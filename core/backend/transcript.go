package backend

import (
	"context"
	"fmt"
<<<<<<< HEAD:core/backend/transcript.go
	"time"

	"github.com/mudler/LocalAI/core/config"
	"github.com/mudler/LocalAI/core/schema"

	"github.com/mudler/LocalAI/pkg/grpc/proto"
	"github.com/mudler/LocalAI/pkg/model"
)

func ModelTranscription(audio, language string, translate bool, ml *model.ModelLoader, backendConfig config.BackendConfig, appConfig *config.ApplicationConfig) (*schema.TranscriptionResult, error) {

	if backendConfig.Backend == "" {
		backendConfig.Backend = model.WhisperBackend
	}

	opts := ModelOptions(backendConfig, appConfig)

	transcriptionModel, err := ml.Load(opts...)
	if err != nil {
		return nil, err
	}
	defer ml.Close()

	if transcriptionModel == nil {
		return nil, fmt.Errorf("could not load transcription model")
	}

	r, err := transcriptionModel.AudioTranscription(context.Background(), &proto.TranscriptRequest{
		Dst:       audio,
		Language:  language,
		Translate: translate,
		Threads:   uint32(*backendConfig.Threads),
	})
	if err != nil {
		return nil, err
	}
	tr := &schema.TranscriptionResult{
		Text: r.Text,
	}
	for _, s := range r.Segments {
		var tks []int
		for _, t := range s.Tokens {
			tks = append(tks, int(t))
		}
		tr.Segments = append(tr.Segments,
			schema.TranscriptionSegment{
				Text:   s.Text,
				Id:     int(s.Id),
				Start:  time.Duration(s.Start),
				End:    time.Duration(s.End),
				Tokens: tks,
			})
	}
	return tr, err
=======

	config "github.com/go-skynet/LocalAI/core/config"
	"github.com/go-skynet/LocalAI/core/schema"

	"github.com/go-skynet/LocalAI/core/options"
	"github.com/go-skynet/LocalAI/pkg/grpc/proto"
	model "github.com/go-skynet/LocalAI/pkg/model"
)

func ModelTranscription(audio, language string, loader *model.ModelLoader, c config.Config, o *options.Option) (*schema.Result, error) {

	opts := modelOpts(c, o, []model.Option{
		model.WithBackendString(model.WhisperBackend),
		model.WithModel(c.Model),
		model.WithContext(o.Context),
		model.WithThreads(uint32(c.Threads)),
		model.WithAssetDir(o.AssetsDestination),
	})

	whisperModel, err := o.Loader.BackendLoader(opts...)
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
>>>>>>> 5f2b87fa (Revert "[Refactor]: Core/API Split" (#1550)):api/backend/transcript.go
}
