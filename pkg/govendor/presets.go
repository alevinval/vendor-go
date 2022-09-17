package govendor

var _ Preset = (*DefaultPreset)(nil)

// Preset interface used to customize the behaviour of the vendor library.
// It allows customizing anything you need, like the names of the spec and
// lock file, also allows customizing the targeted and ignored paths.
type Preset interface {
	GetPresetName() string
	GetSpecFilename() string
	GetSpecLockFilename() string
	GetFilters() *Filters
	GetFiltersForDependency(*Dependency) *Filters
}

// DefaultPreset provides the default configuration for the vendor library.
type DefaultPreset struct{}

func (dp *DefaultPreset) GetPresetName() string {
	return "default"
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
