package govendor

var _ Preset = (*DefaultPreset)(nil)

// Preset interface used to customize the behaviour of the vendor library.
// It allows customizing anything you need, like the names of the spec and
// lock file, also allows customizing the targeted and ignored paths.
type Preset interface {
	GetSpecFilename() string
	GetSpecLockFilename() string
	GetExtensions() []string
	GetTargets(dep *Dependency) []string
	GetIgnores(dep *Dependency) []string
}

// DefaultPreset provides the default configuration for the vendor library.
type DefaultPreset struct{}

func (dp *DefaultPreset) GetSpecFilename() string {
	return SPEC_FILENAME
}

func (dp *DefaultPreset) GetSpecLockFilename() string {
	return SPEC_LOCK_FILENAME
}

func (dp *DefaultPreset) GetExtensions() []string {
	return []string{}
}

func (dp *DefaultPreset) GetTargets(dep *Dependency) []string {
	return []string{}
}

func (dp *DefaultPreset) GetIgnores(dep *Dependency) []string {
	return []string{}
}
