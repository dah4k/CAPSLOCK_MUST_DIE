package main

// References:
// - https://gist.github.com/obonyojimmy/52d836a1b31e2fc914d19a81bd2e0a1b
// - https://gist.github.com/jordansissel/1e08b1c65157bde0f30a87c4fb569237
// - https://github.com/susam/uncap/blob/master/uncap.c
// - https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowshookexw

import (
	//"fmt"
	"syscall"
	//"unsafe"
	"golang.org/x/sys/windows"
)

var (
	user32                = windows.NewLazySystemDLL("user32.dll")
	procSetWindowsHookExW = user32.NewProc("SetWindowsHookExW")
	procLowLevelKeyboard  = user32.NewProc("LowLevelKeyboardProc")
)

const (
	WH_KEYBOARD_LL = 13
)

type (
	DWORD     uint32
	WPARAM    uintptr
	LPARAM    uintptr
	LRESULT   uintptr
	HANDLE    uintptr
	HINSTANCE HANDLE
	HHOOK     HANDLE
)

type HOOKPROC func(int, WPARAM, LPARAM) LRESULT

func SetWindowsHookExW(idHook int, lpfn HOOKPROC, hMod HINSTANCE, dwThreadId DWORD) HHOOK {
	ret, _, _ := procSetWindowsHookExW.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId),
	)
	return HHOOK(ret)
}

func LowLevelKeyboardProc(nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procLowLevelKeyboard.Call(
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return LRESULT(ret)
}

func Start() {
	// TODO
}

func main() {
	go Start()
}
