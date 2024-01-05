package backend

import (
	"context"
	"fmt"

	"github.com/go-skynet/LocalAI/core/services"
	"github.com/go-skynet/LocalAI/pkg/grpc/proto"
	"github.com/go-skynet/LocalAI/pkg/model"
	"github.com/go-skynet/LocalAI/pkg/schema"
)

func ModelTranscription(audio, language string, loader *model.ModelLoader, c schema.Config, o *schema.StartupOptions) (*schema.WhisperResult, error) {

	opts := modelOpts(c, o, []model.Option{
		model.WithBackendString(model.WhisperBackend),
		model.WithModel(c.Model),
		model.WithContext(o.Context),
		model.WithThreads(uint32(c.Threads)),
		model.WithAssetDir(o.AssetsDestination),
		model.WithExternalBackends(o.ExternalGRPCBackends, false),
	})

	whisperModel, err := loader.BackendLoader(opts...)
	if err != nil {
		return nil, err
	}

	if whisperModel == nil {
		return nil, fmt.Errorf("could not load whisper model")
	}

	return whisperModel.AudioTranscription(context.Background(), &proto.TranscriptRequest{
		Dst:      audio,
		Language: language,
		Threads:  uint32(c.Threads),
	})
}

func TranscriptionOpenAIRequest(modelName string, input *schema.OpenAIRequest, audioFilePath string, cl *services.ConfigLoader, ml *model.ModelLoader, startupOptions *schema.StartupOptions) (*schema.WhisperResult, error) {
	config, input, err := ReadConfigFromFileAndCombineWithOpenAIRequest(modelName, input, cl, startupOptions)
	if err != nil {
		return nil, fmt.Errorf("failed reading parameters from request:%w", err)
	}

	tr, err := ModelTranscription(audioFilePath, input.Language, ml, *config, startupOptions)
	if err != nil {
		return nil, err
	}

	return tr, nil
}
