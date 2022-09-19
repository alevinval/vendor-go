package vending

import (
	"sort"
)

type Filters struct {
	Extensions []string `yaml:"extensions,omitempty"`
	Targets    []string `yaml:"targets,omitempty"`
	Ignores    []string `yaml:"ignores,omitempty"`
}

func NewFilters() *Filters {
	return &Filters{
		Extensions: []string{},
		Targets:    []string{},
		Ignores:    []string{},
	}
}

func (f *Filters) AddExtension(extension ...string) *Filters {
	f.Extensions = append(f.Extensions, extension...)
	return f
}

func (f *Filters) AddTarget(target ...string) *Filters {
	f.Targets = append(f.Targets, target...)
	return f
}

func (f *Filters) AddIgnore(ignore ...string) *Filters {
	f.Ignores = append(f.Ignores, ignore...)
	return f
}

func (f *Filters) ApplyFilters(filters *Filters) *Filters {
	f.Extensions = filtersUnion(f.Extensions, filters.Extensions)
	f.Targets = filtersUnion(f.Targets, filters.Targets)
	f.Ignores = filtersUnion(f.Ignores, filters.Ignores)
	return f
}

func (f *Filters) ApplyPreset(preset Preset) *Filters {
	return f.
		ApplyFilters(preset.GetFilters())
}

func (f *Filters) ApplyDep(preset Preset, dep *Dependency) *Filters {
	return f.
		ApplyFilters(dep.Filters).
		ApplyFilters(preset.GetFiltersForDependency(dep))
}

func (f *Filters) Clone() *Filters {
	return NewFilters().
		AddExtension(f.Extensions...).
		AddTarget(f.Targets...).
		AddIgnore(f.Ignores...)
}

func filtersUnion(a, b []string) []string {
	union := map[string]struct{}{}
	for i := range a {
		union[a[i]] = struct{}{}
	}
	for i := range b {
		union[b[i]] = struct{}{}
	}
	list := make([]string, 0, len(union))
	for key := range union {
		list = append(list, key)
	}
	sort.Strings(list)
	return list
}
