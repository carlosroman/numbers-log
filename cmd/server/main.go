package main

import (
	"encoding/json"
	"github.com/carlosroman/numbers-log/internal/pkg/repo"
	"github.com/carlosroman/numbers-log/internal/pkg/server"
	"go.uber.org/zap"
	"os"
)

func main() {

	r := repo.NewRepo()
	l := getLogger()
	defer l.Sync()
	h := server.NewHandler(r, l)
	s := server.NewServer(5, "localhost", 4000, h)
	if err := s.Start(); err != nil {
		os.Exit(2)
	}
	for {
		if err := s.Process(); err != nil {
			os.Exit(1)
		}
	}
}

func getLogger() *zap.Logger {
	rawJSON := []byte(`{
	  "level": "info",
	  "encoding": "console",
	  "outputPaths": ["numbers.log"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}
