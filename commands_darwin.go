package main

import (
	// #cgo LDFLAGS: -framework IOKit -framework ApplicationServices
	// #cgo CFLAGS: -Wno-deprecated
	// #include "brightness_darwin.h"
	//
	// #cgo LDFLAGS: -framework AudioToolbox
	// #cgo CFLAGS: -Wno-deprecated -Wno-nonnull
	// #include "volume_darwin.h"
	"C"
	"math"

	"github.com/ian-kent/go-log/log"
)

// SetVolume takes a float value between [0,1] and sets the global output volume.
func SetVolume(vNorm float64) {
	v := C.float(constrain(vNorm, 0, 1))
	log.Debug("Setting volume to %v (%v)", vNorm, v)

	if res := C.setVolume(v); int(res) != 0 {
		log.Fatalf("Failed to set volume: %d", res)
	}
}

// SetBrightness takes a float value betweem [0,1] and sets the global screen brightness.
func SetBrightness(vNorm float64) {
	v := C.float(constrain(vNorm, 0, 1))
	log.Debug("Setting brightness to %v (%v)", vNorm, v)

	if res := C.setBrightness(v); int(res) != 0 {
		log.Fatalf("Failed to set brightness: %d", res)
	}
}

func constrain(v float64, min float64, max float64) float64 {
	return math.Min(max, math.Max(min, v))
}
