package openai

import (
<<<<<<< HEAD:core/http/endpoints/openai/inference.go
	"github.com/mudler/LocalAI/core/backend"
	"github.com/mudler/LocalAI/core/config"

	"github.com/mudler/LocalAI/core/schema"
	model "github.com/mudler/LocalAI/pkg/model"
=======
	"github.com/go-skynet/LocalAI/api/backend"
	config "github.com/go-skynet/LocalAI/api/config"
	"github.com/go-skynet/LocalAI/api/options"
	"github.com/go-skynet/LocalAI/api/schema"
	model "github.com/go-skynet/LocalAI/pkg/model"
>>>>>>> 5f2b87fa (Revert "[Refactor]: Core/API Split" (#1550)):api/openai/inference.go
)

func ComputeChoices(
	req *schema.OpenAIRequest,
	predInput string,
<<<<<<< HEAD:core/http/endpoints/openai/inference.go
	config *config.BackendConfig,
	bcl *config.BackendConfigLoader,
	o *config.ApplicationConfig,
=======
	config *config.Config,
	o *options.Option,
>>>>>>> 5f2b87fa (Revert "[Refactor]: Core/API Split" (#1550)):api/openai/inference.go
	loader *model.ModelLoader,
	cb func(string, *[]schema.Choice),
	tokenCallback func(string, backend.TokenUsage) bool) ([]schema.Choice, backend.TokenUsage, error) {
	n := req.N // number of completions to return
	result := []schema.Choice{}

	if n == 0 {
		n = 1
	}

	images := []string{}
	for _, m := range req.Messages {
		images = append(images, m.StringImages...)
	}
<<<<<<< HEAD:core/http/endpoints/openai/inference.go
	videos := []string{}
	for _, m := range req.Messages {
		videos = append(videos, m.StringVideos...)
	}
	audios := []string{}
	for _, m := range req.Messages {
		audios = append(audios, m.StringAudios...)
	}

	// get the model function to call for the result
	predFunc, err := backend.ModelInference(req.Context, predInput, req.Messages, images, videos, audios, loader, config, bcl, o, tokenCallback)
=======

	// get the model function to call for the result
	predFunc, err := backend.ModelInference(req.Context, predInput, images, loader, *config, o, tokenCallback)
>>>>>>> 5f2b87fa (Revert "[Refactor]: Core/API Split" (#1550)):api/openai/inference.go
	if err != nil {
		return result, backend.TokenUsage{}, err
	}

	tokenUsage := backend.TokenUsage{}

	for i := 0; i < n; i++ {
		prediction, err := predFunc()
		if err != nil {
			return result, backend.TokenUsage{}, err
		}

		tokenUsage.Prompt += prediction.Usage.Prompt
		tokenUsage.Completion += prediction.Usage.Completion
<<<<<<< HEAD:core/http/endpoints/openai/inference.go
		tokenUsage.TimingPromptProcessing += prediction.Usage.TimingPromptProcessing
		tokenUsage.TimingTokenGeneration += prediction.Usage.TimingTokenGeneration
=======
>>>>>>> 5f2b87fa (Revert "[Refactor]: Core/API Split" (#1550)):api/openai/inference.go

		finetunedResponse := backend.Finetune(*config, predInput, prediction.Response)
		cb(finetunedResponse, &result)

		//result = append(result, Choice{Text: prediction})

	}
	return result, tokenUsage, err
}
