package server

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

type add interface {
	IsUnique(n uint32) (unique bool)
}

func TestAddOkay(t *testing.T) {
	c := newChecker()
	testAddOkay(t, c)
}

func TestAddDuplicate(t *testing.T) {
	c := newChecker()
	testAddDuplicate(t, c)
}

func TestAddOkayAlt(t *testing.T) {
	c := NewNumberChecker()
	testAddOkay(t, c)
}

func TestAddDuplicateAlt(t *testing.T) {
	c := NewNumberChecker()
	testAddDuplicate(t, c)
}

func testAddOkay(t *testing.T, a add) {
	assert.Equal(t, true, a.IsUnique(1337))
}

func testAddDuplicate(t *testing.T, a add) {
	a.IsUnique(1337)
	assert.Equal(t, false, a.IsUnique(1337))
}

func BenchmarkAdd(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := newChecker()
		for a := uint32(0); a < (1000000000 / 100); a++ {
			assert.Equal(b, true, c.IsUnique(a))
		}
	}
}

func BenchmarkAddAlt(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := NewNumberChecker()
		for a := uint32(0); a < (1000000000 / 100); a++ {
			assert.Equal(b, true, c.IsUnique(a))
		}
	}
}

func BenchmarkRandomAdd(b *testing.B) {
	s := rand.NewSource(42)
	r := rand.New(s)
	c := newChecker()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a := r.Int31n(1000000000)
		c.IsUnique(uint32(a))
	}
}
