package cli

<<<<<<< HEAD
import (
	cliContext "github.com/mudler/LocalAI/core/cli/context"
	"github.com/mudler/LocalAI/core/cli/worker"
)

var CLI struct {
	cliContext.Context `embed:""`

	Run             RunCMD             `cmd:"" help:"Run LocalAI, this the default command if no other command is specified. Run 'local-ai run --help' for more information" default:"withargs"`
	Federated       FederatedCLI       `cmd:"" help:"Run LocalAI in federated mode"`
	Models          ModelsCMD          `cmd:"" help:"Manage LocalAI models and definitions"`
	Backends        BackendsCMD        `cmd:"" help:"Manage LocalAI backends and definitions"`
	TTS             TTSCMD             `cmd:"" help:"Convert text to speech"`
	SoundGeneration SoundGenerationCMD `cmd:"" help:"Generates audio files from text or audio"`
	Transcript      TranscriptCMD      `cmd:"" help:"Convert audio to text"`
	Worker          worker.Worker      `cmd:"" help:"Run workers to distribute workload (llama.cpp-only)"`
	Util            UtilCMD            `cmd:"" help:"Utility commands"`
	Explorer        ExplorerCMD        `cmd:"" help:"Run p2p explorer"`
=======
import "embed"

type Context struct {
	Debug    bool    `env:"LOCALAI_DEBUG,DEBUG" default:"false" hidden:"" help:"DEPRECATED, use --log-level=debug instead. Enable debug logging"`
	LogLevel *string `env:"LOCALAI_LOG_LEVEL" enum:"error,warn,info,debug" help:"Set the level of logs to output [${enum}]"`

	// This field is not a command line argument/flag, the struct tag excludes it from the parsed CLI
	BackendAssets embed.FS `kong:"-"`
}

var CLI struct {
	Context `embed:""`

	Run        RunCMD        `cmd:"" help:"Run LocalAI, this the default command if no other command is specified. Run 'local-ai run --help' for more information" default:"withargs"`
	Models     ModelsCMD     `cmd:"" help:"Manage LocalAI models and definitions"`
	TTS        TTSCMD        `cmd:"" help:"Convert text to speech"`
	Transcript TranscriptCMD `cmd:"" help:"Convert audio to text"`
>>>>>>> e16d5918 (feat: kong cli refactor fixes #1955 (#1974))
}
