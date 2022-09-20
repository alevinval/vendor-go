package internal

import (
	"testing"

	"github.com/alevinval/vendor-go/pkg/vending"
	"github.com/stretchr/testify/assert"
)

func TestCache_RepositoryPath(t *testing.T) {
	dep := vending.NewDependency("some-url", "some-branch")

	actual := getRepositoryPath("some-cache-root", dep)

	assert.Equal(t,
		"some-cache-root/repos/736f6d652d75726ce3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		actual,
	)
}

func TestCache_RepositoryLockPath(t *testing.T) {
	dep := vending.NewDependency("some-url", "some-branch")

	actual := getRepositoryLockPath("some-cache-root", dep)

	assert.Equal(t,
		"some-cache-root/locks/736f6d652d75726ce3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		actual,
	)
}
