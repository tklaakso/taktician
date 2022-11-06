// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: tak/proto/corpus_entry.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CorpusEntry_InTak int32

const (
	CorpusEntry_UNSET      CorpusEntry_InTak = 0
	CorpusEntry_NOT_IN_TAK CorpusEntry_InTak = 1
	CorpusEntry_IN_TAK     CorpusEntry_InTak = 2
)

// Enum value maps for CorpusEntry_InTak.
var (
	CorpusEntry_InTak_name = map[int32]string{
		0: "UNSET",
		1: "NOT_IN_TAK",
		2: "IN_TAK",
	}
	CorpusEntry_InTak_value = map[string]int32{
		"UNSET":      0,
		"NOT_IN_TAK": 1,
		"IN_TAK":     2,
	}
)

func (x CorpusEntry_InTak) Enum() *CorpusEntry_InTak {
	p := new(CorpusEntry_InTak)
	*p = x
	return p
}

func (x CorpusEntry_InTak) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CorpusEntry_InTak) Descriptor() protoreflect.EnumDescriptor {
	return file_tak_proto_corpus_entry_proto_enumTypes[0].Descriptor()
}

func (CorpusEntry_InTak) Type() protoreflect.EnumType {
	return &file_tak_proto_corpus_entry_proto_enumTypes[0]
}

func (x CorpusEntry_InTak) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CorpusEntry_InTak.Descriptor instead.
func (CorpusEntry_InTak) EnumDescriptor() ([]byte, []int) {
	return file_tak_proto_corpus_entry_proto_rawDescGZIP(), []int{0, 0}
}

type CorpusEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Day      string            `protobuf:"bytes,1,opt,name=day,proto3" json:"day,omitempty"`
	Id       int32             `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Ply      int32             `protobuf:"varint,3,opt,name=ply,proto3" json:"ply,omitempty"`
	Tps      string            `protobuf:"bytes,4,opt,name=tps,proto3" json:"tps,omitempty"`
	Move     string            `protobuf:"bytes,5,opt,name=move,proto3" json:"move,omitempty"`
	Value    float32           `protobuf:"fixed32,6,opt,name=value,proto3" json:"value,omitempty"`
	Plies    int32             `protobuf:"varint,7,opt,name=plies,proto3" json:"plies,omitempty"`
	Features []int64           `protobuf:"varint,8,rep,packed,name=features,proto3" json:"features,omitempty"`
	InTak    CorpusEntry_InTak `protobuf:"varint,9,opt,name=in_tak,json=inTak,proto3,enum=tak.proto.CorpusEntry_InTak" json:"in_tak,omitempty"`
}

func (x *CorpusEntry) Reset() {
	*x = CorpusEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tak_proto_corpus_entry_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CorpusEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CorpusEntry) ProtoMessage() {}

func (x *CorpusEntry) ProtoReflect() protoreflect.Message {
	mi := &file_tak_proto_corpus_entry_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CorpusEntry.ProtoReflect.Descriptor instead.
func (*CorpusEntry) Descriptor() ([]byte, []int) {
	return file_tak_proto_corpus_entry_proto_rawDescGZIP(), []int{0}
}

func (x *CorpusEntry) GetDay() string {
	if x != nil {
		return x.Day
	}
	return ""
}

func (x *CorpusEntry) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CorpusEntry) GetPly() int32 {
	if x != nil {
		return x.Ply
	}
	return 0
}

func (x *CorpusEntry) GetTps() string {
	if x != nil {
		return x.Tps
	}
	return ""
}

func (x *CorpusEntry) GetMove() string {
	if x != nil {
		return x.Move
	}
	return ""
}

func (x *CorpusEntry) GetValue() float32 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *CorpusEntry) GetPlies() int32 {
	if x != nil {
		return x.Plies
	}
	return 0
}

func (x *CorpusEntry) GetFeatures() []int64 {
	if x != nil {
		return x.Features
	}
	return nil
}

func (x *CorpusEntry) GetInTak() CorpusEntry_InTak {
	if x != nil {
		return x.InTak
	}
	return CorpusEntry_UNSET
}

var File_tak_proto_corpus_entry_proto protoreflect.FileDescriptor

var file_tak_proto_corpus_entry_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x74, 0x61, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72, 0x70,
	0x75, 0x73, 0x5f, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09,
	0x74, 0x61, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x94, 0x02, 0x0a, 0x0b, 0x43, 0x6f,
	0x72, 0x70, 0x75, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x61, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x64, 0x61, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x70,
	0x6c, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x70, 0x6c, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x74, 0x70, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x70, 0x73, 0x12,
	0x12, 0x0a, 0x04, 0x6d, 0x6f, 0x76, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6d,
	0x6f, 0x76, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x02, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x6c, 0x69,
	0x65, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x70, 0x6c, 0x69, 0x65, 0x73, 0x12,
	0x1a, 0x0a, 0x08, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28,
	0x03, 0x52, 0x08, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x12, 0x33, 0x0a, 0x06, 0x69,
	0x6e, 0x5f, 0x74, 0x61, 0x6b, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x74, 0x61,
	0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x72, 0x70, 0x75, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x2e, 0x49, 0x6e, 0x54, 0x61, 0x6b, 0x52, 0x05, 0x69, 0x6e, 0x54, 0x61, 0x6b,
	0x22, 0x2e, 0x0a, 0x05, 0x49, 0x6e, 0x54, 0x61, 0x6b, 0x12, 0x09, 0x0a, 0x05, 0x55, 0x4e, 0x53,
	0x45, 0x54, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x4e, 0x4f, 0x54, 0x5f, 0x49, 0x4e, 0x5f, 0x54,
	0x41, 0x4b, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x49, 0x4e, 0x5f, 0x54, 0x41, 0x4b, 0x10, 0x02,
	0x42, 0x21, 0x5a, 0x1f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e,
	0x65, 0x6c, 0x68, 0x61, 0x67, 0x65, 0x2f, 0x74, 0x61, 0x6b, 0x74, 0x69, 0x63, 0x69, 0x61, 0x6e,
	0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tak_proto_corpus_entry_proto_rawDescOnce sync.Once
	file_tak_proto_corpus_entry_proto_rawDescData = file_tak_proto_corpus_entry_proto_rawDesc
)

func file_tak_proto_corpus_entry_proto_rawDescGZIP() []byte {
	file_tak_proto_corpus_entry_proto_rawDescOnce.Do(func() {
		file_tak_proto_corpus_entry_proto_rawDescData = protoimpl.X.CompressGZIP(file_tak_proto_corpus_entry_proto_rawDescData)
	})
	return file_tak_proto_corpus_entry_proto_rawDescData
}

var file_tak_proto_corpus_entry_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_tak_proto_corpus_entry_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_tak_proto_corpus_entry_proto_goTypes = []interface{}{
	(CorpusEntry_InTak)(0), // 0: tak.proto.CorpusEntry.InTak
	(*CorpusEntry)(nil),    // 1: tak.proto.CorpusEntry
}
var file_tak_proto_corpus_entry_proto_depIdxs = []int32{
	0, // 0: tak.proto.CorpusEntry.in_tak:type_name -> tak.proto.CorpusEntry.InTak
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_tak_proto_corpus_entry_proto_init() }
func file_tak_proto_corpus_entry_proto_init() {
	if File_tak_proto_corpus_entry_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tak_proto_corpus_entry_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CorpusEntry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_tak_proto_corpus_entry_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_tak_proto_corpus_entry_proto_goTypes,
		DependencyIndexes: file_tak_proto_corpus_entry_proto_depIdxs,
		EnumInfos:         file_tak_proto_corpus_entry_proto_enumTypes,
		MessageInfos:      file_tak_proto_corpus_entry_proto_msgTypes,
	}.Build()
	File_tak_proto_corpus_entry_proto = out.File
	file_tak_proto_corpus_entry_proto_rawDesc = nil
	file_tak_proto_corpus_entry_proto_goTypes = nil
	file_tak_proto_corpus_entry_proto_depIdxs = nil
}
