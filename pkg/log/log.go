package log

import (
	"encoding/json"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Level is the atomic level that will be used for logging with the default
	// logger of the library. When injecting a new zap logger ensure the Level
	// field points to this one.
	//
	// Level is used by the cmd package when creating the root command, which
	// supports a flag to enable debug log level.
	Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	// logger keeps reference to the current logger that will be used by the
	// vendoring tool.
	logger *zap.Logger
)

func init() {
	logger = zap.Must(getDefaultCfg().Build())
	defer logger.Sync()
}

// S returns the zap suggared logger from the current logger instance.
func S() *zap.SugaredLogger {
	return logger.Sugar()
}

// SetLogger replaces the logger instance.
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
