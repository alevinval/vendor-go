package cache

import "os"

var _ fsOps = (*defaultFs)(nil)

type (
	fsOps interface {
		MkdirAll(path string, perm os.FileMode) error
		RemoveAll(path string) error
	}
	defaultFs struct{}
)

func (defaultFs) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (defaultFs) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
