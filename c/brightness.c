/* gcc -std=c99 -o brightness brightness.c -framework IOKit -framework ApplicationServices */

#include <stdio.h>
#include <unistd.h>
#include <IOKit/graphics/IOGraphicsLib.h>
#include <ApplicationServices/ApplicationServices.h>

const int kMaxDisplays = 16;
const CFStringRef kDisplayBrightness = CFSTR(kIODisplayBrightnessKey);

void errexit(const char *fmt, ...) {
  va_list ap;
  va_start(ap, fmt);
  fprintf(stderr, "brightness: ");
  vfprintf(stderr, fmt, ap);
  fprintf(stderr, "\n");
  exit(1);
}

void usage() {
  fprintf(stderr, "usage: brightness [-m|-d display] [-v] <brightness>\n");
  exit(1);
}

int main(int argc, char *const argv[]) {
  if (argc == 1)
    usage();

  unsigned long displayToSet = 0;
  enum { ACTION_SET_ALL, ACTION_SET_ONE } action = ACTION_SET_ALL;
  extern char *optarg;
  extern int optind;
  int ch;

  while ( (ch = getopt(argc, argv, "md:")) != -1) {
    switch (ch) {
    case 'm':
      if (action != ACTION_SET_ALL)
        usage();
      action = ACTION_SET_ONE;
      displayToSet = (unsigned long)CGMainDisplayID();
      break;
    case 'd':
      if (action != ACTION_SET_ALL)
        usage();
      action = ACTION_SET_ONE;
      errno = 0;
      displayToSet = strtoul(optarg, NULL, 0);
      if (errno == EINVAL || errno == ERANGE)
        errexit("display must be an integer index (0) or a hexadecimal ID (0x4270a80)");
      break;
    default:
      usage();
    }
  }
 
  argc -= optind;
  argv += optind;

  if (argc != 1)
    usage();

  errno = 0;
  float brightness = strtof(argv[0], NULL);
  if (errno == ERANGE)
    usage();
  if (brightness < 0 || brightness > 1)
    errexit("brightness must be between 0 and 1");

  CGDirectDisplayID display[kMaxDisplays];
  CGDisplayCount numDisplays;
  CGDisplayErr err;
  err = CGGetActiveDisplayList(kMaxDisplays, display, &numDisplays);
  if (err != CGDisplayNoErr)
    errexit("cannot get list of displays (error %d)", err);

  for (CGDisplayCount i = 0; i < numDisplays; ++i) {
    CGDirectDisplayID displayID = display[i];
    #pragma GCC diagnostic push
    #pragma GCC diagnostic ignored "-Wdeprecated-declarations"
    CFDictionaryRef originalMode = CGDisplayCurrentMode(displayID);
    if (originalMode == NULL)
      continue;

    io_service_t service = CGDisplayIOServicePort(displayID);
    #pragma GCC diagnostic pop
    switch (action) {
      case ACTION_SET_ONE:
        if ((CGDirectDisplayID)displayToSet != displayID && displayToSet != i)
          continue;
      case ACTION_SET_ALL:
        err = IODisplaySetFloatParameter(service, kNilOptions, kDisplayBrightness, brightness);
        if (err != kIOReturnSuccess) {
          errexit("failed to set brightness of display 0x%x (error %d)", (unsigned int)displayID, err);
          continue;
        }
    }
  }

  return 0;
}
