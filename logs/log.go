package logs
import (
	"fmt"
	"runtime"
	"path"
)

const (
	LevelDebug 	= iota //= "Debug"
	LevelInfo 		//= "Info"
	LevelWarn 		//= "Warn"
	LevelPanic		//= "Panic"
	LevelFatal		//= "Fatal"
)

var levelShortNames = []string{
	"[D]",
	"[I]",
	"[W]",
	"[P]",
	"[F]",
}

type ILog interface {
	WriteMessage(level int, file string, line int, prefix string, msg string)
}


type Log struct {
	level int
	handlers map[string]ILog
}

func NewLog() *Log {
	var log = &Log{}
	log.handlers = make(map[string]ILog)
	return log
}

func (this *Log) SetLogLevel(level int) {
	this.level = level
}

func (this *Log) GetLogLevel() int {
	return this.level
}

func (this *Log) SetOutput(key string, out ILog) {
	this.handlers[key] = out
}

func (this *Log) GetOutput(key string) ILog {
	return this.handlers[key]
}

func (this *Log) RemoveOutput(key string) {
	delete(this.handlers, key)
}

func (this *Log)writeMessage(level int, format string, v ...interface{}) {
	var _, file, line, ok = runtime.Caller(2)
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

func (this *Log)Debug(format string, v ...interface{}) {
	this.writeMessage(LevelDebug, format, v...)
}

func (this *Log)Print(format string, v ...interface{}) {
	this.Debug(format, v...)
}

//func Info(format string, v ...interface{}) {
//	var levelShortNames
//}

func Print(v ...interface{}) {
//	Std.Output("", Linfo, 2, fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) {
}

func Debug(v ...interface{}) {
}

func Infof(format string, v ...interface{}) {
}

func Info(v ...interface{}) {
}

func Warnf(format string, v ...interface{}) {
}

func Warn(v ...interface{}) {
}

func Panicf(format string, v ...interface{}) {
}

func Panic(v ...interface{}) {
}

func Fatalf(format string, v ...interface{}) {
}

func Fatal(v ...interface{}) {
}