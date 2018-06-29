package serverlog

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/core/logging/api"
	"github.com/sirupsen/logrus"
	"log"
	"fmt"
)


type consoleLogProvider struct {
	logrusLogger *logrus.Logger
}

func GetConsoleLogProvider() api.LoggerProvider {
	logger:= logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:true,
	}

	provider:= &consoleLogProvider{
		logrusLogger: logger,
	}

	return provider
}

func (l *consoleLogProvider) GetLogger(module string) api.Logger {

	clogger := log.New(l.logrusLogger.Writer(), fmt.Sprintf("[%s] ", module), log.Ldate|log.Ltime|log.LUTC)
	//clogger.SetOutput(l.logrusLogger.Writer())
	return &ConsoleLogger{logger: clogger, loggerus:l.logrusLogger, module: module}
}

type ConsoleLogger struct {
	logger *log.Logger
	module string
	loggerus *logrus.Logger
}


//Fatal serverlog
func (l *ConsoleLogger) Fatal(v ...interface{}) {
	l.logger.Fatalln(v...)
}

//Fatalf serverlog
func (l *ConsoleLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatalf(format, v...)
}

//Fatalln serverlog
func (l *ConsoleLogger) Fatalln(v ...interface{}) {
	l.logger.Fatalln(v...)
}

//Panic serverlog
func (l *ConsoleLogger) Panic(v ...interface{}) {
	l.logger.Panicln(v...)
}

//Panicf serverlog
func (l *ConsoleLogger) Panicf(format string, v ...interface{}) {
	l.logger.Panicf(format, v...)
}

//Panicln serverlog
func (l *ConsoleLogger) Panicln(v ...interface{}) {
	l.logger.Panicln(v...)
}

//Print serverlog
func (l *ConsoleLogger) Print(v ...interface{}) {
	l.logger.Println(v...)
}

//Printf serverlog
func (l *ConsoleLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

//Println serverlog
func (l *ConsoleLogger) Println(v ...interface{}) {
	l.logger.Println(v...)
}

//Debug serverlog
func (l *ConsoleLogger) Debug(args ...interface{}) {
	l.logger.Print(args...)
}

//Debugf serverlog
func (l *ConsoleLogger) Debugf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

//Debugln serverlog
func (l *ConsoleLogger) Debugln(args ...interface{}) {
	l.logger.Println(args...)
}

//Info serverlog
func (l *ConsoleLogger) Info(args ...interface{}) {
	l.logger.Print(args...)
}

//Infof serverlog
func (l *ConsoleLogger) Infof(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

//Infoln serverlog
func (l *ConsoleLogger) Infoln(args ...interface{}) {
	l.logger.Println(args...)
}

//Warn serverlog
func (l *ConsoleLogger) Warn(args ...interface{}) {
	l.logger.Print(args...)
}

//Warnf serverlog
func (l *ConsoleLogger) Warnf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

//Warnln serverlog
func (l *ConsoleLogger) Warnln(args ...interface{}) {
	l.logger.Println(args...)
}

//Error serverlog
func (l *ConsoleLogger) Error(args ...interface{}) {
	l.logger.Print(args...)
}

//Errorf serverlog
func (l *ConsoleLogger) Errorf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

//Errorln serverlog
func (l *ConsoleLogger) Errorln(args ...interface{}) {
	l.logger.Println(args...)
}

func (l *ConsoleLogger) WithError(err error, args ...interface{}) {
	l.loggerus.WithError(err).Error(args...)
}

func (l *ConsoleLogger) WithErrorf(err error,format string, args ...interface{}) {
	l.loggerus.WithError(err).Errorf(format, args...)
}
