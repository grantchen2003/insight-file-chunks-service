// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.3
// source: file_chunks_service.proto

package protobufs

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	FileChunksService_CreateFileChunks_FullMethodName                           = "/FileChunksService/CreateFileChunks"
	FileChunksService_GetSortedFileChunksContent_FullMethodName                 = "/FileChunksService/GetSortedFileChunksContent"
	FileChunksService_DeleteFileChunksByRepositoryId_FullMethodName             = "/FileChunksService/DeleteFileChunksByRepositoryId"
	FileChunksService_DeleteFileChunksByRepositoryIdAndFilePaths_FullMethodName = "/FileChunksService/DeleteFileChunksByRepositoryIdAndFilePaths"
)

// FileChunksServiceClient is the client API for FileChunksService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileChunksServiceClient interface {
	CreateFileChunks(ctx context.Context, in *CreateFileChunksRequest, opts ...grpc.CallOption) (*CreateFileChunksResponse, error)
	GetSortedFileChunksContent(ctx context.Context, in *GetSortedFileChunksContentRequest, opts ...grpc.CallOption) (FileChunksService_GetSortedFileChunksContentClient, error)
	DeleteFileChunksByRepositoryId(ctx context.Context, in *DeleteFileChunksByRepositoryIdRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteFileChunksByRepositoryIdAndFilePaths(ctx context.Context, in *DeleteFileChunksByRepositoryIdAndFilePathsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type fileChunksServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFileChunksServiceClient(cc grpc.ClientConnInterface) FileChunksServiceClient {
	return &fileChunksServiceClient{cc}
}

func (c *fileChunksServiceClient) CreateFileChunks(ctx context.Context, in *CreateFileChunksRequest, opts ...grpc.CallOption) (*CreateFileChunksResponse, error) {
	out := new(CreateFileChunksResponse)
	err := c.cc.Invoke(ctx, FileChunksService_CreateFileChunks_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileChunksServiceClient) GetSortedFileChunksContent(ctx context.Context, in *GetSortedFileChunksContentRequest, opts ...grpc.CallOption) (FileChunksService_GetSortedFileChunksContentClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileChunksService_ServiceDesc.Streams[0], FileChunksService_GetSortedFileChunksContent_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &fileChunksServiceGetSortedFileChunksContentClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FileChunksService_GetSortedFileChunksContentClient interface {
	Recv() (*FileChunkContent, error)
	grpc.ClientStream
}

type fileChunksServiceGetSortedFileChunksContentClient struct {
	grpc.ClientStream
}

func (x *fileChunksServiceGetSortedFileChunksContentClient) Recv() (*FileChunkContent, error) {
	m := new(FileChunkContent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileChunksServiceClient) DeleteFileChunksByRepositoryId(ctx context.Context, in *DeleteFileChunksByRepositoryIdRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, FileChunksService_DeleteFileChunksByRepositoryId_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileChunksServiceClient) DeleteFileChunksByRepositoryIdAndFilePaths(ctx context.Context, in *DeleteFileChunksByRepositoryIdAndFilePathsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, FileChunksService_DeleteFileChunksByRepositoryIdAndFilePaths_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileChunksServiceServer is the server API for FileChunksService service.
// All implementations must embed UnimplementedFileChunksServiceServer
// for forward compatibility
type FileChunksServiceServer interface {
	CreateFileChunks(context.Context, *CreateFileChunksRequest) (*CreateFileChunksResponse, error)
	GetSortedFileChunksContent(*GetSortedFileChunksContentRequest, FileChunksService_GetSortedFileChunksContentServer) error
	DeleteFileChunksByRepositoryId(context.Context, *DeleteFileChunksByRepositoryIdRequest) (*emptypb.Empty, error)
	DeleteFileChunksByRepositoryIdAndFilePaths(context.Context, *DeleteFileChunksByRepositoryIdAndFilePathsRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedFileChunksServiceServer()
}

// UnimplementedFileChunksServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFileChunksServiceServer struct {
}

func (UnimplementedFileChunksServiceServer) CreateFileChunks(context.Context, *CreateFileChunksRequest) (*CreateFileChunksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFileChunks not implemented")
}
func (UnimplementedFileChunksServiceServer) GetSortedFileChunksContent(*GetSortedFileChunksContentRequest, FileChunksService_GetSortedFileChunksContentServer) error {
	return status.Errorf(codes.Unimplemented, "method GetSortedFileChunksContent not implemented")
}
func (UnimplementedFileChunksServiceServer) DeleteFileChunksByRepositoryId(context.Context, *DeleteFileChunksByRepositoryIdRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFileChunksByRepositoryId not implemented")
}
func (UnimplementedFileChunksServiceServer) DeleteFileChunksByRepositoryIdAndFilePaths(context.Context, *DeleteFileChunksByRepositoryIdAndFilePathsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFileChunksByRepositoryIdAndFilePaths not implemented")
}
func (UnimplementedFileChunksServiceServer) mustEmbedUnimplementedFileChunksServiceServer() {}

// UnsafeFileChunksServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileChunksServiceServer will
// result in compilation errors.
type UnsafeFileChunksServiceServer interface {
	mustEmbedUnimplementedFileChunksServiceServer()
}

func RegisterFileChunksServiceServer(s grpc.ServiceRegistrar, srv FileChunksServiceServer) {
	s.RegisterService(&FileChunksService_ServiceDesc, srv)
}

func _FileChunksService_CreateFileChunks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFileChunksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileChunksServiceServer).CreateFileChunks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileChunksService_CreateFileChunks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileChunksServiceServer).CreateFileChunks(ctx, req.(*CreateFileChunksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileChunksService_GetSortedFileChunksContent_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetSortedFileChunksContentRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FileChunksServiceServer).GetSortedFileChunksContent(m, &fileChunksServiceGetSortedFileChunksContentServer{stream})
}

type FileChunksService_GetSortedFileChunksContentServer interface {
	Send(*FileChunkContent) error
	grpc.ServerStream
}

type fileChunksServiceGetSortedFileChunksContentServer struct {
	grpc.ServerStream
}

func (x *fileChunksServiceGetSortedFileChunksContentServer) Send(m *FileChunkContent) error {
	return x.ServerStream.SendMsg(m)
}

func _FileChunksService_DeleteFileChunksByRepositoryId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFileChunksByRepositoryIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileChunksServiceServer).DeleteFileChunksByRepositoryId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileChunksService_DeleteFileChunksByRepositoryId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileChunksServiceServer).DeleteFileChunksByRepositoryId(ctx, req.(*DeleteFileChunksByRepositoryIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileChunksService_DeleteFileChunksByRepositoryIdAndFilePaths_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFileChunksByRepositoryIdAndFilePathsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileChunksServiceServer).DeleteFileChunksByRepositoryIdAndFilePaths(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileChunksService_DeleteFileChunksByRepositoryIdAndFilePaths_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileChunksServiceServer).DeleteFileChunksByRepositoryIdAndFilePaths(ctx, req.(*DeleteFileChunksByRepositoryIdAndFilePathsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FileChunksService_ServiceDesc is the grpc.ServiceDesc for FileChunksService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileChunksService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "FileChunksService",
	HandlerType: (*FileChunksServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateFileChunks",
			Handler:    _FileChunksService_CreateFileChunks_Handler,
		},
		{
			MethodName: "DeleteFileChunksByRepositoryId",
			Handler:    _FileChunksService_DeleteFileChunksByRepositoryId_Handler,
		},
		{
			MethodName: "DeleteFileChunksByRepositoryIdAndFilePaths",
			Handler:    _FileChunksService_DeleteFileChunksByRepositoryIdAndFilePaths_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetSortedFileChunksContent",
			Handler:       _FileChunksService_GetSortedFileChunksContent_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "file_chunks_service.proto",
}
