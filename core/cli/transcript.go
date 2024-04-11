package cli

import (
	"context"
	"errors"
	"fmt"

<<<<<<< HEAD
	"github.com/mudler/LocalAI/core/backend"
	cliContext "github.com/mudler/LocalAI/core/cli/context"
	"github.com/mudler/LocalAI/core/config"
	"github.com/mudler/LocalAI/pkg/model"
	"github.com/rs/zerolog/log"
=======
	"github.com/go-skynet/LocalAI/core/backend"
	"github.com/go-skynet/LocalAI/core/config"
	"github.com/go-skynet/LocalAI/pkg/model"
>>>>>>> e16d5918 (feat: kong cli refactor fixes #1955 (#1974))
)

type TranscriptCMD struct {
	Filename string `arg:""`

	Backend           string `short:"b" default:"whisper" help:"Backend to run the transcription model"`
	Model             string `short:"m" required:"" help:"Model name to run the TTS"`
	Language          string `short:"l" help:"Language of the audio file"`
<<<<<<< HEAD
	Translate         bool   `short:"c" help:"Translate the transcription to english"`
=======
>>>>>>> e16d5918 (feat: kong cli refactor fixes #1955 (#1974))
	Threads           int    `short:"t" default:"1" help:"Number of threads used for parallel computation"`
	ModelsPath        string `env:"LOCALAI_MODELS_PATH,MODELS_PATH" type:"path" default:"${basepath}/models" help:"Path containing models used for inferencing" group:"storage"`
	BackendAssetsPath string `env:"LOCALAI_BACKEND_ASSETS_PATH,BACKEND_ASSETS_PATH" type:"path" default:"/tmp/localai/backend_data" help:"Path used to extract libraries that are required by some of the backends in runtime" group:"storage"`
}

<<<<<<< HEAD
func (t *TranscriptCMD) Run(ctx *cliContext.Context) error {
=======
func (t *TranscriptCMD) Run(ctx *Context) error {
>>>>>>> e16d5918 (feat: kong cli refactor fixes #1955 (#1974))
	opts := &config.ApplicationConfig{
		ModelPath:         t.ModelsPath,
		Context:           context.Background(),
		AssetsDestination: t.BackendAssetsPath,
	}

<<<<<<< HEAD
	cl := config.NewBackendConfigLoader(t.ModelsPath)
	ml := model.NewModelLoader(opts.ModelPath, opts.SingleBackend)
=======
	cl := config.NewBackendConfigLoader()
	ml := model.NewModelLoader(opts.ModelPath)
>>>>>>> e16d5918 (feat: kong cli refactor fixes #1955 (#1974))
	if err := cl.LoadBackendConfigsFromPath(t.ModelsPath); err != nil {
		return err
	}

	c, exists := cl.GetBackendConfig(t.Model)
	if !exists {
		return errors.New("model not found")
	}

	c.Threads = &t.Threads

<<<<<<< HEAD
	defer func() {
		err := ml.StopAllGRPC()
		if err != nil {
			log.Error().Err(err).Msg("unable to stop all grpc processes")
		}
	}()

	tr, err := backend.ModelTranscription(t.Filename, t.Language, t.Translate, ml, c, opts)
=======
	defer ml.StopAllGRPC()

	tr, err := backend.ModelTranscription(t.Filename, t.Language, ml, c, opts)
>>>>>>> e16d5918 (feat: kong cli refactor fixes #1955 (#1974))
	if err != nil {
		return err
	}
	for _, segment := range tr.Segments {
		fmt.Println(segment.Start.String(), "-", segment.Text)
	}
	return nil
}
