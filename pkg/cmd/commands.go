package cmd

import (
	"fmt"

	"github.com/alevinval/vendor-go/internal/log"
	"github.com/alevinval/vendor-go/pkg/govendor"

	"github.com/spf13/cobra"
)

var (
	logger = log.GetLogger()

	isDebugEnabled bool
)

func newRootCmd(commandName string) *cobra.Command {
	return &cobra.Command{
		Use:   commandName,
		Short: fmt.Sprintf("[%s] %s is a flexible and customizable vendoring tool", govendor.VERSION, commandName),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if isDebugEnabled {
				log.EnableDebug()
			}
		},
	}
}

func newInitCmd(co *CmdOrchestrator) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "initialises the current directory",
		Run: func(cmd *cobra.Command, args []string) {
			err := co.Init()
			if err != nil {
				logger.Errorf("%s", err)
			}
		},
	}
}

func newAddCmd(co *CmdOrchestrator) *cobra.Command {
	return &cobra.Command{
		Use:   "add [url] [branch]",
		Short: "Add a new dependency to the spec",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			branch := args[1]
			err := co.AddDependency(url, branch)
			if err != nil {
				logger.Errorf("%s", err)
			}
		},
	}
}

func newInstallCmd(co *CmdOrchestrator) *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Installs dependencies respectring the lockfile",
		Run: func(cmd *cobra.Command, args []string) {
			err := co.Install()
			if err != nil {
				logger.Errorf("%s", err)
			}
		},
	}
}

func newUpdateCmd(co *CmdOrchestrator) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "update dependencies to the latest commit from the branch of the spec",
		Run: func(cmd *cobra.Command, args []string) {
			err := co.Update()
			if err != nil {
				logger.Errorf("%s", err)
			}
		},
	}
}

func NewVendorCmd(commandName string, preset govendor.Preset) *cobra.Command {
	rootCmd := newRootCmd(commandName)
	rootCmd.PersistentFlags().BoolVarP(&isDebugEnabled, "debug", "d", false, "enable debug logging")

	orchestrator := NewOrchestrator(preset)
	rootCmd.AddCommand(newInitCmd(orchestrator))
	rootCmd.AddCommand(newAddCmd(orchestrator))
	rootCmd.AddCommand(newInstallCmd(orchestrator))
	rootCmd.AddCommand(newUpdateCmd(orchestrator))
	return rootCmd
}
