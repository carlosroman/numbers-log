package server

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"os"
)

type Writer interface {
	Info(msg string, fields ...zap.Field)
	Sync() error
}

func GetWriter(file string) Writer {
	cfgJson := fmt.Sprintf(`{
	  "level": "info",
	  "encoding": "console",
	  "development": false,
	  "disableStacktrace": true,
	  "disableCaller": true,
	  "outputPaths": ["%s"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelEncoder": "lowercase",
	    "lineEnding": "\n"
	  }
	}`, url.QueryEscape(file))
	rawJSON := []byte(cfgJson)
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	if _, err := os.Stat(file); err == nil {
		if err := os.Remove(file); err != nil {
			panic(err)
		}
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}
