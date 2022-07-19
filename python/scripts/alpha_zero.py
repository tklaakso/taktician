import argparse
import time
import typing as T  # noqa
import os.path
import yaml

import torch
from torch import multiprocessing

import grpc
from tak.proto import analysis_pb2_grpc
import asyncio

import xformer
from xformer import data, model, train, loading
from xformer.train import hooks, lr_schedules
from tak.model import batches, heads, losses

import tak.model.server
from tak import self_play, mcts
from tak.alphazero import model_process
from tak import alphazero
from xformer import yaml_ext  # noqa


def parse_args():
    parser = argparse.ArgumentParser(description="Train a Tak player using self-play")

    parser.add_argument("--load-model", type=str, help="Initial model to load")
    parser.add_argument("--layers", type=int, default=2, help="Number of layers")
    parser.add_argument("--d_model", type=int, default=None, help="embedding dimension")
    parser.add_argument("--d_head", type=int, default=32, help="head dimension")
    parser.add_argument(
        "--n_ctx", type=int, default=1024, help="maximum context length"
    )
    parser.add_argument(
        "--pe", type=str, default=None, help="positional encoding (sin, learned, none)"
    )

    parser.add_argument("--batch", type=int, default=64, help="batch size")

    parser.add_argument(
        "--device", type=str, choices=("cpu", "cuda"), default="cuda", help="device"
    )

    parser.add_argument("--lr", type=float, default=5e-4, help="learning rate")

    parser.add_argument("--steps", type=int, default=10)

    parser.add_argument("--size", type=int, default=3)

    parser.add_argument("--rollouts-per-step", type=int, default=1000)
    parser.add_argument("--replay-buffer-steps", type=int, default=4)
    parser.add_argument("--train-positions", type=int, default=1024)

    parser.add_argument("--rollout-workers", type=int, default=50)
    parser.add_argument("--rollout-simulations", type=int, default=25)

    parser.add_argument("--save-dir", type=str, metavar="PATH")
    parser.add_argument("--save-freq", type=int, metavar="STEPS", default=10)

    parser.add_argument("--progress", default=True, action="store_true")
    parser.add_argument("--no-progress", dest="progress", action="store_false")

    parser.add_argument("--job-name", type=str, default=None, help="job name for wandb")
    parser.add_argument("--wandb", action="store_true", default=False)
    parser.add_argument("--no-wandb", action="store_false", dest="wandb")

    return parser.parse_args()


def check_and_clear_save_request(save_root) -> bool:
    flagpath = os.path.join(save_root, "SAVE_NOW")
    if os.path.exists(flagpath):
        os.unlink(flagpath)
        return True
    return False


def main():
    multiprocessing.set_start_method("spawn")

    args = parse_args()

    if args.load_model:
        model_cfg = loading.load_config(args.load_model)
    else:
        model_cfg = xformer.Config(
            n_layer=args.layers,
            d_model=args.d_model or 128 * args.layers,
            d_head=args.d_head,
            n_ctx=args.n_ctx,
            n_vocab=256,
            autoregressive_mask=False,
            output_head=heads.PolicyValue,
        )
        if args.pe is not None:
            model_cfg.positional_encoding = args.pe

    config = alphazero.Config(
        device=args.device,
        server_port=5001,
        size=args.size,
        load_model=args.load_model,
        rollout_workers=args.rollout_workers,
        rollouts_per_step=args.rollouts_per_step,
        rollout_resignation_threshold=0.99,
        rollout_ply_limit=20,
        replay_buffer_steps=args.replay_buffer_steps,
        train_batch=args.batch,
        train_positions=args.train_positions,
        lr=args.lr,
        save_path=args.save_dir,
        save_freq=args.save_freq,
        train_steps=args.steps,
        wandb=args.wandb,
        job_name=args.job_name,
    )
    config.rollout_config.simulation_limit = args.rollout_simulations
    config.rollout_config.time_limit = 0

    if config.save_path:
        config_path = os.path.join(config.save_path, "run.yaml")
        os.makedirs(os.path.dirname(config_path), exist_ok=True)
        with open(config_path, "w") as fh:
            yaml.dump(config, fh)

    srv = model_process.create_server(config=config, model_config=model_cfg)
    srv.start()

    rollout_engine = self_play.MultiprocessSelfPlayEngine(
        config=self_play.SelfPlayConfig(
            size=config.size,
            workers=config.rollout_workers,
            resignation_threshold=config.rollout_resignation_threshold,
            ply_limit=config.rollout_ply_limit,
            engine_factory=self_play.BuildRemoteMCTS(
                host="localhost",
                port=config.server_port,
                config=config.rollout_config,
            ),
        )
    )

    for step in range(config.train_steps):
        logs = rollout_engine.play_many(
            config.rollouts_per_step, progress=args.progress
        )

        batch = self_play.encode_games(logs)
        batch["positions"] = batch["positions"].to(torch.long)
        srv.train_step({k: v.share_memory_() for (k, v) in batch.items()})

        if config.save_path and (
            step % config.save_freq == 0
            or step == config.train_steps - 1
            or check_and_clear_save_request(config.save_path)
        ):
            save_dir = os.path.join(config.save_path, f"step_{step:06d}")
            print(f"Saving snapshot to {save_dir}...")
            srv.save_model(save_dir)
            latest_link = os.path.join(config.save_path, "latest")
            try:
                os.unlink(latest_link)
            except FileNotFoundError:
                pass
            os.symlink(
                os.path.basename(save_dir),
                latest_link,
            )

    rollout_engine.stop()
    srv.stop()


if __name__ == "__main__":
    main()
