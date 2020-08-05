package main

import (
	"flag"
	"log"
)

type configType struct {
	logLevel int
}

var (
	config configType
)

func init() {

	flag.IntVar(&config.logLevel, "log-level", 1, `Level log: 
		0 - silent, 
		1 - error, start and end, 
		2 - error, start and end, warning, 
		3 - error, start and end, warning, info,
		4 - error, start and end, warning, info and more info`)
	flag.Parse()
	if config.logLevel > 4 {
		config.logLevel = 4
	}
	if config.logLevel != 0 {
		log.SetFlags(log.Ldate | log.Ltime)
	} // If logLevel not specified - silent mode
}

func main() {
}
