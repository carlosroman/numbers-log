package server

import (
	"fmt"
	"go.uber.org/atomic"
)

type Recorder interface {
	markUnique()
	markDuplicate()
	getReport() string
}

func NewRecorder() Recorder {
	return &recorder{}
}

type recorder struct {
	u atomic.Uint32
	d atomic.Uint32
	t atomic.Uint32
}

func (r *recorder) markUnique() {
	r.u.Inc()
	r.t.Inc()
}
func (r *recorder) markDuplicate() {
	r.d.Inc()
}
func (r *recorder) getReport() string {
	return fmt.Sprintf(
		"Received %v unique numbers, %v duplicates. Unique total: %v",
		r.u.Swap(0), r.d.Swap(0), r.t.Load())
}
