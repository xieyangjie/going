package logs

import (
	"fmt"
	"path"
	"runtime"
	"sync"
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

var messagePool *sync.Pool

////////////////////////////////////////////////////////////////////////////////
type ILogWriter interface {
	SetLevel(level int)
	GetLevel() int
	WriteMessage(level int, file string, line int, prefix string, msg string)

	Close()
	Flush()
}

////////////////////////////////////////////////////////////////////////////////
type logMessage struct {
	level          int
	file           string
	line           int
	levelShortName string
	message        string
}

////////////////////////////////////////////////////////////////////////////////
type Logger struct {
	level        int
	enableLogger bool
	enableStack  bool
	stackLevel   int
	outputs      map[string]ILogWriter
	messageChan  chan *logMessage
	signalChan   chan string
	waitGroup    sync.WaitGroup
}

func NewLogger() *Logger {
	return NewLoggerWithChannel(256)
}

func NewLoggerWithChannel(channelLen int64) *Logger {
	var log = &Logger{}
	log.outputs = make(map[string]ILogWriter)
	log.SetLogLevel(LOG_LEVEL_DEBUG)
	log.SetEnableLogger(true)
	log.SetEnableStack(false)
	log.SetStackLevel(LOG_LEVEL_PANIC)
	log.messageChan = make(chan *logMessage, channelLen)
	log.signalChan  = make(chan string, 1)
	log.waitGroup.Add(1)

	if messagePool == nil {
		messagePool = &sync.Pool{}
		messagePool.New = func() interface{} {
			return &logMessage{}
		}
	}

	go log.startLogger()
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

func (this *Logger) startLogger() {
	var end = false
	for {
		select {
		case msg := <- this.messageChan:
			this._writeMessage(msg.level, msg.file, msg.line, msg.levelShortName, msg.message)
			messagePool.Put(msg)
		case sc := <- this.signalChan:
			this.flush()
			if sc == "close" {
				for _, o := range this.outputs {
					o.Close()
				}
				this.outputs = nil
				end = true
			}
			this.waitGroup.Done()
		}

		if end {
			break
		}
	}
}

func (this *Logger) writeMessage(level int, format string, v ...interface{}) {
	if !this.enableLogger {
		return
	}

	var skip = 2
	if this == DefaultLogger {
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

	var msg = messagePool.Get().(*logMessage)
	msg.level = level
	msg.file = file
	msg.line = line
	msg.levelShortName = levelShortName
	msg.message = message
	this.messageChan <- msg
}

func (this *Logger) _writeMessage(level int, file string, line int, levelShortName, message string) {
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
	this.signalChan <- "close"
	this.waitGroup.Wait()
	close(this.messageChan)
	close(this.signalChan)
}

func (this *Logger) Flush() {
	this.signalChan <- "flush"
	this.waitGroup.Wait()
	this.waitGroup.Add(1)
}

func (this *Logger) flush() {
	for {
		if len(this.messageChan) > 0 {
			var mc = <- this.messageChan
			this._writeMessage(mc.level, mc.file, mc.line, mc.levelShortName, mc.message)
			messagePool.Put(mc)
			continue
		}
		break
	}
	for _, o := range this.outputs {
		o.Flush()
	}
}

////////////////////////////////////////////////////////////////////////////////
var DefaultLogger *Logger

func SharedLogger() *Logger {
	if DefaultLogger == nil {
		DefaultLogger = NewLogger()
		DefaultLogger.AddOutput("console", NewConsoleWriter(LOG_LEVEL_DEBUG))
	}
	return DefaultLogger
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
