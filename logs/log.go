package logs

import (
	"fmt"
	"runtime"
	"path"
)

const (
	LOG_LEVEL_DEBUG = iota	//= "Debug"
	LOG_LEVEL_INFO        	//= "Info"
	LOG_LEVEL_WARN         	//= "Warn"
	LOG_LEVEL_PANIC        	//= "Panic"
	LOG_LEVEL_FATAL        	//= "Fatal"
)

var levelShortNames = []string{
	"[D]",
	"[I]",
	"[W]",
	"[P]",
	"[F]",
}

var sharedLogger *Logger

type ILog interface {
	SetLevel(level int)
	GetLevel() int
	WriteMessage(level int, file string, line int, prefix string, msg string)
}

type Logger struct {
	level int
	handlers map[string]ILog
}

func NewLogger() *Logger {
	var log = &Logger{}
	log.handlers = make(map[string]ILog)
	return log
}

func SharedLogger() *Logger {
	if sharedLogger == nil {
		sharedLogger = NewLogger()
		sharedLogger.AddOutput("console", NewConsole(LOG_LEVEL_DEBUG))
	}
	return sharedLogger
}

func (this *Logger) SetLogLevel(level int) {
	this.level = level
}

func (this *Logger) GetLogLevel() int {
	return this.level
}

func (this *Logger) AddOutput(key string, out ILog) {
	this.handlers[key] = out
}

func (this *Logger) GetOutput(key string) ILog {
	return this.handlers[key]
}

func (this *Logger) RemoveOutput(key string) {
	delete(this.handlers, key)
}

func (this *Logger)writeMessage(level int, format string, v ...interface{}) {
	var skip = 2
	if this == sharedLogger {
		skip = 3
	}

	var _, file, line, ok = runtime.Caller(skip)
	if !ok {
		file = "???"
		line = -1
	} else {
		_, file = path.Split(file)
	}

	var levelShortName = levelShortNames[level]
	var message = fmt.Sprintf(format, v...)

	for _, value := range this.handlers {
		value.WriteMessage(level, file, line, levelShortName, message)
	}
}

func (this *Logger)Debug(format string, v ...interface{}) {
	this.writeMessage(LOG_LEVEL_DEBUG, format, v...)
}

func (this *Logger)Print(format string, v ...interface{}) {
	this.Debug(format, v...)
}

func (this *Logger)Info(format string, v ...interface{}) {
	this.writeMessage(LOG_LEVEL_INFO, format, v...)
}

func (this *Logger)Warn(format string, v ...interface{}) {
	this.writeMessage(LOG_LEVEL_WARN, format, v...)
}

func (this *Logger)Panic(format string, v ...interface{}) {
	this.writeMessage(LOG_LEVEL_PANIC, format, v...)
}

func (this *Logger)Fatal(format string, v ...interface{}) {
	this.writeMessage(LOG_LEVEL_FATAL, format, v...)
}

func Debug(format string, v ...interface{}) {
	SharedLogger().Debug(format, v...)
}

func Print(format string, v ...interface{}) {
	SharedLogger().Print(format, v...)
}

func Info(format string, v ...interface{}) {
	SharedLogger().Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	SharedLogger().Warn(format, v...)
}

func Panic(format string, v ...interface{}) {
	SharedLogger().Panic(format, v...)
}

func Fatal(format string, v ...interface{}) {
	SharedLogger().Fatal(format, v...)
}