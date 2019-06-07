package server

import (
	"sync"
)

type NumberChecker interface {
	IsUnique(n uint32) (unique bool)
	GetReport() string
}

type checker struct {
	mu sync.Mutex
	tm map[uint32]bool
	r  Recorder
}

func newChecker(r Recorder) NumberChecker {
	return &checker{
		tm: make(map[uint32]bool),
		r:  r,
	}
}

func (c *checker) GetReport() string {
	return c.r.getReport()
}

func (c *checker) IsUnique(n uint32) (unique bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.tm[n]; !ok {
		c.tm[n] = true
		c.r.markUnique()
		return true
	}
	c.r.markDuplicate()
	return false
}

type checkerImplList struct {
	mu sync.Mutex
	tm []bool
	r  Recorder
}

func NewNumberChecker(r Recorder) NumberChecker {
	return &checkerImplList{
		tm: make([]bool, 1000000000),
		r:  r,
	}
}

func (c *checkerImplList) IsUnique(n uint32) (unique bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ok := c.tm[n]; !ok {
		c.tm[n] = true
		c.r.markUnique()
		return true
	}
	c.r.markDuplicate()
	return false
}

func (c *checkerImplList) GetReport() string {
	return c.r.getReport()
}
