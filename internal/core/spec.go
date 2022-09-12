package core

import (
	"strings"
)

const VERSION = "0.0.1"
const SPEC_FILENAME = ".vendor.yml"
const SPEC_LOCK_FILENAME = ".vendor-lock.yml"

type Spec struct {
	Version    string
	Preset     string        `yaml:"preset,omitempty"`
	VendorDir  string        `yaml:"vendor_dir,omitempty"`
	Extensions []string      `yaml:"extensions,omitempty"`
	Targets    []string      `yaml:"targets,omitempty"`
	Ignores    []string      `yaml:"ignores,omitempty"`
	Deps       []*Dependency `yaml:"deps"`
}

func NewSpec() *Spec {
	return &Spec{
		Version:    VERSION,
		VendorDir:  "vendor/",
		Extensions: []string{},
		Targets:    []string{},
		Ignores:    []string{},
		Deps:       []*Dependency{},
	}
}
func (s *Spec) Add(dependency *Dependency) {
	if dep, ok := s.findDep(dependency); ok {
		dep.Update(dependency)
	} else {
		s.Deps = append(s.Deps, dependency)
	}
}

func (s *Spec) ToYaml() []byte {
	return toYaml(s)
}

func (s *Spec) LoadPreset() {
	if s.Preset == "" {
		return
	}
	preset := presets[s.Preset]
	s.Extensions = preset.GetExtensions()
	for _, dep := range s.Deps {
		dep.Targets = preset.GetTargets(dep)
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
