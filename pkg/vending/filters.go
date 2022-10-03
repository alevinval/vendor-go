package vending

import (
	"sort"
)

// Filters is a collection type. When vendoring dependencies we look if the
// file paths are copied or ignored, or wether the extension is supported.
// Filters are directly serialized into the output YAML.
//
// Filters supposed to make life easier proving a builder pattern so it's
// convenient to create, but also merge operations with other Filters or
// even Presets and Dependencies. This is really handy because this tool
// supports both global customization, and specific granular configurations
// for each dependency. For this reason, we often need to combine multiple
// Filters before we're able to vendor a dependency.
type Filters struct {
	Extensions []string `yaml:"extensions,omitempty"`
	Targets    []string `yaml:"targets,omitempty"`
	Ignores    []string `yaml:"ignores,omitempty"`
}

// NewFilters allocates a Filters instance with empty lists initialization.
func NewFilters() *Filters {
	return &Filters{
		Extensions: []string{},
		Targets:    []string{},
		Ignores:    []string{},
	}
}

// AddExtension adds one or several extensions to the collection.
func (f *Filters) AddExtension(extension ...string) *Filters {
	f.Extensions = append(f.Extensions, extension...)
	return f
}

// AddTarget adds one or several target paths to the collection.
func (f *Filters) AddTarget(target ...string) *Filters {
	f.Targets = append(f.Targets, target...)
	return f
}

// AddIgnore adds one or several ignored paths to the collection.
func (f *Filters) AddIgnore(ignore ...string) *Filters {
	f.Ignores = append(f.Ignores, ignore...)
	return f
}

// ApplyFilters merges the current Filters instance with another one.
// Because this method combines instances, it stabilizes the results
// by doing the intersection, but also sorting the entries.
func (f *Filters) ApplyFilters(filters *Filters) *Filters {
	f.Extensions = sortedUnion(f.Extensions, filters.Extensions)
	f.Targets = sortedUnion(f.Targets, filters.Targets)
	f.Ignores = sortedUnion(f.Ignores, filters.Ignores)
	return f
}

// ApplyPreset applies the Filters of a Preset.
func (f *Filters) ApplyPreset(preset Preset) *Filters {
	if preset.ForceFilters() {
		f.Clear()
	}
	return f.
		ApplyFilters(preset.GetFilters())
}

// ApplyPresetForDependency applies the filters of a Preset, and a Dependency.
func (f *Filters) ApplyPresetForDependency(preset Preset, dep *Dependency) *Filters {
	if preset.ForceFilters() {
		f.Clear()
	}
	return f.
		ApplyFilters(dep.Filters).
		ApplyFilters(preset.GetFiltersForDependency(dep))
}

// Clear resets the lists.
func (f *Filters) Clear() *Filters {
	f.Extensions = []string{}
	f.Targets = []string{}
	f.Ignores = []string{}
	return f
}

// Clone allocates a new instance, and clones the contents of the current one.
func (f *Filters) Clone() *Filters {
	return NewFilters().
		AddExtension(f.Extensions...).
		AddTarget(f.Targets...).
		AddIgnore(f.Ignores...)
}

func sortedUnion(a, b []string) []string {
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
