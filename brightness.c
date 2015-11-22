#include <stdio.h>
#include <unistd.h>
#include <IOKit/graphics/IOGraphicsLib.h>
#include <ApplicationServices/ApplicationServices.h>

const int kMaxDisplays = 16;
const CFStringRef kDisplayBrightness = CFSTR(kIODisplayBrightnessKey);

void set_brightness(float v) {
  CGDirectDisplayID display[kMaxDisplays];
  CGDisplayCount numDisplays;
  CGDisplayErr err;

  err = CGGetActiveDisplayList(kMaxDisplays, display, &numDisplays);
  if (err != CGDisplayNoErr)
    printf("cannot get list of displays (error %d)\n",err);
  
  for (CGDisplayCount i = 0; i < numDisplays; ++i) {
    CGDirectDisplayID dspy = display[i];

#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wdeprecated-declarations"  
    CFDictionaryRef originalMode = CGDisplayCurrentMode(dspy);
    if (originalMode == NULL)
      continue;
    io_service_t service = CGDisplayIOServicePort(dspy);
#pragma GCC diagnostic pop

    float brightness;
    err= IODisplayGetFloatParameter(service, kNilOptions, kDisplayBrightness, &brightness);
    if (err != kIOReturnSuccess) {
      fprintf(stderr, "failed to get brightness of display 0x%x (error %d)", (unsigned int)dspy, err);
      continue;
    }

    err = IODisplaySetFloatParameter(service, kNilOptions, kDisplayBrightness, v);
    if (err != kIOReturnSuccess) {
      fprintf(stderr, "Failed to set brightness of display 0x%x (error %d)", (unsigned int)dspy, err);
      continue;
    }
  }
}
