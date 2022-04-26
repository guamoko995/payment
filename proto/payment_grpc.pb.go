// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.0
// source: proto/payment.proto

package payment

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

// PaymentClient is the client API for Payment service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PaymentClient interface {
	// TopUpAccount - RPC, которая пополняет аккаунт
	// переданного ID на переданную Amount.
	UpSum(ctx context.Context, in *UpRequest, opts ...grpc.CallOption) (*Enum, error)
	// TopUpAccount - RPC, которая переводит Amount средств
	// с первого ID на второй.
	SumTransfer(ctx context.Context, in *TransferRequest, opts ...grpc.CallOption) (*Enum, error)
}

type paymentClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentClient(cc grpc.ClientConnInterface) PaymentClient {
	return &paymentClient{cc}
}

func (c *paymentClient) UpSum(ctx context.Context, in *UpRequest, opts ...grpc.CallOption) (*Enum, error) {
	out := new(Enum)
	err := c.cc.Invoke(ctx, "/proto.Payment/UpSum", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) SumTransfer(ctx context.Context, in *TransferRequest, opts ...grpc.CallOption) (*Enum, error) {
	out := new(Enum)
	err := c.cc.Invoke(ctx, "/proto.Payment/SumTransfer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentServer is the server API for Payment service.
// All implementations must embed UnimplementedPaymentServer
// for forward compatibility
type PaymentServer interface {
	// TopUpAccount - RPC, которая пополняет аккаунт
	// переданного ID на переданную Amount.
	UpSum(context.Context, *UpRequest) (*Enum, error)
	// TopUpAccount - RPC, которая переводит Amount средств
	// с первого ID на второй.
	SumTransfer(context.Context, *TransferRequest) (*Enum, error)
	mustEmbedUnimplementedPaymentServer()
}

// UnimplementedPaymentServer must be embedded to have forward compatible implementations.
type UnimplementedPaymentServer struct {
}

func (UnimplementedPaymentServer) UpSum(context.Context, *UpRequest) (*Enum, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpSum not implemented")
}
func (UnimplementedPaymentServer) SumTransfer(context.Context, *TransferRequest) (*Enum, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SumTransfer not implemented")
}
func (UnimplementedPaymentServer) mustEmbedUnimplementedPaymentServer() {}

// UnsafePaymentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PaymentServer will
// result in compilation errors.
type UnsafePaymentServer interface {
	mustEmbedUnimplementedPaymentServer()
}

func RegisterPaymentServer(s grpc.ServiceRegistrar, srv PaymentServer) {
	s.RegisterService(&Payment_ServiceDesc, srv)
}

func _Payment_UpSum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).UpSum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Payment/UpSum",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).UpSum(ctx, req.(*UpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_SumTransfer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).SumTransfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Payment/SumTransfer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).SumTransfer(ctx, req.(*TransferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Payment_ServiceDesc is the grpc.ServiceDesc for Payment service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Payment_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Payment",
	HandlerType: (*PaymentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpSum",
			Handler:    _Payment_UpSum_Handler,
		},
		{
			MethodName: "SumTransfer",
			Handler:    _Payment_SumTransfer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/payment.proto",
}
