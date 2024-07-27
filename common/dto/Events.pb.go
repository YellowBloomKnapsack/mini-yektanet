// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.21.12
// source: Events.proto

package dto

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

type ClickEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PublisherUsername string `protobuf:"bytes,1,opt,name=publisher_username,json=publisherUsername,proto3" json:"publisher_username,omitempty"`
	EventTime         string `protobuf:"bytes,2,opt,name=event_time,json=eventTime,proto3" json:"event_time,omitempty"`
	AdId              uint32 `protobuf:"varint,3,opt,name=ad_id,json=adId,proto3" json:"ad_id,omitempty"`
}

func (x *ClickEvent) Reset() {
	*x = ClickEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Events_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClickEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClickEvent) ProtoMessage() {}

func (x *ClickEvent) ProtoReflect() protoreflect.Message {
	mi := &file_Events_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClickEvent.ProtoReflect.Descriptor instead.
func (*ClickEvent) Descriptor() ([]byte, []int) {
	return file_Events_proto_rawDescGZIP(), []int{0}
}

func (x *ClickEvent) GetPublisherUsername() string {
	if x != nil {
		return x.PublisherUsername
	}
	return ""
}

func (x *ClickEvent) GetEventTime() string {
	if x != nil {
		return x.EventTime
	}
	return ""
}

func (x *ClickEvent) GetAdId() uint32 {
	if x != nil {
		return x.AdId
	}
	return 0
}

type ImpressionEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PublisherUsername string `protobuf:"bytes,1,opt,name=publisher_username,json=publisherUsername,proto3" json:"publisher_username,omitempty"`
	EventTime         string `protobuf:"bytes,2,opt,name=event_time,json=eventTime,proto3" json:"event_time,omitempty"`
	AdId              uint32 `protobuf:"varint,3,opt,name=ad_id,json=adId,proto3" json:"ad_id,omitempty"`
}

func (x *ImpressionEvent) Reset() {
	*x = ImpressionEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Events_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImpressionEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImpressionEvent) ProtoMessage() {}

func (x *ImpressionEvent) ProtoReflect() protoreflect.Message {
	mi := &file_Events_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImpressionEvent.ProtoReflect.Descriptor instead.
func (*ImpressionEvent) Descriptor() ([]byte, []int) {
	return file_Events_proto_rawDescGZIP(), []int{1}
}

func (x *ImpressionEvent) GetPublisherUsername() string {
	if x != nil {
		return x.PublisherUsername
	}
	return ""
}

func (x *ImpressionEvent) GetEventTime() string {
	if x != nil {
		return x.EventTime
	}
	return ""
}

func (x *ImpressionEvent) GetAdId() uint32 {
	if x != nil {
		return x.AdId
	}
	return 0
}

var File_Events_proto protoreflect.FileDescriptor

var file_Events_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03,
	0x64, 0x74, 0x6f, 0x22, 0x6f, 0x0a, 0x0a, 0x43, 0x6c, 0x69, 0x63, 0x6b, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x12, 0x2d, 0x0a, 0x12, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x72, 0x5f, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x70,
	0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x13, 0x0a, 0x05, 0x61, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04,
	0x61, 0x64, 0x49, 0x64, 0x22, 0x74, 0x0a, 0x0f, 0x49, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x2d, 0x0a, 0x12, 0x70, 0x75, 0x62, 0x6c, 0x69,
	0x73, 0x68, 0x65, 0x72, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x11, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x72, 0x55, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x13, 0x0a, 0x05, 0x61, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x61, 0x64, 0x49, 0x64, 0x42, 0x2e, 0x5a, 0x2c, 0x59, 0x65,
	0x6c, 0x6c, 0x6f, 0x77, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x4b, 0x6e, 0x61, 0x70, 0x73, 0x61, 0x63,
	0x6b, 0x2f, 0x6d, 0x69, 0x6e, 0x69, 0x2d, 0x79, 0x65, 0x6b, 0x74, 0x61, 0x6e, 0x65, 0x74, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x64, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_Events_proto_rawDescOnce sync.Once
	file_Events_proto_rawDescData = file_Events_proto_rawDesc
)

func file_Events_proto_rawDescGZIP() []byte {
	file_Events_proto_rawDescOnce.Do(func() {
		file_Events_proto_rawDescData = protoimpl.X.CompressGZIP(file_Events_proto_rawDescData)
	})
	return file_Events_proto_rawDescData
}

var file_Events_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_Events_proto_goTypes = []any{
	(*ClickEvent)(nil),      // 0: dto.ClickEvent
	(*ImpressionEvent)(nil), // 1: dto.ImpressionEvent
}
var file_Events_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_Events_proto_init() }
func file_Events_proto_init() {
	if File_Events_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_Events_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ClickEvent); i {
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
		file_Events_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ImpressionEvent); i {
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
			RawDescriptor: file_Events_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_Events_proto_goTypes,
		DependencyIndexes: file_Events_proto_depIdxs,
		MessageInfos:      file_Events_proto_msgTypes,
	}.Build()
	File_Events_proto = out.File
	file_Events_proto_rawDesc = nil
	file_Events_proto_goTypes = nil
	file_Events_proto_depIdxs = nil
}
