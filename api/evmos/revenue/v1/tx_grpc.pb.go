// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:LGPL-3.0-only

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: evmos/revenue/v1/tx.proto

package revenuev1

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

const (
	Msg_RegisterRevenue_FullMethodName = "/evmos.revenue.v1.Msg/RegisterRevenue"
	Msg_UpdateRevenue_FullMethodName   = "/evmos.revenue.v1.Msg/UpdateRevenue"
	Msg_CancelRevenue_FullMethodName   = "/evmos.revenue.v1.Msg/CancelRevenue"
	Msg_UpdateParams_FullMethodName    = "/evmos.revenue.v1.Msg/UpdateParams"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	// RegisterRevenue registers a new contract for receiving transaction fees
	RegisterRevenue(ctx context.Context, in *MsgRegisterRevenue, opts ...grpc.CallOption) (*MsgRegisterRevenueResponse, error)
	// UpdateRevenue updates the withdrawer address of a revenue
	UpdateRevenue(ctx context.Context, in *MsgUpdateRevenue, opts ...grpc.CallOption) (*MsgUpdateRevenueResponse, error)
	// CancelRevenue cancels a contract's fee registration and further receival
	// of transaction fees
	CancelRevenue(ctx context.Context, in *MsgCancelRevenue, opts ...grpc.CallOption) (*MsgCancelRevenueResponse, error)
	// UpdateParams defined a governance operation for updating the x/revenue module parameters.
	// The authority is hard-coded to the Cosmos SDK x/gov module account
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) RegisterRevenue(ctx context.Context, in *MsgRegisterRevenue, opts ...grpc.CallOption) (*MsgRegisterRevenueResponse, error) {
	out := new(MsgRegisterRevenueResponse)
	err := c.cc.Invoke(ctx, Msg_RegisterRevenue_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateRevenue(ctx context.Context, in *MsgUpdateRevenue, opts ...grpc.CallOption) (*MsgUpdateRevenueResponse, error) {
	out := new(MsgUpdateRevenueResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateRevenue_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) CancelRevenue(ctx context.Context, in *MsgCancelRevenue, opts ...grpc.CallOption) (*MsgCancelRevenueResponse, error) {
	out := new(MsgCancelRevenueResponse)
	err := c.cc.Invoke(ctx, Msg_CancelRevenue_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateParams_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	// RegisterRevenue registers a new contract for receiving transaction fees
	RegisterRevenue(context.Context, *MsgRegisterRevenue) (*MsgRegisterRevenueResponse, error)
	// UpdateRevenue updates the withdrawer address of a revenue
	UpdateRevenue(context.Context, *MsgUpdateRevenue) (*MsgUpdateRevenueResponse, error)
	// CancelRevenue cancels a contract's fee registration and further receival
	// of transaction fees
	CancelRevenue(context.Context, *MsgCancelRevenue) (*MsgCancelRevenueResponse, error)
	// UpdateParams defined a governance operation for updating the x/revenue module parameters.
	// The authority is hard-coded to the Cosmos SDK x/gov module account
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) RegisterRevenue(context.Context, *MsgRegisterRevenue) (*MsgRegisterRevenueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterRevenue not implemented")
}
func (UnimplementedMsgServer) UpdateRevenue(context.Context, *MsgUpdateRevenue) (*MsgUpdateRevenueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRevenue not implemented")
}
func (UnimplementedMsgServer) CancelRevenue(context.Context, *MsgCancelRevenue) (*MsgCancelRevenueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelRevenue not implemented")
}
func (UnimplementedMsgServer) UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_RegisterRevenue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterRevenue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RegisterRevenue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_RegisterRevenue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RegisterRevenue(ctx, req.(*MsgRegisterRevenue))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateRevenue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateRevenue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateRevenue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateRevenue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateRevenue(ctx, req.(*MsgUpdateRevenue))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CancelRevenue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCancelRevenue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CancelRevenue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_CancelRevenue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CancelRevenue(ctx, req.(*MsgCancelRevenue))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateParams_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "evmos.revenue.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterRevenue",
			Handler:    _Msg_RegisterRevenue_Handler,
		},
		{
			MethodName: "UpdateRevenue",
			Handler:    _Msg_UpdateRevenue_Handler,
		},
		{
			MethodName: "CancelRevenue",
			Handler:    _Msg_CancelRevenue_Handler,
		},
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "evmos/revenue/v1/tx.proto",
}