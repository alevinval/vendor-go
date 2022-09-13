package core

type Preset interface {
	GetSpecFilename() string
	GetSpecLockFilename() string
	GetExtensions() []string
	GetTargets(dep *Dependency) []string
	GetIgnores(dep *Dependency) []string
}

type PresetWrapper struct {
	Preset
}

func WrapPreset(preset Preset) *PresetWrapper {
	return &PresetWrapper{preset}
}

func (w *PresetWrapper) GetSpecFilename() string {
	if w.Preset == nil {
		return SPEC_FILENAME
	} else {
		return w.Preset.GetSpecFilename()
	}
}

func (w *PresetWrapper) GetSpecLockFilename() string {
	if w.Preset == nil {
		return SPEC_FILENAME
	} else {
		return w.Preset.GetSpecLockFilename()
	}
}
