package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rid-lin/go-template/config"
)

type configType struct {
	logLevel int
}

var (
	conf configType
)

func init() {

	flag.IntVar(&conf.logLevel, "log-level", 1, `Level log: 
		0 - silent, 
		1 - error, start and end, 
		2 - error, start and end, warning, 
		3 - error, start and end, warning, info,
		4 - error, start and end, warning, info and more info`)
	flag.Parse()
	if conf.logLevel > 4 {
		conf.logLevel = 4
	}
	if conf.logLevel != 0 {
		log.SetFlags(log.Ldate | log.Ltime)
	} // If logLevel not specified - silent mode
	if len(os.Args) == 1 {
		// loads values from .env into the system
		if err := godotenv.Load(); err != nil {
			log.Print("No .env file found")
		}
	}
}

func main() {
	conf := config.New()

	// Print out environment variables
	fmt.Println(conf.GitHub.Username)
	fmt.Println(conf.GitHub.APIKey)
	fmt.Println(conf.DebugMode)
	fmt.Println(conf.MaxUsers)

	// Print out each role
	for _, role := range conf.UserRoles {
		fmt.Println(role)
	}
}
