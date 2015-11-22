package main

import (
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/codegangsta/cli"
	"github.com/ian-kent/go-log/layout"
	"github.com/ian-kent/go-log/levels"
	"github.com/ian-kent/go-log/log"
)

var version string

func init() {
	logger := log.Logger()
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		logger.Appender().SetLayout(layout.Pattern("%p: %m"))
	} else {
		logger.Appender().SetLayout(layout.Pattern("%d %p: %m"))
	}
	logger.SetLevel(levels.DEBUG)
}

func main() {
	app := cli.NewApp()
	app.Name = "necd"
	app.Usage = "Network Environment Change Detector"
	app.Version = version
	app.Action = actionMain

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "daemon, d",
			EnvVar: "DAEMON",
			Usage:  "Daemonize",
		},
	}

	app.RunAndExitOnError()
}

func actionMain(c *cli.Context) {
	if c.Bool("daemon") {
		if err := Daemonize("necd"); err != nil {
			log.Fatalf("Failed to daemonize: %v", err)
		}
	} else {
		KeepAlive(runLoop)
	}
}

func runLoop() {
	conf, err := readConf()
	AssertNoErr(err, "Failed to read config file")
	log.Debug("Global conf: %#v", conf)

	nif := conf["config"]["if"]
	interval := atoi(conf["config"]["interval"], 5)

	for name := range notifyOnChange(nif, interval) {
		log.Info("Network changed: %s", name)

		if section, ok := conf["ssid:"+name]; ok {
			log.Debug("Found section: %v", section)
			ApplyCmds(section)
		} else {
			log.Debug("Undefined section for: %s", name)
		}
	}
}

func notifyOnChange(nif string, interval int) <-chan (string) {
	previousName := NetworkName(nif)
	c := make(chan string, 1)
	c <- previousName

	go func() {
		for _ = range time.Tick(time.Duration(interval) * time.Second) {
			currentName := NetworkName(nif)
			log.Debug("Current name: %s", currentName)
			if previousName != currentName {
				previousName = currentName
				c <- currentName
			}
		}
	}()

	return c
}

func atoi(a string, def int) int {
	if a == "" {
		return def
	}
	i, err := strconv.Atoi(a)
	if err != nil {
		log.Fatalf("Failed to parse value: %s", a)
	}
	return i
}
