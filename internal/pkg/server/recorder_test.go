package server

import "github.com/stretchr/testify/mock"

type mockRecorder struct {
	mock.Mock
}

func (mr *mockRecorder) markUnique() {
	mr.Called()
}

func (mr *mockRecorder) markDuplicate() {
	mr.Called()
}

func (mr *mockRecorder) printStats() string {
	return mr.Called().String(0)
}
