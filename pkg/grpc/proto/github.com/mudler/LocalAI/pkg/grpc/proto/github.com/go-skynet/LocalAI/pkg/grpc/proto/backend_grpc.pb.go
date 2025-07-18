// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.21.12
// source: backend.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Backend_Health_FullMethodName             = "/backend.Backend/Health"
	Backend_Predict_FullMethodName            = "/backend.Backend/Predict"
	Backend_LoadModel_FullMethodName          = "/backend.Backend/LoadModel"
	Backend_PredictStream_FullMethodName      = "/backend.Backend/PredictStream"
	Backend_Embedding_FullMethodName          = "/backend.Backend/Embedding"
	Backend_GenerateImage_FullMethodName      = "/backend.Backend/GenerateImage"
	Backend_GenerateVideo_FullMethodName      = "/backend.Backend/GenerateVideo"
	Backend_AudioTranscription_FullMethodName = "/backend.Backend/AudioTranscription"
	Backend_TTS_FullMethodName                = "/backend.Backend/TTS"
	Backend_SoundGeneration_FullMethodName    = "/backend.Backend/SoundGeneration"
	Backend_TokenizeString_FullMethodName     = "/backend.Backend/TokenizeString"
	Backend_Status_FullMethodName             = "/backend.Backend/Status"
	Backend_StoresSet_FullMethodName          = "/backend.Backend/StoresSet"
	Backend_StoresDelete_FullMethodName       = "/backend.Backend/StoresDelete"
	Backend_StoresGet_FullMethodName          = "/backend.Backend/StoresGet"
	Backend_StoresFind_FullMethodName         = "/backend.Backend/StoresFind"
	Backend_Rerank_FullMethodName             = "/backend.Backend/Rerank"
	Backend_GetMetrics_FullMethodName         = "/backend.Backend/GetMetrics"
	Backend_VAD_FullMethodName                = "/backend.Backend/VAD"
)

// BackendClient is the client API for Backend service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BackendClient interface {
	Health(ctx context.Context, in *HealthMessage, opts ...grpc.CallOption) (*Reply, error)
	Predict(ctx context.Context, in *PredictOptions, opts ...grpc.CallOption) (*Reply, error)
	LoadModel(ctx context.Context, in *ModelOptions, opts ...grpc.CallOption) (*Result, error)
	PredictStream(ctx context.Context, in *PredictOptions, opts ...grpc.CallOption) (Backend_PredictStreamClient, error)
	Embedding(ctx context.Context, in *PredictOptions, opts ...grpc.CallOption) (*EmbeddingResult, error)
	GenerateImage(ctx context.Context, in *GenerateImageRequest, opts ...grpc.CallOption) (*Result, error)
	GenerateVideo(ctx context.Context, in *GenerateVideoRequest, opts ...grpc.CallOption) (*Result, error)
	AudioTranscription(ctx context.Context, in *TranscriptRequest, opts ...grpc.CallOption) (*TranscriptResult, error)
	TTS(ctx context.Context, in *TTSRequest, opts ...grpc.CallOption) (*Result, error)
	SoundGeneration(ctx context.Context, in *SoundGenerationRequest, opts ...grpc.CallOption) (*Result, error)
	TokenizeString(ctx context.Context, in *PredictOptions, opts ...grpc.CallOption) (*TokenizationResponse, error)
	Status(ctx context.Context, in *HealthMessage, opts ...grpc.CallOption) (*StatusResponse, error)
	StoresSet(ctx context.Context, in *StoresSetOptions, opts ...grpc.CallOption) (*Result, error)
	StoresDelete(ctx context.Context, in *StoresDeleteOptions, opts ...grpc.CallOption) (*Result, error)
	StoresGet(ctx context.Context, in *StoresGetOptions, opts ...grpc.CallOption) (*StoresGetResult, error)
	StoresFind(ctx context.Context, in *StoresFindOptions, opts ...grpc.CallOption) (*StoresFindResult, error)
	Rerank(ctx context.Context, in *RerankRequest, opts ...grpc.CallOption) (*RerankResult, error)
	GetMetrics(ctx context.Context, in *MetricsRequest, opts ...grpc.CallOption) (*MetricsResponse, error)
	VAD(ctx context.Context, in *VADRequest, opts ...grpc.CallOption) (*VADResponse, error)
}

type backendClient struct {
	cc grpc.ClientConnInterface
}

func NewBackendClient(cc grpc.ClientConnInterface) BackendClient {
	return &backendClient{cc}
}

func (c *backendClient) Health(ctx context.Context, in *HealthMessage, opts ...grpc.CallOption) (*Reply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Reply)
	err := c.cc.Invoke(ctx, Backend_Health_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) Predict(ctx context.Context, in *PredictOptions, opts ...grpc.CallOption) (*Reply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Reply)
	err := c.cc.Invoke(ctx, Backend_Predict_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) LoadModel(ctx context.Context, in *ModelOptions, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, Backend_LoadModel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) PredictStream(ctx context.Context, in *PredictOptions, opts ...grpc.CallOption) (Backend_PredictStreamClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Backend_ServiceDesc.Streams[0], Backend_PredictStream_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &backendPredictStreamClient{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Backend_PredictStreamClient interface {
	Recv() (*Reply, error)
	grpc.ClientStream
}

type backendPredictStreamClient struct {
	grpc.ClientStream
}

func (x *backendPredictStreamClient) Recv() (*Reply, error) {
	m := new(Reply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *backendClient) Embedding(ctx context.Context, in *PredictOptions, opts ...grpc.CallOption) (*EmbeddingResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EmbeddingResult)
	err := c.cc.Invoke(ctx, Backend_Embedding_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) GenerateImage(ctx context.Context, in *GenerateImageRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, Backend_GenerateImage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) GenerateVideo(ctx context.Context, in *GenerateVideoRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, Backend_GenerateVideo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) AudioTranscription(ctx context.Context, in *TranscriptRequest, opts ...grpc.CallOption) (*TranscriptResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TranscriptResult)
	err := c.cc.Invoke(ctx, Backend_AudioTranscription_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) TTS(ctx context.Context, in *TTSRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, Backend_TTS_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) SoundGeneration(ctx context.Context, in *SoundGenerationRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, Backend_SoundGeneration_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) TokenizeString(ctx context.Context, in *PredictOptions, opts ...grpc.CallOption) (*TokenizationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TokenizationResponse)
	err := c.cc.Invoke(ctx, Backend_TokenizeString_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) Status(ctx context.Context, in *HealthMessage, opts ...grpc.CallOption) (*StatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, Backend_Status_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) StoresSet(ctx context.Context, in *StoresSetOptions, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, Backend_StoresSet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) StoresDelete(ctx context.Context, in *StoresDeleteOptions, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, Backend_StoresDelete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) StoresGet(ctx context.Context, in *StoresGetOptions, opts ...grpc.CallOption) (*StoresGetResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StoresGetResult)
	err := c.cc.Invoke(ctx, Backend_StoresGet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) StoresFind(ctx context.Context, in *StoresFindOptions, opts ...grpc.CallOption) (*StoresFindResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StoresFindResult)
	err := c.cc.Invoke(ctx, Backend_StoresFind_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) Rerank(ctx context.Context, in *RerankRequest, opts ...grpc.CallOption) (*RerankResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RerankResult)
	err := c.cc.Invoke(ctx, Backend_Rerank_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) GetMetrics(ctx context.Context, in *MetricsRequest, opts ...grpc.CallOption) (*MetricsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MetricsResponse)
	err := c.cc.Invoke(ctx, Backend_GetMetrics_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) VAD(ctx context.Context, in *VADRequest, opts ...grpc.CallOption) (*VADResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VADResponse)
	err := c.cc.Invoke(ctx, Backend_VAD_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BackendServer is the server API for Backend service.
// All implementations must embed UnimplementedBackendServer
// for forward compatibility
type BackendServer interface {
	Health(context.Context, *HealthMessage) (*Reply, error)
	Predict(context.Context, *PredictOptions) (*Reply, error)
	LoadModel(context.Context, *ModelOptions) (*Result, error)
	PredictStream(*PredictOptions, Backend_PredictStreamServer) error
	Embedding(context.Context, *PredictOptions) (*EmbeddingResult, error)
	GenerateImage(context.Context, *GenerateImageRequest) (*Result, error)
	GenerateVideo(context.Context, *GenerateVideoRequest) (*Result, error)
	AudioTranscription(context.Context, *TranscriptRequest) (*TranscriptResult, error)
	TTS(context.Context, *TTSRequest) (*Result, error)
	SoundGeneration(context.Context, *SoundGenerationRequest) (*Result, error)
	TokenizeString(context.Context, *PredictOptions) (*TokenizationResponse, error)
	Status(context.Context, *HealthMessage) (*StatusResponse, error)
	StoresSet(context.Context, *StoresSetOptions) (*Result, error)
	StoresDelete(context.Context, *StoresDeleteOptions) (*Result, error)
	StoresGet(context.Context, *StoresGetOptions) (*StoresGetResult, error)
	StoresFind(context.Context, *StoresFindOptions) (*StoresFindResult, error)
	Rerank(context.Context, *RerankRequest) (*RerankResult, error)
	GetMetrics(context.Context, *MetricsRequest) (*MetricsResponse, error)
	VAD(context.Context, *VADRequest) (*VADResponse, error)
	mustEmbedUnimplementedBackendServer()
}

// UnimplementedBackendServer must be embedded to have forward compatible implementations.
type UnimplementedBackendServer struct {
}

func (UnimplementedBackendServer) Health(context.Context, *HealthMessage) (*Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Health not implemented")
}
func (UnimplementedBackendServer) Predict(context.Context, *PredictOptions) (*Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Predict not implemented")
}
func (UnimplementedBackendServer) LoadModel(context.Context, *ModelOptions) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadModel not implemented")
}
func (UnimplementedBackendServer) PredictStream(*PredictOptions, Backend_PredictStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method PredictStream not implemented")
}
func (UnimplementedBackendServer) Embedding(context.Context, *PredictOptions) (*EmbeddingResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Embedding not implemented")
}
func (UnimplementedBackendServer) GenerateImage(context.Context, *GenerateImageRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateImage not implemented")
}
func (UnimplementedBackendServer) GenerateVideo(context.Context, *GenerateVideoRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateVideo not implemented")
}
func (UnimplementedBackendServer) AudioTranscription(context.Context, *TranscriptRequest) (*TranscriptResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AudioTranscription not implemented")
}
func (UnimplementedBackendServer) TTS(context.Context, *TTSRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TTS not implemented")
}
func (UnimplementedBackendServer) SoundGeneration(context.Context, *SoundGenerationRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SoundGeneration not implemented")
}
func (UnimplementedBackendServer) TokenizeString(context.Context, *PredictOptions) (*TokenizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TokenizeString not implemented")
}
func (UnimplementedBackendServer) Status(context.Context, *HealthMessage) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedBackendServer) StoresSet(context.Context, *StoresSetOptions) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoresSet not implemented")
}
func (UnimplementedBackendServer) StoresDelete(context.Context, *StoresDeleteOptions) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoresDelete not implemented")
}
func (UnimplementedBackendServer) StoresGet(context.Context, *StoresGetOptions) (*StoresGetResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoresGet not implemented")
}
func (UnimplementedBackendServer) StoresFind(context.Context, *StoresFindOptions) (*StoresFindResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoresFind not implemented")
}
func (UnimplementedBackendServer) Rerank(context.Context, *RerankRequest) (*RerankResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rerank not implemented")
}
func (UnimplementedBackendServer) GetMetrics(context.Context, *MetricsRequest) (*MetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetrics not implemented")
}
func (UnimplementedBackendServer) VAD(context.Context, *VADRequest) (*VADResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VAD not implemented")
}
func (UnimplementedBackendServer) mustEmbedUnimplementedBackendServer() {}

// UnsafeBackendServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BackendServer will
// result in compilation errors.
type UnsafeBackendServer interface {
	mustEmbedUnimplementedBackendServer()
}

func RegisterBackendServer(s grpc.ServiceRegistrar, srv BackendServer) {
	s.RegisterService(&Backend_ServiceDesc, srv)
}

func _Backend_Health_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).Health(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_Health_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).Health(ctx, req.(*HealthMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_Predict_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PredictOptions)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).Predict(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_Predict_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).Predict(ctx, req.(*PredictOptions))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_LoadModel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModelOptions)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).LoadModel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_LoadModel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).LoadModel(ctx, req.(*ModelOptions))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_PredictStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PredictOptions)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BackendServer).PredictStream(m, &backendPredictStreamServer{ServerStream: stream})
}

type Backend_PredictStreamServer interface {
	Send(*Reply) error
	grpc.ServerStream
}

type backendPredictStreamServer struct {
	grpc.ServerStream
}

func (x *backendPredictStreamServer) Send(m *Reply) error {
	return x.ServerStream.SendMsg(m)
}

func _Backend_Embedding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PredictOptions)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).Embedding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_Embedding_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).Embedding(ctx, req.(*PredictOptions))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_GenerateImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).GenerateImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_GenerateImage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).GenerateImage(ctx, req.(*GenerateImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_GenerateVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateVideoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).GenerateVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_GenerateVideo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).GenerateVideo(ctx, req.(*GenerateVideoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_AudioTranscription_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TranscriptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).AudioTranscription(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_AudioTranscription_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).AudioTranscription(ctx, req.(*TranscriptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_TTS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TTSRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).TTS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_TTS_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).TTS(ctx, req.(*TTSRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_SoundGeneration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SoundGenerationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).SoundGeneration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_SoundGeneration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).SoundGeneration(ctx, req.(*SoundGenerationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_TokenizeString_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PredictOptions)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).TokenizeString(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_TokenizeString_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).TokenizeString(ctx, req.(*PredictOptions))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_Status_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).Status(ctx, req.(*HealthMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_StoresSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoresSetOptions)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).StoresSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_StoresSet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).StoresSet(ctx, req.(*StoresSetOptions))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_StoresDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoresDeleteOptions)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).StoresDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_StoresDelete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).StoresDelete(ctx, req.(*StoresDeleteOptions))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_StoresGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoresGetOptions)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).StoresGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_StoresGet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).StoresGet(ctx, req.(*StoresGetOptions))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_StoresFind_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoresFindOptions)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).StoresFind(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_StoresFind_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).StoresFind(ctx, req.(*StoresFindOptions))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_Rerank_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RerankRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).Rerank(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_Rerank_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).Rerank(ctx, req.(*RerankRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_GetMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).GetMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_GetMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).GetMetrics(ctx, req.(*MetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_VAD_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VADRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).VAD(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Backend_VAD_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).VAD(ctx, req.(*VADRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Backend_ServiceDesc is the grpc.ServiceDesc for Backend service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Backend_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "backend.Backend",
	HandlerType: (*BackendServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Health",
			Handler:    _Backend_Health_Handler,
		},
		{
			MethodName: "Predict",
			Handler:    _Backend_Predict_Handler,
		},
		{
			MethodName: "LoadModel",
			Handler:    _Backend_LoadModel_Handler,
		},
		{
			MethodName: "Embedding",
			Handler:    _Backend_Embedding_Handler,
		},
		{
			MethodName: "GenerateImage",
			Handler:    _Backend_GenerateImage_Handler,
		},
		{
			MethodName: "GenerateVideo",
			Handler:    _Backend_GenerateVideo_Handler,
		},
		{
			MethodName: "AudioTranscription",
			Handler:    _Backend_AudioTranscription_Handler,
		},
		{
			MethodName: "TTS",
			Handler:    _Backend_TTS_Handler,
		},
		{
			MethodName: "SoundGeneration",
			Handler:    _Backend_SoundGeneration_Handler,
		},
		{
			MethodName: "TokenizeString",
			Handler:    _Backend_TokenizeString_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _Backend_Status_Handler,
		},
		{
			MethodName: "StoresSet",
			Handler:    _Backend_StoresSet_Handler,
		},
		{
			MethodName: "StoresDelete",
			Handler:    _Backend_StoresDelete_Handler,
		},
		{
			MethodName: "StoresGet",
			Handler:    _Backend_StoresGet_Handler,
		},
		{
			MethodName: "StoresFind",
			Handler:    _Backend_StoresFind_Handler,
		},
		{
			MethodName: "Rerank",
			Handler:    _Backend_Rerank_Handler,
		},
		{
			MethodName: "GetMetrics",
			Handler:    _Backend_GetMetrics_Handler,
		},
		{
			MethodName: "VAD",
			Handler:    _Backend_VAD_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PredictStream",
			Handler:       _Backend_PredictStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "backend.proto",
}
