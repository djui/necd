#include "volume_darwin.h"

// setting system volume. Mutes if under threshhold.
// Courtesy: https://gist.github.com/atr000/205140
int setVolume(float newVolume) {
  if (newVolume < 0.0 || newVolume > 1.0) {
    fprintf(stderr, "ERROR: Requested volume out of range (%.2f)", newVolume);
    return 1;
  }

  // get output device
  UInt32 propertySize = 0;
  OSStatus status = noErr;
  AudioObjectPropertyAddress propertyAOPA;
  propertyAOPA.mElement = kAudioObjectPropertyElementMaster;
  propertyAOPA.mScope = kAudioDevicePropertyScopeOutput;

  if (newVolume < 0.001) {
    propertyAOPA.mSelector = kAudioDevicePropertyMute;
  } else {
    propertyAOPA.mSelector = kAudioHardwareServiceDeviceProperty_VirtualMasterVolume;
  }

  AudioDeviceID outputDeviceID = defaultOutputDeviceID();

  if (outputDeviceID == kAudioObjectUnknown) {
    fprintf(stderr, "ERROR: Unknown device");
    return 2;
  }

  if (!AudioHardwareServiceHasProperty(outputDeviceID, &propertyAOPA)) {
    fprintf(stderr, "ERROR: Device 0x%0x does not support volume control", outputDeviceID);
    return 3;
  }

  Boolean canSetVolume = false;
  status = AudioHardwareServiceIsPropertySettable(outputDeviceID, &propertyAOPA, &canSetVolume);

  if (status || !canSetVolume) {
    fprintf(stderr, "ERROR: Device 0x%0x does not support volume control", outputDeviceID);
    return 4;
  }

  if (propertyAOPA.mSelector == kAudioDevicePropertyMute) {
    propertySize = sizeof(UInt32);
    UInt32 mute = 1;
    status = AudioHardwareServiceSetPropertyData(outputDeviceID, &propertyAOPA, 0, NULL, propertySize, &mute);
  } else {
    propertySize = sizeof(Float32);
    status = AudioHardwareServiceSetPropertyData(outputDeviceID, &propertyAOPA, 0, NULL, propertySize, &newVolume);

    if (status) {
      fprintf(stderr, "ERROR: Unable to set volume for device 0x%0x", outputDeviceID);
      return 5;
    }

    // make sure we're not muted
    propertyAOPA.mSelector = kAudioDevicePropertyMute;
    propertySize = sizeof(UInt32);
    UInt32 mute = 0;

    if (!AudioHardwareServiceHasProperty(outputDeviceID, &propertyAOPA)) {
      fprintf(stderr, "ERROR: Device 0x%0x does not support muting", outputDeviceID);
      return 6;
    }

    Boolean canSetMute = false;

    status = AudioHardwareServiceIsPropertySettable(outputDeviceID, &propertyAOPA, &canSetMute);

    if (status || !canSetMute) {
      fprintf(stderr, "ERROR: Device 0x%0x does not support muting", outputDeviceID);
      return 7;
    }

    status = AudioHardwareServiceSetPropertyData(outputDeviceID, &propertyAOPA, 0, NULL, propertySize, &mute);
  }

  if (status) {
    fprintf(stderr, "ERROR: Unable to set volume for device 0x%0x", outputDeviceID);
    return 8;
  }

  return 0;
}

AudioDeviceID defaultOutputDeviceID() {
  AudioDeviceID outputDeviceID = kAudioObjectUnknown;

  // get output device device
  UInt32 propertySize = 0;
  OSStatus status = noErr;
  AudioObjectPropertyAddress propertyAOPA;
  propertyAOPA.mScope = kAudioObjectPropertyScopeGlobal;
  propertyAOPA.mElement = kAudioObjectPropertyElementMaster;
  propertyAOPA.mSelector = kAudioHardwarePropertyDefaultOutputDevice;

  if (!AudioHardwareServiceHasProperty(kAudioObjectSystemObject, &propertyAOPA)) {
    fprintf(stderr, "ERROR: Cannot find default output device!");
    return outputDeviceID;
  }

  propertySize = sizeof(AudioDeviceID);

  status = AudioHardwareServiceGetPropertyData(kAudioObjectSystemObject, &propertyAOPA, 0, NULL, &propertySize, &outputDeviceID);
  if (status) {
    fprintf(stderr, "ERROR: Cannot find default output device!");
  }
  
  return outputDeviceID;
}
