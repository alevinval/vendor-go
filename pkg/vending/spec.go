package vending

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/alevinval/vendor-go/pkg/log"
	"gopkg.in/yaml.v3"
)

// VERSION of the current tool
const VERSION = "v0.4.0"

// Spec holds relevant information related to the specification of what
// versions need to be fetched when updating dependencies.
//
// This model directly maps to the serialized YAML of the spec lock file.
type Spec struct {
	Version    string        `yaml:"version"`
	PresetName string        `yaml:"preset"`
	VendorDir  string        `yaml:"vendor_dir,omitempty"`
	Filters    *Filters      `yaml:",inline"`
	Deps       []*Dependency `yaml:"deps"`
	preset     Preset        `yaml:"-"`
}

// NewSpec allocates a new Spec instance with the default initialization.
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

// AddDependency adds a Dependency to the list of dependencies to vendor.
func (s *Spec) AddDependency(dependency *Dependency) {
	if dep, ok := s.findDep(dependency); ok {
		dep.Update(dependency)
	} else {
		s.Deps = append(s.Deps, dependency)
	}
	s.applyPreset(s.preset)
}

// Load Spec from the filesystem.
func (s *Spec) Load() error {
	preset := s.preset

	filename := preset.GetSpecFilename()
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read: %w", err)
	}

	err = yaml.Unmarshal(data, s)
	if err != nil {
		return fmt.Errorf("cannot parse: %w", err)
	}

	s.applyPreset(preset)
	return nil
}

// Save converts the Spec to YAML, and writes the data in the spec file,
// as specified by the Preset.
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

	if s.Version < VERSION {
		s.Version = VERSION
	}

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

func toYaml(obj interface{}) ([]byte, error) {
	b := &bytes.Buffer{}
	encoder := yaml.NewEncoder(b)
	encoder.SetIndent(2)
	err := encoder.Encode(obj)
	if err != nil {
		return nil, err
	}
	encoder.Close()
	return b.Bytes(), nil
}
