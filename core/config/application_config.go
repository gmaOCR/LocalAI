package config

import (
	"context"
<<<<<<< HEAD
	"encoding/json"
	"regexp"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/mudler/LocalAI/pkg/xsysinfo"
=======
	"embed"
	"encoding/json"
	"time"

	"github.com/go-skynet/LocalAI/pkg/gallery"
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
	"github.com/rs/zerolog/log"
)

type ApplicationConfig struct {
	Context                             context.Context
	ConfigFile                          string
	ModelPath                           string
<<<<<<< HEAD
	BackendsPath                        string
	ExternalBackends                    []string
	LibPath                             string
	UploadLimitMB, Threads, ContextSize int
	F16                                 bool
	Debug                               bool
	GeneratedContentDir                 string

	ConfigsDir string
	UploadDir  string

	DynamicConfigsDir             string
	DynamicConfigsDirPollInterval time.Duration
	CORS                          bool
	CSRF                          bool
	PreloadJSONModels             string
	PreloadModelsFromPath         string
	CORSAllowOrigins              string
	ApiKeys                       []string
	P2PToken                      string
	P2PNetworkID                  string

	DisableWebUI                       bool
	EnforcePredownloadScans            bool
	OpaqueErrors                       bool
	UseSubtleKeyComparison             bool
	DisableApiKeyRequirementForHttpGet bool
	DisableMetrics                     bool
	HttpGetExemptedEndpoints           []*regexp.Regexp
	DisableGalleryEndpoint             bool
	LoadToMemory                       []string

	Galleries        []Gallery
	BackendGalleries []Gallery

	BackendAssets     *rice.Box
=======
	UploadLimitMB, Threads, ContextSize int
	F16                                 bool
	Debug, DisableMessage               bool
	ImageDir                            string
	AudioDir                            string
	UploadDir                           string
	CORS                                bool
	PreloadJSONModels                   string
	PreloadModelsFromPath               string
	CORSAllowOrigins                    string
	ApiKeys                             []string

	ModelLibraryURL string

	Galleries []gallery.Gallery

	BackendAssets     embed.FS
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
	AssetsDestination string

	ExternalGRPCBackends map[string]string

<<<<<<< HEAD
	AutoloadGalleries, AutoloadBackendGalleries bool
=======
	AutoloadGalleries bool
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))

	SingleBackend           bool
	ParallelBackendRequests bool

	WatchDogIdle bool
	WatchDogBusy bool
	WatchDog     bool

	ModelsURL []string

	WatchDogBusyTimeout, WatchDogIdleTimeout time.Duration
<<<<<<< HEAD

	MachineTag string
=======
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
}

type AppOption func(*ApplicationConfig)

func NewApplicationConfig(o ...AppOption) *ApplicationConfig {
	opt := &ApplicationConfig{
<<<<<<< HEAD
		Context:       context.Background(),
		UploadLimitMB: 15,
		ContextSize:   512,
		Debug:         true,
=======
		Context:        context.Background(),
		UploadLimitMB:  15,
		Threads:        1,
		ContextSize:    512,
		Debug:          true,
		DisableMessage: true,
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
	}
	for _, oo := range o {
		oo(opt)
	}
	return opt
}

func WithModelsURL(urls ...string) AppOption {
	return func(o *ApplicationConfig) {
		o.ModelsURL = urls
	}
}

func WithModelPath(path string) AppOption {
	return func(o *ApplicationConfig) {
		o.ModelPath = path
	}
}

<<<<<<< HEAD
func WithBackendsPath(path string) AppOption {
	return func(o *ApplicationConfig) {
		o.BackendsPath = path
	}
}

func WithExternalBackends(backends ...string) AppOption {
	return func(o *ApplicationConfig) {
		o.ExternalBackends = backends
	}
}

func WithMachineTag(tag string) AppOption {
	return func(o *ApplicationConfig) {
		o.MachineTag = tag
	}
}

=======
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
func WithCors(b bool) AppOption {
	return func(o *ApplicationConfig) {
		o.CORS = b
	}
}

<<<<<<< HEAD
func WithP2PNetworkID(s string) AppOption {
	return func(o *ApplicationConfig) {
		o.P2PNetworkID = s
	}
}

func WithCsrf(b bool) AppOption {
	return func(o *ApplicationConfig) {
		o.CSRF = b
	}
}

func WithP2PToken(s string) AppOption {
	return func(o *ApplicationConfig) {
		o.P2PToken = s
	}
}

func WithLibPath(path string) AppOption {
	return func(o *ApplicationConfig) {
		o.LibPath = path
=======
func WithModelLibraryURL(url string) AppOption {
	return func(o *ApplicationConfig) {
		o.ModelLibraryURL = url
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
	}
}

var EnableWatchDog = func(o *ApplicationConfig) {
	o.WatchDog = true
}

var EnableWatchDogIdleCheck = func(o *ApplicationConfig) {
	o.WatchDog = true
	o.WatchDogIdle = true
}

<<<<<<< HEAD
var DisableGalleryEndpoint = func(o *ApplicationConfig) {
	o.DisableGalleryEndpoint = true
}

=======
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
var EnableWatchDogBusyCheck = func(o *ApplicationConfig) {
	o.WatchDog = true
	o.WatchDogBusy = true
}

<<<<<<< HEAD
var DisableWebUI = func(o *ApplicationConfig) {
	o.DisableWebUI = true
}

=======
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
func SetWatchDogBusyTimeout(t time.Duration) AppOption {
	return func(o *ApplicationConfig) {
		o.WatchDogBusyTimeout = t
	}
}

func SetWatchDogIdleTimeout(t time.Duration) AppOption {
	return func(o *ApplicationConfig) {
		o.WatchDogIdleTimeout = t
	}
}

var EnableSingleBackend = func(o *ApplicationConfig) {
	o.SingleBackend = true
}

var EnableParallelBackendRequests = func(o *ApplicationConfig) {
	o.ParallelBackendRequests = true
}

var EnableGalleriesAutoload = func(o *ApplicationConfig) {
	o.AutoloadGalleries = true
}

<<<<<<< HEAD
var EnableBackendGalleriesAutoload = func(o *ApplicationConfig) {
	o.AutoloadBackendGalleries = true
}

=======
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
func WithExternalBackend(name string, uri string) AppOption {
	return func(o *ApplicationConfig) {
		if o.ExternalGRPCBackends == nil {
			o.ExternalGRPCBackends = make(map[string]string)
		}
		o.ExternalGRPCBackends[name] = uri
	}
}

func WithCorsAllowOrigins(b string) AppOption {
	return func(o *ApplicationConfig) {
		o.CORSAllowOrigins = b
	}
}

func WithBackendAssetsOutput(out string) AppOption {
	return func(o *ApplicationConfig) {
		o.AssetsDestination = out
	}
}

<<<<<<< HEAD
func WithBackendAssets(f *rice.Box) AppOption {
=======
func WithBackendAssets(f embed.FS) AppOption {
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
	return func(o *ApplicationConfig) {
		o.BackendAssets = f
	}
}

func WithStringGalleries(galls string) AppOption {
	return func(o *ApplicationConfig) {
		if galls == "" {
<<<<<<< HEAD
			o.Galleries = []Gallery{}
			return
		}
		var galleries []Gallery
		if err := json.Unmarshal([]byte(galls), &galleries); err != nil {
			log.Error().Err(err).Msg("failed loading galleries")
=======
			log.Debug().Msgf("no galleries to load")
			o.Galleries = []gallery.Gallery{}
			return
		}
		var galleries []gallery.Gallery
		if err := json.Unmarshal([]byte(galls), &galleries); err != nil {
			log.Error().Msgf("failed loading galleries: %s", err.Error())
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
		}
		o.Galleries = append(o.Galleries, galleries...)
	}
}

<<<<<<< HEAD
func WithBackendGalleries(galls string) AppOption {
	return func(o *ApplicationConfig) {
		if galls == "" {
			o.BackendGalleries = []Gallery{}
			return
		}
		var galleries []Gallery
		if err := json.Unmarshal([]byte(galls), &galleries); err != nil {
			log.Error().Err(err).Msg("failed loading galleries")
		}
		o.BackendGalleries = append(o.BackendGalleries, galleries...)
	}
}

func WithGalleries(galleries []Gallery) AppOption {
=======
func WithGalleries(galleries []gallery.Gallery) AppOption {
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
	return func(o *ApplicationConfig) {
		o.Galleries = append(o.Galleries, galleries...)
	}
}

func WithContext(ctx context.Context) AppOption {
	return func(o *ApplicationConfig) {
		o.Context = ctx
	}
}

func WithYAMLConfigPreload(configFile string) AppOption {
	return func(o *ApplicationConfig) {
		o.PreloadModelsFromPath = configFile
	}
}

func WithJSONStringPreload(configFile string) AppOption {
	return func(o *ApplicationConfig) {
		o.PreloadJSONModels = configFile
	}
}
func WithConfigFile(configFile string) AppOption {
	return func(o *ApplicationConfig) {
		o.ConfigFile = configFile
	}
}

func WithUploadLimitMB(limit int) AppOption {
	return func(o *ApplicationConfig) {
		o.UploadLimitMB = limit
	}
}

func WithThreads(threads int) AppOption {
	return func(o *ApplicationConfig) {
<<<<<<< HEAD
		if threads == 0 { // 0 is not allowed
			threads = xsysinfo.CPUPhysicalCores()
		}
=======
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
		o.Threads = threads
	}
}

func WithContextSize(ctxSize int) AppOption {
	return func(o *ApplicationConfig) {
		o.ContextSize = ctxSize
	}
}

func WithF16(f16 bool) AppOption {
	return func(o *ApplicationConfig) {
		o.F16 = f16
	}
}

func WithDebug(debug bool) AppOption {
	return func(o *ApplicationConfig) {
		o.Debug = debug
	}
}

<<<<<<< HEAD
func WithGeneratedContentDir(generatedContentDir string) AppOption {
	return func(o *ApplicationConfig) {
		o.GeneratedContentDir = generatedContentDir
=======
func WithDisableMessage(disableMessage bool) AppOption {
	return func(o *ApplicationConfig) {
		o.DisableMessage = disableMessage
	}
}

func WithAudioDir(audioDir string) AppOption {
	return func(o *ApplicationConfig) {
		o.AudioDir = audioDir
	}
}

func WithImageDir(imageDir string) AppOption {
	return func(o *ApplicationConfig) {
		o.ImageDir = imageDir
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
	}
}

func WithUploadDir(uploadDir string) AppOption {
	return func(o *ApplicationConfig) {
		o.UploadDir = uploadDir
	}
}

<<<<<<< HEAD
func WithConfigsDir(configsDir string) AppOption {
	return func(o *ApplicationConfig) {
		o.ConfigsDir = configsDir
	}
}

func WithDynamicConfigDir(dynamicConfigsDir string) AppOption {
	return func(o *ApplicationConfig) {
		o.DynamicConfigsDir = dynamicConfigsDir
	}
}

func WithDynamicConfigDirPollInterval(interval time.Duration) AppOption {
	return func(o *ApplicationConfig) {
		o.DynamicConfigsDirPollInterval = interval
	}
}

=======
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
func WithApiKeys(apiKeys []string) AppOption {
	return func(o *ApplicationConfig) {
		o.ApiKeys = apiKeys
	}
}

<<<<<<< HEAD
func WithEnforcedPredownloadScans(enforced bool) AppOption {
	return func(o *ApplicationConfig) {
		o.EnforcePredownloadScans = enforced
	}
}

func WithOpaqueErrors(opaque bool) AppOption {
	return func(o *ApplicationConfig) {
		o.OpaqueErrors = opaque
	}
}

func WithLoadToMemory(models []string) AppOption {
	return func(o *ApplicationConfig) {
		o.LoadToMemory = models
	}
}

func WithSubtleKeyComparison(subtle bool) AppOption {
	return func(o *ApplicationConfig) {
		o.UseSubtleKeyComparison = subtle
	}
}

func WithDisableApiKeyRequirementForHttpGet(required bool) AppOption {
	return func(o *ApplicationConfig) {
		o.DisableApiKeyRequirementForHttpGet = required
	}
}

var DisableMetricsEndpoint AppOption = func(o *ApplicationConfig) {
	o.DisableMetrics = true
}

func WithHttpGetExemptedEndpoints(endpoints []string) AppOption {
	return func(o *ApplicationConfig) {
		o.HttpGetExemptedEndpoints = []*regexp.Regexp{}
		for _, epr := range endpoints {
			r, err := regexp.Compile(epr)
			if err == nil && r != nil {
				o.HttpGetExemptedEndpoints = append(o.HttpGetExemptedEndpoints, r)
			} else {
				log.Warn().Err(err).Str("regex", epr).Msg("Error while compiling HTTP Get Exemption regex, skipping this entry.")
			}
		}
	}
}

// ToConfigLoaderOptions returns a slice of ConfigLoader Option.
// Some options defined at the application level are going to be passed as defaults for
// all the configuration for the models.
// This includes for instance the context size or the number of threads.
// If a model doesn't set configs directly to the config model file
// it will use the defaults defined here.
func (o *ApplicationConfig) ToConfigLoaderOptions() []ConfigLoaderOption {
	return []ConfigLoaderOption{
		LoadOptionContextSize(o.ContextSize),
		LoadOptionDebug(o.Debug),
		LoadOptionF16(o.F16),
		LoadOptionThreads(o.Threads),
		ModelPath(o.ModelPath),
	}
}

=======
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
// func WithMetrics(meter *metrics.Metrics) AppOption {
// 	return func(o *StartupOptions) {
// 		o.Metrics = meter
// 	}
// }
