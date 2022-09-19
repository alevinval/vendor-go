package importer

import (
	"path/filepath"
	"strings"

	"github.com/alevinval/vendor-go/pkg/vending"
)

type Selector struct {
	filters *vending.Filters
}

func NewSelector(spec *vending.Spec, dep *vending.Dependency) *Selector {
	filters := spec.Filters.Clone().ApplyFilters(dep.Filters)
	return &Selector{
		filters,
	}
}

func (sel *Selector) Select(path string) bool {
	return sel.isTarget(path) && sel.hasExt(path) && !sel.isIgnored(path)
}

func (sel *Selector) isTarget(path string) bool {
	if len(sel.filters.Targets) == 0 {
		return true
	}

	for _, target := range sel.filters.Targets {
		if hasPrefix(path, target) {
			return true
		}
	}
	return false
}

func (sel *Selector) hasExt(path string) bool {
	ext := filepath.Ext(path)
	if ext == "" {
		return false
	}
	for _, targetExt := range sel.filters.Extensions {
		// Ignore initial dot that filepath.Ext returns
		if strings.EqualFold(ext[1:], targetExt) {
			return true
		}
	}

	return false
}

func (sel *Selector) isIgnored(path string) bool {
	for _, prefix := range sel.filters.Ignores {
		if hasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func hasPrefix(path, prefix string) bool {
	path = strings.TrimPrefix(path, "/")
	prefix = strings.TrimPrefix(prefix, "/")
	return strings.HasPrefix(path, prefix)
}
