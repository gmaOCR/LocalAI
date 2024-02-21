package grpc

import (
	"context"
<<<<<<< HEAD
<<<<<<< HEAD

	pb "github.com/mudler/LocalAI/pkg/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
=======
	"github.com/go-skynet/LocalAI/api/schema"
	pb "github.com/go-skynet/LocalAI/pkg/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
>>>>>>> d6352300 (feat(grpc): backend SPI pluggable in embedding mode (#1621))
=======
	"time"

	"github.com/go-skynet/LocalAI/core/schema"
	pb "github.com/go-skynet/LocalAI/pkg/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
>>>>>>> 68598ebe (MQTT Startup Refactoring Part 1: core/ packages part 1 (#1728))
)

var _ Backend = new(embedBackend)
var _ pb.Backend_PredictStreamServer = new(embedBackendServerStream)

type embedBackend struct {
	s *server
}

func (e *embedBackend) IsBusy() bool {
	return e.s.llm.Busy()
}

func (e *embedBackend) HealthCheck(ctx context.Context) (bool, error) {
	return true, nil
}

func (e *embedBackend) Embeddings(ctx context.Context, in *pb.PredictOptions, opts ...grpc.CallOption) (*pb.EmbeddingResult, error) {
	return e.s.Embedding(ctx, in)
}

func (e *embedBackend) Predict(ctx context.Context, in *pb.PredictOptions, opts ...grpc.CallOption) (*pb.Reply, error) {
	return e.s.Predict(ctx, in)
}

func (e *embedBackend) LoadModel(ctx context.Context, in *pb.ModelOptions, opts ...grpc.CallOption) (*pb.Result, error) {
	return e.s.LoadModel(ctx, in)
}

<<<<<<< HEAD
func (e *embedBackend) PredictStream(ctx context.Context, in *pb.PredictOptions, f func(reply *pb.Reply), opts ...grpc.CallOption) error {
=======
func (e *embedBackend) PredictStream(ctx context.Context, in *pb.PredictOptions, f func(s []byte), opts ...grpc.CallOption) error {
>>>>>>> d6352300 (feat(grpc): backend SPI pluggable in embedding mode (#1621))
	bs := &embedBackendServerStream{
		ctx: ctx,
		fn:  f,
	}
	return e.s.PredictStream(in, bs)
}

func (e *embedBackend) GenerateImage(ctx context.Context, in *pb.GenerateImageRequest, opts ...grpc.CallOption) (*pb.Result, error) {
	return e.s.GenerateImage(ctx, in)
}

<<<<<<< HEAD
func (e *embedBackend) GenerateVideo(ctx context.Context, in *pb.GenerateVideoRequest, opts ...grpc.CallOption) (*pb.Result, error) {
	return e.s.GenerateVideo(ctx, in)
}

=======
>>>>>>> d6352300 (feat(grpc): backend SPI pluggable in embedding mode (#1621))
func (e *embedBackend) TTS(ctx context.Context, in *pb.TTSRequest, opts ...grpc.CallOption) (*pb.Result, error) {
	return e.s.TTS(ctx, in)
}

<<<<<<< HEAD
func (e *embedBackend) SoundGeneration(ctx context.Context, in *pb.SoundGenerationRequest, opts ...grpc.CallOption) (*pb.Result, error) {
	return e.s.SoundGeneration(ctx, in)
}

func (e *embedBackend) AudioTranscription(ctx context.Context, in *pb.TranscriptRequest, opts ...grpc.CallOption) (*pb.TranscriptResult, error) {
	return e.s.AudioTranscription(ctx, in)
=======
func (e *embedBackend) AudioTranscription(ctx context.Context, in *pb.TranscriptRequest, opts ...grpc.CallOption) (*schema.Result, error) {
	r, err := e.s.AudioTranscription(ctx, in)
	if err != nil {
		return nil, err
	}
	tr := &schema.Result{}
	for _, s := range r.Segments {
		var tks []int
		for _, t := range s.Tokens {
			tks = append(tks, int(t))
		}
		tr.Segments = append(tr.Segments,
			schema.Segment{
				Text:   s.Text,
				Id:     int(s.Id),
				Start:  time.Duration(s.Start),
				End:    time.Duration(s.End),
				Tokens: tks,
			})
	}
	tr.Text = r.Text
	return tr, err
>>>>>>> d6352300 (feat(grpc): backend SPI pluggable in embedding mode (#1621))
}

func (e *embedBackend) TokenizeString(ctx context.Context, in *pb.PredictOptions, opts ...grpc.CallOption) (*pb.TokenizationResponse, error) {
	return e.s.TokenizeString(ctx, in)
}

func (e *embedBackend) Status(ctx context.Context) (*pb.StatusResponse, error) {
	return e.s.Status(ctx, &pb.HealthMessage{})
}

<<<<<<< HEAD
func (e *embedBackend) StoresSet(ctx context.Context, in *pb.StoresSetOptions, opts ...grpc.CallOption) (*pb.Result, error) {
	return e.s.StoresSet(ctx, in)
}

func (e *embedBackend) StoresDelete(ctx context.Context, in *pb.StoresDeleteOptions, opts ...grpc.CallOption) (*pb.Result, error) {
	return e.s.StoresDelete(ctx, in)
}

func (e *embedBackend) StoresGet(ctx context.Context, in *pb.StoresGetOptions, opts ...grpc.CallOption) (*pb.StoresGetResult, error) {
	return e.s.StoresGet(ctx, in)
}

func (e *embedBackend) StoresFind(ctx context.Context, in *pb.StoresFindOptions, opts ...grpc.CallOption) (*pb.StoresFindResult, error) {
	return e.s.StoresFind(ctx, in)
}

func (e *embedBackend) Rerank(ctx context.Context, in *pb.RerankRequest, opts ...grpc.CallOption) (*pb.RerankResult, error) {
	return e.s.Rerank(ctx, in)
}

func (e *embedBackend) VAD(ctx context.Context, in *pb.VADRequest, opts ...grpc.CallOption) (*pb.VADResponse, error) {
	return e.s.VAD(ctx, in)
}

func (e *embedBackend) GetTokenMetrics(ctx context.Context, in *pb.MetricsRequest, opts ...grpc.CallOption) (*pb.MetricsResponse, error) {
	return e.s.GetMetrics(ctx, in)
}

type embedBackendServerStream struct {
	ctx context.Context
	fn  func(reply *pb.Reply)
}

func (e *embedBackendServerStream) Send(reply *pb.Reply) error {
	e.fn(reply)
=======
type embedBackendServerStream struct {
	ctx context.Context
	fn  func(s []byte)
}

func (e *embedBackendServerStream) Send(reply *pb.Reply) error {
	e.fn(reply.GetMessage())
>>>>>>> d6352300 (feat(grpc): backend SPI pluggable in embedding mode (#1621))
	return nil
}

func (e *embedBackendServerStream) SetHeader(md metadata.MD) error {
	return nil
}

func (e *embedBackendServerStream) SendHeader(md metadata.MD) error {
	return nil
}

func (e *embedBackendServerStream) SetTrailer(md metadata.MD) {
}

func (e *embedBackendServerStream) Context() context.Context {
	return e.ctx
}

func (e *embedBackendServerStream) SendMsg(m any) error {
	if x, ok := m.(*pb.Reply); ok {
		return e.Send(x)
	}
	return nil
}

func (e *embedBackendServerStream) RecvMsg(m any) error {
	return nil
}
