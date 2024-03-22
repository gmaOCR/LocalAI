package e2e_test

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"testing"

<<<<<<< HEAD
	"github.com/docker/go-connections/nat"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sashabaranov/go-openai"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var container testcontainers.Container
=======
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sashabaranov/go-openai"
)

var pool *dockertest.Pool
var resource *dockertest.Resource
>>>>>>> f1b4a748 (feat(aio): add tests, update model definitions (#1880))
var client *openai.Client

var containerImage = os.Getenv("LOCALAI_IMAGE")
var containerImageTag = os.Getenv("LOCALAI_IMAGE_TAG")
var modelsDir = os.Getenv("LOCALAI_MODELS_DIR")
<<<<<<< HEAD
var apiEndpoint = os.Getenv("LOCALAI_API_ENDPOINT")
var apiKey = os.Getenv("LOCALAI_API_KEY")

const (
	defaultApiPort = "8080"
)
=======
var apiPort = os.Getenv("LOCALAI_API_PORT")
>>>>>>> f1b4a748 (feat(aio): add tests, update model definitions (#1880))

func TestLocalAI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LocalAI E2E test suite")
}

var _ = BeforeSuite(func() {

<<<<<<< HEAD
	var defaultConfig openai.ClientConfig
	if apiEndpoint == "" {
		startDockerImage()
		apiPort, err := container.MappedPort(context.Background(), nat.Port(defaultApiPort))
		Expect(err).To(Not(HaveOccurred()))

		defaultConfig = openai.DefaultConfig(apiKey)
		apiEndpoint = "http://localhost:" + apiPort.Port() + "/v1" // So that other tests can reference this value safely.
		defaultConfig.BaseURL = apiEndpoint
	} else {
		GinkgoWriter.Printf("docker apiEndpoint set from env: %q\n", apiEndpoint)
		defaultConfig = openai.DefaultConfig(apiKey)
		defaultConfig.BaseURL = apiEndpoint
	}

	// Wait for API to be ready
	client = openai.NewClientWithConfig(defaultConfig)

	Eventually(func() error {
		_, err := client.ListModels(context.TODO())
		return err
	}, "50m").ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	if container != nil {
		Expect(container.Terminate(context.Background())).To(Succeed())
	}
})

var _ = AfterEach(func() {
	// Add any cleanup needed after each test
})

type logConsumer struct {
}

func (l *logConsumer) Accept(log testcontainers.Log) {
	GinkgoWriter.Write([]byte(log.Content))
}

func startDockerImage() {
=======
	if containerImage == "" {
		Fail("LOCALAI_IMAGE is not set")
	}
	if containerImageTag == "" {
		Fail("LOCALAI_IMAGE_TAG is not set")
	}
	if apiPort == "" {
		apiPort = "8080"
	}

	p, err := dockertest.NewPool("")
	Expect(err).To(Not(HaveOccurred()))
	Expect(p.Client.Ping()).To(Succeed())

	pool = p

>>>>>>> f1b4a748 (feat(aio): add tests, update model definitions (#1880))
	// get cwd
	cwd, err := os.Getwd()
	Expect(err).To(Not(HaveOccurred()))
	md := cwd + "/models"

	if modelsDir != "" {
		md = modelsDir
	}

	proc := runtime.NumCPU()
<<<<<<< HEAD

	req := testcontainers.ContainerRequest{

		Image:        fmt.Sprintf("%s:%s", containerImage, containerImageTag),
		ExposedPorts: []string{defaultApiPort},
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Consumers: []testcontainers.LogConsumer{
				&logConsumer{},
			},
		},
		Env: map[string]string{
			"MODELS_PATH":                   "/models",
			"DEBUG":                         "true",
			"THREADS":                       fmt.Sprint(proc),
			"LOCALAI_SINGLE_ACTIVE_BACKEND": "true",
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      md,
				ContainerFilePath: "/models",
				FileMode:          0o755,
			},
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(defaultApiPort)),
		//	wait.ForHTTP("/v1/models").WithPort(nat.Port(apiPort)).WithStartupTimeout(50*time.Minute),
		),
	}

	GinkgoWriter.Printf("Launching Docker Container %s:%s\n", containerImage, containerImageTag)

	ctx := context.Background()
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	Expect(err).To(Not(HaveOccurred()))

	container = c
}
=======
	options := &dockertest.RunOptions{
		Repository: containerImage,
		Tag:        containerImageTag,
		//	Cmd:        []string{"server", "/data"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"8080/tcp": []docker.PortBinding{{HostPort: apiPort}},
		},
		Env:    []string{"MODELS_PATH=/models", "DEBUG=true", "THREADS=" + fmt.Sprint(proc)},
		Mounts: []string{md + ":/models"},
	}

	r, err := pool.RunWithOptions(options)
	Expect(err).To(Not(HaveOccurred()))

	resource = r

	defaultConfig := openai.DefaultConfig("")
	defaultConfig.BaseURL = "http://localhost:" + apiPort + "/v1"

	// Wait for API to be ready
	client = openai.NewClientWithConfig(defaultConfig)

	Eventually(func() error {
		_, err := client.ListModels(context.TODO())
		return err
	}, "20m").ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	Expect(pool.Purge(resource)).To(Succeed())
	//dat, err := os.ReadFile(resource.Container.LogPath)
	//Expect(err).To(Not(HaveOccurred()))
	//Expect(string(dat)).To(ContainSubstring("GRPC Service Ready"))
	//fmt.Println(string(dat))
})

var _ = AfterEach(func() {
	//Expect(dbClient.Clear()).To(Succeed())
})
>>>>>>> f1b4a748 (feat(aio): add tests, update model definitions (#1880))
