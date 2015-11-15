package main

import (
	"time"

	"github.com/codegangsta/cli"
	"github.com/ian-kent/go-log/layout"
	"github.com/ian-kent/go-log/levels"
	"github.com/ian-kent/go-log/log"
)

var version string

func init() {
	//log.SetFlags(0)
	logger := log.Logger()
	logger.Appender().SetLayout(layout.Pattern("%p: %m"))
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

	log.Info("Global conf: %#v", conf)

	nif := conf["config"]["if"]

	for name := range notifyOnChange(nif) {
		log.Info("Network changed: %s", name)

		if section, ok := conf[name]; ok {
			log.Debug("Found section: %v", section)
			ApplyCmds(section)
		}
	}
}

func notifyOnChange(nif string) <-chan (string) {
	previousName := NetworkName(nif)
	c := make(chan string, 1)
	c <- previousName

	go func() {
		for _ = range time.Tick(5 * time.Second) {
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
