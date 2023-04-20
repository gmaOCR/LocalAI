package model

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/rs/zerolog/log"

	gptj "github.com/go-skynet/go-gpt4all-j.cpp"
	llama "github.com/go-skynet/go-llama.cpp"
)

type ModelLoader struct {
	modelPath        string
	mu               sync.Mutex
	models           map[string]*llama.LLama
	gptmodels        map[string]*gptj.GPTJ
	promptsTemplates map[string]*template.Template
}

func NewModelLoader(modelPath string) *ModelLoader {
	return &ModelLoader{modelPath: modelPath, gptmodels: make(map[string]*gptj.GPTJ), models: make(map[string]*llama.LLama), promptsTemplates: make(map[string]*template.Template)}
}

func (ml *ModelLoader) ExistsInModelPath(s string) bool {
	_, err := os.Stat(filepath.Join(ml.modelPath, s))
	return err == nil
}

func (ml *ModelLoader) ListModels() ([]string, error) {
	files, err := ioutil.ReadDir(ml.modelPath)
	if err != nil {
		return []string{}, err
	}

	models := []string{}
	for _, file := range files {
		// Skip templates, YAML and .keep files
		if strings.HasSuffix(file.Name(), ".tmpl") || strings.HasSuffix(file.Name(), ".keep") || strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml") {
			continue
		}

		models = append(models, file.Name())
	}

	return models, nil
}

func (ml *ModelLoader) TemplatePrefix(modelName string, in interface{}) (string, error) {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	m, ok := ml.promptsTemplates[modelName]
	if !ok {
		return "", fmt.Errorf("no prompt template available")
	}

	var buf bytes.Buffer

	if err := m.Execute(&buf, in); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (ml *ModelLoader) loadTemplateIfExists(modelName, modelFile string) error {
	// Check if the template was already loaded
	if _, ok := ml.promptsTemplates[modelName]; ok {
		return nil
	}

	// Check if the model path exists
	// skip any error here - we run anyway if a template is not exist
	modelTemplateFile := fmt.Sprintf("%s.tmpl", modelName)

	if !ml.ExistsInModelPath(modelTemplateFile) {
		return nil
	}

	dat, err := os.ReadFile(filepath.Join(ml.modelPath, modelTemplateFile))
	if err != nil {
		return err
	}

	// Parse the template
	tmpl, err := template.New("prompt").Parse(string(dat))
	if err != nil {
		return err
	}
	ml.promptsTemplates[modelName] = tmpl

	return nil
}

func (ml *ModelLoader) LoadGPTJModel(modelName string) (*gptj.GPTJ, error) {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	// Check if we already have a loaded model
	if !ml.ExistsInModelPath(modelName) {
		return nil, fmt.Errorf("model does not exist")
	}

	if m, ok := ml.gptmodels[modelName]; ok {
		log.Debug().Msgf("Model already loaded in memory: %s", modelName)
		return m, nil
	}

	// Load the model and keep it in memory for later use
	modelFile := filepath.Join(ml.modelPath, modelName)
	log.Debug().Msgf("Loading model in memory from file: %s", modelFile)

	model, err := gptj.New(modelFile)
	if err != nil {
		return nil, err
	}

	// If there is a prompt template, load it
	if err := ml.loadTemplateIfExists(modelName, modelFile); err != nil {
		return nil, err
	}

	ml.gptmodels[modelName] = model
	return model, err
}

func (ml *ModelLoader) LoadLLaMAModel(modelName string, opts ...llama.ModelOption) (*llama.LLama, error) {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	log.Debug().Msgf("Loading model name: %s", modelName)

	// Check if we already have a loaded model
	if !ml.ExistsInModelPath(modelName) {
		return nil, fmt.Errorf("model does not exist")
	}

	if m, ok := ml.models[modelName]; ok {
		log.Debug().Msgf("Model already loaded in memory: %s", modelName)
		return m, nil
	}

	// TODO: This needs refactoring, it's really bad to have it in here
	// Check if we have a GPTJ model loaded instead - if we do we return an error so the API tries with GPTJ
	if _, ok := ml.gptmodels[modelName]; ok {
		log.Debug().Msgf("Model is GPTJ: %s", modelName)
		return nil, fmt.Errorf("this model is a GPTJ one")
	}

	// Load the model and keep it in memory for later use
	modelFile := filepath.Join(ml.modelPath, modelName)
	log.Debug().Msgf("Loading model in memory from file: %s", modelFile)

	model, err := llama.New(modelFile, opts...)
	if err != nil {
		return nil, err
	}

	// If there is a prompt template, load it
	if err := ml.loadTemplateIfExists(modelName, modelFile); err != nil {
		return nil, err
	}

	ml.models[modelName] = model
	return model, err
}
