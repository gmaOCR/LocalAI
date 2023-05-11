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

	rwkv "github.com/donomii/go-rwkv.cpp"
	whisper "github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	bloomz "github.com/go-skynet/bloomz.cpp"
	bert "github.com/go-skynet/go-bert.cpp"
	gpt2 "github.com/go-skynet/go-gpt2.cpp"
	gptj "github.com/go-skynet/go-gpt4all-j.cpp"
	llama "github.com/go-skynet/go-llama.cpp"
	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"
)

type ModelLoader struct {
	ModelPath string
	mu        sync.Mutex
	// TODO: this needs generics
	models            map[string]*llama.LLama
	gptmodels         map[string]*gptj.GPTJ
	gpt2models        map[string]*gpt2.GPT2
	gptstablelmmodels map[string]*gpt2.StableLM
	dollymodels       map[string]*gpt2.Dolly
	redpajama         map[string]*gpt2.RedPajama
	rwkv              map[string]*rwkv.RwkvState
	bloomz            map[string]*bloomz.Bloomz
	bert              map[string]*bert.Bert
	promptsTemplates  map[string]*template.Template
	whisperModels     map[string]whisper.Model
}

func NewModelLoader(modelPath string) *ModelLoader {
	return &ModelLoader{
		ModelPath:         modelPath,
		gpt2models:        make(map[string]*gpt2.GPT2),
		gptmodels:         make(map[string]*gptj.GPTJ),
		gptstablelmmodels: make(map[string]*gpt2.StableLM),
		dollymodels:       make(map[string]*gpt2.Dolly),
		redpajama:         make(map[string]*gpt2.RedPajama),
		models:            make(map[string]*llama.LLama),
		rwkv:              make(map[string]*rwkv.RwkvState),
		bloomz:            make(map[string]*bloomz.Bloomz),
		bert:              make(map[string]*bert.Bert),
		promptsTemplates:  make(map[string]*template.Template),
		whisperModels:     make(map[string]whisper.Model),
	}
}

func (ml *ModelLoader) ExistsInModelPath(s string) bool {
	_, err := os.Stat(filepath.Join(ml.ModelPath, s))
	return err == nil
}

func (ml *ModelLoader) ListModels() ([]string, error) {
	files, err := ioutil.ReadDir(ml.ModelPath)
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
		modelFile := filepath.Join(ml.ModelPath, modelName)
		if err := ml.loadTemplateIfExists(modelName, modelFile); err != nil {
			return "", err
		}

		t, exists := ml.promptsTemplates[modelName]
		if exists {
			m = t
		}
	}
	if m == nil {
		return "", fmt.Errorf("failed loading any template")
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
	// skip any error here - we run anyway if a template does not exist
	modelTemplateFile := fmt.Sprintf("%s.tmpl", modelName)

	if !ml.ExistsInModelPath(modelTemplateFile) {
		return nil
	}

	dat, err := os.ReadFile(filepath.Join(ml.ModelPath, modelTemplateFile))
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

func (ml *ModelLoader) LoadRedPajama(modelName string) (*gpt2.RedPajama, error) {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	// Check if we already have a loaded model
	if !ml.ExistsInModelPath(modelName) {
		return nil, fmt.Errorf("model does not exist")
	}

	if m, ok := ml.redpajama[modelName]; ok {
		log.Debug().Msgf("Model already loaded in memory: %s", modelName)
		return m, nil
	}

	// Load the model and keep it in memory for later use
	modelFile := filepath.Join(ml.ModelPath, modelName)
	log.Debug().Msgf("Loading model in memory from file: %s", modelFile)

	model, err := gpt2.NewRedPajama(modelFile)
	if err != nil {
		return nil, err
	}

	// If there is a prompt template, load it
	if err := ml.loadTemplateIfExists(modelName, modelFile); err != nil {
		return nil, err
	}

	ml.redpajama[modelName] = model
	return model, err
}

func (ml *ModelLoader) LoadDollyModel(modelName string) (*gpt2.Dolly, error) {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	// Check if we already have a loaded model
	if !ml.ExistsInModelPath(modelName) {
		return nil, fmt.Errorf("model does not exist")
	}

	if m, ok := ml.dollymodels[modelName]; ok {
		log.Debug().Msgf("Model already loaded in memory: %s", modelName)
		return m, nil
	}

	// Load the model and keep it in memory for later use
	modelFile := filepath.Join(ml.ModelPath, modelName)
	log.Debug().Msgf("Loading model in memory from file: %s", modelFile)

	model, err := gpt2.NewDolly(modelFile)
	if err != nil {
		return nil, err
	}

	// If there is a prompt template, load it
	if err := ml.loadTemplateIfExists(modelName, modelFile); err != nil {
		return nil, err
	}

	ml.dollymodels[modelName] = model
	return model, err
}

func (ml *ModelLoader) LoadStableLMModel(modelName string) (*gpt2.StableLM, error) {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	// Check if we already have a loaded model
	if !ml.ExistsInModelPath(modelName) {
		return nil, fmt.Errorf("model does not exist")
	}

	if m, ok := ml.gptstablelmmodels[modelName]; ok {
		log.Debug().Msgf("Model already loaded in memory: %s", modelName)
		return m, nil
	}

	// Load the model and keep it in memory for later use
	modelFile := filepath.Join(ml.ModelPath, modelName)
	log.Debug().Msgf("Loading model in memory from file: %s", modelFile)

	model, err := gpt2.NewStableLM(modelFile)
	if err != nil {
		return nil, err
	}

	// If there is a prompt template, load it
	if err := ml.loadTemplateIfExists(modelName, modelFile); err != nil {
		return nil, err
	}

	ml.gptstablelmmodels[modelName] = model
	return model, err
}

func (ml *ModelLoader) LoadBERT(modelName string) (*bert.Bert, error) {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	// Check if we already have a loaded model
	if !ml.ExistsInModelPath(modelName) {
		return nil, fmt.Errorf("model does not exist")
	}

	if m, ok := ml.bert[modelName]; ok {
		log.Debug().Msgf("Model already loaded in memory: %s", modelName)
		return m, nil
	}

	// Load the model and keep it in memory for later use
	modelFile := filepath.Join(ml.ModelPath, modelName)
	log.Debug().Msgf("Loading model in memory from file: %s", modelFile)

	model, err := bert.New(modelFile)
	if err != nil {
		return nil, err
	}

	// If there is a prompt template, load it
	if err := ml.loadTemplateIfExists(modelName, modelFile); err != nil {
		return nil, err
	}

	ml.bert[modelName] = model
	return model, err
}

func (ml *ModelLoader) LoadBloomz(modelName string) (*bloomz.Bloomz, error) {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	// Check if we already have a loaded model
	if !ml.ExistsInModelPath(modelName) {
		return nil, fmt.Errorf("model does not exist")
	}

	if m, ok := ml.bloomz[modelName]; ok {
		log.Debug().Msgf("Model already loaded in memory: %s", modelName)
		return m, nil
	}

	// Load the model and keep it in memory for later use
	modelFile := filepath.Join(ml.ModelPath, modelName)
	log.Debug().Msgf("Loading model in memory from file: %s", modelFile)

	model, err := bloomz.New(modelFile)
	if err != nil {
		return nil, err
	}

	// If there is a prompt template, load it
	if err := ml.loadTemplateIfExists(modelName, modelFile); err != nil {
		return nil, err
	}

	ml.bloomz[modelName] = model
	return model, err
}

func (ml *ModelLoader) LoadGPT2Model(modelName string) (*gpt2.GPT2, error) {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	// Check if we already have a loaded model
	if !ml.ExistsInModelPath(modelName) {
		return nil, fmt.Errorf("model does not exist")
	}

	if m, ok := ml.gpt2models[modelName]; ok {
		log.Debug().Msgf("Model already loaded in memory: %s", modelName)
		return m, nil
	}

	// Load the model and keep it in memory for later use
	modelFile := filepath.Join(ml.ModelPath, modelName)
	log.Debug().Msgf("Loading model in memory from file: %s", modelFile)

	model, err := gpt2.New(modelFile)
	if err != nil {
		return nil, err
	}

	// If there is a prompt template, load it
	if err := ml.loadTemplateIfExists(modelName, modelFile); err != nil {
		return nil, err
	}

	ml.gpt2models[modelName] = model
	return model, err
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
	modelFile := filepath.Join(ml.ModelPath, modelName)
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

func (ml *ModelLoader) LoadRWKV(modelName, tokenFile string, threads uint32) (*rwkv.RwkvState, error) {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	log.Debug().Msgf("Loading model name: %s", modelName)

	// Check if we already have a loaded model
	if !ml.ExistsInModelPath(modelName) {
		return nil, fmt.Errorf("model does not exist")
	}

	if m, ok := ml.rwkv[modelName]; ok {
		log.Debug().Msgf("Model already loaded in memory: %s", modelName)
		return m, nil
	}

	// Load the model and keep it in memory for later use
	modelFile := filepath.Join(ml.ModelPath, modelName)
	tokenPath := filepath.Join(ml.ModelPath, tokenFile)
	log.Debug().Msgf("Loading model in memory from file: %s", modelFile)

	model := rwkv.LoadFiles(modelFile, tokenPath, threads)
	if model == nil {
		return nil, fmt.Errorf("could not load model")
	}

	ml.rwkv[modelName] = model
	return model, nil
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

	// Load the model and keep it in memory for later use
	modelFile := filepath.Join(ml.ModelPath, modelName)
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

func (ml *ModelLoader) LoadWhisperModel(modelName string) (whisper.Model, error) {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	// Check if we already have a loaded model
	if !ml.ExistsInModelPath(modelName) {
		return nil, fmt.Errorf("model does not exist -- %s", modelName)
	}

	if m, ok := ml.whisperModels[modelName]; ok {
		log.Debug().Msgf("Model already loaded in memory: %s", modelName)
		return m, nil
	}

	// Load the model and keep it in memory for later use
	modelFile := filepath.Join(ml.ModelPath, modelName)
	log.Debug().Msgf("Loading model in memory from file: %s", modelFile)

	model, err := whisper.New(modelFile)
	if err != nil {
		return nil, err
	}

	ml.whisperModels[modelName] = model
	return model, err
}

const tokenizerSuffix = ".tokenizer.json"

var loadedModels map[string]interface{} = map[string]interface{}{}
var muModels sync.Mutex

func (ml *ModelLoader) BackendLoader(backendString string, modelFile string, llamaOpts []llama.ModelOption, threads uint32) (model interface{}, err error) {
	switch strings.ToLower(backendString) {
	case "llama":
		return ml.LoadLLaMAModel(modelFile, llamaOpts...)
	case "bloomz":
		return ml.LoadBloomz(modelFile)
	case "stablelm":
		return ml.LoadStableLMModel(modelFile)
	case "dolly":
		return ml.LoadDollyModel(modelFile)
	case "redpajama":
		return ml.LoadRedPajama(modelFile)
	case "gpt2":
		return ml.LoadGPT2Model(modelFile)
	case "gptj":
		return ml.LoadGPTJModel(modelFile)
	case "bert-embeddings":
		return ml.LoadBERT(modelFile)
	case "rwkv":
		return ml.LoadRWKV(modelFile, modelFile+tokenizerSuffix, threads)
	default:
		return nil, fmt.Errorf("backend unsupported: %s", backendString)
	}
}

func (ml *ModelLoader) WhisperLoader(backendString string, modelFile string) (model whisper.Model, err error) {
	//TODO expose more whisper options in next PR
	switch strings.ToLower(backendString) {
	case "whisper":
		return ml.LoadWhisperModel(modelFile)
	default:
		return nil, fmt.Errorf("whisper backend unsupported: %s", backendString)
	}
}

func (ml *ModelLoader) GreedyLoader(modelFile string, llamaOpts []llama.ModelOption, threads uint32) (model interface{}, err error) {
	updateModels := func(model interface{}) {
		muModels.Lock()
		defer muModels.Unlock()
		loadedModels[modelFile] = model
	}

	muModels.Lock()
	m, exists := loadedModels[modelFile]
	if exists {
		muModels.Unlock()
		return m, nil
	}
	muModels.Unlock()

	model, modelerr := ml.LoadLLaMAModel(modelFile, llamaOpts...)
	if modelerr == nil {
		updateModels(model)
		return model, nil
	} else {
		err = multierror.Append(err, modelerr)
	}

	model, modelerr = ml.LoadGPTJModel(modelFile)
	if modelerr == nil {
		updateModels(model)
		return model, nil
	} else {
		err = multierror.Append(err, modelerr)
	}

	model, modelerr = ml.LoadGPT2Model(modelFile)
	if modelerr == nil {
		updateModels(model)
		return model, nil
	} else {
		err = multierror.Append(err, modelerr)
	}

	model, modelerr = ml.LoadStableLMModel(modelFile)
	if modelerr == nil {
		updateModels(model)
		return model, nil
	} else {
		err = multierror.Append(err, modelerr)
	}

	model, modelerr = ml.LoadDollyModel(modelFile)
	if modelerr == nil {
		updateModels(model)
		return model, nil
	} else {
		err = multierror.Append(err, modelerr)
	}

	model, modelerr = ml.LoadRedPajama(modelFile)
	if modelerr == nil {
		updateModels(model)
		return model, nil
	} else {
		err = multierror.Append(err, modelerr)
	}

	model, modelerr = ml.LoadBloomz(modelFile)
	if modelerr == nil {
		updateModels(model)
		return model, nil
	} else {
		err = multierror.Append(err, modelerr)
	}

	model, modelerr = ml.LoadRWKV(modelFile, modelFile+tokenizerSuffix, threads)
	if modelerr == nil {
		updateModels(model)
		return model, nil
	} else {
		err = multierror.Append(err, modelerr)
	}

	model, modelerr = ml.LoadBERT(modelFile)
	if modelerr == nil {
		updateModels(model)
		return model, nil
	} else {
		err = multierror.Append(err, modelerr)
	}

	return nil, fmt.Errorf("could not load model - all backends returned error: %s", err.Error())
}
