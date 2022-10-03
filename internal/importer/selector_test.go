package importer

import (
	"testing"

	"github.com/alevinval/vendor-go/pkg/vending"
	"github.com/stretchr/testify/assert"
)

const (
	IS_SELECTED = true
	IS_TARGET   = true
	IS_IGNORED  = true
	HAS_EXT     = true
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

func TestSelectorSelect_WithTargets(t *testing.T) {
	filtersWithTargets := vending.NewFilters().
		AddExtension("proto").
		AddTarget("target/a").
		AddIgnore("ignored/a", "target/a/ignored")

	sut := Selector{
		filters: filtersWithTargets,
	}

	assertSelection(t, sut, "target/a/some-file.proto", IS_SELECTED, IS_TARGET, !IS_IGNORED, HAS_EXT)
	assertSelection(t, sut, "target/a/ignored/ignored.proto", !IS_SELECTED, IS_TARGET, IS_IGNORED, HAS_EXT)
	assertSelection(t, sut, "ignored/a/ignored.proto", !IS_SELECTED, !IS_TARGET, IS_IGNORED, HAS_EXT)
	assertSelection(t, sut, "target/a/readme.md", !IS_SELECTED, IS_TARGET, !IS_IGNORED, !HAS_EXT)
	assertSelection(t, sut, "target/a/no-extension", !IS_SELECTED, IS_TARGET, !IS_IGNORED, !HAS_EXT)
}

func TestSelectorSelect_WithoutTargets(t *testing.T) {
	filtersWithoutTargets := vending.NewFilters().
		AddExtension("proto").
		AddIgnore("ignored/a", "target/a/ignored")

	sut := Selector{
		filters: filtersWithoutTargets,
	}

	assertSelection(t, sut, "target/a/some-file.proto", IS_SELECTED, IS_TARGET, !IS_IGNORED, HAS_EXT)
	assertSelection(t, sut, "target/a/ignored/ignored.proto", !IS_SELECTED, IS_TARGET, IS_IGNORED, HAS_EXT)
	assertSelection(t, sut, "ignored/a/ignored.proto", !IS_SELECTED, IS_TARGET, IS_IGNORED, HAS_EXT)
	assertSelection(t, sut, "target/a/readme.md", !IS_SELECTED, IS_TARGET, !IS_IGNORED, !HAS_EXT)
	assertSelection(t, sut, "target/a/no-extension", !IS_SELECTED, IS_TARGET, !IS_IGNORED, !HAS_EXT)
}

func assertSelection(
	t *testing.T, sut Selector, path string,
	expectedIsSelected, expectedIsTarget, expectedIsIgnored, expectedHasExt bool,
) {
	isSelected, isTarget, isIgnored, hasExt := sut.Select(path)

	assert.Equal(t, expectedIsSelected, isSelected, "isSelected missmatch")
	assert.Equal(t, expectedIsTarget, isTarget, "isTarget missmatch")
	assert.Equal(t, expectedHasExt, hasExt, "hasExt missmatch")
	assert.Equal(t, expectedIsIgnored, isIgnored, "isIgnored missmatch")
}
