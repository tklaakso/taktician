import typing as T  # noqa
import time
import os
import itertools

import torch
import queue
from torch import multiprocessing

import grpc
import asyncio

import wandb

import xformer
from xformer import loading
import yaml

from tak.proto import analysis_pb2_grpc
import tak.model.server
from tak.model import batches, losses
from tak import self_play

from attrs import field, define
import attrs

from .. import Config
from . import data, stats


@define
class TrainingRun:
    model: xformer.Transformer
    config: Config

    wandb: T.Optional["wandb.Run"] = field(default=None, init=False)

    loop: asyncio.BaseEventLoop = field(init=False)
    server: grpc.aio.Server = field(init=False)
    tasks: list[asyncio.Task] = field(init=False, factory=list)
    replay_buffer: list[dict[str, torch.Tensor]] = field(init=False, factory=list)
    train_params: dict[str, torch.Tensor] = field(init=False)
    opt: torch.optim.AdamW = field(init=False)

    elapsed: stats.Elapsed = field(factory=stats.Elapsed, init=False)

    step_start: float = field(init=False)

    def run(self):
        asyncio.run(self.run_async())

    def serve_mode(self):
        self.train_params = {k: v.cpu() for (k, v) in self.model.state_dict().items()}
        self.model.to(device=self.config.device, dtype=self.config.serve_dtype)

    def train_mode(self):
        self.model.to(self.config.train_dtype).load_state_dict(self.train_params)

    def check_and_clear_save_request(self) -> bool:
        run_dir = self.config.run_dir
        if not run_dir:
            return False
        flagpath = os.path.join(run_dir, "SAVE_NOW")
        if os.path.exists(flagpath):
            os.unlink(flagpath)
            return True
        return False

    def train_step(self, batch):
        rollout_time = time.monotonic() - self.step_start

        self.replay_buffer.append(batch)
        if len(self.replay_buffer) > self.config.replay_buffer_steps:
            self.replay_buffer = self.replay_buffer[1:]

        self.train_mode()

        self.elapsed.step += 1

        loss_fn = losses.PolicyValue()
        ds = data.ReplayBufferDataset(
            replay_buffer=self.replay_buffer,
            batch_size=self.config.train_batch,
            device=self.config.device,
        )

        unique = len(
            set(
                tuple(p[m].tolist())
                for (p, m) in zip(batch["positions"], batch["mask"])
            )
        )
        plies = len(batch["positions"])
        stats = {
            "rollout_plies": plies,
            "rollout_games": self.config.rollouts_per_step,
            "rollout_unique_plies": unique,
            "replay_buffer_plies": len(ds.flat_replay_buffer["positions"]),
            "train_step": self.elapsed.step,
            "rollout_time": rollout_time,
        }

        self.elapsed.epoch += 1
        train_start = time.monotonic()

        it = iter(ds)
        for i in range(0, self.config.train_positions, self.config.train_batch):
            try:
                self.elapsed.epoch += 1
                batch = next(it)
            except StopIteration:
                it = iter(ds)
                batch = next(it)

            self.opt.zero_grad()
            out = self.model(batch.inputs, *batch.extra_inputs)
            loss, metrics = loss_fn.loss_and_metrics(batch, out)
            loss.backward()
            self.opt.step()

            self.elapsed.positions += batch.inputs.size(0)

            if self.wandb is not None:
                self.wandb.log(
                    {
                        "train_loss": loss.item(),
                        "train_epoch": self.elapsed.epoch,
                        "positions": self.elapsed.positions,
                    }
                    | stats
                    | metrics
                )

        train_time = time.monotonic() - train_start
        step_time = time.monotonic() - self.step_start
        print(
            f"step={self.elapsed.step}"
            f" games={self.config.rollouts_per_step}"
            f" plies={plies}"
            f" unique={unique}"
            f" rollout_time={rollout_time:0.2f}s"
            f" train_time={train_time:0.2f}s"
            f" step_time={step_time:0.2f}s"
            f" ply/s={plies/(rollout_time):.1f}s"
            f" last_loss={loss.item():0.2f}"
        )

        if self.config.run_dir and (
            self.elapsed.step % self.config.save_freq == 0
            or self.elapsed.step == self.config.train_steps
            or self.check_and_clear_save_request()
        ):
            save_dir = os.path.join(
                self.config.run_dir, f"step_{self.elapsed.step:06d}"
            )
            print(f"Saving snapshot to {save_dir}...")
            self.save_snapshot(save_dir)
            latest_link = os.path.join(self.config.run_dir, "latest")
            try:
                os.unlink(latest_link)
            except FileNotFoundError:
                pass
            os.symlink(
                os.path.basename(save_dir),
                latest_link,
            )

        self.serve_mode()

    def save_snapshot(self, snapshot_path):
        os.makedirs(snapshot_path, exist_ok=True)
        loading.save_model(self.model, snapshot_path)
        torch.save(
            self.opt.state_dict(),
            os.path.join(snapshot_path, "opt.pt"),
        )
        torch.save(
            self.replay_buffer,
            os.path.join(snapshot_path, "replay_buffer.pt"),
        )
        with open(os.path.join(snapshot_path, "elapsed.yaml"), "w") as fh:
            yaml.dump(self.elapsed, fh)

    def should_exit(self):
        return self.elapsed.step >= self.config.train_steps

    def load_or_init_model(self):
        if self.config.run_dir:
            state_dir = os.path.join(self.config.run_dir, "latest")
            if os.path.exists(state_dir):
                loading.load_snapshot(self.model, state_dir)

                self.opt.load_state_dict(
                    torch.load(
                        os.path.join(state_dir, "opt.pt"),
                    )
                )
                self.replay_buffer = torch.load(
                    os.path.join(state_dir, "replay_buffer.pt"),
                )
                with open(os.path.join(state_dir, "elapsed.yaml"), "r") as fh:
                    self.elapsed = yaml.unsafe_load(fh)

                return

        if self.config.load_model:
            loading.load_snapshot(self.model, self.config.load_model)
        else:
            self.model.init_weights()

    async def train_loop(self):
        rollout_engine = self_play.MultiprocessSelfPlayEngine(
            config=self_play.SelfPlayConfig(
                size=self.config.size,
                workers=self.config.rollout_workers,
                resignation_threshold=self.config.rollout_resignation_threshold,
                ply_limit=self.config.rollout_ply_limit,
                engine_factory=self_play.BuildRemoteMCTS(
                    host="localhost",
                    port=self.config.server_port,
                    config=self.config.rollout_config,
                ),
            )
        )

        try:
            while not self.should_exit():
                self.step_start = time.monotonic()
                logs = await self.loop.run_in_executor(
                    None, rollout_engine.play_many, self.config.rollouts_per_step
                )
                batch = self_play.encode_games(logs)
                batch["positions"] = batch["positions"].to(torch.long)
                self.train_step(batch)
        finally:
            rollout_engine.stop()

    async def run_async(self):
        if self.config.wandb:
            self.wandb = wandb.init(
                project=self.config.project,
                name=self.config.job_name,
                id=self.config.job_id,
                resume="allow",
            )
            wandb.config.update(attrs.asdict(self.config), allow_val_change=True)

        self.loop = asyncio.get_event_loop()

        self.opt = torch.optim.AdamW(self.model.parameters(), lr=self.config.lr)

        self.load_or_init_model()
        if self.config.run_dir:
            config_path = os.path.join(self.config.run_dir, "run.yaml")
            os.makedirs(os.path.dirname(config_path), exist_ok=True)
            with open(config_path, "w") as fh:
                yaml.dump(self.config, fh)

        self.serve_mode()

        multiprocessing.set_start_method("spawn")
        os.environ["GRPC_ENABLE_FORK_SUPPORT"] = "0"

        self.server = grpc.aio.server()
        self.server.add_insecure_port(f"localhost:{self.config.server_port}")

        analysis = tak.model.server.Server(model=self.model, device=self.config.device)

        self.tasks.append(asyncio.create_task(analysis.worker_loop()))

        analysis_pb2_grpc.add_AnalysisServicer_to_server(
            analysis,
            self.server,
        )
        await self.server.start()
        train_task = asyncio.create_task(self.train_loop())
        self.tasks.append(train_task)
        try:
            done, pending = await asyncio.wait(
                self.tasks + [self.server.wait_for_termination()],
                return_when=asyncio.FIRST_COMPLETED,
            )
            await self.server.stop(0)
            for task in pending:
                task.cancel()
            for task in done:
                task.result()
        finally:
            if self.wandb:
                self.wandb.finish()
