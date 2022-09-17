package paths

import (
	"testing"

	"github.com/alevinval/vendor-go/pkg/govendor"
	"github.com/stretchr/testify/assert"
)

func TestPathSelector(t *testing.T) {
	spec := govendor.NewSpec(nil)
	spec.Filters = govendor.NewFilters().
		AddExtension("spec-extension").
		AddTarget("spec-target").
		AddIgnore("spec-ignore")

	dep := govendor.NewDependency("some-url", "some-branch")
	dep.Filters = govendor.NewFilters().
		AddExtension("dep-extension").
		AddTarget("dep-target").
		AddIgnore("dep-ignore")

	sut := NewPathSelector(spec, dep)
	assert.Equal(t, spec.Filters.Clone().ApplyFilters(dep.Filters), sut.filters)
}

func TestPathSelectorSelect(t *testing.T) {
	filtersWithTargets := govendor.NewFilters().
		AddExtension("proto").
		AddTarget("target/a").
		AddIgnore("ignored/a", "target/a/ignored")
	sutWithTargets := PathSelector{
		filters: filtersWithTargets,
	}

	filtersWithoutTargets := govendor.NewFilters().
		AddExtension("proto").
		AddIgnore("ignored/a", "target/a/ignored")
	sutWithoutTargets := PathSelector{
		filters: filtersWithoutTargets,
	}

	for _, sut := range []PathSelector{sutWithTargets, sutWithoutTargets} {
		assert.True(t, sut.Select("target/a/some-file.proto"))
		assert.False(t, sut.Select("target/a/ignored/ignored.proto"))
		assert.False(t, sut.Select("ignored/a/ignored.proto"))
		assert.False(t, sut.Select("target/a/readme.md"))
		assert.False(t, sut.Select("target/a/no-extension"))
	}
}
