package index

import "sync"

type StateMonitor struct {
	lock   sync.Mutex
	active bool
	wg     sync.WaitGroup
	ch     chan bool
}

func (s StateMonitor) State() <-chan bool {
	return s.ch
}

func (s StateMonitor) IsActive() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.active
}

func (s *StateMonitor) SetActive() *sync.WaitGroup {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.wg.Add(1)

	if !s.active {
		s.active = true
		s.ch <- true

		go func(wg *sync.WaitGroup) {
			wg.Wait()
			s.setInactive()
		}(&s.wg)
	}
	return &s.wg

}

func (s *StateMonitor) setInactive() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.active = false
	s.ch <- false
}

func NewStateMonitor() *StateMonitor {
	return &StateMonitor{
		ch: make(chan bool, 1),
	}
}
