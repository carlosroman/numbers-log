package server

import (
	"sync"
)

func NewNumberChecker(r Recorder) NumberChecker {
	return newBoolListChecker(r)
}

type NumberChecker interface {
	IsUnique(n uint32) (unique bool)
	GetReport() string
}

type checker struct {
	mu sync.Mutex
	tm map[uint32]bool
	r  Recorder
}

func newMapChecker(r Recorder) NumberChecker {
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

func newAltChecker(r Recorder) NumberChecker {
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

type aBool struct {
	marked bool
	mu     sync.Mutex
}

func (a *aBool) mark() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.marked == false {
		a.marked = true
		return true
	}
	return false

}

func newBoolListChecker(r Recorder) NumberChecker {
	return &checkerImplABoolList{
		tm: make([]aBool, 1000000000),
		r:  r,
	}
}

type checkerImplABoolList struct {
	tm []aBool
	r  Recorder
}

func (c *checkerImplABoolList) IsUnique(n uint32) (unique bool) {
	if ok := c.tm[n].mark(); ok {
		c.r.markUnique()
		return ok
	}
	c.r.markDuplicate()
	return false
}

func (c *checkerImplABoolList) GetReport() string {
	return c.r.getReport()
}
