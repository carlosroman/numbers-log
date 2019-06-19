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
	c := newMapChecker(mr)
	testAddOkay(t, c)
}

func TestAddDuplicate(t *testing.T) {
	mr := &mockRecorder{}
	mr.On("markUnique").Return()
	mr.On("markDuplicate").Return()
	c := newMapChecker(mr)
	testAddDuplicate(t, c)
}

func TestAddOkayAlt(t *testing.T) {
	mr := &mockRecorder{}
	mr.On("markUnique").Return()
	c := newAltChecker(mr)
	testAddOkay(t, c)
}

func TestAddDuplicateAlt(t *testing.T) {
	mr := &mockRecorder{}
	mr.On("markUnique").Return()
	mr.On("markDuplicate").Return()
	c := newAltChecker(mr)
	testAddDuplicate(t, c)
}

func TestAddOkayABool(t *testing.T) {
	mr := &mockRecorder{}
	mr.On("markUnique").Return()
	c := newBoolListChecker(mr)
	testAddOkay(t, c)
}

func TestAddDuplicateABool(t *testing.T) {
	mr := &mockRecorder{}
	mr.On("markUnique").Return()
	mr.On("markDuplicate").Return()
	c := newBoolListChecker(mr)
	testAddDuplicate(t, c)
}

func testAddOkay(t *testing.T, a NumberChecker) {
	assert.Equal(t, true, a.IsUnique(1337))
}

func testAddDuplicate(t *testing.T, a NumberChecker) {
	a.IsUnique(1337)
	assert.Equal(t, false, a.IsUnique(1337))
}

func BenchmarkChecker(b *testing.B) {
	benchmarks := []struct {
		name    string
		checker NumberChecker
	}{
		{
			name:    "Map",
			checker: newMapChecker(&noopRecorder{}),
		},
		{
			name:    "Alt",
			checker: newAltChecker(&noopRecorder{}),
		},
		{
			name:    "Bool",
			checker: newBoolListChecker(&noopRecorder{}),
		},
	}

	max := 10000000
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				wg := sync.WaitGroup{}
				wg.Add(5)
				for r := 0; r < 5; r++ {
					go func(rr int) {
						t := uint32(max * (rr + 1))
						for a := uint32(max * (rr)); a < t; a++ {
							bm.checker.IsUnique(a)
						}
						wg.Done()
					}(r)
				}
				wg.Wait()
			}
		})
	}
}

func BenchmarkRandomAdd(b *testing.B) {
	benchmarks := []struct {
		name    string
		checker NumberChecker
	}{
		{
			name:    "Map",
			checker: newMapChecker(&noopRecorder{}),
		},
		{
			name:    "Alt",
			checker: newAltChecker(&noopRecorder{}),
		},
		{
			name:    "Bool",
			checker: newBoolListChecker(&noopRecorder{}),
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				wg := sync.WaitGroup{}
				wg.Add(5)
				for r := 0; r < 5; r++ {
					go func() {
						s := rand.NewSource(42)
						rn := rand.New(s)
						for do := 0; do < 1000; do++ {
							a := rn.Int63n(100000000)
							bm.checker.IsUnique(uint32(a))
						}
						wg.Done()
					}()
				}
				wg.Wait()
			}
		})
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
