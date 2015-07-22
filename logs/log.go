package logs

import (
	"fmt"
	"path"
	"runtime"
)

const (
	LOG_LEVEL_DEBUG = iota //= "Debug"
	LOG_LEVEL_INFO         //= "Info"
	LOG_LEVEL_WARN         //= "Warn"
	LOG_LEVEL_PANIC        //= "Panic"
	LOG_LEVEL_FATAL        //= "Fatal"
)

var levelShortNames = []string{
	"[D]",
	"[I]",
	"[W]",
	"[P]",
	"[F]",
}

////////////////////////////////////////////////////////////////////////////////
type ILogWriter interface {
	SetLevel(level int)
	GetLevel() int
	WriteMessage(level int, file string, line int, prefix string, msg string)

	Close()
	Flush()
}

////////////////////////////////////////////////////////////////////////////////
type Logger struct {
	level        int
	enableLogger bool
	enableStack  bool
	stackLevel   int
	outputs      map[string]ILogWriter
}

func NewLogger() *Logger {
	var log = &Logger{}
	log.outputs = make(map[string]ILogWriter)
	log.SetLogLevel(LOG_LEVEL_DEBUG)
	log.SetEnableLogger(true)
	log.SetEnableStack(false)
	log.SetStackLevel(LOG_LEVEL_PANIC)
	return log
}

func (this *Logger) SetLogLevel(level int) {
	this.level = level
}

func (this *Logger) GetLogLevel() int {
	return this.level
}

func (this *Logger) SetEnableLogger(e bool) {
	this.enableLogger = e
}

func (this *Logger) GetEnableLogger() bool {
	return this.enableLogger
}

func (this *Logger) SetEnableStack(e bool) {
	this.enableStack = e
}

func (this *Logger) GetEnableStack() bool {
	return this.enableStack
}

func (this *Logger) SetStackLevel(level int) {
	this.stackLevel = level
}

func (this *Logger) GetStackLevel() int {
	return this.stackLevel
}

func (this *Logger) AddOutput(key string, out ILogWriter) {
	this.outputs[key] = out
}

func (this *Logger) GetOutput(key string) ILogWriter {
	return this.outputs[key]
}

func (this *Logger) RemoveOutput(key string) {
	delete(this.outputs, key)
}

func (this *Logger) writeMessage(level int, format string, v ...interface{}) {
	if !this.enableLogger {
		return
	}

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

	if this.enableStack && level >= this.stackLevel {
		message += "\n"
		buf := make([]byte, 1024*1024)
		n := runtime.Stack(buf, true)
		message += string(buf[:n])
		message += "\n"
	}

	for _, output := range this.outputs {
		output.WriteMessage(level, file, line, levelShortName, message)
	}
}

//debug
func (this *Logger) Debugf(format string, args ...interface{}) {
	this.writeMessage(LOG_LEVEL_DEBUG, format, args...)
}

func (this *Logger) Debugln(args ...interface{}) {
	this.writeMessage(LOG_LEVEL_DEBUG, this.println(args...))
}

//print
func (this *Logger) Printf(format string, args ...interface{}) {
	this.writeMessage(LOG_LEVEL_DEBUG, format, args...)
}

func (this *Logger) Println(args ...interface{}) {
	this.writeMessage(LOG_LEVEL_DEBUG, this.println(args...))
}

//info
func (this *Logger) Infof(format string, args ...interface{}) {
	this.writeMessage(LOG_LEVEL_INFO, format, args...)
}

func (this *Logger) Infoln(args ...interface{}) {
	this.writeMessage(LOG_LEVEL_INFO, this.println(args...))
}

//warn
func (this *Logger) Warnf(format string, args ...interface{}) {
	this.writeMessage(LOG_LEVEL_WARN, format, args...)
}

func (this *Logger) Warnln(args ...interface{}) {
	this.writeMessage(LOG_LEVEL_WARN, this.println(args...))
}

//panic
func (this *Logger) Panicf(format string, args ...interface{}) {
	this.writeMessage(LOG_LEVEL_PANIC, format, args...)
}

func (this *Logger) Panicln(args ...interface{}) {
	this.writeMessage(LOG_LEVEL_PANIC, this.println(args...))
}

//fatal
func (this *Logger) Fatalf(format string, args ...interface{}) {
	this.writeMessage(LOG_LEVEL_FATAL, format, args...)
}

func (this *Logger) Fatalln(args ...interface{}) {
	this.writeMessage(LOG_LEVEL_FATAL, this.println(args...))
}

func (this *Logger) println(args ...interface{}) string {
	var msg = fmt.Sprintln(args...)
	return msg
}

func (this *Logger) Close() {
	for _, output := range this.outputs {
		output.Close()
	}
}

func (this *Logger) Flush() {
	for _, output := range this.outputs {
		output.Flush()
	}
}

////////////////////////////////////////////////////////////////////////////////
var sharedLogger *Logger

func SharedLogger() *Logger {
	if sharedLogger == nil {
		sharedLogger = NewLogger()
		sharedLogger.AddOutput("console", NewConsoleWriter(LOG_LEVEL_DEBUG))
	}
	return sharedLogger
}

func Debugf(format string, args ...interface{}) {
	SharedLogger().Debugf(format, args...)
}

func Debugln(args ...interface{}) {
	SharedLogger().Debugln(args...)
}

func Printf(format string, args ...interface{}) {
	SharedLogger().Printf(format, args...)
}

func Println(args ...interface{}) {
	SharedLogger().Println(args...)
}

func Infof(format string, args ...interface{}) {
	SharedLogger().Infof(format, args...)
}

func Infoln(args ...interface{}) {
	SharedLogger().Infoln(args...)
}

func Warnf(format string, args ...interface{}) {
	SharedLogger().Warnf(format, args...)
}

func Warnln(args ...interface{}) {
	SharedLogger().Warnln(args...)
}

func Panicf(format string, args ...interface{}) {
	SharedLogger().Panicf(format, args...)
}

func Panicln(args ...interface{}) {
	SharedLogger().Panicln(args...)
}

func Fatalf(format string, args ...interface{}) {
	SharedLogger().Fatalf(format, args...)
}

func Fatalln(args ...interface{}) {
	SharedLogger().Fatalln(args...)
}
