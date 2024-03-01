<<<<<<< HEAD
package config
=======
package config_test
>>>>>>> 68598ebe (MQTT Startup Refactoring Part 1: core/ packages part 1 (#1728))

import (
	"os"

<<<<<<< HEAD
=======
	. "github.com/go-skynet/LocalAI/core/config"
<<<<<<< HEAD
	"github.com/go-skynet/LocalAI/core/options"
	"github.com/go-skynet/LocalAI/pkg/model"
>>>>>>> 68598ebe (MQTT Startup Refactoring Part 1: core/ packages part 1 (#1728))
=======

>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test cases for config related functions", func() {

	var (
		configFile string
	)

	Context("Test Read configuration functions", func() {
		configFile = os.Getenv("CONFIG_FILE")
<<<<<<< HEAD
		It("Test readConfigFile", func() {
			config, err := readMultipleBackendConfigsFromFile(configFile)
=======
		It("Test ReadConfigFile", func() {
<<<<<<< HEAD
			config, err := ReadConfigFile(configFile)
>>>>>>> 68598ebe (MQTT Startup Refactoring Part 1: core/ packages part 1 (#1728))
=======
			config, err := ReadBackendConfigFile(configFile)
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
			Expect(err).To(BeNil())
			Expect(config).ToNot(BeNil())
			// two configs in config.yaml
			Expect(config[0].Name).To(Equal("list1"))
			Expect(config[1].Name).To(Equal("list2"))
		})

		It("Test LoadConfigs", func() {
<<<<<<< HEAD
<<<<<<< HEAD

			bcl := NewBackendConfigLoader(os.Getenv("MODELS_PATH"))
			err := bcl.LoadBackendConfigsFromPath(os.Getenv("MODELS_PATH"))

			Expect(err).To(BeNil())
			configs := bcl.GetAllBackendConfigs()
			loadedModelNames := []string{}
			for _, v := range configs {
				loadedModelNames = append(loadedModelNames, v.Name)
			}
			Expect(configs).ToNot(BeNil())

			Expect(loadedModelNames).To(ContainElements("code-search-ada-code-001"))

			// config should includes text-embedding-ada-002 models's api.config
			Expect(loadedModelNames).To(ContainElements("text-embedding-ada-002"))

			// config should includes rwkv_test models's api.config
			Expect(loadedModelNames).To(ContainElements("rwkv_test"))

			// config should includes whisper-1 models's api.config
			Expect(loadedModelNames).To(ContainElements("whisper-1"))
		})

		It("Test new loadconfig", func() {

			bcl := NewBackendConfigLoader(os.Getenv("MODELS_PATH"))
			err := bcl.LoadBackendConfigsFromPath(os.Getenv("MODELS_PATH"))
			Expect(err).To(BeNil())
			configs := bcl.GetAllBackendConfigs()
			loadedModelNames := []string{}
			for _, v := range configs {
				loadedModelNames = append(loadedModelNames, v.Name)
			}
			Expect(configs).ToNot(BeNil())
			totalModels := len(loadedModelNames)

			Expect(loadedModelNames).To(ContainElements("code-search-ada-code-001"))

			// config should includes text-embedding-ada-002 models's api.config
			Expect(loadedModelNames).To(ContainElements("text-embedding-ada-002"))

			// config should includes rwkv_test models's api.config
			Expect(loadedModelNames).To(ContainElements("rwkv_test"))

			// config should includes whisper-1 models's api.config
			Expect(loadedModelNames).To(ContainElements("whisper-1"))

			// create a temp directory and store a temporary model
			tmpdir, err := os.MkdirTemp("", "test")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tmpdir)

			// create a temporary model
			model := `name: "test-model"
description: "test model"
options:
- foo
- bar
- baz
`
			modelFile := tmpdir + "/test-model.yaml"
			err = os.WriteFile(modelFile, []byte(model), 0644)
			Expect(err).ToNot(HaveOccurred())

			err = bcl.LoadBackendConfigsFromPath(tmpdir)
			Expect(err).ToNot(HaveOccurred())

			configs = bcl.GetAllBackendConfigs()
			Expect(len(configs)).ToNot(Equal(totalModels))

			loadedModelNames = []string{}
			var testModel BackendConfig
			for _, v := range configs {
				loadedModelNames = append(loadedModelNames, v.Name)
				if v.Name == "test-model" {
					testModel = v
				}
			}
			Expect(loadedModelNames).To(ContainElements("test-model"))
			Expect(testModel.Description).To(Equal("test model"))
			Expect(testModel.Options).To(ContainElements("foo", "bar", "baz"))

=======
			cm := NewConfigLoader()
			opts := options.NewOptions()
			modelLoader := model.NewModelLoader(os.Getenv("MODELS_PATH"))
			options.WithModelLoader(modelLoader)(opts)

			err := cm.LoadConfigs(opts.Loader.ModelPath)
=======
			cm := NewBackendConfigLoader()
			opts := NewApplicationConfig()
			err := cm.LoadBackendConfigsFromPath(opts.ModelPath)
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
			Expect(err).To(BeNil())
			Expect(cm.ListBackendConfigs()).ToNot(BeNil())

			// config should includes gpt4all models's api.config
			Expect(cm.ListBackendConfigs()).To(ContainElements("gpt4all"))

			// config should includes gpt2 models's api.config
			Expect(cm.ListBackendConfigs()).To(ContainElements("gpt4all-2"))

			// config should includes text-embedding-ada-002 models's api.config
			Expect(cm.ListBackendConfigs()).To(ContainElements("text-embedding-ada-002"))

			// config should includes rwkv_test models's api.config
			Expect(cm.ListBackendConfigs()).To(ContainElements("rwkv_test"))

			// config should includes whisper-1 models's api.config
<<<<<<< HEAD
			Expect(cm.ListConfigs()).To(ContainElements("whisper-1"))
>>>>>>> 68598ebe (MQTT Startup Refactoring Part 1: core/ packages part 1 (#1728))
=======
			Expect(cm.ListBackendConfigs()).To(ContainElements("whisper-1"))
>>>>>>> 1ffb92d8 (refactor: move remaining api packages to core (#1731))
		})
	})
})
