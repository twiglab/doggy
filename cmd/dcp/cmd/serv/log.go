package serv

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

func isConsole(logFile string) bool {
	if logFile == "" || logFile == "console" {
		return true
	}
	return false
}

func RootLog(id string, logFile string, level slog.Level) *slog.Logger {
	var out io.Writer = os.Stdout
	if !isConsole(logFile) {
		out = &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    10, // megabytes
			MaxBackups: 10,
			MaxAge:     10, //days
		}
	}
	h := slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level}).WithAttrs([]slog.Attr{slog.String("id", id)})
	logger := slog.New(h)
	slog.SetDefault(logger)
	return logger
}

func buildRootLog(ctx context.Context, v *viper.Viper) (*slog.Logger, context.Context) {
	logFile := v.GetString("log.root.file")
	if logFile == "" {
		logFile = "dcp.log"
	}

	var level slog.Level

	l := v.GetString("log.root.level")

	switch strings.ToLower(l) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	id := v.GetString("id")

	logger := RootLog(id, logFile, level)
	return logger, context.WithValue(ctx, keyRootLog, logger)
}
