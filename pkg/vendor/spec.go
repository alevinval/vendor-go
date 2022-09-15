package vendor

import (
	"os"
	"strings"

	"github.com/alevinval/vendor-go/internal/log"
	"gopkg.in/yaml.v3"
)

const VERSION = "0.1.0"
const SPEC_FILENAME = ".vendor.yml"
const SPEC_LOCK_FILENAME = ".vendor-lock.yml"

var logger = log.GetLogger()

type Spec struct {
	Version    string
	VendorDir  string        `yaml:"vendor_dir,omitempty"`
	Extensions []string      `yaml:"extensions,omitempty"`
	Targets    []string      `yaml:"targets,omitempty"`
	Ignores    []string      `yaml:"ignores,omitempty"`
	Deps       []*Dependency `yaml:"deps"`
	Preset     Preset        `yaml:"-"`
}

func LoadSpec(preset Preset) (*Spec, error) {
	data, err := os.ReadFile(preset.GetSpecFilename())
	if err != nil {
		logger.Errorf("cannot read %s: %s", preset.GetSpecFilename(), err)
		return nil, err
	}

	spec := &Spec{}
	err = yaml.Unmarshal(data, spec)
	if err != nil {
		logger.Errorf("cannot read %s: %s", preset.GetSpecFilename(), err)
		return nil, err
	}

	spec.applyPreset(preset)
	return spec, nil
}

func NewSpec(preset Preset) *Spec {
	spec := &Spec{
		Version:    VERSION,
		VendorDir:  "vendor/",
		Extensions: []string{},
		Targets:    []string{},
		Ignores:    []string{},
		Deps:       []*Dependency{},
	}
	spec.applyPreset(preset)
	return spec
}

func (s *Spec) Add(dependency *Dependency) {
	if dep, ok := s.findDep(dependency); ok {
		dep.Update(dependency)
	} else {
		s.Deps = append(s.Deps, dependency)
	}
}

func (s *Spec) Save() error {
	data, err := toYaml(s)
	if err != nil {
		return err
	}

	return os.WriteFile(s.Preset.GetSpecFilename(), data, os.ModePerm)
}

func (s *Spec) applyPreset(preset Preset) {
	if preset == nil {
		return
	}
	s.Preset = preset
	s.Extensions = s.Preset.GetExtensions()
	for _, dep := range s.Deps {
		dep.Targets = s.Preset.GetTargets(dep)
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
