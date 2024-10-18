// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: apis/recommend.proto

package apis

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

type UserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string            `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Args   map[string]string `protobuf:"bytes,2,rep,name=args,proto3" json:"args,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *UserRequest) Reset() {
	*x = UserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_recommend_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserRequest) ProtoMessage() {}

func (x *UserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_apis_recommend_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserRequest.ProtoReflect.Descriptor instead.
func (*UserRequest) Descriptor() ([]byte, []int) {
	return file_apis_recommend_proto_rawDescGZIP(), []int{0}
}

func (x *UserRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UserRequest) GetArgs() map[string]string {
	if x != nil {
		return x.Args
	}
	return nil
}

type NoteList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TabName string   `protobuf:"bytes,1,opt,name=tabName,proto3" json:"tabName,omitempty"`
	RList   []uint64 `protobuf:"varint,2,rep,packed,name=rList,proto3" json:"rList,omitempty"`
	Size    int32    `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
}

func (x *NoteList) Reset() {
	*x = NoteList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_recommend_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NoteList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoteList) ProtoMessage() {}

func (x *NoteList) ProtoReflect() protoreflect.Message {
	mi := &file_apis_recommend_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoteList.ProtoReflect.Descriptor instead.
func (*NoteList) Descriptor() ([]byte, []int) {
	return file_apis_recommend_proto_rawDescGZIP(), []int{1}
}

func (x *NoteList) GetTabName() string {
	if x != nil {
		return x.TabName
	}
	return ""
}

func (x *NoteList) GetRList() []uint64 {
	if x != nil {
		return x.RList
	}
	return nil
}

func (x *NoteList) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

type NoteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version string    `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	NoteIds *NoteList `protobuf:"bytes,2,opt,name=noteIds,proto3" json:"noteIds,omitempty"`
}

func (x *NoteResponse) Reset() {
	*x = NoteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_recommend_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoteResponse) ProtoMessage() {}

func (x *NoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_apis_recommend_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoteResponse.ProtoReflect.Descriptor instead.
func (*NoteResponse) Descriptor() ([]byte, []int) {
	return file_apis_recommend_proto_rawDescGZIP(), []int{2}
}

func (x *NoteResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *NoteResponse) GetNoteIds() *NoteList {
	if x != nil {
		return x.NoteIds
	}
	return nil
}

var File_apis_recommend_proto protoreflect.FileDescriptor

var file_apis_recommend_proto_rawDesc = []byte{
	0x0a, 0x14, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x72, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x72, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x64, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x22, 0x9b, 0x01, 0x0a, 0x0b, 0x55, 0x73, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x3b, 0x0a, 0x04, 0x61, 0x72, 0x67, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27,
	0x2e, 0x72, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x41, 0x72,
	0x67, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x61, 0x72, 0x67, 0x73, 0x1a, 0x37, 0x0a,
	0x09, 0x41, 0x72, 0x67, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x4e, 0x0a, 0x08, 0x4e, 0x6f, 0x74, 0x65, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x61, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x61, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x72, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x04, 0x52, 0x05, 0x72, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x5e, 0x0a, 0x0c, 0x4e, 0x6f, 0x74, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x34, 0x0a, 0x07, 0x6e, 0x6f, 0x74, 0x65, 0x49, 0x64, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x2e, 0x4e, 0x6f, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x07, 0x6e,
	0x6f, 0x74, 0x65, 0x49, 0x64, 0x73, 0x32, 0x5f, 0x0a, 0x0b, 0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x50, 0x0a, 0x0f, 0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x64, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x1d, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x64, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x6d, 0x6d,
	0x65, 0x6e, 0x64, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x4e, 0x6f, 0x74, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3a, 0x0a, 0x18, 0x69, 0x6f, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x42, 0x0f, 0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x50, 0x01, 0x5a, 0x05, 0x61, 0x70, 0x69, 0x73, 0x2f, 0xa2, 0x02, 0x03,
	0x48, 0x4c, 0x57, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_apis_recommend_proto_rawDescOnce sync.Once
	file_apis_recommend_proto_rawDescData = file_apis_recommend_proto_rawDesc
)

func file_apis_recommend_proto_rawDescGZIP() []byte {
	file_apis_recommend_proto_rawDescOnce.Do(func() {
		file_apis_recommend_proto_rawDescData = protoimpl.X.CompressGZIP(file_apis_recommend_proto_rawDescData)
	})
	return file_apis_recommend_proto_rawDescData
}

var file_apis_recommend_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_apis_recommend_proto_goTypes = []interface{}{
	(*UserRequest)(nil),  // 0: recommend.server.UserRequest
	(*NoteList)(nil),     // 1: recommend.server.NoteList
	(*NoteResponse)(nil), // 2: recommend.server.NoteResponse
	nil,                  // 3: recommend.server.UserRequest.ArgsEntry
}
var file_apis_recommend_proto_depIdxs = []int32{
	3, // 0: recommend.server.UserRequest.args:type_name -> recommend.server.UserRequest.ArgsEntry
	1, // 1: recommend.server.NoteResponse.noteIds:type_name -> recommend.server.NoteList
	0, // 2: recommend.server.Recommender.RecommendServer:input_type -> recommend.server.UserRequest
	2, // 3: recommend.server.Recommender.RecommendServer:output_type -> recommend.server.NoteResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_apis_recommend_proto_init() }
func file_apis_recommend_proto_init() {
	if File_apis_recommend_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_apis_recommend_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserRequest); i {
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
		file_apis_recommend_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NoteList); i {
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
		file_apis_recommend_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NoteResponse); i {
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
			RawDescriptor: file_apis_recommend_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_apis_recommend_proto_goTypes,
		DependencyIndexes: file_apis_recommend_proto_depIdxs,
		MessageInfos:      file_apis_recommend_proto_msgTypes,
	}.Build()
	File_apis_recommend_proto = out.File
	file_apis_recommend_proto_rawDesc = nil
	file_apis_recommend_proto_goTypes = nil
	file_apis_recommend_proto_depIdxs = nil
}
