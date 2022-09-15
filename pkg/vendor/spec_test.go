package vendor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecAdd(t *testing.T) {
	sut := NewSpec(nil)
	assert.Empty(t, sut.Deps)

	dep := NewDependency("some-url", "some-branch")
	sut.Add(dep)

	assert.Equal(t, []*Dependency{dep}, sut.Deps)
}

func TestSpecAddUpdates(t *testing.T) {
	dep := NewDependency("some-url", "some-branch")
	sut := NewSpec(nil)
	sut.Deps = append(sut.Deps, dep)

	assert.Equal(t, []string{}, sut.Deps[0].Targets)

	other := NewDependency("some-url", "other-branch")
	other.Targets = []string{"other-target"}
	sut.Add(other)

	assert.Equal(t, []*Dependency{dep}, sut.Deps)
	assert.Equal(t, []string{"other-target"}, sut.Deps[0].Targets)
}
