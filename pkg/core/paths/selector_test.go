package paths

import (
	"testing"

	"github.com/alevinval/vendor-go/pkg/core"

	"github.com/stretchr/testify/assert"
)

func TestPathSelector(t *testing.T) {
	spec := core.NewSpec()
	spec.Targets = []string{"spec-target"}
	spec.Ignores = []string{"spec-ignores"}
	spec.Extensions = []string{"spec-extensions"}

	dep := core.NewDependency("some-url", "some-branch")
	dep.Targets = []string{"dep-target"}
	dep.Ignores = []string{"dep-ignores"}
	dep.Extensions = []string{"dep-extensions"}

	sut := NewPathSelector(spec, dep)
	assert.Equal(t, []string{"spec-target", "dep-target"}, sut.Targets)
	assert.Equal(t, []string{"spec-ignores", "dep-ignores"}, sut.Ignores)
	assert.Equal(t, []string{"spec-extensions", "dep-extensions"}, sut.Extensions)
}

func TestPathSelectorSelect(t *testing.T) {
	sutWithTargets := PathSelector{
		Targets:    []string{"target/a"},
		Ignores:    []string{"ignored/a", "target/a/ignored"},
		Extensions: []string{"proto"},
	}

	sutWithoutTargets := PathSelector{
		Targets:    []string{},
		Ignores:    []string{"ignored/a", "target/a/ignored"},
		Extensions: []string{"proto"},
	}

	for _, sut := range []PathSelector{sutWithTargets, sutWithoutTargets} {
		assert.True(t, sut.Select("target/a/some-file.proto"))
		assert.False(t, sut.Select("target/a/ignored/ignored.proto"))
		assert.False(t, sut.Select("ignored/a/ignored.proto"))
		assert.False(t, sut.Select("target/a/readme.md"))
		assert.False(t, sut.Select("target/a/no-extension"))
	}
}
