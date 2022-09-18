package vendor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDependencyUpdate(t *testing.T) {
	dep := NewDependency("some-url", "some-branch")

	assert.Empty(t, dep.Filters.Targets)
	assert.Empty(t, dep.Filters.Ignores)
	assert.Empty(t, dep.Filters.Extensions)

	other := NewDependency("some-url", "some-branch")
	other.Filters = NewFilters().
		AddExtension("some-extension").
		AddTarget("some-target").
		AddIgnore("some-ignore")

	dep.Update(other)

	assert.Equal(t, other.Filters.Extensions, dep.Filters.Extensions)
	assert.Equal(t, other.Filters.Targets, dep.Filters.Targets)
	assert.Equal(t, other.Filters.Ignores, dep.Filters.Ignores)
}

func TestDependencyID(t *testing.T) {
	dep := NewDependency("some-url", "some-branch")
	assert.Equal(t, "736f6d652d75726ce3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", dep.ID())
}
