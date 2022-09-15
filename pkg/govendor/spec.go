package govendor

import (
	"os"
	"strings"

	"github.com/alevinval/vendor-go/internal/utils"
	"github.com/alevinval/vendor-go/pkg/govendor/log"
	"gopkg.in/yaml.v3"
)

const (
	VERSION            = "0.1.0"
	SPEC_FILENAME      = ".vendor.yml"
	SPEC_LOCK_FILENAME = ".vendor-lock.yml"
)

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
	if preset == nil {
		logger.Warnf("no preset has been provided, using default preset")
		preset = &DefaultPreset{}
	}

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
	if preset == nil {
		logger.Warnf("no preset has been provided, using default preset")
		preset = &DefaultPreset{}
	}

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
	s.applyPreset(s.Preset)
}

func (s *Spec) Save() error {
	data, err := toYaml(s)
	if err != nil {
		return err
	}

	return os.WriteFile(s.Preset.GetSpecFilename(), data, os.ModePerm)
}

func (s *Spec) applyPreset(preset Preset) {
	s.Preset = preset
	s.Extensions = utils.Union(s.Extensions, s.Preset.GetExtensions())
	for _, dep := range s.Deps {
		dep.Targets = utils.Union(dep.Targets, s.Preset.GetTargets(dep))
		dep.Ignores = utils.Union(dep.Ignores, s.Preset.GetIgnores(dep))
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
