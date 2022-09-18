package log

import (
	"encoding/json"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	Level  = zap.NewAtomicLevelAt(zapcore.InfoLevel)
)

func init() {
	logger = zap.Must(getDefaultCfg().Build())
	defer logger.Sync()
}

func S() *zap.SugaredLogger {
	return logger.Sugar()
}

func SetLogger(newLogger *zap.Logger) {
	logger = newLogger
}

func getDefaultCfg() zap.Config {
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

	cfg.Level = Level
	return cfg
}
