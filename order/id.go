package order

import "sync"

type IDGen struct {
	ID int
	mtx sync.Mutex
	
}

func (id *IDGen) NextGen() int {
	id.mtx.Lock()
	defer id.mtx.Unlock()
	
	id.ID++
	v := id.ID
	return v
}
