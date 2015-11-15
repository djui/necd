package main

import (
	"strconv"

	"github.com/ian-kent/go-log/log"
)

// ApplyCmds applies a map of known commands and their values.
func ApplyCmds(cmds map[string]string) {
	for cmd, val := range cmds {
		switch cmd {
		case "volume":
			v, err := strconv.ParseFloat(val, 32)
			if err != nil {
				log.Warn("Config parse error: %v", err)
				continue
			}
			SetVolume(v)
		case "brightness":
			v, err := strconv.ParseFloat(val, 32)
			if err != nil {
				log.Warn("Config parse error: %v", err)
				continue
			}
			SetBrightness(v)
		}
	}
}
