package server

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGetWriter(t *testing.T) {
	logger := getLogger()
	dir, err := ioutil.TempDir("", "Writer")
	defer os.RemoveAll(dir)

	if err != nil {
		panic(err)
	}
	type args struct {
		file string
	}
	tests := []struct {
		name       string
		args       args
		assertThat func(args args)
	}{
		{
			name: "WriteToFile",
			args: args{
				file: filepath.Join(dir, "WriteToFile"),
			},
			assertThat: func(args args) {
				w := GetWriter(args.file)
				w.Info("expected text")
				bs, err := ioutil.ReadFile(args.file)
				assert.NoError(t, err)
				assert.Equal(t, "expected text\n", string(bs))
			},
		},
		{
			name: "ClearExistingFile",
			args: args{
				file: filepath.Join(dir, "ClearExistingFile"),
			},
			assertThat: func(args args) {
				err := ioutil.WriteFile(args.file, []byte("this"), os.ModePerm)
				assert.NoError(t, err)
				w := GetWriter(args.file)
				w.Info("expected text")
				bs, err := ioutil.ReadFile(args.file)
				assert.NoError(t, err)
				assert.Equal(t, "expected text\n", string(bs))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.Debug("Running test", zap.String("File", tt.args.file))
			tt.assertThat(tt.args)
		})
	}
}
