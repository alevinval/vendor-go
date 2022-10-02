package importer

import (
	"testing"

	"github.com/alevinval/vendor-go/pkg/vending"
	"github.com/stretchr/testify/assert"
)

func TestSelector(t *testing.T) {
	spec := vending.NewSpec(nil)
	spec.Filters = vending.NewFilters().
		AddExtension("spec-extension").
		AddTarget("spec-target").
		AddIgnore("spec-ignore")

	dep := vending.NewDependency("some-url", "some-branch")
	dep.Filters = vending.NewFilters().
		AddExtension("dep-extension").
		AddTarget("dep-target").
		AddIgnore("dep-ignore")

	sut := newSelector(spec, dep)
	assert.Equal(t, spec.Filters.Clone().ApplyFilters(dep.Filters), sut.filters)
}

func TestSelectorSelect(t *testing.T) {
	filtersWithTargets := vending.NewFilters().
		AddExtension("proto").
		AddTarget("target/a").
		AddIgnore("ignored/a", "target/a/ignored")
	sutWithTargets := Selector{
		filters: filtersWithTargets,
	}

	filtersWithoutTargets := vending.NewFilters().
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
