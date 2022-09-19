package importer

import (
	"testing"

	"github.com/alevinval/vendor-go/pkg/vendoring"
	"github.com/stretchr/testify/assert"
)

func TestSelector(t *testing.T) {
	spec := vendoring.NewSpec(nil)
	spec.Filters = vendoring.NewFilters().
		AddExtension("spec-extension").
		AddTarget("spec-target").
		AddIgnore("spec-ignore")

	dep := vendoring.NewDependency("some-url", "some-branch")
	dep.Filters = vendoring.NewFilters().
		AddExtension("dep-extension").
		AddTarget("dep-target").
		AddIgnore("dep-ignore")

	sut := NewSelector(spec, dep)
	assert.Equal(t, spec.Filters.Clone().ApplyFilters(dep.Filters), sut.filters)
}

func TestSelectorSelect(t *testing.T) {
	filtersWithTargets := vendoring.NewFilters().
		AddExtension("proto").
		AddTarget("target/a").
		AddIgnore("ignored/a", "target/a/ignored")
	sutWithTargets := Selector{
		filters: filtersWithTargets,
	}

	filtersWithoutTargets := vendoring.NewFilters().
		AddExtension("proto").
		AddIgnore("ignored/a", "target/a/ignored")
	sutWithoutTargets := Selector{
		filters: filtersWithoutTargets,
	}

	for _, sut := range []Selector{sutWithTargets, sutWithoutTargets} {
		assert.True(t, sut.Select("target/a/some-file.proto"))
		assert.False(t, sut.Select("target/a/ignored/ignored.proto"))
		assert.False(t, sut.Select("ignored/a/ignored.proto"))
		assert.False(t, sut.Select("target/a/readme.md"))
		assert.False(t, sut.Select("target/a/no-extension"))
	}
}
