package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/lmj/mqtt-clients-demo/lumberjack"
	"github.com/sirupsen/logrus"
)

type fileConfig struct {
	FileName   string
	MaxSize    int
	MagAge     int
	MaxBackups int
	LocalTime  bool
	Compress   bool
	Level      logrus.Level
	Formatter  logrus.Formatter
}

type fileHook struct {
	config fileConfig
	writer io.Writer
}

func newFileHook(config fileConfig) (logrus.Hook, error) {
	hook := fileHook{
		config: config,
	}

	var zeroLevel logrus.Level
	if hook.config.Level == zeroLevel {
		hook.config.Level = logrus.InfoLevel
	}
	var zeroFormatter logrus.Formatter
	if hook.config.Formatter == zeroFormatter {
		hook.config.Formatter = new(logrus.TextFormatter)
	}

	hook.writer = &lumberjack.Logger{
		Filename:   fmt.Sprintf(config.FileName, time.Now().Format("2006-01-02")),
		MaxSize:    config.MaxSize,
		MaxAge:     config.MagAge,
		MaxBackups: config.MaxBackups,
		LocalTime:  config.LocalTime,
		Compress:   config.Compress,
	}

	return &hook, nil
}

func NewFileLogger(c Config, fileds ...string) (Logger, error) {
	logLevel, err := logrus.ParseLevel(c.Level)
	if err != nil {
		log.Printf("failed to parse log level (%s), use default level (info)\n", c.Level)
		logLevel = logrus.InfoLevel
	}
	var fileHook logrus.Hook
	if c.Path != "" {
		err = os.MkdirAll(filepath.Dir(c.Path), 0777)
		if err != nil {
			return nil, fmt.Errorf("failed to create log directory: %s", err.Error())
		} else {
			fileHook, err = newFileHook(fileConfig{
				FileName:   c.Path,
				Formatter:  newFormatter(c.Format),
				Level:      logLevel,
				MagAge:     c.Age.Max,
				MaxSize:    c.Size.Max,
				MaxBackups: c.Backup.Max,
				Compress:   true,
				LocalTime:  true,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create log hook: %s", err.Error())
			}
		}
	}

	entry := logrus.NewEntry(logrus.New())
	entry.Level = logLevel
	entry.Logger.Level = logLevel
	entry.Logger.Formatter = newFormatter(c.Format)
	if fileHook != nil {
		entry.Logger.Hooks.Add(fileHook)
	}
	logrusFields := logrus.Fields{}
	for index := 0; index < len(fileds)-1; index = index + 2 {
		logrusFields[fileds[index]] = fileds[index+1]
	}
	return &logger{entry.WithFields(logrusFields)}, nil
}

func (hook *fileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *fileHook) Fire(entry *logrus.Entry) (err error) {
	if hook.config.Level < entry.Level {
		return nil
	}
	b, err := hook.config.Formatter.Format(entry)
	if err != nil {
		return err
	}
	hook.writer.Write(b)
	return nil
}
