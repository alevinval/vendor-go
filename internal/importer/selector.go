package importer

import (
	"path/filepath"
	"strings"

	"github.com/alevinval/vendor-go/pkg/vending"
)

type Selector struct {
	filters *vending.Filters
}

func newSelector(spec *vending.Spec, dep *vending.Dependency) *Selector {
	filters := spec.Filters.Clone().ApplyFilters(dep.Filters)

	return &Selector{
		filters,
	}
}

func (sel *Selector) Select(path string) bool {
	return sel.isTarget(path) && sel.hasExt(path) && !sel.isIgnored(path)
}

func (sel *Selector) isTarget(path string) bool {
	return hasPrefix(path, sel.filters.Targets) || len(sel.filters.Targets) == 0
}

func (sel *Selector) isIgnored(path string) bool {
	return hasPrefix(path, sel.filters.Ignores)
}

func (sel *Selector) hasExt(path string) bool {
	ext := filepath.Ext(path)
	if ext == "" {
		return false
	}

	// Ignore initial dot that filepath.Ext returns
	ext = ext[1:]

	for _, targetExt := range sel.filters.Extensions {
		if strings.EqualFold(ext, targetExt) {
			return true
		}
	}

	return false
}

func hasPrefix(path string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
