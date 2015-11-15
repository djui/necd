package main

import (
	"os"
	"time"

	"github.com/ian-kent/go-log/layout"
	"github.com/ian-kent/go-log/levels"
	"github.com/ian-kent/go-log/log"
)

func init() {
	//log.SetFlags(0)
	logger := log.Logger()
	logger.Appender().SetLayout(layout.Pattern("%p: %m"))
	logger.SetLevel(levels.DEBUG)
}

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "-d" {
		if err := Daemonize("necd"); err != nil {
			log.Fatal("Failed to daemonize: %v", err)
		}
	} else {
		KeepAlive(runLoop)
	}
}

func runLoop() {
	conf, err := readConf()
	AssertNoErr(err, "Failed to read config file")

	log.Info("Global conf: %#v", conf)

	nif := "en0"

	for name := range notifyOnChange(nif) {
		log.Info("Network changed: %s", name)

		if section, ok := conf[name]; ok {
			log.Debug("Found section: %v", section)
			ApplyCmds(section)
		}
	}
}

func notifyOnChange(nif string) <-chan (string) {
	previousName := CurrentWifiName(nif)
	c := make(chan string, 1)
	c <- previousName

	go func() {
		for _ = range time.Tick(5 * time.Second) {
			currentName := CurrentWifiName(nif)
			log.Debug("Current name: %s", currentName)
			if previousName != currentName {
				previousName = currentName
				c <- currentName
			}
		}
	}()

	return c
}
