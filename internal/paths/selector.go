package paths

import (
	"path/filepath"
	"strings"

	"github.com/alevinval/vendor-go/pkg/vendor"
)

type PathSelector struct {
	Targets    []string
	Ignores    []string
	Extensions []string
}

func NewPathSelector(spec *vendor.Spec, dep *vendor.Dependency) *PathSelector {
	return &PathSelector{
		Targets:    append(spec.Targets, dep.Targets...),
		Ignores:    append(spec.Ignores, dep.Ignores...),
		Extensions: append(spec.Extensions, dep.Extensions...),
	}
}

func (sel *PathSelector) Select(path string) bool {
	return isTarget(path, sel.Targets) && hasExt(path, sel.Extensions) && !isIgnored(path, sel.Ignores)
}

func isTarget(path string, targets []string) bool {
	if len(targets) == 0 {
		return true
	}

	for _, target := range targets {
		if hasPrefix(path, target) {
			return true
		}
	}
	return false
}

func hasExt(path string, extensions []string) bool {
	ext := filepath.Ext(path)
	if ext == "" {
		return false
	}
	for _, targetExt := range extensions {
		// Ignore initial dot that filepath.Ext returns
		if strings.EqualFold(ext[1:], targetExt) {
			return true
		}
	}

	return false
}

func isIgnored(path string, ignores []string) bool {
	for _, prefix := range ignores {
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
