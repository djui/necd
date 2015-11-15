package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// CurrentWifiName returns the SSID of the wifi network.
//
// Many other options exist, like `airport` or
// `networksetup -getinfo Wi-Fi`.
//
// Things to improve are:
//
//   - Run `networksetup -listallhardwareports` or
//     `networksetup -listnetworkserviceorder` first to ensure we got
//     the correct network interface name.
//   - Check if the network interface is active, otherwise try to find
//     another one.
func CurrentWifiName(nif string) string {
	if !nifExists(nif) {
		return ""
	}

	if nifStatus(nif) != "active" {
		return ""
	}

	name := wifiName(nif)
	return name
}

func nifExists(nif string) bool {
	cmd := []string{"/sbin/ifconfig", "-l"}
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	AssertNoErr(err, "Failed to list network interfaces")
	nifs := strings.Fields(string(out))
	return Contains(nif, nifs)
}

func nifStatus(nif string) string {
	cmd := []string{"/sbin/ifconfig", nif}
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	AssertNoErr(err, "Failed to obtain network interface status")
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		kv := strings.Fields(line)
		if len(kv) >= 2 && kv[0] == "status:" {
			name := strings.Join(kv[1:], " ")
			return name
		}
	}
	return ""
}

func wifiName(nif string) string {
	cmd := []string{"/usr/sbin/networksetup", "-getairportnetwork", nif}
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	AssertNoErr(err, "Failed to obtain airport network")
	parts := strings.Fields(string(out))
	Assert(len(parts) >= 4, fmt.Sprintf("Failed to parse WiFi name: %v", parts))
	name := strings.Join(parts[3:], " ")
	return name
}
