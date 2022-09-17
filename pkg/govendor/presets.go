package govendor

var _ Preset = (*DefaultPreset)(nil)

// Preset interface used to customize the behaviour of the vendor library.
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

	// ForceFilters flag returns wether the preset will force the overriding of
	// the spec or dependency filters to the preset ones
	ForceFilters() bool
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
	return SPEC_FILENAME
}

func (dp *DefaultPreset) GetSpecLockFilename() string {
	return SPEC_LOCK_FILENAME
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
