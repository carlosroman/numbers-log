package server

type Recorder interface {
	markUnique()
	markDuplicate()
	printStats() string
}

func NewRecorder() Recorder {
	return nil
}

type recorder struct {
}

func (r *recorder) markUnique()        {}
func (r *recorder) markDuplicate()     {}
func (r *recorder) printStats() string { return "" }
