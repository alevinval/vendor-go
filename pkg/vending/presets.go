package vending

import (
	"os"
	"path"

	"github.com/alevinval/vendor-go/pkg/log"
)

var _ Preset = (*DefaultPreset)(nil)

// Preset interface used to customize the behavior of the vendor library.
// It allows customizing anything you need, like the names of the spec and
// lock file, also allows customizing the targeted and ignored paths.
type Preset interface {
	// GetPresetName returns the name for this preset
	GetPresetName() string

	// GetVendorDir returns the name for the vendor folder
	GetVendorDir() string

	// GetSpecFilename returns the name of the spec file
	GetSpecFilename() string

	// GetSpecLockFilename returns the name of the spec lock file
	GetSpecLockFilename() string

	// GetFilters returns the global filters of the preset
	GetFilters() *Filters

	// GetFiltersForDependency returns the specific filters for a dependency
	GetFiltersForDependency(*Dependency) *Filters

	// ForceFilters flag returns whether the preset will force the overriding of
	// the spec or dependency filters to the preset ones
	ForceFilters() bool

	// GetCacheDir returns the path where the repository cache will be kept
	GetCacheDir() string
}

// DefaultPreset provides the default configuration for the vendor library.
type DefaultPreset struct{}

func (dp *DefaultPreset) GetPresetName() string {
	return "default"
}

func (dp *DefaultPreset) GetVendorDir() string {
	return "vendor/"
}

func (dp *DefaultPreset) GetSpecFilename() string {
	return ".vendor.yml"
}

func (dp *DefaultPreset) GetSpecLockFilename() string {
	return ".vendor-lock.yml"
}

func (dp *DefaultPreset) GetFilters() *Filters {
	return NewFilters()
}

func (dp *DefaultPreset) GetFiltersForDependency(*Dependency) *Filters {
	return NewFilters()
}

func (dp *DefaultPreset) ForceFilters() bool {
	return false
}

// GetCacheDir default implementation tries to return a path under the user
// home dir. When this operation fails, it resorts back to a temporary dir..
func (dp *DefaultPreset) GetCacheDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.S().Warnf("Cannot find user home directory, using tempdir as home")
		homeDir = os.TempDir()
	}
	return path.Join(homeDir, ".cache", "vending")
}
