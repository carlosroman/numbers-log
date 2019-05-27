package repo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddOkay(t *testing.T) {
	c := newChecker()
	assert.Equal(t, true, c.Add(1337))
}

func TestAddDuplicate(t *testing.T) {
	c := newChecker()
	c.Add(1337)
	assert.Equal(t, false, c.Add(1337))
}

func BenchmarkAdd(b *testing.B) {
	//for i := 0; i < b.N; i++ {
		c := newChecker()
		for a := uint32(0); a < 1000000000; a++ {
			assert.Equal(b, true, c.Add(a))
		}
	//}
}
