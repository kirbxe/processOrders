package order

import "sync"

type Stats struct {
	ordersReadyCount int
	mtx              sync.RWMutex
}

func (s *Stats) Inc() {
	s.mtx.Lock()
	s.ordersReadyCount++
	s.mtx.Unlock()
}


func (s *Stats) GetStats() int {
		
	s.mtx.RLock()
	count := s.ordersReadyCount
	s.mtx.RUnlock()
	
	return count

}
