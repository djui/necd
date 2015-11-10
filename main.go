package main

import (
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
	KeepAlive(runLoop)
}

func runLoop() {
	_, err := readConf()
	AssertNoErr(err, "Failed to read config file")

	nif := "en0"

	for name := range notifyOnChange(nif) {
		log.Info("Network changed: %s", name)

		// TODO: Use config for these configuration settings
		switch name {
		case "Wall-E":
			SetVolume(0.5)
			SetBrightness(0.5)

		case "ScraperWiki":
			SetVolume(0.1)
			SetBrightness(0)
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
			log.Debug(currentName)
			if previousName != currentName {
				previousName = currentName
				c <- currentName
			}
		}
	}()

	return c
}
