package main

import (
	// #cgo CFLAGS: -x objective-c
	// #cgo LDFLAGS: -fobjc-arc -framework CoreWLAN
	// #include "wifi_darwin.h"
	"C"

	"github.com/ian-kent/go-log/log"
)

// NetworkName returns the SSID of the wifi network.
func NetworkName(nif string) string {
	if nif == "" {
		nif = C.GoString(C.guessWifiInterfaceName())
	}

	if nif == "" {
		log.Debug("Could not find Wi-Fi network interface")
		return ""
	}

	active := C.getWifiActive(C.CString(nif))
	powerOn := C.getWifiPowerOn(C.CString(nif))

	if !active {
		log.Debug("Wi-Fi network interface is not active")
		return ""
	}

	if !powerOn {
		log.Debug("Wi-Fi network interface is not powered on")
		return ""
	}

	ssid := C.GoString(C.getWifiSSID(C.CString(nif)))

	if ssid == "" {
		log.Error("Wi-Fi network interface ssid empty")
	}

	return ssid
}
