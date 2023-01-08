// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.6
// source: payment.proto

package paymentpb

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

// PaymentServiceClient is the client API for PaymentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PaymentServiceClient interface {
	CreatePayment(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*Statement, error)
	CapturePayment(ctx context.Context, in *Payment, opts ...grpc.CallOption) (*Statement, error)
	CancelPayment(ctx context.Context, in *Payment, opts ...grpc.CallOption) (*Statement, error)
	RefundPayment(ctx context.Context, in *Payment, opts ...grpc.CallOption) (*Statement, error)
}

type paymentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentServiceClient(cc grpc.ClientConnInterface) PaymentServiceClient {
	return &paymentServiceClient{cc}
}

func (c *paymentServiceClient) CreatePayment(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*Statement, error) {
	out := new(Statement)
	err := c.cc.Invoke(ctx, "/payment.PaymentService/CreatePayment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentServiceClient) CapturePayment(ctx context.Context, in *Payment, opts ...grpc.CallOption) (*Statement, error) {
	out := new(Statement)
	err := c.cc.Invoke(ctx, "/payment.PaymentService/CapturePayment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentServiceClient) CancelPayment(ctx context.Context, in *Payment, opts ...grpc.CallOption) (*Statement, error) {
	out := new(Statement)
	err := c.cc.Invoke(ctx, "/payment.PaymentService/CancelPayment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentServiceClient) RefundPayment(ctx context.Context, in *Payment, opts ...grpc.CallOption) (*Statement, error) {
	out := new(Statement)
	err := c.cc.Invoke(ctx, "/payment.PaymentService/RefundPayment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentServiceServer is the server API for PaymentService service.
// All implementations must embed UnimplementedPaymentServiceServer
// for forward compatibility
type PaymentServiceServer interface {
	CreatePayment(context.Context, *CreateRequest) (*Statement, error)
	CapturePayment(context.Context, *Payment) (*Statement, error)
	CancelPayment(context.Context, *Payment) (*Statement, error)
	RefundPayment(context.Context, *Payment) (*Statement, error)
	mustEmbedUnimplementedPaymentServiceServer()
}

// UnimplementedPaymentServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPaymentServiceServer struct {
}

func (UnimplementedPaymentServiceServer) CreatePayment(context.Context, *CreateRequest) (*Statement, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePayment not implemented")
}
func (UnimplementedPaymentServiceServer) CapturePayment(context.Context, *Payment) (*Statement, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CapturePayment not implemented")
}
func (UnimplementedPaymentServiceServer) CancelPayment(context.Context, *Payment) (*Statement, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelPayment not implemented")
}
func (UnimplementedPaymentServiceServer) RefundPayment(context.Context, *Payment) (*Statement, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefundPayment not implemented")
}
func (UnimplementedPaymentServiceServer) mustEmbedUnimplementedPaymentServiceServer() {}

// UnsafePaymentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PaymentServiceServer will
// result in compilation errors.
type UnsafePaymentServiceServer interface {
	mustEmbedUnimplementedPaymentServiceServer()
}

func RegisterPaymentServiceServer(s grpc.ServiceRegistrar, srv PaymentServiceServer) {
	s.RegisterService(&PaymentService_ServiceDesc, srv)
}

func _PaymentService_CreatePayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).CreatePayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/payment.PaymentService/CreatePayment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).CreatePayment(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentService_CapturePayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Payment)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).CapturePayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/payment.PaymentService/CapturePayment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).CapturePayment(ctx, req.(*Payment))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentService_CancelPayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Payment)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).CancelPayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/payment.PaymentService/CancelPayment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).CancelPayment(ctx, req.(*Payment))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentService_RefundPayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Payment)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).RefundPayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/payment.PaymentService/RefundPayment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).RefundPayment(ctx, req.(*Payment))
	}
	return interceptor(ctx, in, info, handler)
}

// PaymentService_ServiceDesc is the grpc.ServiceDesc for PaymentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PaymentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "payment.PaymentService",
	HandlerType: (*PaymentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePayment",
			Handler:    _PaymentService_CreatePayment_Handler,
		},
		{
			MethodName: "CapturePayment",
			Handler:    _PaymentService_CapturePayment_Handler,
		},
		{
			MethodName: "CancelPayment",
			Handler:    _PaymentService_CancelPayment_Handler,
		},
		{
			MethodName: "RefundPayment",
			Handler:    _PaymentService_RefundPayment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "payment.proto",
}
