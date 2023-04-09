#ifndef MOUSE_H
#define MOUSE_H
#include <windows.h>

int sendInputMove(const int dx, const int dy);
int sendInputKey(const int key);
int sendInputScroll(const int scrollDir, const int size);
#endif