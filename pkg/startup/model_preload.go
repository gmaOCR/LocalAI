package startup

import (
	"errors"
<<<<<<< HEAD
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mudler/LocalAI/core/config"
	"github.com/mudler/LocalAI/core/gallery"
	"github.com/mudler/LocalAI/pkg/downloader"
	"github.com/mudler/LocalAI/pkg/utils"
	"github.com/rs/zerolog/log"
)

// InstallModels will preload models from the given list of URLs and galleries
// It will download the model if it is not already present in the model path
// It will also try to resolve if the model is an embedded model YAML configuration
func InstallModels(galleries, backendGalleries []config.Gallery, modelPath, backendBasePath string, enforceScan, autoloadBackendGalleries bool, downloadStatus func(string, string, string, float64), models ...string) error {
	// create an error that groups all errors
	var err error

	for _, url := range models {
		// As a best effort, try to resolve the model from the remote library
		// if it's not resolved we try with the other method below

		uri := downloader.URI(url)

		switch {
		case uri.LooksLikeOCI():
			log.Debug().Msgf("[startup] resolved OCI model to download: %s", url)

			// convert OCI image name to a file name.
			ociName := strings.TrimPrefix(url, downloader.OCIPrefix)
			ociName = strings.TrimPrefix(ociName, downloader.OllamaPrefix)
			ociName = strings.ReplaceAll(ociName, "/", "__")
			ociName = strings.ReplaceAll(ociName, ":", "__")

			// check if file exists
			if _, e := os.Stat(filepath.Join(modelPath, ociName)); errors.Is(e, os.ErrNotExist) {
				modelDefinitionFilePath := filepath.Join(modelPath, ociName)
				e := uri.DownloadFile(modelDefinitionFilePath, "", 0, 0, func(fileName, current, total string, percent float64) {
					utils.DisplayDownloadFunction(fileName, current, total, percent)
				})
				if e != nil {
					log.Error().Err(e).Str("url", url).Str("filepath", modelDefinitionFilePath).Msg("error downloading model")
					err = errors.Join(err, e)
				}
			}

			log.Info().Msgf("[startup] installed model from OCI repository: %s", ociName)
		case uri.LooksLikeURL():
			log.Debug().Msgf("[startup] downloading %s", url)

			// Extract filename from URL
			fileName, e := uri.FilenameFromUrl()
			if e != nil {
				log.Warn().Err(e).Str("url", url).Msg("error extracting filename from URL")
				err = errors.Join(err, e)
				continue
			}

			modelPath := filepath.Join(modelPath, fileName)

			if e := utils.VerifyPath(fileName, modelPath); e != nil {
				log.Error().Err(e).Str("filepath", modelPath).Msg("error verifying path")
				err = errors.Join(err, e)
				continue
			}

			// check if file exists
			if _, e := os.Stat(modelPath); errors.Is(e, os.ErrNotExist) {
				e := uri.DownloadFile(modelPath, "", 0, 0, func(fileName, current, total string, percent float64) {
					utils.DisplayDownloadFunction(fileName, current, total, percent)
				})
				if e != nil {
					log.Error().Err(e).Str("url", url).Str("filepath", modelPath).Msg("error downloading model")
					err = errors.Join(err, e)
				}
			}
		default:
			if _, e := os.Stat(url); e == nil {
				log.Debug().Msgf("[startup] resolved local model: %s", url)
				// copy to modelPath
				md5Name := utils.MD5(url)

				modelYAML, e := os.ReadFile(url)
				if e != nil {
					log.Error().Err(e).Str("filepath", url).Msg("error reading model definition")
					err = errors.Join(err, e)
					continue
				}

				modelDefinitionFilePath := filepath.Join(modelPath, md5Name) + ".yaml"
				if e := os.WriteFile(modelDefinitionFilePath, modelYAML, 0600); e != nil {
					log.Error().Err(err).Str("filepath", modelDefinitionFilePath).Msg("error loading model: %s")
					err = errors.Join(err, e)
				}
			} else {
				// Check if it's a model gallery, or print a warning
				e, found := installModel(galleries, backendGalleries, url, modelPath, backendBasePath, downloadStatus, enforceScan, autoloadBackendGalleries)
				if e != nil && found {
					log.Error().Err(err).Msgf("[startup] failed installing model '%s'", url)
					err = errors.Join(err, e)
				} else if !found {
					log.Warn().Msgf("[startup] failed resolving model '%s'", url)
					err = errors.Join(err, fmt.Errorf("failed resolving model '%s'", url))
				}
			}
		}
	}
	return err
}

func installModel(galleries, backendGalleries []config.Gallery, modelName, modelPath, backendBasePath string, downloadStatus func(string, string, string, float64), enforceScan, autoloadBackendGalleries bool) (error, bool) {
	models, err := gallery.AvailableGalleryModels(galleries, modelPath)
	if err != nil {
		return err, false
	}

	model := gallery.FindGalleryElement(models, modelName, modelPath)
	if model == nil {
		return err, false
	}

	if downloadStatus == nil {
		downloadStatus = utils.DisplayDownloadFunction
	}

	log.Info().Str("model", modelName).Str("license", model.License).Msg("installing model")
	err = gallery.InstallModelFromGallery(galleries, backendGalleries, modelName, modelPath, backendBasePath, gallery.GalleryModel{}, downloadStatus, enforceScan, autoloadBackendGalleries)
	if err != nil {
		return err, true
	}

	return nil, true
=======
	"os"
	"path/filepath"

	"github.com/go-skynet/LocalAI/embedded"
	"github.com/go-skynet/LocalAI/pkg/downloader"
	"github.com/go-skynet/LocalAI/pkg/utils"
	"github.com/rs/zerolog/log"
)

// PreloadModelsConfigurations will preload models from the given list of URLs
// It will download the model if it is not already present in the model path
// It will also try to resolve if the model is an embedded model YAML configuration
func PreloadModelsConfigurations(modelLibraryURL string, modelPath string, models ...string) {
	for _, url := range models {

		// As a best effort, try to resolve the model from the remote library
		// if it's not resolved we try with the other method below
		if modelLibraryURL != "" {
			lib, err := embedded.GetRemoteLibraryShorteners(modelLibraryURL)
			if err == nil {
				if lib[url] != "" {
					log.Debug().Msgf("[startup] model configuration is defined remotely: %s (%s)", url, lib[url])
					url = lib[url]
				}
			}
		}

		url = embedded.ModelShortURL(url)
		switch {
		case embedded.ExistsInModelsLibrary(url):
			modelYAML, err := embedded.ResolveContent(url)
			// If we resolve something, just save it to disk and continue
			if err != nil {
				log.Error().Msgf("error loading model: %s", err.Error())
				continue
			}

			log.Debug().Msgf("[startup] resolved embedded model: %s", url)
			md5Name := utils.MD5(url)
			if err := os.WriteFile(filepath.Join(modelPath, md5Name)+".yaml", modelYAML, os.ModePerm); err != nil {
				log.Error().Msgf("error loading model: %s", err.Error())
			}
		case downloader.LooksLikeURL(url):
			log.Debug().Msgf("[startup] resolved model to download: %s", url)

			// md5 of model name
			md5Name := utils.MD5(url)

			// check if file exists
			if _, err := os.Stat(filepath.Join(modelPath, md5Name)); errors.Is(err, os.ErrNotExist) {
				err := downloader.DownloadFile(url, filepath.Join(modelPath, md5Name)+".yaml", "", func(fileName, current, total string, percent float64) {
					utils.DisplayDownloadFunction(fileName, current, total, percent)
				})
				if err != nil {
					log.Error().Msgf("error loading model: %s", err.Error())
				}
			}
		default:
			log.Warn().Msgf("[startup] failed resolving model '%s'", url)
		}
	}
>>>>>>> fda6bf56 (feat: embedded model configurations, add popular model examples, refactoring (#1532))
}
