// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AcknowledgementServiceClient is the client API for AcknowledgementService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AcknowledgementServiceClient interface {
	Ack(ctx context.Context, opts ...grpc.CallOption) (AcknowledgementService_AckClient, error)
}

type acknowledgementServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAcknowledgementServiceClient(cc grpc.ClientConnInterface) AcknowledgementServiceClient {
	return &acknowledgementServiceClient{cc}
}

func (c *acknowledgementServiceClient) Ack(ctx context.Context, opts ...grpc.CallOption) (AcknowledgementService_AckClient, error) {
	stream, err := c.cc.NewStream(ctx, &AcknowledgementService_ServiceDesc.Streams[0], "/proto.acknowledge.AcknowledgementService/Ack", opts...)
	if err != nil {
		return nil, err
	}
	x := &acknowledgementServiceAckClient{stream}
	return x, nil
}

type AcknowledgementService_AckClient interface {
	Send(*AckRequest) error
	Recv() (*AckResponse, error)
	grpc.ClientStream
}

type acknowledgementServiceAckClient struct {
	grpc.ClientStream
}

func (x *acknowledgementServiceAckClient) Send(m *AckRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *acknowledgementServiceAckClient) Recv() (*AckResponse, error) {
	m := new(AckResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AcknowledgementServiceServer is the server API for AcknowledgementService service.
// All implementations must embed UnimplementedAcknowledgementServiceServer
// for forward compatibility
type AcknowledgementServiceServer interface {
	Ack(AcknowledgementService_AckServer) error
	mustEmbedUnimplementedAcknowledgementServiceServer()
}

// UnimplementedAcknowledgementServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAcknowledgementServiceServer struct {
}

func (UnimplementedAcknowledgementServiceServer) Ack(AcknowledgementService_AckServer) error {
	return status.Errorf(codes.Unimplemented, "method Ack not implemented")
}
func (UnimplementedAcknowledgementServiceServer) mustEmbedUnimplementedAcknowledgementServiceServer() {
}

// UnsafeAcknowledgementServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AcknowledgementServiceServer will
// result in compilation errors.
type UnsafeAcknowledgementServiceServer interface {
	mustEmbedUnimplementedAcknowledgementServiceServer()
}

func RegisterAcknowledgementServiceServer(s grpc.ServiceRegistrar, srv AcknowledgementServiceServer) {
	s.RegisterService(&AcknowledgementService_ServiceDesc, srv)
}

func _AcknowledgementService_Ack_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AcknowledgementServiceServer).Ack(&acknowledgementServiceAckServer{stream})
}

type AcknowledgementService_AckServer interface {
	Send(*AckResponse) error
	Recv() (*AckRequest, error)
	grpc.ServerStream
}

type acknowledgementServiceAckServer struct {
	grpc.ServerStream
}

func (x *acknowledgementServiceAckServer) Send(m *AckResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *acknowledgementServiceAckServer) Recv() (*AckRequest, error) {
	m := new(AckRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AcknowledgementService_ServiceDesc is the grpc.ServiceDesc for AcknowledgementService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AcknowledgementService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.acknowledge.AcknowledgementService",
	HandlerType: (*AcknowledgementServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Ack",
			Handler:       _AcknowledgementService_Ack_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/acknowledge.proto",
}
