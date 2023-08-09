package api

import (
	"errors"
	"strings"

	config "github.com/go-skynet/LocalAI/api/config"
	"github.com/go-skynet/LocalAI/api/localai"
	"github.com/go-skynet/LocalAI/api/openai"
	"github.com/go-skynet/LocalAI/api/options"
	"github.com/go-skynet/LocalAI/internal"
	"github.com/go-skynet/LocalAI/pkg/assets"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func App(opts ...options.AppOption) (*fiber.App, error) {
	options := options.NewOptions(opts...)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if options.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Return errors as JSON responses
	app := fiber.New(fiber.Config{
		BodyLimit:             options.UploadLimitMB * 1024 * 1024, // this is the default limit of 4MB
		DisableStartupMessage: options.DisableMessage,
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Send custom error page
			return ctx.Status(code).JSON(
				openai.ErrorResponse{
					Error: &openai.APIError{Message: err.Error(), Code: code},
				},
			)
		},
	})

	if options.Debug {
		app.Use(logger.New(logger.Config{
			Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
		}))
	}

	log.Info().Msgf("Starting LocalAI using %d threads, with models path: %s", options.Threads, options.Loader.ModelPath)
	log.Info().Msgf("LocalAI version: %s", internal.PrintableVersion())

	cm := config.NewConfigLoader()
	if err := cm.LoadConfigs(options.Loader.ModelPath); err != nil {
		log.Error().Msgf("error loading config files: %s", err.Error())
	}

	if options.ConfigFile != "" {
		if err := cm.LoadConfigFile(options.ConfigFile); err != nil {
			log.Error().Msgf("error loading config file: %s", err.Error())
		}
	}

	if options.Debug {
		for _, v := range cm.ListConfigs() {
			cfg, _ := cm.GetConfig(v)
			log.Debug().Msgf("Model: %s (config: %+v)", v, cfg)
		}
	}

	if options.AssetsDestination != "" {
		// Extract files from the embedded FS
		err := assets.ExtractFiles(options.BackendAssets, options.AssetsDestination)
		log.Debug().Msgf("Extracting backend assets files to %s", options.AssetsDestination)
		if err != nil {
			log.Warn().Msgf("Failed extracting backend assets files: %s (might be required for some backends to work properly, like gpt4all)", err)
		}
	}

	// Default middleware config
	app.Use(recover.New())

	// Auth middleware checking if API key is valid. If no API key is set, no auth is required.
	auth := func(c *fiber.Ctx) error {
		if len(options.ApiKeys) > 0 {
			authHeader := c.Get("Authorization")
			if authHeader == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Authorization header missing"})
			}
			authHeaderParts := strings.Split(authHeader, " ")
			if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Authorization header format"})
			}

			apiKey := authHeaderParts[1]
			validApiKey := false
			for _, key := range options.ApiKeys {
				if apiKey == key {
					validApiKey = true
				}
			}
			if !validApiKey {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid API key"})
			}
		}
		return c.Next()
	}

	if options.PreloadJSONModels != "" {
		if err := localai.ApplyGalleryFromString(options.Loader.ModelPath, options.PreloadJSONModels, cm, options.Galleries); err != nil {
			return nil, err
		}
	}

	if options.PreloadModelsFromPath != "" {
		if err := localai.ApplyGalleryFromFile(options.Loader.ModelPath, options.PreloadModelsFromPath, cm, options.Galleries); err != nil {
			return nil, err
		}
	}

	if options.CORS {
		var c func(ctx *fiber.Ctx) error
		if options.CORSAllowOrigins == "" {
			c = cors.New()
		} else {
			c = cors.New(cors.Config{AllowOrigins: options.CORSAllowOrigins})
		}

		app.Use(c)
	}

	// LocalAI API endpoints
	galleryService := localai.NewGalleryService(options.Loader.ModelPath)
	galleryService.Start(options.Context, cm)

	app.Get("/version", auth, func(c *fiber.Ctx) error {
		return c.JSON(struct {
			Version string `json:"version"`
		}{Version: internal.PrintableVersion()})
	})

	app.Post("/models/apply", auth, localai.ApplyModelGalleryEndpoint(options.Loader.ModelPath, cm, galleryService.C, options.Galleries))
	app.Get("/models/available", auth, localai.ListModelFromGalleryEndpoint(options.Galleries, options.Loader.ModelPath))
	app.Get("/models/jobs/:uuid", auth, localai.GetOpStatusEndpoint(galleryService))

	// openAI compatible API endpoint

	// chat
	app.Post("/v1/chat/completions", auth, openai.ChatEndpoint(cm, options))
	app.Post("/chat/completions", auth, openai.ChatEndpoint(cm, options))

	// edit
	app.Post("/v1/edits", auth, openai.EditEndpoint(cm, options))
	app.Post("/edits", auth, openai.EditEndpoint(cm, options))

	// completion
	app.Post("/v1/completions", auth, openai.CompletionEndpoint(cm, options))
	app.Post("/completions", auth, openai.CompletionEndpoint(cm, options))
	app.Post("/v1/engines/:model/completions", auth, openai.CompletionEndpoint(cm, options))

	// embeddings
	app.Post("/v1/embeddings", auth, openai.EmbeddingsEndpoint(cm, options))
	app.Post("/embeddings", auth, openai.EmbeddingsEndpoint(cm, options))
	app.Post("/v1/engines/:model/embeddings", auth, openai.EmbeddingsEndpoint(cm, options))

	// audio
	app.Post("/v1/audio/transcriptions", auth, openai.TranscriptEndpoint(cm, options))
	app.Post("/tts", auth, localai.TTSEndpoint(cm, options))

	// images
	app.Post("/v1/images/generations", auth, openai.ImageEndpoint(cm, options))

	if options.ImageDir != "" {
		app.Static("/generated-images", options.ImageDir)
	}

	if options.AudioDir != "" {
		app.Static("/generated-audio", options.AudioDir)
	}

	ok := func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	}

	// Kubernetes health checks
	app.Get("/healthz", ok)
	app.Get("/readyz", ok)

	// models
	app.Get("/v1/models", auth, openai.ListModelsEndpoint(options.Loader, cm))
	app.Get("/models", auth, openai.ListModelsEndpoint(options.Loader, cm))

	// turn off any process that was started by GRPC if the context is canceled
	go func() {
		<-options.Context.Done()
		log.Debug().Msgf("Context canceled, shutting down")
		options.Loader.StopGRPC()
	}()

	return app, nil
}
