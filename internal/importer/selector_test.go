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

func TestSelectorSelect_WithTargets(t *testing.T) {
	filtersWithTargets := vending.NewFilters().
		AddExtension("proto").
		AddTarget("target/a").
		AddIgnore("ignored/a", "target/a/ignored")

	sut := Selector{
		filters: filtersWithTargets,
	}

	// Filepaths
	assertSelection(t, sut, "target/a/some-file.proto", true, true)
	assertSelection(t, sut, "target/a/ignored/ignored.proto", false, false)
	assertSelection(t, sut, "ignored/a/ignored.proto", false, false)
	assertSelection(t, sut, "target/a/readme.md", false, true)
	assertSelection(t, sut, "target/a/no-extension", false, true)

	// Dirs
	assertSelection(t, sut, "nontarget/a/b", false, false)
	assertSelection(t, sut, "ignored/a/b", false, false)
	assertSelection(t, sut, "target", false, true)
}

func TestSelectorSelect_WithoutTargets(t *testing.T) {
	filtersWithoutTargets := vending.NewFilters().
		AddExtension("proto").
		AddIgnore("ignored/a", "target/a/ignored")

	sut := Selector{
		filters: filtersWithoutTargets,
	}

	// Filepaths
	assertSelection(t, sut, "target/a/some-file.proto", true, true)
	assertSelection(t, sut, "target/a/ignored/ignored.proto", false, false)
	assertSelection(t, sut, "ignored/a/ignored.proto", false, false)
	assertSelection(t, sut, "target/a/readme.md", false, true)
	assertSelection(t, sut, "target/a/no-extension", false, true)

	// Dirs
	assertSelection(t, sut, "nontarget/a/b", false, true)
	assertSelection(t, sut, "ignored/a/b", false, false)
}

func assertSelection(
	t *testing.T, sut Selector, path string,
	expectedIsSelected, expectedShouldEnterDir bool,
) {
	isSelected := sut.SelectPath(path)
	shouldEnterDir := sut.SelectDir(path)

	assert.Equal(t, expectedIsSelected, isSelected, "isSelected missmatch")
	assert.Equal(t, expectedShouldEnterDir, shouldEnterDir, "shouldEnterDir missmatch")
}
