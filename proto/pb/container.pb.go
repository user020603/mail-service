// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.12.4
// source: proto/container.proto

package pb

import (
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

type GetContainerInformationRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	StartTime     int64                  `protobuf:"varint,1,opt,name=startTime,proto3" json:"startTime,omitempty"`
	EndTime       int64                  `protobuf:"varint,2,opt,name=endTime,proto3" json:"endTime,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetContainerInformationRequest) Reset() {
	*x = GetContainerInformationRequest{}
	mi := &file_proto_container_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetContainerInformationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetContainerInformationRequest) ProtoMessage() {}

func (x *GetContainerInformationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_container_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetContainerInformationRequest.ProtoReflect.Descriptor instead.
func (*GetContainerInformationRequest) Descriptor() ([]byte, []int) {
	return file_proto_container_proto_rawDescGZIP(), []int{0}
}

func (x *GetContainerInformationRequest) GetStartTime() int64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *GetContainerInformationRequest) GetEndTime() int64 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

type GetContainerInformationResponse struct {
	state                protoimpl.MessageState `protogen:"open.v1"`
	NumContainers        int64                  `protobuf:"varint,1,opt,name=numContainers,proto3" json:"numContainers,omitempty"`
	NumRunningContainers int64                  `protobuf:"varint,2,opt,name=numRunningContainers,proto3" json:"numRunningContainers,omitempty"`
	NumStoppedContainers int64                  `protobuf:"varint,3,opt,name=numStoppedContainers,proto3" json:"numStoppedContainers,omitempty"`
	MeanUptimeRatio      float32                `protobuf:"fixed32,4,opt,name=meanUptimeRatio,proto3" json:"meanUptimeRatio,omitempty"`
	unknownFields        protoimpl.UnknownFields
	sizeCache            protoimpl.SizeCache
}

func (x *GetContainerInformationResponse) Reset() {
	*x = GetContainerInformationResponse{}
	mi := &file_proto_container_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetContainerInformationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetContainerInformationResponse) ProtoMessage() {}

func (x *GetContainerInformationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_container_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetContainerInformationResponse.ProtoReflect.Descriptor instead.
func (*GetContainerInformationResponse) Descriptor() ([]byte, []int) {
	return file_proto_container_proto_rawDescGZIP(), []int{1}
}

func (x *GetContainerInformationResponse) GetNumContainers() int64 {
	if x != nil {
		return x.NumContainers
	}
	return 0
}

func (x *GetContainerInformationResponse) GetNumRunningContainers() int64 {
	if x != nil {
		return x.NumRunningContainers
	}
	return 0
}

func (x *GetContainerInformationResponse) GetNumStoppedContainers() int64 {
	if x != nil {
		return x.NumStoppedContainers
	}
	return 0
}

func (x *GetContainerInformationResponse) GetMeanUptimeRatio() float32 {
	if x != nil {
		return x.MeanUptimeRatio
	}
	return 0
}

var File_proto_container_proto protoreflect.FileDescriptor

const file_proto_container_proto_rawDesc = "" +
	"\n" +
	"\x15proto/container.proto\x12\x15container_adm_service\"X\n" +
	"\x1eGetContainerInformationRequest\x12\x1c\n" +
	"\tstartTime\x18\x01 \x01(\x03R\tstartTime\x12\x18\n" +
	"\aendTime\x18\x02 \x01(\x03R\aendTime\"\xd9\x01\n" +
	"\x1fGetContainerInformationResponse\x12$\n" +
	"\rnumContainers\x18\x01 \x01(\x03R\rnumContainers\x122\n" +
	"\x14numRunningContainers\x18\x02 \x01(\x03R\x14numRunningContainers\x122\n" +
	"\x14numStoppedContainers\x18\x03 \x01(\x03R\x14numStoppedContainers\x12(\n" +
	"\x0fmeanUptimeRatio\x18\x04 \x01(\x02R\x0fmeanUptimeRatio2\xa0\x01\n" +
	"\x13ContainerAdmService\x12\x88\x01\n" +
	"\x17GetContainerInformation\x125.container_adm_service.GetContainerInformationRequest\x1a6.container_adm_service.GetContainerInformationResponseB\fZ\n" +
	"./proto/pbb\x06proto3"

var (
	file_proto_container_proto_rawDescOnce sync.Once
	file_proto_container_proto_rawDescData []byte
)

func file_proto_container_proto_rawDescGZIP() []byte {
	file_proto_container_proto_rawDescOnce.Do(func() {
		file_proto_container_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_container_proto_rawDesc), len(file_proto_container_proto_rawDesc)))
	})
	return file_proto_container_proto_rawDescData
}

var file_proto_container_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_container_proto_goTypes = []any{
	(*GetContainerInformationRequest)(nil),  // 0: container_adm_service.GetContainerInformationRequest
	(*GetContainerInformationResponse)(nil), // 1: container_adm_service.GetContainerInformationResponse
}
var file_proto_container_proto_depIdxs = []int32{
	0, // 0: container_adm_service.ContainerAdmService.GetContainerInformation:input_type -> container_adm_service.GetContainerInformationRequest
	1, // 1: container_adm_service.ContainerAdmService.GetContainerInformation:output_type -> container_adm_service.GetContainerInformationResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_container_proto_init() }
func file_proto_container_proto_init() {
	if File_proto_container_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_container_proto_rawDesc), len(file_proto_container_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_container_proto_goTypes,
		DependencyIndexes: file_proto_container_proto_depIdxs,
		MessageInfos:      file_proto_container_proto_msgTypes,
	}.Build()
	File_proto_container_proto = out.File
	file_proto_container_proto_goTypes = nil
	file_proto_container_proto_depIdxs = nil
}
