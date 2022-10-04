package cmd

import (
	"fmt"

	"github.com/alevinval/vendor-go/internal/control"
	"github.com/alevinval/vendor-go/pkg/log"
	"github.com/alevinval/vendor-go/pkg/vending"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
)

// NewCobraCommand returns a configured cobra command that serves as entry point
// to the vending CLI.
func NewCobraCommand(opts ...Option) *cobra.Command {
	b := &builder{
		debugFlag: new(bool),
	}
	for _, opt := range opts {
		opt(b)
	}
	return b.buildCobra()
}

type builder struct {
	preset      vending.Preset
	commandName string
	debugFlag   *bool
}

func (b *builder) buildCobra() *cobra.Command {
	rootCmd := newRootCmd(b.commandName, b.debugFlag)
	rootCmd.PersistentFlags().BoolVarP(b.debugFlag, "debug", "d", false, "enable debug logging")

	controller := control.New(
		control.WithPreset(b.preset),
	)

	rootCmd.AddCommand(newInitCmd(controller))
	rootCmd.AddCommand(newAddCmd(controller))
	rootCmd.AddCommand(newInstallCmd(controller))
	rootCmd.AddCommand(newUpdateCmd(controller))
	rootCmd.AddCommand(newCleanCacheCmd(controller))
	return rootCmd
}

func newRootCmd(commandName string, debugFlag *bool) *cobra.Command {
	return &cobra.Command{
		Use:   commandName,
		Short: fmt.Sprintf("%s is a flexible and customizable vending tool (%s)", commandName, vending.VERSION),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if *debugFlag {
				log.Level.SetLevel(zapcore.DebugLevel)
			}
		},
	}
}

func newInitCmd(controller *control.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "initializes the current directory",
		Run: func(cmd *cobra.Command, args []string) {
			err := controller.Init()
			if err != nil {
				log.S().Errorf("%s", err)
			}
		},
	}
}

func newAddCmd(controller *control.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "add [url] [branch]",
		Short: "Add a new dependency to the spec",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			branch := args[1]
			err := controller.AddDependency(url, branch)
			if err != nil {
				log.S().Errorf("%s", err)
			}
		},
	}
}

func newInstallCmd(controller *control.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Installs dependencies respecting the lockfile",
		Run: func(cmd *cobra.Command, args []string) {
			err := controller.Install()
			if err != nil {
				log.S().Errorf("%s", err)
			}
		},
	}
}

func newUpdateCmd(controller *control.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "update dependencies to the latest commit from the branch of the spec",
		Run: func(cmd *cobra.Command, args []string) {
			err := controller.Update()
			if err != nil {
				log.S().Errorf("%s", err)
			}
		},
	}
}

func newCleanCacheCmd(controller *control.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "cleancache",
		Short: "resets the repository cache",
		Run: func(cmd *cobra.Command, args []string) {
			err := controller.CleanCache()
			if err != nil {
				log.S().Errorf("%s", err)
			}
		},
	}
}
