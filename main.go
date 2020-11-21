package main

import (
	"bufio"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

/*Begin defining an array for the command line*/
type arrayFlags []string

func (i *arrayFlags) String() string {
	return "List of strings"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

/*End defining an array for the command line*/

type Config struct {
	SubNets             arrayFlags `yaml:"SubNets" toml:"subnets" env:"SUBNETS"`
	IgnorList           arrayFlags `yaml:"IgnorList" toml:"ignorlist" env:"IGNORLIST"`
	LogLevel            string     `yaml:"LogLevel" toml:"loglevel" env:"LOG_LEVEL"`
	ProcessingDirection string     `yaml:"ProcessingDirection" toml:"direct" env:"DIRECT" env-default:"both"`
	FlowAddr            string     `yaml:"FlowAddr" toml:"flowaddr" env:"FLOW_ADDR"`
	FlowPort            int        `yaml:"FlowPort" toml:"flowport" env:"FLOW_PORT" env-default:"2055"`
	NameFileToLog       string     `yaml:"FileToLog" toml:"log" env:"FLOW_LOG"`
}

var (
	cfg                Config
	SubNets, IgnorList arrayFlags
	writer             *bufio.Writer
	FileToLog          *os.File
	err                error
	configFilename     string = "/etc/go/config.toml"
)

func init() {
	flag.StringVar(&cfg.FlowAddr, "addr", "", "NetFlow/IPFIX listening address")
	flag.IntVar(&cfg.FlowPort, "port", 2055, "NetFlow/IPFIX listening port")
	flag.StringVar(&cfg.LogLevel, "loglevel", "info", "Log level")
	flag.Var(&cfg.SubNets, "subnet", "List of internal subnets")
	flag.Var(&cfg.IgnorList, "ignorlist", "List of ignored words/parameters per string")
	flag.StringVar(&cfg.ProcessingDirection, "direct", "both", "")
	flag.StringVar(&cfg.NameFileToLog, "log", "", "The file where logs will be written in the format of squid logs")
	flag.Parse()

	var config_source string
	if SubNets == nil && IgnorList == nil {
		// err := cleanenv.ReadConfig("goflow.toml", &cfg)
		err := cleanenv.ReadConfig(configFilename, &cfg)
		if err != nil {
			log.Warningf("No config file(%v) found: %v", configFilename, err)
		}
		lvl, err2 := log.ParseLevel(cfg.LogLevel)
		if err2 != nil {
			log.Errorf("Error in determining the level of logs (%v). Installed by default = Info", cfg.LogLevel)
			lvl, _ = log.ParseLevel("info")
		}
		log.SetLevel(lvl)
		config_source = "ENV/CFG"
	} else {
		config_source = "CLI"
	}
	log.Debugf("Config read from %s: IgnorList=(%v), SubNets=(%v), FlowAddr=(%v), FlowPort=(%v), LogLevel=(%v), ProcessingDirection=(%v)",
		config_source,
		cfg.IgnorList,
		cfg.SubNets,
		cfg.FlowAddr,
		cfg.FlowPort,
		cfg.LogLevel,
		cfg.ProcessingDirection)

}

// getExitSignalsChannel Intercepts program termination signals
func getExitSignalsChannel() chan os.Signal {

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		// syscall.SIGKILL, // "always fatal", "SIGKILL and SIGSTOP may not be caught by a program"
		syscall.SIGHUP, // "terminal is disconnected"
	)
	return c

}

func openOutputDevice(filename string) *bufio.Writer {
	if filename == "" {
		writer = bufio.NewWriter(os.Stdout)
		log.Debug("Output in os.Stdout")
		return writer

	} else {
		FileToLog, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Errorf("Error, the '%v' file could not be created (there are not enough premissions or it is busy with another program): %v", cfg.NameFileToLog, err)
			writer = bufio.NewWriter(os.Stdout)
			FileToLog.Close()
			log.Debug("Output in os.Stdout with error open file")
			return writer

		} else {
			defer FileToLog.Close()
			writer = bufio.NewWriter(FileToLog)
			log.Debugf("Output in file (%v)(%v)", filename, FileToLog)
			return writer

		}
	}

}

func main() {

	/*Creating a channel to intercept the program end signal*/
	exitChan := getExitSignalsChannel()

	go func() {
		<-exitChan
		// HERE Insert commands to be executed before the program terminates
		writer.Flush()
		FileToLog.Close()
		log.Println("Shutting down")
		os.Exit(0)

	}()

	writer = openOutputDevice(cfg.NameFileToLog)
}
