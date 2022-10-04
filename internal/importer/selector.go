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

// SelectPath determines if a filepath should be selected or not.
func (sel *Selector) SelectPath(path string) bool {
	return sel.isTarget(path) && sel.hasExt(path) && !sel.isIgnored(path)
}

// SelectDir determines if a directory should be walked or not.
func (sel *Selector) SelectDir(dir string) bool {
	return !sel.isIgnored(dir) && (len(sel.filters.Targets) == 0 ||
		inverseHasPrefix(sel.filters.Targets, dir))
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

	// Reaching this path means the file does not contain a supported extension
	// Check for a perfect match between any of the targets, in which case we want
	// to select the file anyway.
	return hasPerfectMatch(path, sel.filters.Targets)
}

func hasPrefix(path string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func hasPerfectMatch(path string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if path == prefix {
			return true
		}
	}
	return false
}

// inverseHasPrefix is used to determine when the WalkDirFunc should enter inside
// a directory whenever a path has not been selected.
func inverseHasPrefix(paths []string, prefix string) bool {
	prefix = normDir(prefix)
	for _, path := range paths {
		path = normDir(path)
		if len(prefix) > len(path) {
			path, prefix = prefix, path
		}
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func normDir(path string) string {
	pathDir := filepath.Dir(path)
	if pathDir == "." {
		pathDir = path
	}
	return pathDir
}
