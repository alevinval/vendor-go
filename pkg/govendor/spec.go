package govendor

import (
	"os"
	"strings"

	"github.com/alevinval/vendor-go/internal/log"
	"github.com/alevinval/vendor-go/internal/utils"
	"gopkg.in/yaml.v3"
)

const (
	VERSION            = "0.2.0"
	SPEC_FILENAME      = ".vendor.yml"
	SPEC_LOCK_FILENAME = ".vendor-lock.yml"
)

var logger = log.GetLogger()

type Spec struct {
	Version    string
	Preset     string        `yaml:"preset"`
	VendorDir  string        `yaml:"vendor_dir,omitempty"`
	Extensions []string      `yaml:"extensions,omitempty"`
	Targets    []string      `yaml:"targets,omitempty"`
	Ignores    []string      `yaml:"ignores,omitempty"`
	Deps       []*Dependency `yaml:"deps"`
	preset     Preset        `yaml:"-"`
}

func LoadSpec(preset Preset) (*Spec, error) {
	preset = checkPreset(preset, true)

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
	preset = checkPreset(preset, true)

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
	s.applyPreset(s.preset)
}

func (s *Spec) Save() error {
	data, err := toYaml(s)
	if err != nil {
		return err
	}

	return os.WriteFile(s.preset.GetSpecFilename(), data, os.ModePerm)
}

func (s *Spec) applyPreset(preset Preset) {
	s.preset = preset
	s.Preset = preset.GetPresetName()
	s.Extensions = utils.Union(s.Extensions, s.preset.GetExtensions())
	s.Targets = utils.Union(s.Targets, s.preset.GetTargets())
	s.Ignores = utils.Union(s.Ignores, s.preset.GetIgnores())
	for _, dep := range s.Deps {
		dep.Extensions = utils.Union(dep.Extensions, s.preset.GetDepExtensions(dep))
		dep.Targets = utils.Union(dep.Targets, s.preset.GetDepTargets(dep))
		dep.Ignores = utils.Union(dep.Ignores, s.preset.GetDepIgnores(dep))
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
			logger.Warnf("no preset has been provided, using default preset")
		}
		return &DefaultPreset{}
	}

	if warn {
		switch p := preset.(type) {
		case *DefaultPreset:
			break
		default:
			if p.GetPresetName() == "default" {
				logger.Warnf("custom preset injected, but the name is \"default\" which is used for the default preset name, this will be confusing for users, consider a different name for your preset")
			}
		}
	}

	return preset
}
