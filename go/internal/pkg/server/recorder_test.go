package server

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"testing"

	"github.com/stretchr/testify/mock"
)

type mockRecorder struct {
	mock.Mock
}

func (mr *mockRecorder) markUnique() {
	mr.Called()
}

func (mr *mockRecorder) markDuplicate() {
	mr.Called()
}

func (mr *mockRecorder) getReport() string {
	return mr.Called().String(0)
}

func Test_recorder_getReport(t *testing.T) {

	type args struct {
		unique int
		dupe   int
	}
	tests := []struct {
		name string
		r    Recorder
		want string
		args args
	}{
		{
			name: "OneOfEach",
			r:    NewRecorder(),
			args: args{unique: 1, dupe: 1},
		},
		{
			name: "Start",
			r:    NewRecorder(),
			args: args{unique: 0, dupe: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.args.unique; i++ {
				tt.r.markUnique()
			}
			for i := 0; i < tt.args.unique; i++ {
				tt.r.markDuplicate()
			}
			exp := fmt.Sprintf("Received %v unique numbers, %v duplicates. Unique total: %v", tt.args.unique, tt.args.dupe, tt.args.unique)
			assert.Equal(t, exp, tt.r.getReport())
			exp = fmt.Sprintf("Received %v unique numbers, %v duplicates. Unique total: %v", 0, 0, tt.args.unique)
			assert.Equal(t, exp, tt.r.getReport())
		})
	}
}
