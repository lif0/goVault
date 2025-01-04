package concurrency

import (
	"testing"
)

func TestSemaphore_AcquireRelease(t *testing.T) {
	sem := NewSemaphore(2)

	sem.Acquire()
	sem.Acquire()

	select {
	case sem.tickets <- struct{}{}:
		t.Errorf("Semaphore allowed more tickets than expected")
	default:
	}

	sem.Release()
	sem.Release()

	select {
	case sem.tickets <- struct{}{}:
		// nice - semaphone work is correct
	default:
		t.Errorf("Semaphore did not release a ticket as expected")
	}
}

func TestSemaphore_ZeroTickets(t *testing.T) {
	sem := NewSemaphore(0)

	select {
	case sem.tickets <- struct{}{}:
		t.Errorf("Semaphore with 0 tickets should not allow Acquire")
	default:
	}
}

func TestSemaphore_NilSemaphore(t *testing.T) {
	var sem *Semaphore

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Calling methods on nil Semaphore caused panic")
		}
	}()
	sem.Acquire()
	sem.Release()
}
