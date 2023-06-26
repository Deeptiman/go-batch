package logger

import (
	"github.com/sirupsen/logrus"
)

type LogLevel int

const (
	Error LogLevel = iota
	Info
	Debug		
)

type Logger struct {
	log 	*logrus.Logger
}

func NewLogger() *Logger {
	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		ForceColors:     true,
		DisableTimestamp: true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	// only log errors by default
	log.SetLevel(logrus.ErrorLevel)

	return &Logger{
		log: log,
	}
}

func (l *Logger) SetLogLevel(level LogLevel) {
	if level == Debug {
		l.log.Level = logrus.DebugLevel
	} else if level == Info {
		l.log.Level = logrus.InfoLevel
	} else {
		l.log.Level = logrus.ErrorLevel
	}
}

func (l *Logger) Trace(format string, args ...interface{}) {

}

func (l *Logger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *Logger) Debugln(args ...interface{}) {
	l.log.Debugln(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.log.Infoln(args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log.Warn(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *Logger) Warnln(format string, args ...interface{}) {
	l.log.Warnln(args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log.Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *Logger) Fatalln(format string, args ...interface{}) {
	l.log.Fatalln(args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *Logger) Errorln(format string, args ...interface{}) {
	l.log.Errorln(args...)
}

func (l *Logger) WithField(key string, value interface{}) {
	l.log.WithField(key, value)
}

func (l *Logger) WithFields(fields logrus.Fields) {
	l.log.WithFields(fields)
}
