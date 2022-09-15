package internal

import "github.com/alevinval/vendor-go/pkg/govendor"

type PresetWrapper struct {
	govendor.Preset
}

func WrapPreset(preset govendor.Preset) *PresetWrapper {
	return &PresetWrapper{preset}
}

func (w *PresetWrapper) LoadSpec() (*govendor.Spec, error) {
	return govendor.LoadSpec(w)
}

func (w *PresetWrapper) LoadSpecLock() (*govendor.SpecLock, error) {
	return govendor.LoadSpecLock(w)
}

func (w *PresetWrapper) NewSpec() *govendor.Spec {
	return govendor.NewSpec(w)
}
