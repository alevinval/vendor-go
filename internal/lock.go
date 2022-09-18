package internal

import (
	"fmt"
	"time"

	"github.com/alevinval/vendor-go/pkg/log"
	"github.com/rogpeppe/go-internal/lockedfile"
)

type Lock struct {
	file *lockedfile.File
	path string
	warn string
}

func NewLock(path string) *Lock {
	return &Lock{
		path: path,
	}
}

func (l *Lock) Acquire() (err error) {
	ch := createLock(l.path)
	l.file, err = l.doPoll(ch)
	return
}

func (l *Lock) Release() error {
	return l.file.Close()
}

func (l *Lock) WithWarn(warn string) *Lock {
	l.warn = warn
	return l
}

func (p *Lock) doPoll(ch <-chan *lockedfile.File) (*lockedfile.File, error) {
	warn := p.warn != ""
	for {
		select {
		case <-time.After(1 * time.Second):
			if warn {
				log.S().Warnf(p.warn)
				warn = false
			}
		case lock, ok := <-ch:
			if ok {
				return lock, nil
			} else {
				return nil, fmt.Errorf("cannot acquire lock")
			}
		}
	}
}

func createLock(path string) <-chan *lockedfile.File {
	ch := make(chan *lockedfile.File)
	go func() {
		lock, err := lockedfile.Create(path)
		if err != nil {
			close(ch)
			return
		}
		ch <- lock
	}()
	return ch
}
