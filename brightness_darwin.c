#include "brightness_darwin.h"

const int kMaxDisplays = 16;
const CFStringRef kDisplayBrightness = CFSTR(kIODisplayBrightnessKey);

// Courtesy: https://github.com/nriley/brightness/blob/master/brightness.c
int setBrightness(float v) {
  CGDirectDisplayID display[kMaxDisplays];
  CGDisplayCount numDisplays;
  CGDisplayErr err;

  err = CGGetActiveDisplayList(kMaxDisplays, display, &numDisplays);
  if (err != CGDisplayNoErr) {
    fprintf(stderr, "ERROR: Cannot get list of displays (error %d)\n", err);
    return 1;
  }

  for (CGDisplayCount i = 0; i < numDisplays; ++i) {
    CGDirectDisplayID dspy = display[i];

    CFDictionaryRef originalMode = CGDisplayCurrentMode(dspy);
    if (originalMode == NULL)
      continue;

    io_service_t service = CGDisplayIOServicePort(dspy);

    err = IODisplaySetFloatParameter(service, kNilOptions, kDisplayBrightness, v);
    if (err != kIOReturnSuccess) {
      fprintf(stderr, "ERROR: Failed to set brightness of display 0x%x (error %d)\n", (unsigned int)dspy, err);
      return 2;
    }
  }

  return 0;
}
