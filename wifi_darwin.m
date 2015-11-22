#include "wifi_darwin.h"

const char* guessWifiInterfaceName() {
  CWInterface* nif = CWWiFiClient.sharedWiFiClient.interface;
  return nif.interfaceName.UTF8String;
}

const char* getWifiSSID(const char* name) {
  NSString* nsName = [[NSString alloc] initWithUTF8String: name];
  CWInterface* nif = [[CWInterface alloc] initWithInterfaceName: nsName];
  return nif.ssid.UTF8String;
}

bool getWifiActive(const char* name) {
  NSString* nsName = [[NSString alloc] initWithUTF8String: name];
  CWInterface* nif = [[CWInterface alloc] initWithInterfaceName: nsName];
  return nif.serviceActive;
}

bool getWifiPowerOn(const char* name) {
  NSString* nsName = [[NSString alloc] initWithUTF8String: name];
  CWInterface* nif = [[CWInterface alloc] initWithInterfaceName: nsName];
  return nif.powerOn;
}
