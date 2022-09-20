package lock

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func cleanUp(locks ...string) {
	for _, lock := range locks {
		os.Remove(lock)
	}
}

func TestLock_NoConcurrency(t *testing.T) {
	defer cleanUp("LOCK")

	sut := New("LOCK")

	err := sut.Acquire()
	assert.NoError(t, err)

	err = sut.Release()
	assert.NoError(t, err)
}

func TestLock_WithConcurrency_DoesBlock(t *testing.T) {
	defer cleanUp("LOCK")

	one := New("LOCK")
	two := New("LOCK")
	two.WithWarn("not able to lock within period as expected").WithPeriod(90 * time.Millisecond)

	err := one.Acquire()
	assert.NoError(t, err)

	select {
	case <-tryAcquire(t, two):
		t.Errorf("should not have been able to acquire lock")
	case <-time.After(100 * time.Millisecond):
	}
}

func tryAcquire(t *testing.T, l *Lock) chan struct{} {
	out := make(chan struct{})

	go func() {
		err := l.Acquire()
		assert.NoError(t, err)

		out <- struct{}{}
	}()

	return out
}
