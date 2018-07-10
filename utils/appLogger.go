package utils

import (
	"github.com/sirupsen/logrus"
	"sync"
	"io"
	"fmt"
)

type AppLogger struct {
	logger *logrus.Logger
	module string
	prefix string
	once sync.Once
}

func NewAppLogger(module string, prefix string) *AppLogger{
	return &AppLogger{
		module:module,
		prefix:prefix,
	}
}

func (l *AppLogger) Writer() io.Writer {
	return l.logger.Writer()
}

func (l *AppLogger) withPrefixf(format string, args ...interface{}) string {
	if len(l.prefix) == 0 {
		return fmt.Sprintf("[%v] %v", l.module, fmt.Sprintf(format, args...))
	}
	return fmt.Sprintf("%v %v", l.prefix, fmt.Sprintf(format, args...))
}

func (l *AppLogger) withPrefix(args ...interface{}) string {
	if len(l.prefix) == 0 {
		return fmt.Sprintf("[%v] %v", l.module, fmt.Sprintf("%v", args...))
	}
	return fmt.Sprintf("%v %v", l.prefix, fmt.Sprintf("%v", args...))
}

//Fatal serverlog
func (l *AppLogger) Fatal(v ...interface{}) {
	l.getLogger().Fatalln(l.withPrefix(v...))
}

//Fatalf serverlog
func (l *AppLogger) Fatalf(format string, v ...interface{}) {
	l.getLogger().Fatal(l.withPrefixf(format, v...))
}

//Fatalln serverlog
func (l *AppLogger) Fatalln(v ...interface{}) {
	l.getLogger().Fatalln(l.withPrefix(v...))
}

//Panic serverlog
func (l *AppLogger) Panic(v ...interface{}) {
	l.getLogger().Panicln(l.withPrefix(v...))
}

//Panicf serverlog
func (l *AppLogger) Panicf(format string, v ...interface{}) {
	l.getLogger().Panicf(format, v...)
}

//Panicln serverlog
func (l *AppLogger) Panicln(v ...interface{}) {
	l.getLogger().Panicln(l.withPrefix(v...))
}

//Print serverlog
func (l *AppLogger) Print(v ...interface{}) {
	l.getLogger().Println(l.withPrefix(v...))
}

//Printf serverlog
func (l *AppLogger) Printf(format string, v ...interface{}) {
	l.getLogger().Printf(l.withPrefixf(format, v...))
}

//Println serverlog
func (l *AppLogger) Println(v ...interface{}) {
	l.getLogger().Println(l.withPrefix(v...))
}

//Debug serverlog
func (l *AppLogger) Debug(args ...interface{}) {
	l.getLogger().Debug(l.withPrefix(args...))
}

//Debugf serverlog
func (l *AppLogger) Debugf(format string, args ...interface{}) {
	l.getLogger().Debugf(l.withPrefixf(format, args...))
}

//Debugln serverlog
func (l *AppLogger) Debugln(args ...interface{}) {
	l.getLogger().Debugln(l.withPrefix(args...))
}

//Info serverlog
func (l *AppLogger) Info(args ...interface{}) {
	l.getLogger().Info(l.withPrefix(args...))
}

//Infof serverlog
func (l *AppLogger) Infof(format string, args ...interface{}) {
	l.getLogger().Infof(l.withPrefixf(format, args...))
}

//Infoln serverlog
func (l *AppLogger) Infoln(args ...interface{}) {
	l.getLogger().Infoln(l.withPrefix(args...))
}

//Warn serverlog
func (l *AppLogger) Warn(args ...interface{}) {
	l.getLogger().Warn(l.withPrefix(args...))
}

//Warnf serverlog
func (l *AppLogger) Warnf(format string, args ...interface{}) {
	l.getLogger().Warnf(l.withPrefixf(format, args...))
}

//Warnln serverlog
func (l *AppLogger) Warnln(args ...interface{}) {
	l.getLogger().Warnln(l.withPrefix(args...))
}

//Error serverlog
func (l *AppLogger) Error(args ...interface{}) {
	l.getLogger().Errorln(l.withPrefix(args...))
}

//Errorf serverlog
func (l *AppLogger) Errorf(format string, args ...interface{}) {
	l.getLogger().Errorf(l.withPrefixf(format, args...))
}

//Errorln serverlog
func (l *AppLogger) Errorln(args ...interface{}) {
	l.getLogger().Errorln(l.withPrefix(args...))
}

func (l *AppLogger) getLogger() *logrus.Logger  {
	l.once.Do(func(){
		ruslogger := logrus.New()
		ruslogger.Formatter = &logrus.TextFormatter{
			ForceColors: true,
			FullTimestamp:true,
		}
		ruslogger.SetLevel(logrus.DebugLevel)
		l.logger = ruslogger
	})

	return l.logger
}
