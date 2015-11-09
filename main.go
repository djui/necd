package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/ian-kent/go-log/layout"
	"github.com/ian-kent/go-log/log"
)

func init() {
	//log.SetFlags()
	log.Logger().Appender().SetLayout(layout.Pattern("%p: %m"))
}

func main() {
	KeepAlive(runLoop)
}

func runLoop() {
	_, err := readConf()
	AssertNoErr(err, "Failed to read config file")

	for name := range notifyOnChange() {
		log.Info("WiFi network changed: %s", name)
	}
}

func notifyOnChange() <-chan (string) {
	previousName := currentWifiName()
	c := make(chan string, 1)
	c <- previousName

	go func() {
		for _ = range time.Tick(5 * time.Second) {
			currentName := currentWifiName()
			if previousName != currentName {
				previousName = currentName
				c <- currentName
			}
		}
	}()

	return c
}

// wifiName returns the SSID of the wifi network.
//
// Many other options exist, like `airport` or
// `networksetup -getinfo Wi-Fi`.
//
// Things to improve are:

//   - Run `networksetup -listallhardwareports` or
//     `networksetup -listnetworkserviceorder` first to ensure we got
//     the correct network interface name.
//   - Check if the network interface is active, otherwise try to find
//     another one.
func currentWifiName() string {
	//cmd := []string{"networksetup", "-getairportnetwork", "en0"}
	//out, err := exec.Command(cmd...).Output()
	out, err := exec.Command("networksetup", "-getairportnetwork", "en0").Output()
	AssertNoErr(err, "Failed to obtain airport network")

	parts := strings.Fields(string(out))
	Assert(len(parts) >= 4, fmt.Sprintf("Failed to parse WiFi name: %v", parts))

	name := strings.Join(parts[3:], " ")
	log.Debug(name)
	return name
}
