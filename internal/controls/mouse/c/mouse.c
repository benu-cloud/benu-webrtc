#include "mouse.h"

int sendInputMove(const int dx, const int dy) {
  INPUT input = {0};
  input.type = INPUT_MOUSE;
  input.mi.dwFlags = MOUSEEVENTF_MOVE;
  input.mi.dx = dx;
  input.mi.dy = dy;
  UINT sent = SendInput(1, &input, sizeof(INPUT));
  if (!sent) return HRESULT_FROM_WIN32(GetLastError());
  return ERROR_SUCCESS;
}
int sendInputKey(const int key) {
  INPUT input = {0};
  input.type = INPUT_MOUSE;
  input.mi.dwFlags = key;
  UINT sent = SendInput(1, &input, sizeof(INPUT));
  if (!sent) return HRESULT_FROM_WIN32(GetLastError());
  return ERROR_SUCCESS;
}
int sendInputScroll(const int scrollDir, const int size) {
  INPUT input = {0};
  input.type = INPUT_MOUSE;
  input.mi.dwFlags = scrollDir;
  input.mi.mouseData = size;
  UINT sent = SendInput(1, &input, sizeof(INPUT));
  if (!sent) return HRESULT_FROM_WIN32(GetLastError());
  return ERROR_SUCCESS;
}