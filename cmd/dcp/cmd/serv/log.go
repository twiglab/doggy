package serv

import (
	"log/slog"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

func RootLog(logFile string, level slog.Level) *slog.Logger {
	out := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     10, //days
	}
	h := slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level})
	logger := slog.New(h)
	slog.SetDefault(logger)
	return logger
}

func BuildRootLog(conf AppConf) *slog.Logger {
	logFile := conf.LoggerConf.LogFile
	if logFile == "" {
		logFile = "dcp.log"
	}

	var level slog.Level

	l := strings.ToUpper(conf.LoggerConf.Level)
	switch l {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	return RootLog(logFile, level)
}
