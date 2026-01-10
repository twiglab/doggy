package be

import (
	"context"
	"errors"
	"log/slog"
	"path/filepath"

	"github.com/twiglab/doggy/pkg/human"
	"gopkg.in/natefinch/lumberjack.v2"
)

var ErrUnimplType = errors.New("unsupport type")

const (
	TAOS = "taos"
	MQTT = "mqtt" // mqtt 3.11
	LOG  = "log"

	MQTT5 = "mqtt5" // mqtt5 保留
	HTTP  = "http"  // 保留，暂不提供

	NOOP = "noop" // 保留，仅作占位符
)

func HasHuman(in, out int) bool {
	return in != 0 || out != 0
}

func HasCount(count int) bool {
	return count != 0
}

type DataHandler interface {
	HandleData(ctx context.Context, data human.DataMix) error
	Name() string
}

type MutiAction struct {
	actions []DataHandler
}

func NewMutiAction(actions ...DataHandler) *MutiAction {
	return &MutiAction{actions: actions}
}

func (a MutiAction) Name() string {
	return "muti"
}

func (a *MutiAction) Add(h DataHandler) *MutiAction {
	if h != nil {
		a.actions = append(a.actions, h)
	}
	return a
}

func (a MutiAction) HandleData(ctx context.Context, data human.DataMix) error {
	for _, action := range a.actions {
		if err := action.HandleData(ctx, data); err != nil {
			return err
		}
	}
	return nil
}

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
