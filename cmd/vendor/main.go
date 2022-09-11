package main

import (
	"io/fs"
	"os"
	"vendor-go/internal/core"
	"vendor-go/internal/core/installers"
	"vendor-go/internal/core/log"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var isDebugEnabled bool

var logger = log.GetLogger()

var RootCmd = &cobra.Command{
	Use:   "vendor",
	Short: "Vendor is a flexible and customizable vendoring tool",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if isDebugEnabled {
			log.EnableDebug()
		}
	},
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialises vendor in the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.ReadFile(core.SPEC_FILENAME)
		if err == nil {
			logger.Warnf("%s already exists", core.SPEC_FILENAME)
			return
		}

		spec := core.NewSpec()

		err = saveFile(core.SPEC_FILENAME, spec)
		if err != nil {
			logger.Errorf("failed initializing: %s", err)
			return
		}

		logger.Infof("%s has been created", core.SPEC_FILENAME)
	},
}

var AddCmd = &cobra.Command{
	Use:   "add [url] [branch]",
	Short: "Adds a new dependency to the vendor spec",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		branch := args[1]

		spec, err := loadSpec()
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

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs dependencies respectring the lock spec",
	Run: func(cmd *cobra.Command, args []string) {
		spec, err := loadSpec()
		if err != nil {
			return
		}

		specLock, err := loadSpecLock()
		if err != nil {
			return
		}

		cache := os.TempDir()
		logger.Debugf("repository cache located at %s", cache)

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

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update dependencies respectring the spec",
	Run: func(cmd *cobra.Command, args []string) {
		spec, err := loadSpec()
		if err != nil {
			return
		}

		specLock, err := loadSpecLock()
		if err != nil {
			return
		}

		cache := os.TempDir()
		logger.Debugf("repository cache located at %s", cache)

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

func loadSpec() (*core.Spec, error) {
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

func main() {
	RootCmd.PersistentFlags().BoolVarP(&isDebugEnabled, "debug", "d", false, "enable debug logging")
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(AddCmd)
	RootCmd.AddCommand(InstallCmd)
	RootCmd.AddCommand(UpdateCmd)
	RootCmd.Execute()
}
