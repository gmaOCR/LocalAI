package localai

import (
	"github.com/go-skynet/LocalAI/core/backend"
	"github.com/go-skynet/LocalAI/core/services"
	"github.com/go-skynet/LocalAI/pkg/model"
	"github.com/go-skynet/LocalAI/pkg/schema"
	"github.com/gofiber/fiber/v2"
)

func TTSEndpoint(cl *services.ConfigLoader, ml *model.ModelLoader, so *schema.StartupOptions) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		input := new(schema.TTSRequest)
		// Get input data from the request body
		if err := c.BodyParser(input); err != nil {
			return err
		}

		filePath, _, err := backend.ModelTTS(input.Backend, input.Input, input.Model, ml, so)
		if err != nil {
			return err
		}
		return c.Download(filePath)
	}
}
