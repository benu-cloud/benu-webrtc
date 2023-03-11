#include "keyboard.h"

int sendInputKeyCode(const int key, bool keyDown) {
  INPUT input = {0};
  input.type = INPUT_KEYBOARD;
  input.ki.wVk = key;
  if (!keyDown) input.ki.dwFlags = KEYEVENTF_KEYUP;
  UINT sent = SendInput(1, &input, sizeof(INPUT));
  if (!sent) return HRESULT_FROM_WIN32(GetLastError());
  return ERROR_SUCCESS;
}
int sendInputKeyChar(const char key, bool keyDown) {
  int code;
  code = VkKeyScan(key);
  if (code == 0xFFFF) {
    return ERROR_BAD_ARGUMENTS;
  }
  return sendInputKeyCode(code, keyDown);
}