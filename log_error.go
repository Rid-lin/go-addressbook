package main

import (
	"log"
)

func checkErrorFatal(lvIn, logLevel int, message string, err error) {
	if err != nil {
		toLog(lvIn, logLevel, "%v %v", message, err)
	}
}

func toLog(level int, lvIn int, format string, v ...interface{}) {
	// level - level logging
	// 0 - silent
	// 1 - only panic
	// 2 - panic, warning
	// 3 - panic, warning, some info
	// 4 - debug info
	// lvIn - level log for this event
	if lvIn <= level {
		log.Printf(format, v...)
	}
}
