package lock

import (
	"fmt"
	"time"

	"github.com/alevinval/vendor-go/pkg/log"
	"github.com/rogpeppe/go-internal/lockedfile"
)

// Lock is used to ensure only one agent has permission to do certain operations.
// It relies on creating files
type Lock struct {
	file   *lockedfile.File
	path   string
	warn   string
	period time.Duration

	acquired bool
}

// New returns a Lock instance, it does not try to acquire the lock yet, so no
// file will be created until Acquire() is called.
func New(path string) *Lock {
	return &Lock{
		path:   path,
		period: time.Duration(1 * time.Second),
	}
}

// WithWarn configures the warning message that will be printed on slow lock acquisition.
func (l *Lock) WithWarn(warn string) *Lock {
	l.warn = warn
	return l
}

// WithPeriod configures the polling period during lock acquisition, for purposes of
// printing the warning message.
func (l *Lock) WithPeriod(period time.Duration) *Lock {
	l.period = period
	return l
}

// Acquire attempts to lock the file. This operation blocks as long as needed
// until it can be acquired.
func (l *Lock) Acquire() (err error) {
	if l.file, err = l.acquireLock(); err != nil {
		return fmt.Errorf("cannot create lockfile: %w", err)
	} else {
		l.acquired = true
		return nil
	}
}

// Release the lock that was acquired.
func (l *Lock) Release() error {
	if !l.acquired {
		panic("attempted to release a lock that was never acquired")
	}
	return l.file.Close()
}

// acquireLock creates the lockfile in a go-routine and waits for the process to complete.
// If the process takes too long and warning is enabled, it will print a warning message.
func (l *Lock) acquireLock() (*lockedfile.File, error) {
	fileCh := make(chan *lockedfile.File)
	errCh := make(chan error)
	go createLockFile(l.path, fileCh, errCh)
	return waitForAcquire(fileCh, errCh, l.period, l.warn)
}

// waitForAcquire polls the channel and displays a warning message in case the lock
// operation takes longer than expected.
func waitForAcquire(
	fileIn <-chan *lockedfile.File,
	errIn <-chan error,
	period time.Duration,
	warnMsg string,
) (*lockedfile.File, error) {
	warn := warnMsg != ""
	for {
		select {
		case <-time.After(period):
			if warn {
				log.S().Warnf(warnMsg)
				warn = false
			}
		case e := <-errIn:
			return nil, e
		case lock := <-fileIn:
			return lock, nil
		}
	}
}

// createLockFile creates the lockfile and sends it to the output channel. This is done
// so we can poll the channel and show a warning message, rather than blocking indefinetely
// in case the lock is already acquired by someone else.
func createLockFile(path string, fileOut chan *lockedfile.File, errOut chan error) {
	if file, err := lockedfile.Create(path); err != nil {
		errOut <- err
	} else {
		fileOut <- file
	}
}
