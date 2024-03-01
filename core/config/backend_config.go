package config

import (
<<<<<<< HEAD
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/mudler/LocalAI/core/schema"
	"github.com/mudler/LocalAI/pkg/downloader"
	"github.com/mudler/LocalAI/pkg/functions"
	"gopkg.in/yaml.v3"
)

const (
	RAND_SEED = -1
)

type TTSConfig struct {

	// Voice wav path or id
	Voice string `yaml:"voice"`

	AudioPath string `yaml:"audio_path"`
}

=======
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-skynet/LocalAI/core/schema"
	"github.com/go-skynet/LocalAI/pkg/downloader"
	"github.com/go-skynet/LocalAI/pkg/utils"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
type BackendConfig struct {
	schema.PredictionOptions `yaml:"parameters"`
	Name                     string `yaml:"name"`

<<<<<<< HEAD
	F16                 *bool                  `yaml:"f16"`
	Threads             *int                   `yaml:"threads"`
	Debug               *bool                  `yaml:"debug"`
	Roles               map[string]string      `yaml:"roles"`
	Embeddings          *bool                  `yaml:"embeddings"`
	Backend             string                 `yaml:"backend"`
	TemplateConfig      TemplateConfig         `yaml:"template"`
	KnownUsecaseStrings []string               `yaml:"known_usecases"`
	KnownUsecases       *BackendConfigUsecases `yaml:"-"`
	Pipeline            Pipeline               `yaml:"pipeline"`

	PromptStrings, InputStrings                []string               `yaml:"-"`
	InputToken                                 [][]int                `yaml:"-"`
	functionCallString, functionCallNameString string                 `yaml:"-"`
	ResponseFormat                             string                 `yaml:"-"`
	ResponseFormatMap                          map[string]interface{} `yaml:"-"`

	FunctionsConfig functions.FunctionsConfig `yaml:"function"`
=======
	F16            bool              `yaml:"f16"`
	Threads        int               `yaml:"threads"`
	Debug          bool              `yaml:"debug"`
	Roles          map[string]string `yaml:"roles"`
	Embeddings     bool              `yaml:"embeddings"`
	Backend        string            `yaml:"backend"`
	TemplateConfig TemplateConfig    `yaml:"template"`

	PromptStrings, InputStrings                []string `yaml:"-"`
	InputToken                                 [][]int  `yaml:"-"`
	functionCallString, functionCallNameString string   `yaml:"-"`

	FunctionsConfig Functions `yaml:"function"`
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))

	FeatureFlag FeatureFlag `yaml:"feature_flags"` // Feature Flag registry. We move fast, and features may break on a per model/backend basis. Registry for (usually temporary) flags that indicate aborting something early.
	// LLM configs (GPT4ALL, Llama.cpp, ...)
	LLMConfig `yaml:",inline"`

<<<<<<< HEAD
=======
	// AutoGPTQ specifics
	AutoGPTQ AutoGPTQ `yaml:"autogptq"`

>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
	// Diffusers
	Diffusers Diffusers `yaml:"diffusers"`
	Step      int       `yaml:"step"`

	// GRPC Options
	GRPC GRPC `yaml:"grpc"`

<<<<<<< HEAD
	// TTS specifics
	TTSConfig `yaml:"tts"`
=======
	// Vall-e-x
	VallE VallE `yaml:"vall-e"`
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))

	// CUDA
	// Explicitly enable CUDA or not (some backends might need it)
	CUDA bool `yaml:"cuda"`

	DownloadFiles []File `yaml:"download_files"`

	Description string `yaml:"description"`
	Usage       string `yaml:"usage"`
<<<<<<< HEAD

	Options   []string `yaml:"options"`
	Overrides []string `yaml:"overrides"`
}

// Pipeline defines other models to use for audio-to-audio
type Pipeline struct {
	TTS           string `yaml:"tts"`
	LLM           string `yaml:"llm"`
	Transcription string `yaml:"transcription"`
	VAD           string `yaml:"vad"`
}

type File struct {
	Filename string         `yaml:"filename" json:"filename"`
	SHA256   string         `yaml:"sha256" json:"sha256"`
	URI      downloader.URI `yaml:"uri" json:"uri"`
=======
}

type File struct {
	Filename string `yaml:"filename" json:"filename"`
	SHA256   string `yaml:"sha256" json:"sha256"`
	URI      string `yaml:"uri" json:"uri"`
}

type VallE struct {
	AudioPath string `yaml:"audio_path"`
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
}

type FeatureFlag map[string]*bool

func (ff FeatureFlag) Enabled(s string) bool {
	v, exist := ff[s]
	return exist && v != nil && *v
}

type GRPC struct {
	Attempts          int `yaml:"attempts"`
	AttemptsSleepTime int `yaml:"attempts_sleep_time"`
}

type Diffusers struct {
<<<<<<< HEAD
	CUDA             bool   `yaml:"cuda"`
	PipelineType     string `yaml:"pipeline_type"`
	SchedulerType    string `yaml:"scheduler_type"`
	EnableParameters string `yaml:"enable_parameters"` // A list of comma separated parameters to specify
	IMG2IMG          bool   `yaml:"img2img"`           // Image to Image Diffuser
	ClipSkip         int    `yaml:"clip_skip"`         // Skip every N frames
	ClipModel        string `yaml:"clip_model"`        // Clip model to use
	ClipSubFolder    string `yaml:"clip_subfolder"`    // Subfolder to use for clip model
	ControlNet       string `yaml:"control_net"`
}

// LLMConfig is a struct that holds the configuration that are
// generic for most of the LLM backends.
=======
	CUDA             bool    `yaml:"cuda"`
	PipelineType     string  `yaml:"pipeline_type"`
	SchedulerType    string  `yaml:"scheduler_type"`
	EnableParameters string  `yaml:"enable_parameters"` // A list of comma separated parameters to specify
	CFGScale         float32 `yaml:"cfg_scale"`         // Classifier-Free Guidance Scale
	IMG2IMG          bool    `yaml:"img2img"`           // Image to Image Diffuser
	ClipSkip         int     `yaml:"clip_skip"`         // Skip every N frames
	ClipModel        string  `yaml:"clip_model"`        // Clip model to use
	ClipSubFolder    string  `yaml:"clip_subfolder"`    // Subfolder to use for clip model
	ControlNet       string  `yaml:"control_net"`
}

>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
type LLMConfig struct {
	SystemPrompt    string   `yaml:"system_prompt"`
	TensorSplit     string   `yaml:"tensor_split"`
	MainGPU         string   `yaml:"main_gpu"`
	RMSNormEps      float32  `yaml:"rms_norm_eps"`
	NGQA            int32    `yaml:"ngqa"`
	PromptCachePath string   `yaml:"prompt_cache_path"`
	PromptCacheAll  bool     `yaml:"prompt_cache_all"`
	PromptCacheRO   bool     `yaml:"prompt_cache_ro"`
<<<<<<< HEAD
	MirostatETA     *float64 `yaml:"mirostat_eta"`
	MirostatTAU     *float64 `yaml:"mirostat_tau"`
	Mirostat        *int     `yaml:"mirostat"`
	NGPULayers      *int     `yaml:"gpu_layers"`
	MMap            *bool    `yaml:"mmap"`
	MMlock          *bool    `yaml:"mmlock"`
	LowVRAM         *bool    `yaml:"low_vram"`
	Reranking       *bool    `yaml:"reranking"`
	Grammar         string   `yaml:"grammar"`
	StopWords       []string `yaml:"stopwords"`
	Cutstrings      []string `yaml:"cutstrings"`
	ExtractRegex    []string `yaml:"extract_regex"`
	TrimSpace       []string `yaml:"trimspace"`
	TrimSuffix      []string `yaml:"trimsuffix"`

	ContextSize          *int             `yaml:"context_size"`
	NUMA                 bool             `yaml:"numa"`
	LoraAdapter          string           `yaml:"lora_adapter"`
	LoraBase             string           `yaml:"lora_base"`
	LoraAdapters         []string         `yaml:"lora_adapters"`
	LoraScales           []float32        `yaml:"lora_scales"`
	LoraScale            float32          `yaml:"lora_scale"`
	NoMulMatQ            bool             `yaml:"no_mulmatq"`
	DraftModel           string           `yaml:"draft_model"`
	NDraft               int32            `yaml:"n_draft"`
	Quantization         string           `yaml:"quantization"`
	LoadFormat           string           `yaml:"load_format"`
	GPUMemoryUtilization float32          `yaml:"gpu_memory_utilization"` // vLLM
	TrustRemoteCode      bool             `yaml:"trust_remote_code"`      // vLLM
	EnforceEager         bool             `yaml:"enforce_eager"`          // vLLM
	SwapSpace            int              `yaml:"swap_space"`             // vLLM
	MaxModelLen          int              `yaml:"max_model_len"`          // vLLM
	TensorParallelSize   int              `yaml:"tensor_parallel_size"`   // vLLM
	DisableLogStatus     bool             `yaml:"disable_log_stats"`      // vLLM
	DType                string           `yaml:"dtype"`                  // vLLM
	LimitMMPerPrompt     LimitMMPerPrompt `yaml:"limit_mm_per_prompt"`    // vLLM
	MMProj               string           `yaml:"mmproj"`

	FlashAttention bool   `yaml:"flash_attention"`
	NoKVOffloading bool   `yaml:"no_kv_offloading"`
	CacheTypeK     string `yaml:"cache_type_k"`
	CacheTypeV     string `yaml:"cache_type_v"`
=======
	MirostatETA     float64  `yaml:"mirostat_eta"`
	MirostatTAU     float64  `yaml:"mirostat_tau"`
	Mirostat        int      `yaml:"mirostat"`
	NGPULayers      int      `yaml:"gpu_layers"`
	MMap            bool     `yaml:"mmap"`
	MMlock          bool     `yaml:"mmlock"`
	LowVRAM         bool     `yaml:"low_vram"`
	Grammar         string   `yaml:"grammar"`
	StopWords       []string `yaml:"stopwords"`
	Cutstrings      []string `yaml:"cutstrings"`
	TrimSpace       []string `yaml:"trimspace"`
	TrimSuffix      []string `yaml:"trimsuffix"`

	ContextSize  int     `yaml:"context_size"`
	NUMA         bool    `yaml:"numa"`
	LoraAdapter  string  `yaml:"lora_adapter"`
	LoraBase     string  `yaml:"lora_base"`
	LoraScale    float32 `yaml:"lora_scale"`
	NoMulMatQ    bool    `yaml:"no_mulmatq"`
	DraftModel   string  `yaml:"draft_model"`
	NDraft       int32   `yaml:"n_draft"`
	Quantization string  `yaml:"quantization"`
	MMProj       string  `yaml:"mmproj"`
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))

	RopeScaling string `yaml:"rope_scaling"`
	ModelType   string `yaml:"type"`

	YarnExtFactor  float32 `yaml:"yarn_ext_factor"`
	YarnAttnFactor float32 `yaml:"yarn_attn_factor"`
	YarnBetaFast   float32 `yaml:"yarn_beta_fast"`
	YarnBetaSlow   float32 `yaml:"yarn_beta_slow"`
<<<<<<< HEAD

	CFGScale float32 `yaml:"cfg_scale"` // Classifier-Free Guidance Scale
}

// LimitMMPerPrompt is a struct that holds the configuration for the limit-mm-per-prompt config in vLLM
type LimitMMPerPrompt struct {
	LimitImagePerPrompt int `yaml:"image"`
	LimitVideoPerPrompt int `yaml:"video"`
	LimitAudioPerPrompt int `yaml:"audio"`
}

// TemplateConfig is a struct that holds the configuration of the templating system
type TemplateConfig struct {
	// Chat is the template used in the chat completion endpoint
	Chat string `yaml:"chat"`

	// ChatMessage is the template used for chat messages
	ChatMessage string `yaml:"chat_message"`

	// Completion is the template used for completion requests
	Completion string `yaml:"completion"`

	// Edit is the template used for edit completion requests
	Edit string `yaml:"edit"`

	// Functions is the template used when tools are present in the client requests
	Functions string `yaml:"function"`

	// UseTokenizerTemplate is a flag that indicates if the tokenizer template should be used.
	// Note: this is mostly consumed for backends such as vllm and transformers
	// that can use the tokenizers specified in the JSON config files of the models
	UseTokenizerTemplate bool `yaml:"use_tokenizer_template"`

	// JoinChatMessagesByCharacter is a string that will be used to join chat messages together.
	// It defaults to \n
	JoinChatMessagesByCharacter *string `yaml:"join_chat_messages_by_character"`

	Multimodal string `yaml:"multimodal"`

	JinjaTemplate bool `yaml:"jinja_template"`

	ReplyPrefix string `yaml:"reply_prefix"`
}

func (c *BackendConfig) UnmarshalYAML(value *yaml.Node) error {
	type BCAlias BackendConfig
	var aux BCAlias
	if err := value.Decode(&aux); err != nil {
		return err
	}
	*c = BackendConfig(aux)

	c.KnownUsecases = GetUsecasesFromYAML(c.KnownUsecaseStrings)
	// Make sure the usecases are valid, we rewrite with what we identified
	c.KnownUsecaseStrings = []string{}
	for k, usecase := range GetAllBackendConfigUsecases() {
		if c.HasUsecases(usecase) {
			c.KnownUsecaseStrings = append(c.KnownUsecaseStrings, k)
		}
	}
	return nil
=======
}

type AutoGPTQ struct {
	ModelBaseName    string `yaml:"model_base_name"`
	Device           string `yaml:"device"`
	Triton           bool   `yaml:"triton"`
	UseFastTokenizer bool   `yaml:"use_fast_tokenizer"`
}

type Functions struct {
	DisableNoAction         bool   `yaml:"disable_no_action"`
	NoActionFunctionName    string `yaml:"no_action_function_name"`
	NoActionDescriptionName string `yaml:"no_action_description_name"`
	ParallelCalls           bool   `yaml:"parallel_calls"`
}

type TemplateConfig struct {
	Chat        string `yaml:"chat"`
	ChatMessage string `yaml:"chat_message"`
	Completion  string `yaml:"completion"`
	Edit        string `yaml:"edit"`
	Functions   string `yaml:"function"`
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
}

func (c *BackendConfig) SetFunctionCallString(s string) {
	c.functionCallString = s
}

func (c *BackendConfig) SetFunctionCallNameString(s string) {
	c.functionCallNameString = s
}

func (c *BackendConfig) ShouldUseFunctions() bool {
	return ((c.functionCallString != "none" || c.functionCallString == "") || c.ShouldCallSpecificFunction())
}

func (c *BackendConfig) ShouldCallSpecificFunction() bool {
	return len(c.functionCallNameString) > 0
}

<<<<<<< HEAD
// MMProjFileName returns the filename of the MMProj file
// If the MMProj is a URL, it will return the MD5 of the URL which is the filename
func (c *BackendConfig) MMProjFileName() string {
	uri := downloader.URI(c.MMProj)
	if uri.LooksLikeURL() {
		f, _ := uri.FilenameFromUrl()
		return f
	}

	return c.MMProj
}

func (c *BackendConfig) IsMMProjURL() bool {
	uri := downloader.URI(c.MMProj)
	return uri.LooksLikeURL()
}

func (c *BackendConfig) IsModelURL() bool {
	uri := downloader.URI(c.Model)
	return uri.LooksLikeURL()
}

// ModelFileName returns the filename of the model
// If the model is a URL, it will return the MD5 of the URL which is the filename
func (c *BackendConfig) ModelFileName() string {
	uri := downloader.URI(c.Model)
	if uri.LooksLikeURL() {
		f, _ := uri.FilenameFromUrl()
		return f
	}

	return c.Model
}

func (c *BackendConfig) FunctionToCall() string {
	if c.functionCallNameString != "" &&
		c.functionCallNameString != "none" && c.functionCallNameString != "auto" {
		return c.functionCallNameString
	}

	return c.functionCallString
}

func (cfg *BackendConfig) SetDefaults(opts ...ConfigLoaderOption) {
	lo := &LoadOptions{}
	lo.Apply(opts...)

	ctx := lo.ctxSize
	threads := lo.threads
	f16 := lo.f16
	debug := lo.debug
	// https://github.com/ggerganov/llama.cpp/blob/75cd4c77292034ecec587ecb401366f57338f7c0/common/sampling.h#L22
	defaultTopP := 0.95
	defaultTopK := 40
	defaultTemp := 0.9
	// https://github.com/mudler/LocalAI/issues/2780
	defaultMirostat := 0
	defaultMirostatTAU := 5.0
	defaultMirostatETA := 0.1
	defaultTypicalP := 1.0
	defaultTFZ := 1.0
	defaultZero := 0

	trueV := true
	falseV := false

	if cfg.Seed == nil {
		//  random number generator seed
		defaultSeed := RAND_SEED
		cfg.Seed = &defaultSeed
	}

	if cfg.TopK == nil {
		cfg.TopK = &defaultTopK
	}

	if cfg.TypicalP == nil {
		cfg.TypicalP = &defaultTypicalP
	}

	if cfg.TFZ == nil {
		cfg.TFZ = &defaultTFZ
	}

	if cfg.MMap == nil {
		// MMap is enabled by default

		// Only exception is for Intel GPUs
		if os.Getenv("XPU") != "" {
			cfg.MMap = &falseV
		} else {
			cfg.MMap = &trueV
		}
	}

	if cfg.MMlock == nil {
		// MMlock is disabled by default
		cfg.MMlock = &falseV
	}

	if cfg.TopP == nil {
		cfg.TopP = &defaultTopP
	}
	if cfg.Temperature == nil {
		cfg.Temperature = &defaultTemp
	}

	if cfg.Maxtokens == nil {
		cfg.Maxtokens = &defaultZero
	}

	if cfg.Mirostat == nil {
		cfg.Mirostat = &defaultMirostat
	}

	if cfg.MirostatETA == nil {
		cfg.MirostatETA = &defaultMirostatETA
	}

	if cfg.MirostatTAU == nil {
		cfg.MirostatTAU = &defaultMirostatTAU
	}

	if cfg.LowVRAM == nil {
		cfg.LowVRAM = &falseV
	}

	if cfg.Embeddings == nil {
		cfg.Embeddings = &falseV
	}

	if cfg.Reranking == nil {
		cfg.Reranking = &falseV
	}

	if threads == 0 {
		// Threads can't be 0
		threads = 4
	}

	if cfg.Threads == nil {
		cfg.Threads = &threads
	}

	if cfg.F16 == nil {
		cfg.F16 = &f16
	}

	if cfg.Debug == nil {
		cfg.Debug = &falseV
	}

	if debug {
		cfg.Debug = &trueV
	}

	guessDefaultsFromFile(cfg, lo.modelPath, ctx)
}

func (c *BackendConfig) Validate() bool {
	downloadedFileNames := []string{}
	for _, f := range c.DownloadFiles {
		downloadedFileNames = append(downloadedFileNames, f.Filename)
	}
	validationTargets := []string{c.Backend, c.Model, c.MMProj}
	validationTargets = append(validationTargets, downloadedFileNames...)
	// Simple validation to make sure the model can be correctly loaded
	for _, n := range validationTargets {
		if n == "" {
			continue
		}
		if strings.HasPrefix(n, string(os.PathSeparator)) ||
			strings.Contains(n, "..") {
			return false
		}
	}

	if c.Backend != "" {
		// a regex that checks that is a string name with no special characters, except '-' and '_'
		re := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
		return re.MatchString(c.Backend)
	}

	return true
}

func (c *BackendConfig) HasTemplate() bool {
	return c.TemplateConfig.Completion != "" || c.TemplateConfig.Edit != "" || c.TemplateConfig.Chat != "" || c.TemplateConfig.ChatMessage != ""
}

type BackendConfigUsecases int

const (
	FLAG_ANY              BackendConfigUsecases = 0b000000000000
	FLAG_CHAT             BackendConfigUsecases = 0b000000000001
	FLAG_COMPLETION       BackendConfigUsecases = 0b000000000010
	FLAG_EDIT             BackendConfigUsecases = 0b000000000100
	FLAG_EMBEDDINGS       BackendConfigUsecases = 0b000000001000
	FLAG_RERANK           BackendConfigUsecases = 0b000000010000
	FLAG_IMAGE            BackendConfigUsecases = 0b000000100000
	FLAG_TRANSCRIPT       BackendConfigUsecases = 0b000001000000
	FLAG_TTS              BackendConfigUsecases = 0b000010000000
	FLAG_SOUND_GENERATION BackendConfigUsecases = 0b000100000000
	FLAG_TOKENIZE         BackendConfigUsecases = 0b001000000000
	FLAG_VAD              BackendConfigUsecases = 0b010000000000
	FLAG_VIDEO            BackendConfigUsecases = 0b100000000000

	// Common Subsets
	FLAG_LLM BackendConfigUsecases = FLAG_CHAT | FLAG_COMPLETION | FLAG_EDIT
)

func GetAllBackendConfigUsecases() map[string]BackendConfigUsecases {
	return map[string]BackendConfigUsecases{
		"FLAG_ANY":              FLAG_ANY,
		"FLAG_CHAT":             FLAG_CHAT,
		"FLAG_COMPLETION":       FLAG_COMPLETION,
		"FLAG_EDIT":             FLAG_EDIT,
		"FLAG_EMBEDDINGS":       FLAG_EMBEDDINGS,
		"FLAG_RERANK":           FLAG_RERANK,
		"FLAG_IMAGE":            FLAG_IMAGE,
		"FLAG_TRANSCRIPT":       FLAG_TRANSCRIPT,
		"FLAG_TTS":              FLAG_TTS,
		"FLAG_SOUND_GENERATION": FLAG_SOUND_GENERATION,
		"FLAG_TOKENIZE":         FLAG_TOKENIZE,
		"FLAG_VAD":              FLAG_VAD,
		"FLAG_LLM":              FLAG_LLM,
		"FLAG_VIDEO":            FLAG_VIDEO,
	}
}

func stringToFlag(s string) string {
	return "FLAG_" + strings.ToUpper(s)
}

func GetUsecasesFromYAML(input []string) *BackendConfigUsecases {
	if len(input) == 0 {
		return nil
	}
	result := FLAG_ANY
	flags := GetAllBackendConfigUsecases()
	for _, str := range input {
		flag, exists := flags[stringToFlag(str)]
		if exists {
			result |= flag
		}
	}
	return &result
}

// HasUsecases examines a BackendConfig and determines which endpoints have a chance of success.
func (c *BackendConfig) HasUsecases(u BackendConfigUsecases) bool {
	if (c.KnownUsecases != nil) && ((u & *c.KnownUsecases) == u) {
		return true
	}
	return c.GuessUsecases(u)
}

// GuessUsecases is a **heuristic based** function, as the backend in question may not be loaded yet, and the config may not record what it's useful at.
// In its current state, this function should ideally check for properties of the config like templates, rather than the direct backend name checks for the lower half.
// This avoids the maintenance burden of updating this list for each new backend - but unfortunately, that's the best option for some services currently.
func (c *BackendConfig) GuessUsecases(u BackendConfigUsecases) bool {
	if (u & FLAG_CHAT) == FLAG_CHAT {
		if c.TemplateConfig.Chat == "" && c.TemplateConfig.ChatMessage == "" {
			return false
		}
	}
	if (u & FLAG_COMPLETION) == FLAG_COMPLETION {
		if c.TemplateConfig.Completion == "" {
			return false
		}
	}
	if (u & FLAG_EDIT) == FLAG_EDIT {
		if c.TemplateConfig.Edit == "" {
			return false
		}
	}
	if (u & FLAG_EMBEDDINGS) == FLAG_EMBEDDINGS {
		if c.Embeddings == nil || !*c.Embeddings {
			return false
		}
	}
	if (u & FLAG_IMAGE) == FLAG_IMAGE {
		imageBackends := []string{"diffusers", "stablediffusion", "stablediffusion-ggml"}
		if !slices.Contains(imageBackends, c.Backend) {
			return false
		}

		if c.Backend == "diffusers" && c.Diffusers.PipelineType == "" {
			return false
		}

	}
	if (u & FLAG_VIDEO) == FLAG_VIDEO {
		videoBackends := []string{"diffusers", "stablediffusion"}
		if !slices.Contains(videoBackends, c.Backend) {
			return false
		}

		if c.Backend == "diffusers" && c.Diffusers.PipelineType == "" {
			return false
		}

	}
	if (u & FLAG_RERANK) == FLAG_RERANK {
		if c.Backend != "rerankers" {
			return false
		}
	}
	if (u & FLAG_TRANSCRIPT) == FLAG_TRANSCRIPT {
		if c.Backend != "whisper" {
			return false
		}
	}
	if (u & FLAG_TTS) == FLAG_TTS {
		ttsBackends := []string{"bark-cpp", "piper", "transformers-musicgen"}
		if !slices.Contains(ttsBackends, c.Backend) {
			return false
		}
	}

	if (u & FLAG_SOUND_GENERATION) == FLAG_SOUND_GENERATION {
		if c.Backend != "transformers-musicgen" {
			return false
		}
	}

	if (u & FLAG_TOKENIZE) == FLAG_TOKENIZE {
		tokenizeCapableBackends := []string{"llama.cpp", "rwkv"}
		if !slices.Contains(tokenizeCapableBackends, c.Backend) {
			return false
		}
	}

	if (u & FLAG_VAD) == FLAG_VAD {
		if c.Backend != "silero-vad" {
			return false
		}
	}

	return true
=======
func (c *BackendConfig) FunctionToCall() string {
	return c.functionCallNameString
}

func defaultPredictOptions(modelFile string) schema.PredictionOptions {
	return schema.PredictionOptions{
		TopP:        0.7,
		TopK:        80,
		Maxtokens:   512,
		Temperature: 0.9,
		Model:       modelFile,
	}
}

func DefaultConfig(modelFile string) *BackendConfig {
	return &BackendConfig{
		PredictionOptions: defaultPredictOptions(modelFile),
	}
}

////// Config Loader ////////

type BackendConfigLoader struct {
	configs map[string]BackendConfig
	sync.Mutex
}

// Load a config file for a model
func LoadBackendConfigFileByName(modelName, modelPath string, cl *BackendConfigLoader, debug bool, threads, ctx int, f16 bool) (*BackendConfig, error) {
	// Load a config file if present after the model name
	modelConfig := filepath.Join(modelPath, modelName+".yaml")

	var cfg *BackendConfig

	defaults := func() {
		cfg = DefaultConfig(modelName)
		cfg.ContextSize = ctx
		cfg.Threads = threads
		cfg.F16 = f16
		cfg.Debug = debug
	}

	cfgExisting, exists := cl.GetBackendConfig(modelName)
	if !exists {
		if _, err := os.Stat(modelConfig); err == nil {
			if err := cl.LoadBackendConfig(modelConfig); err != nil {
				return nil, fmt.Errorf("failed loading model config (%s) %s", modelConfig, err.Error())
			}
			cfgExisting, exists = cl.GetBackendConfig(modelName)
			if exists {
				cfg = &cfgExisting
			} else {
				defaults()
			}
		} else {
			defaults()
		}
	} else {
		cfg = &cfgExisting
	}

	// Set the parameters for the language model prediction
	//updateConfig(cfg, input)

	// Don't allow 0 as setting
	if cfg.Threads == 0 {
		if threads != 0 {
			cfg.Threads = threads
		} else {
			cfg.Threads = 4
		}
	}

	// Enforce debug flag if passed from CLI
	if debug {
		cfg.Debug = true
	}

	return cfg, nil
}

func NewBackendConfigLoader() *BackendConfigLoader {
	return &BackendConfigLoader{
		configs: make(map[string]BackendConfig),
	}
}
func ReadBackendConfigFile(file string) ([]*BackendConfig, error) {
	c := &[]*BackendConfig{}
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}
	if err := yaml.Unmarshal(f, c); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config file: %w", err)
	}

	return *c, nil
}

func ReadBackendConfig(file string) (*BackendConfig, error) {
	c := &BackendConfig{}
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}
	if err := yaml.Unmarshal(f, c); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config file: %w", err)
	}

	return c, nil
}

func (cm *BackendConfigLoader) LoadBackendConfigFile(file string) error {
	cm.Lock()
	defer cm.Unlock()
	c, err := ReadBackendConfigFile(file)
	if err != nil {
		return fmt.Errorf("cannot load config file: %w", err)
	}

	for _, cc := range c {
		cm.configs[cc.Name] = *cc
	}
	return nil
}

func (cl *BackendConfigLoader) LoadBackendConfig(file string) error {
	cl.Lock()
	defer cl.Unlock()
	c, err := ReadBackendConfig(file)
	if err != nil {
		return fmt.Errorf("cannot read config file: %w", err)
	}

	cl.configs[c.Name] = *c
	return nil
}

func (cl *BackendConfigLoader) GetBackendConfig(m string) (BackendConfig, bool) {
	cl.Lock()
	defer cl.Unlock()
	v, exists := cl.configs[m]
	return v, exists
}

func (cl *BackendConfigLoader) GetAllBackendConfigs() []BackendConfig {
	cl.Lock()
	defer cl.Unlock()
	var res []BackendConfig
	for _, v := range cl.configs {
		res = append(res, v)
	}
	return res
}

func (cl *BackendConfigLoader) ListBackendConfigs() []string {
	cl.Lock()
	defer cl.Unlock()
	var res []string
	for k := range cl.configs {
		res = append(res, k)
	}
	return res
}

// Preload prepare models if they are not local but url or huggingface repositories
func (cl *BackendConfigLoader) Preload(modelPath string) error {
	cl.Lock()
	defer cl.Unlock()

	status := func(fileName, current, total string, percent float64) {
		utils.DisplayDownloadFunction(fileName, current, total, percent)
	}

	log.Info().Msgf("Preloading models from %s", modelPath)

	for i, config := range cl.configs {

		// Download files and verify their SHA
		for _, file := range config.DownloadFiles {
			log.Debug().Msgf("Checking %q exists and matches SHA", file.Filename)

			if err := utils.VerifyPath(file.Filename, modelPath); err != nil {
				return err
			}
			// Create file path
			filePath := filepath.Join(modelPath, file.Filename)

			if err := downloader.DownloadFile(file.URI, filePath, file.SHA256, status); err != nil {
				return err
			}
		}

		modelURL := config.PredictionOptions.Model
		modelURL = downloader.ConvertURL(modelURL)

		if downloader.LooksLikeURL(modelURL) {
			// md5 of model name
			md5Name := utils.MD5(modelURL)

			// check if file exists
			if _, err := os.Stat(filepath.Join(modelPath, md5Name)); errors.Is(err, os.ErrNotExist) {
				err := downloader.DownloadFile(modelURL, filepath.Join(modelPath, md5Name), "", status)
				if err != nil {
					return err
				}
			}

			cc := cl.configs[i]
			c := &cc
			c.PredictionOptions.Model = md5Name
			cl.configs[i] = *c
		}
		if cl.configs[i].Name != "" {
			log.Info().Msgf("Model name: %s", cl.configs[i].Name)
		}
		if cl.configs[i].Description != "" {
			log.Info().Msgf("Model description: %s", cl.configs[i].Description)
		}
		if cl.configs[i].Usage != "" {
			log.Info().Msgf("Model usage: \n%s", cl.configs[i].Usage)
		}
	}
	return nil
}

func (cm *BackendConfigLoader) LoadBackendConfigsFromPath(path string) error {
	cm.Lock()
	defer cm.Unlock()
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	files := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}
		files = append(files, info)
	}
	for _, file := range files {
		// Skip templates, YAML and .keep files
		if !strings.Contains(file.Name(), ".yaml") && !strings.Contains(file.Name(), ".yml") {
			continue
		}
		c, err := ReadBackendConfig(filepath.Join(path, file.Name()))
		if err == nil {
			cm.configs[c.Name] = *c
		}
	}

	return nil
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
}
