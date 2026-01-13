package be

import (
	"context"
	"log/slog"
	"path/filepath"

	"github.com/twiglab/doggy/pkg/human"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogAction struct {
	logDensity *slog.Logger // density 密度
	logQueue   *slog.Logger // queue 排队长度
	logCount   *slog.Logger // count 人数
	logDir     string
}

const dir = "belogs"

func NewLogAction(logDir string) LogAction {
	path := dir
	if logDir == "" {
		path = logDir
	}

	return LogAction{
		logDensity: newLog(logfile(path, human.DENSITY)),
		logQueue:   newLog(logfile(path, human.QUEUE)),
		logCount:   newLog(logfile(path, human.COUNT)),
		logDir:     path,
	}
}

func (d LogAction) Name() string {
	return LOG
}

func (d LogAction) HandleData(ctx context.Context, data human.DataMix) error {
	switch data.Type {
	case human.COUNT:
		d.logCount.InfoContext(ctx, human.COUNT, slog.Any("data", data))
	case human.QUEUE:
		d.logQueue.InfoContext(ctx, human.QUEUE, slog.Any("data", data))
	case human.DENSITY:
		d.logDensity.InfoContext(ctx, human.DENSITY, slog.Any("data", data))
	}
	return nil
}

func logfile(dir, file string) string {
	logf := file + ".log"
	return filepath.Join(dir, logf)
}

func newLog(logFile string) *slog.Logger {
	out := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // megabytes
		MaxBackups: 30,
		MaxAge:     30, //days
	}
	return slog.New(slog.NewJSONHandler(out, nil))
}
