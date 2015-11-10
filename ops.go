package main

import (
	"fmt"
	"os/exec"
)

// SetVolume takes a float value between [0,1] and sets the global output volume.
func SetVolume(vNorm float32) {
	v := 7 * vNorm
	cmd := []string{"/usr/bin/osascript", "-e", fmt.Sprintf("set volume %f", v)}
	out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	AssertNoErr(err, fmt.Sprintf("Failed execute cmd: %v: %s", cmd, string(out)))
}

// SetBrightness takes a float value betweem [0,1] and sets the global screen brightness.
func SetBrightness(vNorm float32) {
	cmd := []string{"brightness", fmt.Sprintf("%f", vNorm)}
	out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	AssertNoErr(err, fmt.Sprintf("Failed execute cmd: %v: %s", cmd, string(out)))
	return

	v := int(16 * vNorm)

	for i := 1; i <= 16; i++ {
		cmd := []string{"/usr/bin/osascript", "-e", "tell application \"System Events\" to key code 107"}
		out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
		AssertNoErr(err, fmt.Sprintf("Failed execute cmd: %v: %s", cmd, string(out)))
	}

	for i := 1; i <= v; i++ {
		cmd := []string{"/usr/bin/osascript", "-e", "tell application \"System Events\" to key code 113"}
		out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
		AssertNoErr(err, fmt.Sprintf("Failed execute cmd: %v: %s", cmd, string(out)))
	}
}
