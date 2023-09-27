package logger

import (
	"log"
	"strings"
	"time"

	"github.com/lmj/mqtt-clients-demo/config"
	"github.com/sirupsen/logrus"
)

const (
	persistenceFile = "file"
	persistencePG   = "postgres"
)

var Log Logger

type Logger interface {
	WithField(key string, value interface{}) Logger
	WithError(err error) Logger
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
}

type logger struct {
	entry *logrus.Entry
}

func (l *logger) WithField(key string, value interface{}) Logger {
	return &logger{l.entry.WithField(key, value)}
}
func (l *logger) WithError(err error) Logger {
	return &logger{l.entry.WithError(err)}
}
func (l *logger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}
func (l *logger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}
func (l *logger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}
func (l *logger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}
func (l *logger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}
func (l *logger) Debugln(args ...interface{}) {
	l.entry.Debugln(args...)
}
func (l *logger) Infoln(args ...interface{}) {
	l.entry.Infoln(args...)
}
func (l *logger) Warnln(args ...interface{}) {
	l.entry.Warnln(args...)
}
func (l *logger) Errorln(args ...interface{}) {
	l.entry.Errorln(args...)
}
func (l *logger) Fatalln(args ...interface{}) {
	l.entry.Fatalln(args...)
}

type Config struct {
	Path   string `yaml:"path" json:"path"`
	Level  string `yaml:"level" json:"level" default:"info" validate:"regexp=^(info|debug|warn|error|fatal)$"`
	Format string `yaml:"format" json:"format" default:"text" validate:"regexp=^(text|json)$"`
	Age    struct {
		Max int `yaml:"max" json:"max" default:"15" validate:"min=1"`
	} `yaml:"age" json:"age"` //天
	Size struct {
		Max int `yaml:"max" json:"max" default:"50" validate:"min=1"`
	} `yaml:"size" json:"size"` //MB
	Backup struct {
		Max int `yaml:"max" json:"max" default:"15" validate:"min=1"`
	} `yaml:"backup" json:"backup"` //MB
	Store string `yaml:"store" json:"store"`
}

func newFormatter(format string) logrus.Formatter {
	var formatter logrus.Formatter
	if strings.ToLower(format) == "json" {
		formatter = &logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano}
	} else {
		formatter = &logrus.TextFormatter{TimestampFormat: time.RFC3339Nano, FullTimestamp: true, DisableColors: true}
	}
	return formatter
}

func Init() {
	var err error
	c := Config{
		Path:  config.LOG_PATH,
		Level: config.LOG_LEVEL,
		Age: struct {
			Max int `yaml:"max" json:"max" default:"15" validate:"min=1"`
		}{Max: config.LOG_AGE},
		Size: struct {
			Max int `yaml:"max" json:"max" default:"50" validate:"min=1"`
		}{Max: config.LOG_SIZE},
		Backup: struct {
			Max int `yaml:"max" json:"max" default:"15" validate:"min=1"`
		}{Max: config.LOG_BACKUP_COUNT},
		Store: config.LOG_STORE,
	}
	switch c.Store {
	case persistenceFile:
		Log, err = NewFileLogger(c, "service", "mqttClientsDemo")
		if err != nil {
			log.Panic(err)
		}
	case persistencePG:
		log.Printf("暂未支持psql数据存储方式")
	default:
		log.Printf("存储方式未知: %s", c.Store)
	}
}
