package govendor

import (
	"strings"
)

type SpecLock struct {
	Version string
	Deps    []*DependencyLock `yaml:"deps"`
}

func NewSpecLock() *SpecLock {
	return &SpecLock{
		Version: VERSION,
		Deps:    []*DependencyLock{},
	}
}

func (s *SpecLock) Add(lock *DependencyLock) {
	existing, ok := s.Find(lock.URL)
	if ok {
		existing.Commit = lock.Commit
	} else {
		s.Deps = append(s.Deps, lock)
	}
}

func (s *SpecLock) Find(url string) (*DependencyLock, bool) {
	for _, dep := range s.Deps {
		if strings.EqualFold(dep.URL, url) {
			return dep, true
		}
	}
	return nil, false
}

func (s *SpecLock) ToYaml() []byte {
	return toYaml(s)
}
