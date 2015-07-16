package main

import (
	"github.com/smartwalle/going/logs"
)


func main() {
	logs.Debug("Debug message")
	logs.Info("Info message")
	logs.Warn("Warn message")
	logs.Panic("Panic message")
	logs.Fatal("Fatal message")
}

