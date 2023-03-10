#ifndef KEYBOARD_H
#define KEYBOARD_H
#include <stdbool.h>
#include <windows.h>

int sendInputKeyCode(const int key, bool keyDown);
int sendInputKeyChar(const char key, bool keyDown);

#endif