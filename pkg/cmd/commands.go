package cmd

import (
	"fmt"

	"github.com/alevinval/vendor-go/internal"
	"github.com/alevinval/vendor-go/pkg/log"
	"github.com/alevinval/vendor-go/pkg/vending"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
)

var isDebugEnabled bool

// NewCobraCommand returns a configured cobra command that serves as entry point
// to the vending CLI.
func NewCobraCommand(opts ...Option) *cobra.Command {
	b := &cmdBuilder{}
	for _, opt := range opts {
		opt(b)
	}
	return b.build()
}

type cmdBuilder struct {
	preset      vending.Preset
	commandName string
}

func (cb *cmdBuilder) build() *cobra.Command {
	rootCmd := newRootCmd(cb.commandName)
	rootCmd.PersistentFlags().BoolVarP(&isDebugEnabled, "debug", "d", false, "enable debug logging")

	orchestrator := internal.NewOrchestrator(cb.preset)
	rootCmd.AddCommand(newInitCmd(orchestrator))
	rootCmd.AddCommand(newAddCmd(orchestrator))
	rootCmd.AddCommand(newInstallCmd(orchestrator))
	rootCmd.AddCommand(newUpdateCmd(orchestrator))
	rootCmd.AddCommand(newCleanCacheCmd(orchestrator))
	return rootCmd
}

func newRootCmd(commandName string) *cobra.Command {
	return &cobra.Command{
		Use:   commandName,
		Short: fmt.Sprintf("[%s] %s is a flexible and customizable vending tool", vending.VERSION, commandName),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if isDebugEnabled {
				log.Level.SetLevel(zapcore.DebugLevel)
			}
		},
	}
}

func newInitCmd(co *internal.CmdOrchestrator) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "initializes the current directory",
		Run: func(cmd *cobra.Command, args []string) {
			err := co.Init()
			if err != nil {
				log.S().Errorf("%s", err)
			}
		},
	}
}

func newAddCmd(co *internal.CmdOrchestrator) *cobra.Command {
	return &cobra.Command{
		Use:   "add [url] [branch]",
		Short: "Add a new dependency to the spec",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			branch := args[1]
			err := co.AddDependency(url, branch)
			if err != nil {
				log.S().Errorf("%s", err)
			}
		},
	}
}

func newInstallCmd(co *internal.CmdOrchestrator) *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Installs dependencies respecting the lockfile",
		Run: func(cmd *cobra.Command, args []string) {
			err := co.Install()
			if err != nil {
				log.S().Errorf("%s", err)
			}
		},
	}
}

func newUpdateCmd(co *internal.CmdOrchestrator) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "update dependencies to the latest commit from the branch of the spec",
		Run: func(cmd *cobra.Command, args []string) {
			err := co.Update()
			if err != nil {
				log.S().Errorf("%s", err)
			}
		},
	}
}

func newCleanCacheCmd(co *internal.CmdOrchestrator) *cobra.Command {
	return &cobra.Command{
		Use:   "cleancache",
		Short: "resets the repository cache",
		Run: func(cmd *cobra.Command, args []string) {
			err := co.CleanCache()
			if err != nil {
				log.S().Errorf("%s", err)
			}
		},
	}
}
