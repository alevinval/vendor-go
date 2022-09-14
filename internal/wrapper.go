package internal

import "github.com/alevinval/vendor-go/pkg/core"

type PresetWrapper struct {
	core.Preset
}

func WrapPreset(preset core.Preset) *PresetWrapper {
	return &PresetWrapper{preset}
}

func (w *PresetWrapper) GetSpecFilename() string {
	if w.Preset == nil {
		return core.SPEC_FILENAME
	} else {
		return w.Preset.GetSpecFilename()
	}
}

func (w *PresetWrapper) GetSpecLockFilename() string {
	if w.Preset == nil {
		return core.SPEC_FILENAME
	} else {
		return w.Preset.GetSpecLockFilename()
	}
}
