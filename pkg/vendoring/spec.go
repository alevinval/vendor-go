package vendoring

import (
	"fmt"
	"os"
	"strings"

	"github.com/alevinval/vendor-go/pkg/log"
	"gopkg.in/yaml.v3"
)

const VERSION = "0.3.3"

type Spec struct {
	Version    string        `yaml:"version"`
	PresetName string        `yaml:"preset"`
	VendorDir  string        `yaml:"vendor_dir,omitempty"`
	Filters    *Filters      `yaml:",inline"`
	Deps       []*Dependency `yaml:"deps"`
	preset     Preset        `yaml:"-"`
}

func LoadSpec(preset Preset) (*Spec, error) {
	preset = checkPreset(preset, true)

	filename := preset.GetSpecFilename()
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot read: %w", err)
	}

	spec := &Spec{}
	err = yaml.Unmarshal(data, spec)
	if err != nil {
		return nil, fmt.Errorf("cannot parse: %w", err)
	}

	spec.applyPreset(preset)
	return spec, nil
}

func NewSpec(preset Preset) *Spec {
	preset = checkPreset(preset, true)
	spec := &Spec{
		Version: VERSION,
		Filters: NewFilters(),
		Deps:    []*Dependency{},
	}
	spec.applyPreset(preset)
	return spec
}

func (s *Spec) AddDependency(dependency *Dependency) {
	if dep, ok := s.findDep(dependency); ok {
		dep.Update(dependency)
	} else {
		s.Deps = append(s.Deps, dependency)
	}
	s.applyPreset(s.preset)
}

func (s *Spec) Save() error {
	filename := s.preset.GetSpecFilename()
	data, err := toYaml(s)
	if err != nil {
		return fmt.Errorf("cannot save: %w", err)
	}
	err = os.WriteFile(filename, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot save: %w", err)
	}
	return nil
}

func (s *Spec) applyPreset(preset Preset) {
	s.preset = preset
	s.VendorDir = preset.GetVendorDir()
	s.PresetName = preset.GetPresetName()

	if preset.ForceFilters() {
		s.Filters = preset.GetFilters()
	} else {
		if s.Filters == nil {
			s.Filters = NewFilters()
		}
		s.Filters.ApplyPreset(preset)
	}

	for _, dep := range s.Deps {
		if preset.ForceFilters() {
			dep.Filters = preset.GetFiltersForDependency(dep)
		} else {
			if dep.Filters == nil {
				dep.Filters = NewFilters()
			}
			dep.Filters.ApplyDep(preset, dep)
		}
	}
}

func (s *Spec) findDep(dependency *Dependency) (*Dependency, bool) {
	for _, dep := range s.Deps {
		if strings.EqualFold(dep.URL, dependency.URL) {
			return dep, true
		}
	}
	return nil, false
}

func checkPreset(preset Preset, warn bool) Preset {
	if preset == nil {
		if warn {
			log.S().Warnf("no preset has been provided, using default preset")
		}
		return &DefaultPreset{}
	}

	if warn {
		switch p := preset.(type) {
		case *DefaultPreset:
			break
		default:
			if p.GetPresetName() == "default" {
				log.S().Warnf("custom preset injected, but the name is \"default\" which is used for the default preset name, this will be confusing for users, consider a different name for your preset")
			}
		}
	}

	return preset
}
