// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: cache/v1/cache.proto

package v1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetDataRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDataRequest) Reset() {
	*x = GetDataRequest{}
	mi := &file_cache_v1_cache_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDataRequest) ProtoMessage() {}

func (x *GetDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDataRequest.ProtoReflect.Descriptor instead.
func (*GetDataRequest) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{0}
}

func (x *GetDataRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type GetDataReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value         string                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Source        string                 `protobuf:"bytes,3,opt,name=source,proto3" json:"source,omitempty"` // "cache" or "database"
	Pod           string                 `protobuf:"bytes,4,opt,name=pod,proto3" json:"pod,omitempty"`       // 容器IP地址，用于观察负载均衡
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDataReply) Reset() {
	*x = GetDataReply{}
	mi := &file_cache_v1_cache_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDataReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDataReply) ProtoMessage() {}

func (x *GetDataReply) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDataReply.ProtoReflect.Descriptor instead.
func (*GetDataReply) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{1}
}

func (x *GetDataReply) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *GetDataReply) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *GetDataReply) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *GetDataReply) GetPod() string {
	if x != nil {
		return x.Pod
	}
	return ""
}

type SetDataRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value         string                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SetDataRequest) Reset() {
	*x = SetDataRequest{}
	mi := &file_cache_v1_cache_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetDataRequest) ProtoMessage() {}

func (x *SetDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetDataRequest.ProtoReflect.Descriptor instead.
func (*SetDataRequest) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{2}
}

func (x *SetDataRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SetDataRequest) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type SetDataReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SetDataReply) Reset() {
	*x = SetDataReply{}
	mi := &file_cache_v1_cache_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetDataReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetDataReply) ProtoMessage() {}

func (x *SetDataReply) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetDataReply.ProtoReflect.Descriptor instead.
func (*SetDataReply) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{3}
}

func (x *SetDataReply) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *SetDataReply) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type HealthCheckRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HealthCheckRequest) Reset() {
	*x = HealthCheckRequest{}
	mi := &file_cache_v1_cache_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HealthCheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthCheckRequest) ProtoMessage() {}

func (x *HealthCheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthCheckRequest.ProtoReflect.Descriptor instead.
func (*HealthCheckRequest) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{4}
}

type HealthCheckReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        string                 `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Timestamp     string                 `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Version       string                 `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HealthCheckReply) Reset() {
	*x = HealthCheckReply{}
	mi := &file_cache_v1_cache_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HealthCheckReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthCheckReply) ProtoMessage() {}

func (x *HealthCheckReply) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthCheckReply.ProtoReflect.Descriptor instead.
func (*HealthCheckReply) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{5}
}

func (x *HealthCheckReply) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *HealthCheckReply) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *HealthCheckReply) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

var File_cache_v1_cache_proto protoreflect.FileDescriptor

const file_cache_v1_cache_proto_rawDesc = "" +
	"\n" +
	"\x14cache/v1/cache.proto\x12\fapi.cache.v1\x1a\x1cgoogle/api/annotations.proto\"\"\n" +
	"\x0eGetDataRequest\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\"`\n" +
	"\fGetDataReply\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value\x12\x16\n" +
	"\x06source\x18\x03 \x01(\tR\x06source\x12\x10\n" +
	"\x03pod\x18\x04 \x01(\tR\x03pod\"8\n" +
	"\x0eSetDataRequest\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value\"B\n" +
	"\fSetDataReply\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\"\x14\n" +
	"\x12HealthCheckRequest\"b\n" +
	"\x10HealthCheckReply\x12\x16\n" +
	"\x06status\x18\x01 \x01(\tR\x06status\x12\x1c\n" +
	"\ttimestamp\x18\x02 \x01(\tR\ttimestamp\x12\x18\n" +
	"\aversion\x18\x03 \x01(\tR\aversion2\xaf\x02\n" +
	"\fCacheService\x12_\n" +
	"\aGetData\x12\x1c.api.cache.v1.GetDataRequest\x1a\x1a.api.cache.v1.GetDataReply\"\x1a\x82\xd3\xe4\x93\x02\x14\x12\x12/api/v1/data/{key}\x12\\\n" +
	"\aSetData\x12\x1c.api.cache.v1.SetDataRequest\x1a\x1a.api.cache.v1.SetDataReply\"\x17\x82\xd3\xe4\x93\x02\x11:\x01*\"\f/api/v1/data\x12`\n" +
	"\vHealthCheck\x12 .api.cache.v1.HealthCheckRequest\x1a\x1e.api.cache.v1.HealthCheckReply\"\x0f\x82\xd3\xe4\x93\x02\t\x12\a/healthB\x1fZ\x1dcache-service/api/cache/v1;v1b\x06proto3"

var (
	file_cache_v1_cache_proto_rawDescOnce sync.Once
	file_cache_v1_cache_proto_rawDescData []byte
)

func file_cache_v1_cache_proto_rawDescGZIP() []byte {
	file_cache_v1_cache_proto_rawDescOnce.Do(func() {
		file_cache_v1_cache_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_cache_v1_cache_proto_rawDesc), len(file_cache_v1_cache_proto_rawDesc)))
	})
	return file_cache_v1_cache_proto_rawDescData
}

var file_cache_v1_cache_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_cache_v1_cache_proto_goTypes = []any{
	(*GetDataRequest)(nil),     // 0: api.cache.v1.GetDataRequest
	(*GetDataReply)(nil),       // 1: api.cache.v1.GetDataReply
	(*SetDataRequest)(nil),     // 2: api.cache.v1.SetDataRequest
	(*SetDataReply)(nil),       // 3: api.cache.v1.SetDataReply
	(*HealthCheckRequest)(nil), // 4: api.cache.v1.HealthCheckRequest
	(*HealthCheckReply)(nil),   // 5: api.cache.v1.HealthCheckReply
}
var file_cache_v1_cache_proto_depIdxs = []int32{
	0, // 0: api.cache.v1.CacheService.GetData:input_type -> api.cache.v1.GetDataRequest
	2, // 1: api.cache.v1.CacheService.SetData:input_type -> api.cache.v1.SetDataRequest
	4, // 2: api.cache.v1.CacheService.HealthCheck:input_type -> api.cache.v1.HealthCheckRequest
	1, // 3: api.cache.v1.CacheService.GetData:output_type -> api.cache.v1.GetDataReply
	3, // 4: api.cache.v1.CacheService.SetData:output_type -> api.cache.v1.SetDataReply
	5, // 5: api.cache.v1.CacheService.HealthCheck:output_type -> api.cache.v1.HealthCheckReply
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_cache_v1_cache_proto_init() }
func file_cache_v1_cache_proto_init() {
	if File_cache_v1_cache_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_cache_v1_cache_proto_rawDesc), len(file_cache_v1_cache_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_cache_v1_cache_proto_goTypes,
		DependencyIndexes: file_cache_v1_cache_proto_depIdxs,
		MessageInfos:      file_cache_v1_cache_proto_msgTypes,
	}.Build()
	File_cache_v1_cache_proto = out.File
	file_cache_v1_cache_proto_goTypes = nil
	file_cache_v1_cache_proto_depIdxs = nil
}
