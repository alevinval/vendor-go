package main

import (
	"encoding/json"

	"github.com/alevinval/vendor-go/pkg/cmd"
	"github.com/alevinval/vendor-go/pkg/govendor"
	"go.uber.org/zap"
)

func getZapCfg() zap.Config {
	rawJSON := []byte(`{
		"level": "info",
		"encoding": "console",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	cfg.Level = cmd.LogLevel
	return cfg
}

func main() {
	logger := zap.Must(getZapCfg().Build())
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	cmd.NewVendorCmd("vendor", &govendor.DefaultPreset{}).Execute()
}
