// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.27.0--rc3
// source: proto/server.proto

package server

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type CodeEnum int32

const (
	CodeEnum_Success     CodeEnum = 0 // success
	CodeEnum_ParamErr    CodeEnum = 1 // parameter error
	CodeEnum_InternalErr CodeEnum = 2 // internal error
	CodeEnum_UnknowErr   CodeEnum = 3 // unknown error
	CodeEnum_NoDataFound CodeEnum = 4 // No data found
)

// Enum value maps for CodeEnum.
var (
	CodeEnum_name = map[int32]string{
		0: "Success",
		1: "ParamErr",
		2: "InternalErr",
		3: "UnknowErr",
		4: "NoDataFound",
	}
	CodeEnum_value = map[string]int32{
		"Success":     0,
		"ParamErr":    1,
		"InternalErr": 2,
		"UnknowErr":   3,
		"NoDataFound": 4,
	}
)

func (x CodeEnum) Enum() *CodeEnum {
	p := new(CodeEnum)
	*p = x
	return p
}

func (x CodeEnum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CodeEnum) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_server_proto_enumTypes[0].Descriptor()
}

func (CodeEnum) Type() protoreflect.EnumType {
	return &file_proto_server_proto_enumTypes[0]
}

func (x CodeEnum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CodeEnum.Descriptor instead.
func (CodeEnum) EnumDescriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{0}
}

type Book struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title  string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Author string `protobuf:"bytes,3,opt,name=author,proto3" json:"author,omitempty"`
}

func (x *Book) Reset() {
	*x = Book{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Book) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Book) ProtoMessage() {}

func (x *Book) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Book.ProtoReflect.Descriptor instead.
func (*Book) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{0}
}

func (x *Book) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Book) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Book) GetAuthor() string {
	if x != nil {
		return x.Author
	}
	return ""
}

type AddBookRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title  string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Author string `protobuf:"bytes,2,opt,name=author,proto3" json:"author,omitempty"`
}

func (x *AddBookRequest) Reset() {
	*x = AddBookRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddBookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddBookRequest) ProtoMessage() {}

func (x *AddBookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddBookRequest.ProtoReflect.Descriptor instead.
func (*AddBookRequest) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{1}
}

func (x *AddBookRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *AddBookRequest) GetAuthor() string {
	if x != nil {
		return x.Author
	}
	return ""
}

type AddBookReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code CodeEnum `protobuf:"varint,1,opt,name=code,proto3,enum=CodeEnum" json:"code,omitempty"`
	Msg  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Data *Book    `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *AddBookReply) Reset() {
	*x = AddBookReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddBookReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddBookReply) ProtoMessage() {}

func (x *AddBookReply) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddBookReply.ProtoReflect.Descriptor instead.
func (*AddBookReply) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{2}
}

func (x *AddBookReply) GetCode() CodeEnum {
	if x != nil {
		return x.Code
	}
	return CodeEnum_Success
}

func (x *AddBookReply) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *AddBookReply) GetData() *Book {
	if x != nil {
		return x.Data
	}
	return nil
}

type DeleteBookRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteBookRequest) Reset() {
	*x = DeleteBookRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteBookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteBookRequest) ProtoMessage() {}

func (x *DeleteBookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteBookRequest.ProtoReflect.Descriptor instead.
func (*DeleteBookRequest) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteBookRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type DeleteBookReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code CodeEnum `protobuf:"varint,1,opt,name=code,proto3,enum=CodeEnum" json:"code,omitempty"`
	Msg  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *DeleteBookReply) Reset() {
	*x = DeleteBookReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteBookReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteBookReply) ProtoMessage() {}

func (x *DeleteBookReply) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteBookReply.ProtoReflect.Descriptor instead.
func (*DeleteBookReply) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteBookReply) GetCode() CodeEnum {
	if x != nil {
		return x.Code
	}
	return CodeEnum_Success
}

func (x *DeleteBookReply) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type UpdateBookRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title  string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Author string `protobuf:"bytes,3,opt,name=author,proto3" json:"author,omitempty"`
}

func (x *UpdateBookRequest) Reset() {
	*x = UpdateBookRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateBookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateBookRequest) ProtoMessage() {}

func (x *UpdateBookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateBookRequest.ProtoReflect.Descriptor instead.
func (*UpdateBookRequest) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateBookRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateBookRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *UpdateBookRequest) GetAuthor() string {
	if x != nil {
		return x.Author
	}
	return ""
}

type UpdateBookReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code CodeEnum `protobuf:"varint,1,opt,name=code,proto3,enum=CodeEnum" json:"code,omitempty"`
	Msg  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Data *Book    `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *UpdateBookReply) Reset() {
	*x = UpdateBookReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateBookReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateBookReply) ProtoMessage() {}

func (x *UpdateBookReply) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateBookReply.ProtoReflect.Descriptor instead.
func (*UpdateBookReply) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateBookReply) GetCode() CodeEnum {
	if x != nil {
		return x.Code
	}
	return CodeEnum_Success
}

func (x *UpdateBookReply) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *UpdateBookReply) GetData() *Book {
	if x != nil {
		return x.Data
	}
	return nil
}

type GetBookRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetBookRequest) Reset() {
	*x = GetBookRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBookRequest) ProtoMessage() {}

func (x *GetBookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBookRequest.ProtoReflect.Descriptor instead.
func (*GetBookRequest) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{7}
}

func (x *GetBookRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetBookReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code CodeEnum `protobuf:"varint,1,opt,name=code,proto3,enum=CodeEnum" json:"code,omitempty"`
	Msg  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Data *Book    `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *GetBookReply) Reset() {
	*x = GetBookReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBookReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBookReply) ProtoMessage() {}

func (x *GetBookReply) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBookReply.ProtoReflect.Descriptor instead.
func (*GetBookReply) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{8}
}

func (x *GetBookReply) GetCode() CodeEnum {
	if x != nil {
		return x.Code
	}
	return CodeEnum_Success
}

func (x *GetBookReply) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *GetBookReply) GetData() *Book {
	if x != nil {
		return x.Data
	}
	return nil
}

type ListBookRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListBookRequest) Reset() {
	*x = ListBookRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListBookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListBookRequest) ProtoMessage() {}

func (x *ListBookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListBookRequest.ProtoReflect.Descriptor instead.
func (*ListBookRequest) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{9}
}

type ListBookReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code CodeEnum `protobuf:"varint,1,opt,name=code,proto3,enum=CodeEnum" json:"code,omitempty"`
	Msg  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	List []*Book  `protobuf:"bytes,3,rep,name=list,proto3" json:"list,omitempty"`
}

func (x *ListBookReply) Reset() {
	*x = ListBookReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListBookReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListBookReply) ProtoMessage() {}

func (x *ListBookReply) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListBookReply.ProtoReflect.Descriptor instead.
func (*ListBookReply) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{10}
}

func (x *ListBookReply) GetCode() CodeEnum {
	if x != nil {
		return x.Code
	}
	return CodeEnum_Success
}

func (x *ListBookReply) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *ListBookReply) GetList() []*Book {
	if x != nil {
		return x.List
	}
	return nil
}

var File_proto_server_proto protoreflect.FileDescriptor

var file_proto_server_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x44, 0x0a, 0x04, 0x42, 0x6f, 0x6f, 0x6b, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x22, 0x3e, 0x0a, 0x0e, 0x41, 0x64,
	0x64, 0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x22, 0x5a, 0x0a, 0x0c, 0x41, 0x64,
	0x64, 0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1d, 0x0a, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x09, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x45,
	0x6e, 0x75, 0x6d, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x19, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x42, 0x6f, 0x6f, 0x6b,
	0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x23, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x22, 0x42, 0x0a, 0x0f, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1d,
	0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x09, 0x2e, 0x43,
	0x6f, 0x64, 0x65, 0x45, 0x6e, 0x75, 0x6d, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22,
	0x51, 0x0a, 0x11, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x22, 0x5d, 0x0a, 0x0f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x6f, 0x6f, 0x6b,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1d, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x09, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x45, 0x6e, 0x75, 0x6d, 0x52, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x19, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x22, 0x20, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x5a, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x1d, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x09, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x45, 0x6e, 0x75, 0x6d, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6d, 0x73, 0x67, 0x12, 0x19, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x05, 0x2e, 0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0x11, 0x0a, 0x0f, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x5b, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x1d, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x09, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x45, 0x6e, 0x75, 0x6d, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6d, 0x73, 0x67, 0x12, 0x19, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x05, 0x2e, 0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x2a,
	0x56, 0x0a, 0x08, 0x43, 0x6f, 0x64, 0x65, 0x45, 0x6e, 0x75, 0x6d, 0x12, 0x0b, 0x0a, 0x07, 0x53,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x45, 0x72, 0x72, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x45, 0x72, 0x72, 0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x55, 0x6e, 0x6b, 0x6e, 0x6f,
	0x77, 0x45, 0x72, 0x72, 0x10, 0x03, 0x12, 0x0f, 0x0a, 0x0b, 0x4e, 0x6f, 0x44, 0x61, 0x74, 0x61,
	0x46, 0x6f, 0x75, 0x6e, 0x64, 0x10, 0x04, 0x32, 0x81, 0x02, 0x0a, 0x09, 0x42, 0x6f, 0x6f, 0x6b,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x2b, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x42, 0x6f, 0x6f, 0x6b,
	0x12, 0x0f, 0x2e, 0x41, 0x64, 0x64, 0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x0d, 0x2e, 0x41, 0x64, 0x64, 0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x22, 0x00, 0x12, 0x34, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x42, 0x6f, 0x6f, 0x6b,
	0x12, 0x12, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x42, 0x6f, 0x6f,
	0x6b, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x34, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x42, 0x6f, 0x6f, 0x6b, 0x12, 0x12, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42,
	0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x2b,
	0x0a, 0x07, 0x47, 0x65, 0x74, 0x42, 0x6f, 0x6f, 0x6b, 0x12, 0x0f, 0x2e, 0x47, 0x65, 0x74, 0x42,
	0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x47, 0x65, 0x74,
	0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x2e, 0x0a, 0x08, 0x4c,
	0x69, 0x73, 0x74, 0x42, 0x6f, 0x6f, 0x6b, 0x12, 0x10, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x6f,
	0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x42, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x2e,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_server_proto_rawDescOnce sync.Once
	file_proto_server_proto_rawDescData = file_proto_server_proto_rawDesc
)

func file_proto_server_proto_rawDescGZIP() []byte {
	file_proto_server_proto_rawDescOnce.Do(func() {
		file_proto_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_server_proto_rawDescData)
	})
	return file_proto_server_proto_rawDescData
}

var file_proto_server_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_server_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_proto_server_proto_goTypes = []interface{}{
	(CodeEnum)(0),             // 0: CodeEnum
	(*Book)(nil),              // 1: Book
	(*AddBookRequest)(nil),    // 2: AddBookRequest
	(*AddBookReply)(nil),      // 3: AddBookReply
	(*DeleteBookRequest)(nil), // 4: DeleteBookRequest
	(*DeleteBookReply)(nil),   // 5: DeleteBookReply
	(*UpdateBookRequest)(nil), // 6: UpdateBookRequest
	(*UpdateBookReply)(nil),   // 7: UpdateBookReply
	(*GetBookRequest)(nil),    // 8: GetBookRequest
	(*GetBookReply)(nil),      // 9: GetBookReply
	(*ListBookRequest)(nil),   // 10: ListBookRequest
	(*ListBookReply)(nil),     // 11: ListBookReply
}
var file_proto_server_proto_depIdxs = []int32{
	0,  // 0: AddBookReply.code:type_name -> CodeEnum
	1,  // 1: AddBookReply.data:type_name -> Book
	0,  // 2: DeleteBookReply.code:type_name -> CodeEnum
	0,  // 3: UpdateBookReply.code:type_name -> CodeEnum
	1,  // 4: UpdateBookReply.data:type_name -> Book
	0,  // 5: GetBookReply.code:type_name -> CodeEnum
	1,  // 6: GetBookReply.data:type_name -> Book
	0,  // 7: ListBookReply.code:type_name -> CodeEnum
	1,  // 8: ListBookReply.list:type_name -> Book
	2,  // 9: Bookstore.AddBook:input_type -> AddBookRequest
	4,  // 10: Bookstore.DeleteBook:input_type -> DeleteBookRequest
	6,  // 11: Bookstore.UpdateBook:input_type -> UpdateBookRequest
	8,  // 12: Bookstore.GetBook:input_type -> GetBookRequest
	10, // 13: Bookstore.ListBook:input_type -> ListBookRequest
	3,  // 14: Bookstore.AddBook:output_type -> AddBookReply
	5,  // 15: Bookstore.DeleteBook:output_type -> DeleteBookReply
	7,  // 16: Bookstore.UpdateBook:output_type -> UpdateBookReply
	9,  // 17: Bookstore.GetBook:output_type -> GetBookReply
	11, // 18: Bookstore.ListBook:output_type -> ListBookReply
	14, // [14:19] is the sub-list for method output_type
	9,  // [9:14] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_proto_server_proto_init() }
func file_proto_server_proto_init() {
	if File_proto_server_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_server_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Book); i {
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
		file_proto_server_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddBookRequest); i {
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
		file_proto_server_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddBookReply); i {
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
		file_proto_server_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteBookRequest); i {
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
		file_proto_server_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteBookReply); i {
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
		file_proto_server_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateBookRequest); i {
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
		file_proto_server_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateBookReply); i {
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
		file_proto_server_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBookRequest); i {
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
		file_proto_server_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBookReply); i {
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
		file_proto_server_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListBookRequest); i {
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
		file_proto_server_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListBookReply); i {
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
			RawDescriptor: file_proto_server_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_server_proto_goTypes,
		DependencyIndexes: file_proto_server_proto_depIdxs,
		EnumInfos:         file_proto_server_proto_enumTypes,
		MessageInfos:      file_proto_server_proto_msgTypes,
	}.Build()
	File_proto_server_proto = out.File
	file_proto_server_proto_rawDesc = nil
	file_proto_server_proto_goTypes = nil
	file_proto_server_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// BookstoreClient is the client API for Bookstore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BookstoreClient interface {
	AddBook(ctx context.Context, in *AddBookRequest, opts ...grpc.CallOption) (*AddBookReply, error)
	DeleteBook(ctx context.Context, in *DeleteBookRequest, opts ...grpc.CallOption) (*DeleteBookReply, error)
	UpdateBook(ctx context.Context, in *UpdateBookRequest, opts ...grpc.CallOption) (*UpdateBookReply, error)
	GetBook(ctx context.Context, in *GetBookRequest, opts ...grpc.CallOption) (*GetBookReply, error)
	ListBook(ctx context.Context, in *ListBookRequest, opts ...grpc.CallOption) (*ListBookReply, error)
}

type bookstoreClient struct {
	cc grpc.ClientConnInterface
}

func NewBookstoreClient(cc grpc.ClientConnInterface) BookstoreClient {
	return &bookstoreClient{cc}
}

func (c *bookstoreClient) AddBook(ctx context.Context, in *AddBookRequest, opts ...grpc.CallOption) (*AddBookReply, error) {
	out := new(AddBookReply)
	err := c.cc.Invoke(ctx, "/Bookstore/AddBook", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookstoreClient) DeleteBook(ctx context.Context, in *DeleteBookRequest, opts ...grpc.CallOption) (*DeleteBookReply, error) {
	out := new(DeleteBookReply)
	err := c.cc.Invoke(ctx, "/Bookstore/DeleteBook", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookstoreClient) UpdateBook(ctx context.Context, in *UpdateBookRequest, opts ...grpc.CallOption) (*UpdateBookReply, error) {
	out := new(UpdateBookReply)
	err := c.cc.Invoke(ctx, "/Bookstore/UpdateBook", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookstoreClient) GetBook(ctx context.Context, in *GetBookRequest, opts ...grpc.CallOption) (*GetBookReply, error) {
	out := new(GetBookReply)
	err := c.cc.Invoke(ctx, "/Bookstore/GetBook", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookstoreClient) ListBook(ctx context.Context, in *ListBookRequest, opts ...grpc.CallOption) (*ListBookReply, error) {
	out := new(ListBookReply)
	err := c.cc.Invoke(ctx, "/Bookstore/ListBook", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BookstoreServer is the server API for Bookstore service.
type BookstoreServer interface {
	AddBook(context.Context, *AddBookRequest) (*AddBookReply, error)
	DeleteBook(context.Context, *DeleteBookRequest) (*DeleteBookReply, error)
	UpdateBook(context.Context, *UpdateBookRequest) (*UpdateBookReply, error)
	GetBook(context.Context, *GetBookRequest) (*GetBookReply, error)
	ListBook(context.Context, *ListBookRequest) (*ListBookReply, error)
}

// UnimplementedBookstoreServer can be embedded to have forward compatible implementations.
type UnimplementedBookstoreServer struct {
}

func (*UnimplementedBookstoreServer) AddBook(context.Context, *AddBookRequest) (*AddBookReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddBook not implemented")
}
func (*UnimplementedBookstoreServer) DeleteBook(context.Context, *DeleteBookRequest) (*DeleteBookReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteBook not implemented")
}
func (*UnimplementedBookstoreServer) UpdateBook(context.Context, *UpdateBookRequest) (*UpdateBookReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateBook not implemented")
}
func (*UnimplementedBookstoreServer) GetBook(context.Context, *GetBookRequest) (*GetBookReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBook not implemented")
}
func (*UnimplementedBookstoreServer) ListBook(context.Context, *ListBookRequest) (*ListBookReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBook not implemented")
}

func RegisterBookstoreServer(s *grpc.Server, srv BookstoreServer) {
	s.RegisterService(&_Bookstore_serviceDesc, srv)
}

func _Bookstore_AddBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddBookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookstoreServer).AddBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Bookstore/AddBook",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookstoreServer).AddBook(ctx, req.(*AddBookRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bookstore_DeleteBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteBookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookstoreServer).DeleteBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Bookstore/DeleteBook",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookstoreServer).DeleteBook(ctx, req.(*DeleteBookRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bookstore_UpdateBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateBookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookstoreServer).UpdateBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Bookstore/UpdateBook",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookstoreServer).UpdateBook(ctx, req.(*UpdateBookRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bookstore_GetBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookstoreServer).GetBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Bookstore/GetBook",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookstoreServer).GetBook(ctx, req.(*GetBookRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bookstore_ListBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListBookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookstoreServer).ListBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Bookstore/ListBook",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookstoreServer).ListBook(ctx, req.(*ListBookRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Bookstore_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Bookstore",
	HandlerType: (*BookstoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddBook",
			Handler:    _Bookstore_AddBook_Handler,
		},
		{
			MethodName: "DeleteBook",
			Handler:    _Bookstore_DeleteBook_Handler,
		},
		{
			MethodName: "UpdateBook",
			Handler:    _Bookstore_UpdateBook_Handler,
		},
		{
			MethodName: "GetBook",
			Handler:    _Bookstore_GetBook_Handler,
		},
		{
			MethodName: "ListBook",
			Handler:    _Bookstore_ListBook_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/server.proto",
}
