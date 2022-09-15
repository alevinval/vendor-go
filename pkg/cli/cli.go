package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/alevinval/vendor-go/internal"
	"github.com/alevinval/vendor-go/internal/installers"
	"github.com/alevinval/vendor-go/internal/log"
	"github.com/alevinval/vendor-go/pkg/vendor"

	"github.com/spf13/cobra"
)

var (
	logger = log.GetLogger()

	isDebugEnabled bool
)

func newRootCmd(commandName string) *cobra.Command {
	return &cobra.Command{
		Use:   commandName,
		Short: fmt.Sprintf("[%s] %s is a flexible and customizable vendoring tool", vendor.VERSION, commandName),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if isDebugEnabled {
				log.EnableDebug()
			}
		},
	}
}

func newInitCmd(wrapper *internal.PresetWrapper) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "initialises the current directory",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := os.ReadFile(wrapper.GetSpecFilename())
			if err == nil {
				logger.Warnf("%s already exists", vendor.SPEC_FILENAME)
				return
			}

			spec := wrapper.NewSpec()

			err = spec.Save()
			if err != nil {
				logger.Errorf("failed initializing: %s", err)
				return
			}

			logger.Infof("%s has been created", vendor.SPEC_FILENAME)
		},
	}
}

func newAddCmd(wrapper *internal.PresetWrapper) *cobra.Command {
	return &cobra.Command{
		Use:   "add [url] [branch]",
		Short: "Add a new dependency to the spec",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			branch := args[1]

			spec, err := wrapper.LoadSpec()
			if err != nil {
				return
			}

			dep := vendor.NewDependency(url, branch)
			spec.Add(dep)

			spec.Save()
			if err != nil {
				logger.Errorf("failed adding dependency: %s", err)
				return
			}

			logger.Infof("added dependency %s@%s", url, branch)
		},
	}
}

func newInstallCmd(wrapper *internal.PresetWrapper) *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Installs dependencies respectring the lockfile",
		Run: func(cmd *cobra.Command, args []string) {
			spec, err := wrapper.LoadSpec()
			if err != nil {
				return
			}

			specLock, err := wrapper.LoadSpecLock()
			if err != nil {
				return
			}

			cache := path.Join(os.TempDir(), "vendor-go-cache")
			logger.Infof("repository cache located at %s", cache)

			m := installers.NewInstaller(cache, spec, specLock)
			err = m.Install()
			if err != nil {
				logger.Errorf("install failed: %s", err)
				return
			}

			err = spec.Save()
			if err != nil {
				logger.Errorf("install failed: %s", err)
				return
			}

			err = specLock.Save()
			if err != nil {
				logger.Errorf("install failed: %s", err)
				return
			}

			logger.Infof("install success ✅")
		},
	}
}

func newUpdateCmd(wrapper *internal.PresetWrapper) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "update dependencies to the latest commit from the branch of the spec",
		Run: func(cmd *cobra.Command, args []string) {
			spec, err := wrapper.LoadSpec()
			if err != nil {
				return
			}

			specLock, err := wrapper.LoadSpecLock()
			if err != nil {
				return
			}

			cache := path.Join(os.TempDir(), "vendor-go-cache")
			logger.Infof("repository cache located at %s", cache)

			m := installers.NewInstaller(cache, spec, specLock)
			err = m.Update()
			if err != nil {
				logger.Errorf("update failed: %s", err)
				return
			}

			err = spec.Save()
			if err != nil {
				logger.Errorf("update failed: %s", err)
				return
			}

			err = specLock.Save()
			if err != nil {
				logger.Errorf("update failed: %s", err)
				return
			}

			logger.Infof("update success ✅")
		},
	}
}

func NewVendorCmd(commandName string, preset vendor.Preset) *cobra.Command {
	rootCmd := newRootCmd(commandName)
	rootCmd.PersistentFlags().BoolVarP(&isDebugEnabled, "debug", "d", false, "enable debug logging")

	wrapper := internal.WrapPreset(preset)
	rootCmd.AddCommand(newInitCmd(wrapper))
	rootCmd.AddCommand(newAddCmd(wrapper))
	rootCmd.AddCommand(newInstallCmd(wrapper))
	rootCmd.AddCommand(newUpdateCmd(wrapper))
	return rootCmd
}

func Run() error {
	return NewVendorCmd("vendor", nil).Execute()
}
