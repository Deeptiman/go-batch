package batch

type empty struct{}

type Semaphore struct {
	blockch chan empty
	waiter  int
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		blockch: make(chan empty, n),
		waiter:  n,
	}
}

func (s *Semaphore) Acquire(n int) {

	var e empty
	for i := 0; i < n; i++ {
		s.blockch <- e
	}
}

func (s *Semaphore) Release(n int) {

	for i := 0; i < n; i++ {
		<-s.blockch
	}
}

func (s *Semaphore) Lock() {

	s.Acquire(s.waiter)
}

func (s *Semaphore) Unlock() {

	s.Release(s.waiter)
}

func (s *Semaphore) RLock() {

	s.Acquire(1)
}

func (s *Semaphore) RUnlock() {

	s.Release(1)
}
