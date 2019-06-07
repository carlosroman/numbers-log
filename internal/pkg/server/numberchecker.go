package server

import (
	"sync"
)

type NumberChecker interface {
	IsUnique(n uint32) (unique bool)
}

type checker struct {
	mu sync.Mutex
	tm map[uint32]bool
}

func newChecker() *checker {
	return &checker{
		tm: make(map[uint32]bool),
	}
}
func (c *checker) IsUnique(n uint32) (unique bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.tm[n]; !ok {
		c.tm[n] = true
		return true
	}
	return false
}

type checkerImplList struct {
	mu sync.Mutex
	tm []bool
}

func NewNumberChecker() NumberChecker {
	return &checkerImplList{
		tm: make([]bool, 1000000000),
	}
}

func (c *checkerImplList) IsUnique(n uint32) (unique bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ok := c.tm[n]; !ok {
		c.tm[n] = true
		return true
	}
	return false
}
