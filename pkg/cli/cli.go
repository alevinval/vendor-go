package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path"

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

func newInitCmd(preset core.Preset) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "initialises the current directory",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := os.ReadFile(core.SPEC_FILENAME)
			if err == nil {
				logger.Warnf("%s already exists", core.SPEC_FILENAME)
				return
			}

			spec := core.NewSpec()
			spec.Preset = preset

			err = saveFile(core.SPEC_FILENAME, spec)
			if err != nil {
				logger.Errorf("failed initializing: %s", err)
				return
			}

			logger.Infof("%s has been created", core.SPEC_FILENAME)
		},
	}
}

func newAddCmd(preset core.Preset) *cobra.Command {
	return &cobra.Command{
		Use:   "add [url] [branch]",
		Short: "Add a new dependency to the spec",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			branch := args[1]

			spec, err := loadSpec(preset)
			if err != nil {
				return
			}

			dep := core.NewDependency(url, branch)
			spec.Add(dep)

			err = saveFile(core.SPEC_FILENAME, spec)
			if err != nil {
				logger.Errorf("failed adding dependency: %s", err)
				return
			}

			logger.Infof("added dependency %s@%s", url, branch)
		},
	}
}

func newInstallCmd(preset core.Preset) *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Installs dependencies respectring the lockfile",
		Run: func(cmd *cobra.Command, args []string) {
			spec, err := loadSpec(preset)
			if err != nil {
				return
			}

			specLock, err := loadSpecLock()
			if err != nil {
				return
			}

			cache := path.Join(os.TempDir(), "vendor-go-cache")
			logger.Infof("repository cache located at %s", cache)

			m := installers.NewSpecInstaller(cache, spec, specLock)
			err = m.Install()
			if err != nil {
				logger.Errorf("install failed: %s", err)
				return
			}

			err = saveFile(core.SPEC_LOCK_FILENAME, specLock)
			if err != nil {
				logger.Errorf("install failed: %s", err)
				return
			}

			logger.Infof("install success ✅")
		},
	}
}

func newUpdateCmd(preset core.Preset) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "update dependencies to the latest commit from the branch of the spec",
		Run: func(cmd *cobra.Command, args []string) {
			spec, err := loadSpec(preset)
			if err != nil {
				return
			}

			specLock, err := loadSpecLock()
			if err != nil {
				return
			}

			cache := path.Join(os.TempDir(), "vendor-go-cache")
			logger.Infof("repository cache located at %s", cache)

			m := installers.NewSpecInstaller(cache, spec, specLock)
			err = m.Update()
			if err != nil {
				logger.Errorf("update failed: %s", err)
				return
			}

			err = saveFile(core.SPEC_LOCK_FILENAME, specLock)
			if err != nil {
				logger.Errorf("update failed: %s", err)
				return
			}

			logger.Infof("update success ✅")
		},
	}
}

func loadSpec(preset core.Preset) (*core.Spec, error) {
	data, err := os.ReadFile(core.SPEC_FILENAME)
	if err != nil {
		logger.Warnf("cannot read %s: %s", core.SPEC_FILENAME, err)
		return nil, err
	}

	spec := core.NewSpec()
	err = yaml.Unmarshal(data, spec)
	if err != nil {
		logger.Errorf("cannot read %s: %s", core.SPEC_FILENAME, err)
		return nil, err
	}

	spec.Preset = preset
	return spec, nil
}

func loadSpecLock() (*core.SpecLock, error) {
	data, err := os.ReadFile(core.SPEC_LOCK_FILENAME)
	if err != nil {
		return core.NewSpecLock(), nil
	}

	spec := core.NewSpecLock()
	err = yaml.Unmarshal(data, spec)
	if err != nil {
		logger.Errorf("cannot read %s: %s", core.SPEC_LOCK_FILENAME, err)
		return nil, err
	}

	return spec, nil
}

type YamlSerializable interface {
	ToYaml() []byte
}

func saveFile(filename string, s YamlSerializable) error {
	return os.WriteFile(filename, s.ToYaml(), fs.ModePerm)
}

func NewVendorCmd(commandName string, preset core.Preset) *cobra.Command {
	rootCmd := newRootCmd(commandName)
	rootCmd.PersistentFlags().BoolVarP(&isDebugEnabled, "debug", "d", false, "enable debug logging")

	rootCmd.AddCommand(newInitCmd(preset))
	rootCmd.AddCommand(newAddCmd(preset))
	rootCmd.AddCommand(newInstallCmd(preset))
	rootCmd.AddCommand(newUpdateCmd(preset))
	return rootCmd
}

func Run() error {
	return NewVendorCmd("vendor", nil).Execute()
}