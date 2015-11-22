package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/ian-kent/go-log/log"
)

// NetworkName returns the SSID of the wifi network.
func NetworkName(nif string) string {
	ns := new(networksetup)
	if !ns.getAirportPower(nif) {
		return ""
	}
	return ns.getAirportNetwork(nif)
}

const networksetupCmd = "/usr/sbin/networksetup"

type networksetup struct{}

func (n *networksetup) getAirportPower(nif string) bool {
	cmd := exec.Command(networksetupCmd, "-getairportpower", nif)
	out, err := cmd.Output()
	AssertNoErr(err, "Failed to obtain airport power status")
	s := strings.TrimSpace(string(out))
	log.Debug("getairportpower: %s", s)

	return s == fmt.Sprintf("Wi-Fi Power (%s): On", nif)
}

func (n *networksetup) getAirportNetwork(nif string) string {
	cmd := exec.Command(networksetupCmd, "-getairportnetwork", nif)
	out, err := cmd.Output()
	AssertNoErr(err, "Failed to obtain airport power status")
	s := strings.TrimSpace(string(out))
	log.Debug("getairportnetwork: %s", s)

	p := regexp.MustCompile(`Current Wi-Fi Network: (.+)`)
	if m := p.FindStringSubmatch(s); m != nil {
		return string(m[1])
	}
	return ""
}
