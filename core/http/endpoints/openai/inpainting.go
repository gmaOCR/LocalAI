package openai

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/mudler/LocalAI/core/backend"
	"github.com/mudler/LocalAI/core/config"
	"github.com/mudler/LocalAI/core/http/middleware"
	"github.com/mudler/LocalAI/core/schema"
	model "github.com/mudler/LocalAI/pkg/model"
)

// InpaintingEndpoint handles POST /v1/images/inpainting
//
// Swagger / OpenAPI docstring (swaggo):
// @Summary      Image inpainting
// @Description  Perform image inpainting. Accepts multipart/form-data with `image` and `mask` files.
// @Tags         images
// @Accept       multipart/form-data
// @Produce      application/json
// @Param        model   formData  string  true   "Model identifier"
// @Param        prompt  formData  string  true   "Text prompt guiding the generation"
// @Param        steps   formData  int     false  "Number of inference steps (default 25)"
// @Param        image   formData  file    true   "Original image file"
// @Param        mask    formData  file    true   "Mask image file (white = area to inpaint)"
// @Success      200 {object} schema.OpenAIResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /v1/images/inpainting [post]
func InpaintingEndpoint(cl *config.ModelConfigLoader, ml *model.ModelLoader, appConfig *config.ApplicationConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse basic form values
		modelName := c.FormValue("model")
		prompt := c.FormValue("prompt")
		stepsStr := c.FormValue("steps")

		if modelName == "" || prompt == "" {
			log.Error().Msg("Inpainting Endpoint - missing model or prompt")
			return echo.ErrBadRequest
		}

		// steps default
		steps := 25
		if stepsStr != "" {
			if v, err := strconv.Atoi(stepsStr); err == nil {
				steps = v
			}
		}

		// Get uploaded files
		imageFile, err := c.FormFile("image")
		if err != nil {
			log.Error().Err(err).Msg("Inpainting Endpoint - missing image file")
			return echo.NewHTTPError(http.StatusBadRequest, "missing image file")
		}
		maskFile, err := c.FormFile("mask")
		if err != nil {
			log.Error().Err(err).Msg("Inpainting Endpoint - missing mask file")
			return echo.NewHTTPError(http.StatusBadRequest, "missing mask file")
		}

		// Read files into memory (small files expected)
		imgSrc, err := imageFile.Open()
		if err != nil {
			return err
		}
		defer imgSrc.Close()
		imgBytes, err := io.ReadAll(imgSrc)
		if err != nil {
			return err
		}

		maskSrc, err := maskFile.Open()
		if err != nil {
			return err
		}
		defer maskSrc.Close()
		maskBytes, err := io.ReadAll(maskSrc)
		if err != nil {
			return err
		}

		// Create JSON with base64 fields expected by backend
		b64Image := base64.StdEncoding.EncodeToString(imgBytes)
		b64Mask := base64.StdEncoding.EncodeToString(maskBytes)

		// get model config from context (middleware set it)
		cfg, ok := c.Get("MODEL_CONFIG").(*config.ModelConfig)
		if !ok || cfg == nil {
			log.Error().Msg("Inpainting Endpoint - model config not found in context")
			return echo.ErrBadRequest
		}

		// Use the images subdirectory under GeneratedContentDir so the generated
		// PNG is placed where the HTTP static handler serves `/generated-images`.
		tmpDir := filepath.Join(appConfig.GeneratedContentDir, "images")
		id := uuid.New().String()
		jsonName := fmt.Sprintf("inpaint_%s.json", id)
		jsonPath := filepath.Join(tmpDir, jsonName)
		jsonFile := map[string]string{
			"image":      b64Image,
			"mask_image": b64Mask,
		}
		jf, err := os.CreateTemp(tmpDir, "inpaint_")
		if err != nil {
			return err
		}
		// write JSON
		enc := json.NewEncoder(jf)
		if err := enc.Encode(jsonFile); err != nil {
			jf.Close()
			os.Remove(jf.Name())
			return err
		}
		jf.Close()
		// rename to desired name
		if err := os.Rename(jf.Name(), jsonPath); err != nil {
			os.Remove(jf.Name())
			return err
		}
		// prepare dst
		outTmp, err := os.CreateTemp(tmpDir, "out_")
		if err != nil {
			os.Remove(jsonPath)
			return err
		}
		outTmp.Close()
		dst := outTmp.Name() + ".png"
		if err := os.Rename(outTmp.Name(), dst); err != nil {
			os.Remove(jsonPath)
			return err
		}

		// Determine width/height default
		width := 512
		height := 512

		// Call backend image generation via indirection so tests can stub it
		// Note: ImageGenerationFunc will call into the loaded model's GenerateImage which expects src JSON
		fn, err := backend.ImageGenerationFunc(height, width, 0, steps, 0, prompt, "", jsonPath, dst, ml, *cfg, appConfig, nil)
		if err != nil {
			os.Remove(jsonPath)
			return err
		}

		// Execute generation function (blocking)
		if err := fn(); err != nil {
			os.Remove(jsonPath)
			os.Remove(dst)
			return err
		}

		// On success, build response URL using BaseURL middleware helper and
		// the same `generated-images` prefix used by the server static mount.
		baseURL := middleware.BaseURL(c)

		// Build response using url.JoinPath for correct URL escaping
		imgPath, err := url.JoinPath(baseURL, "generated-images", filepath.Base(dst))
		if err != nil {
			return err
		}

		created := int(time.Now().Unix())
		resp := &schema.OpenAIResponse{
			ID:      id,
			Created: created,
			Data: []schema.Item{{
				URL: imgPath,
			}},
		}

		// cleanup json
		defer os.Remove(jsonPath)

		return c.JSON(http.StatusOK, resp)
	}
}
