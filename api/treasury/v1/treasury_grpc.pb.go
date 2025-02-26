// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.3
// source: treasury/v1/treasury.proto

package v1

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
	Treasury_CreateWithdrawClaim_FullMethodName            = "/treasury.v1.Treasury/CreateWithdrawClaim"
	Treasury_ListWithdrawClaims_FullMethodName             = "/treasury.v1.Treasury/ListWithdrawClaims"
	Treasury_GetWithdrawClaim_FullMethodName               = "/treasury.v1.Treasury/GetWithdrawClaim"
	Treasury_ApproveWithdrawClaim_FullMethodName           = "/treasury.v1.Treasury/ApproveWithdrawClaim"
	Treasury_RejectWithdrawClaim_FullMethodName            = "/treasury.v1.Treasury/RejectWithdrawClaim"
	Treasury_ListWithdrawClaimConfirmations_FullMethodName = "/treasury.v1.Treasury/ListWithdrawClaimConfirmations"
)

// TreasuryClient is the client API for Treasury service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TreasuryClient interface {
	CreateWithdrawClaim(ctx context.Context, in *CreateWithdrawClaimRequest, opts ...grpc.CallOption) (*CreateWithdrawClaimReply, error)
	ListWithdrawClaims(ctx context.Context, in *ListWithdrawClaimsRequest, opts ...grpc.CallOption) (*ListWithdrawClaimsReply, error)
	GetWithdrawClaim(ctx context.Context, in *GetWithdrawClaimRequest, opts ...grpc.CallOption) (*GetWithdrawClaimReply, error)
	ApproveWithdrawClaim(ctx context.Context, in *ApproveWithdrawClaimRequest, opts ...grpc.CallOption) (*ApproveWithdrawClaimReply, error)
	RejectWithdrawClaim(ctx context.Context, in *RejectWithdrawClaimRequest, opts ...grpc.CallOption) (*RejectWithdrawClaimReply, error)
	ListWithdrawClaimConfirmations(ctx context.Context, in *ListWithdrawClaimConfirmationsRequest, opts ...grpc.CallOption) (*ListWithdrawClaimConfirmationsReply, error)
}

type treasuryClient struct {
	cc grpc.ClientConnInterface
}

func NewTreasuryClient(cc grpc.ClientConnInterface) TreasuryClient {
	return &treasuryClient{cc}
}

func (c *treasuryClient) CreateWithdrawClaim(ctx context.Context, in *CreateWithdrawClaimRequest, opts ...grpc.CallOption) (*CreateWithdrawClaimReply, error) {
	out := new(CreateWithdrawClaimReply)
	err := c.cc.Invoke(ctx, Treasury_CreateWithdrawClaim_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *treasuryClient) ListWithdrawClaims(ctx context.Context, in *ListWithdrawClaimsRequest, opts ...grpc.CallOption) (*ListWithdrawClaimsReply, error) {
	out := new(ListWithdrawClaimsReply)
	err := c.cc.Invoke(ctx, Treasury_ListWithdrawClaims_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *treasuryClient) GetWithdrawClaim(ctx context.Context, in *GetWithdrawClaimRequest, opts ...grpc.CallOption) (*GetWithdrawClaimReply, error) {
	out := new(GetWithdrawClaimReply)
	err := c.cc.Invoke(ctx, Treasury_GetWithdrawClaim_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *treasuryClient) ApproveWithdrawClaim(ctx context.Context, in *ApproveWithdrawClaimRequest, opts ...grpc.CallOption) (*ApproveWithdrawClaimReply, error) {
	out := new(ApproveWithdrawClaimReply)
	err := c.cc.Invoke(ctx, Treasury_ApproveWithdrawClaim_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *treasuryClient) RejectWithdrawClaim(ctx context.Context, in *RejectWithdrawClaimRequest, opts ...grpc.CallOption) (*RejectWithdrawClaimReply, error) {
	out := new(RejectWithdrawClaimReply)
	err := c.cc.Invoke(ctx, Treasury_RejectWithdrawClaim_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *treasuryClient) ListWithdrawClaimConfirmations(ctx context.Context, in *ListWithdrawClaimConfirmationsRequest, opts ...grpc.CallOption) (*ListWithdrawClaimConfirmationsReply, error) {
	out := new(ListWithdrawClaimConfirmationsReply)
	err := c.cc.Invoke(ctx, Treasury_ListWithdrawClaimConfirmations_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TreasuryServer is the server API for Treasury service.
// All implementations must embed UnimplementedTreasuryServer
// for forward compatibility
type TreasuryServer interface {
	CreateWithdrawClaim(context.Context, *CreateWithdrawClaimRequest) (*CreateWithdrawClaimReply, error)
	ListWithdrawClaims(context.Context, *ListWithdrawClaimsRequest) (*ListWithdrawClaimsReply, error)
	GetWithdrawClaim(context.Context, *GetWithdrawClaimRequest) (*GetWithdrawClaimReply, error)
	ApproveWithdrawClaim(context.Context, *ApproveWithdrawClaimRequest) (*ApproveWithdrawClaimReply, error)
	RejectWithdrawClaim(context.Context, *RejectWithdrawClaimRequest) (*RejectWithdrawClaimReply, error)
	ListWithdrawClaimConfirmations(context.Context, *ListWithdrawClaimConfirmationsRequest) (*ListWithdrawClaimConfirmationsReply, error)
	mustEmbedUnimplementedTreasuryServer()
}

// UnimplementedTreasuryServer must be embedded to have forward compatible implementations.
type UnimplementedTreasuryServer struct {
}

func (UnimplementedTreasuryServer) CreateWithdrawClaim(context.Context, *CreateWithdrawClaimRequest) (*CreateWithdrawClaimReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateWithdrawClaim not implemented")
}
func (UnimplementedTreasuryServer) ListWithdrawClaims(context.Context, *ListWithdrawClaimsRequest) (*ListWithdrawClaimsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListWithdrawClaims not implemented")
}
func (UnimplementedTreasuryServer) GetWithdrawClaim(context.Context, *GetWithdrawClaimRequest) (*GetWithdrawClaimReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWithdrawClaim not implemented")
}
func (UnimplementedTreasuryServer) ApproveWithdrawClaim(context.Context, *ApproveWithdrawClaimRequest) (*ApproveWithdrawClaimReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApproveWithdrawClaim not implemented")
}
func (UnimplementedTreasuryServer) RejectWithdrawClaim(context.Context, *RejectWithdrawClaimRequest) (*RejectWithdrawClaimReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RejectWithdrawClaim not implemented")
}
func (UnimplementedTreasuryServer) ListWithdrawClaimConfirmations(context.Context, *ListWithdrawClaimConfirmationsRequest) (*ListWithdrawClaimConfirmationsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListWithdrawClaimConfirmations not implemented")
}
func (UnimplementedTreasuryServer) mustEmbedUnimplementedTreasuryServer() {}

// UnsafeTreasuryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TreasuryServer will
// result in compilation errors.
type UnsafeTreasuryServer interface {
	mustEmbedUnimplementedTreasuryServer()
}

func RegisterTreasuryServer(s grpc.ServiceRegistrar, srv TreasuryServer) {
	s.RegisterService(&Treasury_ServiceDesc, srv)
}

func _Treasury_CreateWithdrawClaim_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateWithdrawClaimRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TreasuryServer).CreateWithdrawClaim(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Treasury_CreateWithdrawClaim_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TreasuryServer).CreateWithdrawClaim(ctx, req.(*CreateWithdrawClaimRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Treasury_ListWithdrawClaims_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListWithdrawClaimsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TreasuryServer).ListWithdrawClaims(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Treasury_ListWithdrawClaims_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TreasuryServer).ListWithdrawClaims(ctx, req.(*ListWithdrawClaimsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Treasury_GetWithdrawClaim_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWithdrawClaimRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TreasuryServer).GetWithdrawClaim(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Treasury_GetWithdrawClaim_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TreasuryServer).GetWithdrawClaim(ctx, req.(*GetWithdrawClaimRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Treasury_ApproveWithdrawClaim_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApproveWithdrawClaimRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TreasuryServer).ApproveWithdrawClaim(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Treasury_ApproveWithdrawClaim_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TreasuryServer).ApproveWithdrawClaim(ctx, req.(*ApproveWithdrawClaimRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Treasury_RejectWithdrawClaim_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RejectWithdrawClaimRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TreasuryServer).RejectWithdrawClaim(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Treasury_RejectWithdrawClaim_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TreasuryServer).RejectWithdrawClaim(ctx, req.(*RejectWithdrawClaimRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Treasury_ListWithdrawClaimConfirmations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListWithdrawClaimConfirmationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TreasuryServer).ListWithdrawClaimConfirmations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Treasury_ListWithdrawClaimConfirmations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TreasuryServer).ListWithdrawClaimConfirmations(ctx, req.(*ListWithdrawClaimConfirmationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Treasury_ServiceDesc is the grpc.ServiceDesc for Treasury service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Treasury_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "treasury.v1.Treasury",
	HandlerType: (*TreasuryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateWithdrawClaim",
			Handler:    _Treasury_CreateWithdrawClaim_Handler,
		},
		{
			MethodName: "ListWithdrawClaims",
			Handler:    _Treasury_ListWithdrawClaims_Handler,
		},
		{
			MethodName: "GetWithdrawClaim",
			Handler:    _Treasury_GetWithdrawClaim_Handler,
		},
		{
			MethodName: "ApproveWithdrawClaim",
			Handler:    _Treasury_ApproveWithdrawClaim_Handler,
		},
		{
			MethodName: "RejectWithdrawClaim",
			Handler:    _Treasury_RejectWithdrawClaim_Handler,
		},
		{
			MethodName: "ListWithdrawClaimConfirmations",
			Handler:    _Treasury_ListWithdrawClaimConfirmations_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "treasury/v1/treasury.proto",
}
