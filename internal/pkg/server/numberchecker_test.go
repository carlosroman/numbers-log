package server

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
)

func TestAddOkay(t *testing.T) {
	mr := &mockRecorder{}
	mr.On("markUnique").Return()
	c := newChecker(mr)
	testAddOkay(t, c)
}

func TestAddDuplicate(t *testing.T) {
	mr := &mockRecorder{}
	mr.On("markUnique").Return()
	mr.On("markDuplicate").Return()
	c := newChecker(mr)
	testAddDuplicate(t, c)
}

func TestAddOkayAlt(t *testing.T) {
	mr := &mockRecorder{}
	mr.On("markUnique").Return()
	c := NewNumberChecker(mr)
	testAddOkay(t, c)
}

func TestAddDuplicateAlt(t *testing.T) {
	mr := &mockRecorder{}
	mr.On("markUnique").Return()
	mr.On("markDuplicate").Return()
	c := NewNumberChecker(mr)
	testAddDuplicate(t, c)
}

func testAddOkay(t *testing.T, a NumberChecker) {
	assert.Equal(t, true, a.IsUnique(1337))
}

func testAddDuplicate(t *testing.T, a NumberChecker) {
	a.IsUnique(1337)
	assert.Equal(t, false, a.IsUnique(1337))
}

func BenchmarkAdd(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := newChecker(&noopRecorder{})
		wg := sync.WaitGroup{}
		wg.Add(5)
		max := 1000000
		for r := 0; r < 5; r++ {
			go func(rr int) {
				t := uint32(max * (rr + 1))
				for a := uint32(max * (rr)); a < t; a++ {
					assert.Equal(b, true, c.IsUnique(a))
				}
				wg.Done()
			}(r)
		}
		wg.Wait()
	}
}

func BenchmarkAddAlt(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := NewNumberChecker(&noopRecorder{})
		wg := sync.WaitGroup{}
		wg.Add(5)
		max := 1000000
		for r := 0; r < 5; r++ {
			go func(rr int) {
				t := uint32(max * (rr + 1))
				for a := uint32(max * (rr)); a < t; a++ {
					assert.Equal(b, true, c.IsUnique(a))
				}
				wg.Done()
			}(r)
		}
		wg.Wait()
	}
}

func BenchmarkRandomAdd(b *testing.B) {
	s := rand.NewSource(42)
	r := rand.New(s)
	c := newChecker(&noopRecorder{})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a := r.Int31n(1000000000)
		c.IsUnique(uint32(a))
	}
}

type noopRecorder struct {
}

func (n *noopRecorder) markUnique() {

}
func (n *noopRecorder) markDuplicate() {

}
func (n *noopRecorder) getReport() string {
	return "noop"
}
