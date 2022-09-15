package internal

import "github.com/alevinval/vendor-go/pkg/vendor"

type PresetWrapper struct {
	vendor.Preset
}

func WrapPreset(preset vendor.Preset) *PresetWrapper {
	return &PresetWrapper{preset}
}

func (w *PresetWrapper) GetSpecFilename() string {
	if w.Preset == nil {
		return vendor.SPEC_FILENAME
	} else {
		return w.Preset.GetSpecFilename()
	}
}

func (w *PresetWrapper) GetSpecLockFilename() string {
	if w.Preset == nil {
		return vendor.SPEC_FILENAME
	} else {
		return w.Preset.GetSpecLockFilename()
	}
}

func (w *PresetWrapper) LoadSpec() (*vendor.Spec, error) {
	return vendor.LoadSpec(w)
}

func (w *PresetWrapper) LoadSpecLock() (*vendor.SpecLock, error) {
	return vendor.LoadSpecLock(w)
}

func (w *PresetWrapper) NewSpec() *vendor.Spec {
	return vendor.NewSpec(w)
}
