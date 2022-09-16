package govendor

var _ Preset = (*DefaultPreset)(nil)

// Preset interface used to customize the behaviour of the vendor library.
// It allows customizing anything you need, like the names of the spec and
// lock file, also allows customizing the targeted and ignored paths.
type Preset interface {
	GetPresetName() string
	GetSpecFilename() string
	GetSpecLockFilename() string

	GetExtensions() []string
	GetTargets() []string
	GetIgnores() []string

	GetDepExtensions(dep *Dependency) []string
	GetDepTargets(dep *Dependency) []string
	GetDepIgnores(dep *Dependency) []string
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

func (dp *DefaultPreset) GetExtensions() []string {
	return []string{}
}

func (dp *DefaultPreset) GetTargets() []string {
	return []string{}
}

func (dp *DefaultPreset) GetIgnores() []string {
	return []string{}
}

func (dp *DefaultPreset) GetDepExtensions(*Dependency) []string {
	return []string{}
}

func (dp *DefaultPreset) GetDepTargets(*Dependency) []string {
	return []string{}
}

func (dp *DefaultPreset) GetDepIgnores(*Dependency) []string {
	return []string{}
}
