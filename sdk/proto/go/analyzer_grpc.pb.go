// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: pulumi/analyzer.proto

package pulumirpc

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

// AnalyzerClient is the client API for Analyzer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AnalyzerClient interface {
	// Analyze analyzes a single resource object, and returns any errors that it finds.
	// Called with the "inputs" to the resource, before it is updated.
	Analyze(ctx context.Context, in *AnalyzeRequest, opts ...grpc.CallOption) (*AnalyzeResponse, error)
	// AnalyzeStack analyzes all resources within a stack, at the end of a successful
	// preview or update. The provided resources are the "outputs", after any mutations
	// have taken place.
	AnalyzeStack(ctx context.Context, in *AnalyzeStackRequest, opts ...grpc.CallOption) (*AnalyzeResponse, error)
	// GetAnalyzerInfo returns metadata about the analyzer (e.g., list of policies contained).
	GetAnalyzerInfo(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AnalyzerInfo, error)
	// GetPluginInfo returns generic information about this plugin, like its version.
	GetPluginInfo(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PluginInfo, error)
	// Configure configures the analyzer, passing configuration properties for each policy.
	Configure(ctx context.Context, in *ConfigureAnalyzerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type analyzerClient struct {
	cc grpc.ClientConnInterface
}

func NewAnalyzerClient(cc grpc.ClientConnInterface) AnalyzerClient {
	return &analyzerClient{cc}
}

func (c *analyzerClient) Analyze(ctx context.Context, in *AnalyzeRequest, opts ...grpc.CallOption) (*AnalyzeResponse, error) {
	out := new(AnalyzeResponse)
	err := c.cc.Invoke(ctx, "/pulumirpc.Analyzer/Analyze", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analyzerClient) AnalyzeStack(ctx context.Context, in *AnalyzeStackRequest, opts ...grpc.CallOption) (*AnalyzeResponse, error) {
	out := new(AnalyzeResponse)
	err := c.cc.Invoke(ctx, "/pulumirpc.Analyzer/AnalyzeStack", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analyzerClient) GetAnalyzerInfo(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AnalyzerInfo, error) {
	out := new(AnalyzerInfo)
	err := c.cc.Invoke(ctx, "/pulumirpc.Analyzer/GetAnalyzerInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analyzerClient) GetPluginInfo(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PluginInfo, error) {
	out := new(PluginInfo)
	err := c.cc.Invoke(ctx, "/pulumirpc.Analyzer/GetPluginInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analyzerClient) Configure(ctx context.Context, in *ConfigureAnalyzerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/pulumirpc.Analyzer/Configure", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AnalyzerServer is the server API for Analyzer service.
// All implementations must embed UnimplementedAnalyzerServer
// for forward compatibility
type AnalyzerServer interface {
	// Analyze analyzes a single resource object, and returns any errors that it finds.
	// Called with the "inputs" to the resource, before it is updated.
	Analyze(context.Context, *AnalyzeRequest) (*AnalyzeResponse, error)
	// AnalyzeStack analyzes all resources within a stack, at the end of a successful
	// preview or update. The provided resources are the "outputs", after any mutations
	// have taken place.
	AnalyzeStack(context.Context, *AnalyzeStackRequest) (*AnalyzeResponse, error)
	// GetAnalyzerInfo returns metadata about the analyzer (e.g., list of policies contained).
	GetAnalyzerInfo(context.Context, *emptypb.Empty) (*AnalyzerInfo, error)
	// GetPluginInfo returns generic information about this plugin, like its version.
	GetPluginInfo(context.Context, *emptypb.Empty) (*PluginInfo, error)
	// Configure configures the analyzer, passing configuration properties for each policy.
	Configure(context.Context, *ConfigureAnalyzerRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedAnalyzerServer()
}

// UnimplementedAnalyzerServer must be embedded to have forward compatible implementations.
type UnimplementedAnalyzerServer struct{}

func (UnimplementedAnalyzerServer) Analyze(context.Context, *AnalyzeRequest) (*AnalyzeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Analyze not implemented")
}

func (UnimplementedAnalyzerServer) AnalyzeStack(context.Context, *AnalyzeStackRequest) (*AnalyzeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AnalyzeStack not implemented")
}

func (UnimplementedAnalyzerServer) GetAnalyzerInfo(context.Context, *emptypb.Empty) (*AnalyzerInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAnalyzerInfo not implemented")
}

func (UnimplementedAnalyzerServer) GetPluginInfo(context.Context, *emptypb.Empty) (*PluginInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPluginInfo not implemented")
}

func (UnimplementedAnalyzerServer) Configure(context.Context, *ConfigureAnalyzerRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Configure not implemented")
}
func (UnimplementedAnalyzerServer) mustEmbedUnimplementedAnalyzerServer() {}

// UnsafeAnalyzerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AnalyzerServer will
// result in compilation errors.
type UnsafeAnalyzerServer interface {
	mustEmbedUnimplementedAnalyzerServer()
}

func RegisterAnalyzerServer(s grpc.ServiceRegistrar, srv AnalyzerServer) {
	s.RegisterService(&Analyzer_ServiceDesc, srv)
}

func _Analyzer_Analyze_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AnalyzeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyzerServer).Analyze(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pulumirpc.Analyzer/Analyze",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyzerServer).Analyze(ctx, req.(*AnalyzeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Analyzer_AnalyzeStack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AnalyzeStackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyzerServer).AnalyzeStack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pulumirpc.Analyzer/AnalyzeStack",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyzerServer).AnalyzeStack(ctx, req.(*AnalyzeStackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Analyzer_GetAnalyzerInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyzerServer).GetAnalyzerInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pulumirpc.Analyzer/GetAnalyzerInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyzerServer).GetAnalyzerInfo(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Analyzer_GetPluginInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyzerServer).GetPluginInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pulumirpc.Analyzer/GetPluginInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyzerServer).GetPluginInfo(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Analyzer_Configure_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigureAnalyzerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyzerServer).Configure(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pulumirpc.Analyzer/Configure",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyzerServer).Configure(ctx, req.(*ConfigureAnalyzerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Analyzer_ServiceDesc is the grpc.ServiceDesc for Analyzer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Analyzer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pulumirpc.Analyzer",
	HandlerType: (*AnalyzerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Analyze",
			Handler:    _Analyzer_Analyze_Handler,
		},
		{
			MethodName: "AnalyzeStack",
			Handler:    _Analyzer_AnalyzeStack_Handler,
		},
		{
			MethodName: "GetAnalyzerInfo",
			Handler:    _Analyzer_GetAnalyzerInfo_Handler,
		},
		{
			MethodName: "GetPluginInfo",
			Handler:    _Analyzer_GetPluginInfo_Handler,
		},
		{
			MethodName: "Configure",
			Handler:    _Analyzer_Configure_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pulumi/analyzer.proto",
}
