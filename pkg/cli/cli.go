package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/alevinval/vendor-go/internal"
	"github.com/alevinval/vendor-go/pkg/core"
	"github.com/alevinval/vendor-go/pkg/core/installers"
	"github.com/alevinval/vendor-go/pkg/core/log"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	logger = log.GetLogger()

	isDebugEnabled bool
	preset         core.Preset = nil
)

func newRootCmd(commandName string) *cobra.Command {
	return &cobra.Command{
		Use:   commandName,
		Short: fmt.Sprintf("[%s] %s is a flexible and customizable vendoring tool", core.VERSION, commandName),
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
				logger.Warnf("%s already exists", core.SPEC_FILENAME)
				return
			}

			spec := core.NewSpec()
			spec.Preset = wrapper.Preset

			err = saveSpec(wrapper, spec)
			if err != nil {
				logger.Errorf("failed initializing: %s", err)
				return
			}

			logger.Infof("%s has been created", core.SPEC_FILENAME)
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

			spec, err := loadSpec(wrapper)
			if err != nil {
				return
			}

			dep := core.NewDependency(url, branch)
			spec.Add(dep)

			err = saveSpec(wrapper, spec)
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
			spec, err := loadSpec(wrapper)
			if err != nil {
				return
			}

			specLock, err := loadSpecLock(wrapper)
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

			err = saveSpecLock(wrapper, specLock)
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
			spec, err := loadSpec(wrapper)
			if err != nil {
				return
			}

			specLock, err := loadSpecLock(wrapper)
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

			err = saveSpecLock(wrapper, specLock)
			if err != nil {
				logger.Errorf("update failed: %s", err)
				return
			}

			logger.Infof("update success ✅")
		},
	}
}

func loadSpec(pw *internal.PresetWrapper) (*core.Spec, error) {
	data, err := os.ReadFile(pw.GetSpecFilename())
	if err != nil {
		logger.Warnf("cannot read %s: %s", pw.GetSpecFilename(), err)
		return nil, err
	}

	spec := core.NewSpec()
	err = yaml.Unmarshal(data, spec)
	if err != nil {
		logger.Errorf("cannot read %s: %s", core.SPEC_FILENAME, err)
		return nil, err
	}

	spec.Preset = pw.Preset
	return spec, nil
}

func loadSpecLock(pw *internal.PresetWrapper) (*core.SpecLock, error) {
	data, err := os.ReadFile(pw.GetSpecLockFilename())
	if err != nil {
		return core.NewSpecLock(), nil
	}

	spec := core.NewSpecLock()
	err = yaml.Unmarshal(data, spec)
	if err != nil {
		logger.Errorf("cannot read %s: %s", pw.GetSpecLockFilename(), err)
		return nil, err
	}

	return spec, nil
}

type YamlSerializable interface {
	ToYaml() []byte
}

func saveSpec(wrapper *internal.PresetWrapper, spec *core.Spec) error {
	return saveFile(wrapper.GetSpecFilename(), spec)
}

func saveSpecLock(wrapper *internal.PresetWrapper, specLock *core.SpecLock) error {
	return saveFile(wrapper.GetSpecLockFilename(), specLock)
}

func saveFile(filename string, s YamlSerializable) error {
	return os.WriteFile(filename, s.ToYaml(), fs.ModePerm)
}

func NewVendorCmd(commandName string, preset core.Preset) *cobra.Command {
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
