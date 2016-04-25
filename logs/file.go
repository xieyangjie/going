package logs

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"
	"path/filepath"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
type fileWriter struct {
	sync.Mutex
	file *os.File
}

func newFileWriter() *fileWriter {
	var writer = &fileWriter{}
	return writer
}

func (this *fileWriter) setFile(file *os.File) {
	if this.file != nil {
		this.file.Close()
	}
	this.file = file
}

func (this *fileWriter) Write(b []byte) (int, error) {
	this.Lock()
	defer this.Unlock()
	return this.file.Write(b)
}

////////////////////////////////////////////////////////////////////////////////
type FileWriter struct {
	logger  *log.Logger
	writer  *fileWriter
	level   int
	path    string
	maxSize int64
	maxDays int64
	lock    *sync.Mutex
}

func NewFileWriter(level int, path string) *FileWriter {
	var file = &FileWriter{}
	file.maxDays = 30
	file.maxSize = 10 * 1024 * 1024 //10m
	file.level = level
	file.path = path
	file.lock = &sync.Mutex{}
	file.writer = newFileWriter()
	file.logger = log.New(file.writer, "", log.Ldate|log.Ltime)

	file.init()
	return file
}

func (this *FileWriter) init() {
	//首先创建目录
	if _, err := os.Stat(this.path); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(this.path, os.ModeDir|os.ModePerm)
		}
	}
	this.startLogger()
}

func (this *FileWriter) startLogger() {
	var filename = path.Join(this.path, "temp_logs.log")
	var file, _ = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	this.writer.setFile(file)
}

func (this *FileWriter) checkSize() {
	this.lock.Lock()
	defer this.lock.Unlock()
	var fileInfo, err = this.writer.file.Stat()

	if err == nil {
		var size = fileInfo.Size()
		if size >= this.maxSize {
			this.renameFile()

			go this.removeFile()
		}
	}
}

func (this *FileWriter) renameFile() {
	this.writer.Lock()
	defer this.writer.Unlock()

	this.writer.file.Close()
	var filename = path.Join(this.path, "temp_logs.log")
	var now = time.Now()
	var newName = path.Join(this.path, fmt.Sprintf("%s_%.9d.log", now.Format("2006_01_02_15_04_05"), now.Nanosecond()))
	os.Rename(filename, newName)

	this.startLogger()
}

func (this *FileWriter) removeFile() {
	var dir = filepath.Dir(this.path)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) (rErr error) {
		defer func() {
			if r := recover(); r != nil {
			}
		}()

		if !info.IsDir() && info.ModTime().Unix() < (time.Now().Unix()-60*60*24*this.maxDays) {
			if strings.HasSuffix(path, ".log") {
				rErr = os.Remove(path)
			}
		}
		return rErr
	})
}

func (this *FileWriter) SetLevel(level int) {
	this.level = level
}

func (this *FileWriter) GetLevel() int {
	return this.level
}

func (this *FileWriter) SetMaxSize(size int64) {
	this.maxSize = size
}

func (this *FileWriter) GetMaxSize() int64 {
	return this.maxSize
}

func (this *FileWriter) SetMaxDays(days int64) {
	this.maxDays = days
}

func (this *FileWriter) GetMaxDays() int64 {
	return this.maxDays
}

func (this *FileWriter) WriteMessage(level int, file string, line int, prefix string, msg string) {
	if level < this.level {
		return
	}

	this.logger.Printf("%s [%s:%d] %s", prefix, file, line, msg)
	this.checkSize()
}

func (this *FileWriter) Close() {
	this.Flush()
	this.writer.file.Close()
}

func (this *FileWriter) Flush() {
	this.writer.file.Sync()
}