package main

import (
	"github.com/smartwalle/going/logs"
)


func main() {

	var log = logs.NewLog()

	log.SetOutput("1", logs.NewConsole())

	log.Debug("ssss ")
}


