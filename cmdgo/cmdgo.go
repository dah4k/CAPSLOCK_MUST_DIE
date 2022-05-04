package main

// Inspired By:
// - https://gist.github.com/obonyojimmy/52d836a1b31e2fc914d19a81bd2e0a1b
// - https://gist.github.com/jordansissel/1e08b1c65157bde0f30a87c4fb569237
// - https://github.com/susam/uncap/blob/master/uncap.c

import (
	//"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var (
	user32                = windows.NewLazySystemDLL("user32.dll")
	procSetWindowsHookExW = user32.NewProc("SetWindowsHookExW")
	procLowLevelKeyboard  = user32.NewProc("LowLevelKeyboardProc")
	procSendInput         = user32.NewProc("SendInput")
	procGetMessage        = user32.NewProc("GetMessage")
	procTranslateMessage  = user32.NewProc("TranslateMessage")
	procDispatchMessage   = user32.NewProc("DispatchMessage")
)

// - https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowshookexw
// - https://docs.microsoft.com/en-us/windows/win32/inputdev/wm-keydown
// - https://docs.microsoft.com/en-us/windows/win32/inputdev/wm-keyup
// - https://docs.microsoft.com/en-us/windows/win32/inputdev/wm-syskeydown
// - https://docs.microsoft.com/en-us/windows/win32/inputdev/wm-syskeyup
// - https://docs.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes
const (
	WH_KEYBOARD_LL = 13

	WM_KEYDOWN    = 0x0100
	WM_KEYUP      = 0x0101
	WM_SYSKEYDOWN = 0x0104
	WM_SYSKEYUP   = 0x0105

	VK_CAPITAL  = 0x14
	VK_LCONTROL = 0xA3
)

// https://docs.microsoft.com/en-us/windows/win32/winprog/windows-data-types
type (
	BOOL      uint8
	DWORD     uint32
	WORD      uint16
	INT       int32
	LONG      int32
	WPARAM    uintptr
	LPARAM    uintptr
	LRESULT   uintptr
	HANDLE    uintptr
	ULONG_PTR uintptr
	LPINPUT   *INPUT
	HINSTANCE HANDLE
	HHOOK     HANDLE
	HWND      HANDLE
)

// https://docs.microsoft.com/en-us/windows/win32/api/windef/ns-windef-point
type POINT struct {
	x LONG
	y LONG
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-msg
type MSG struct {
	hwnd     HWND
	message  uint
	wParam   WPARAM
	lParam   LPARAM
	time     DWORD
	pt       POINT
	lPrivate DWORD
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-kbdllhookstruct
type KBDLLHOOKSTRUCT struct {
	vkCode      DWORD
	scanCode    DWORD
	flags       DWORD
	time        DWORD
	dwExtraInfo uintptr
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-keybdinput
type KEYBDINPUT struct {
	wVk         WORD
	wScan       WORD
	dwFlags     DWORD
	time        DWORD
	dwExtraInfo ULONG_PTR
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-input
type INPUT struct {
	dwType DWORD
	ki     KEYBDINPUT // XXX: Golang cannot union(MOUSEINPUT, KEYBDINPUT, HARDWAREINPUT)
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nc-winuser-hookproc
type HOOKPROC func(int, WPARAM, LPARAM) LRESULT

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowshookexw
func SetWindowsHookExW(idHook int, lpfn HOOKPROC, hMod HINSTANCE, dwThreadId DWORD) HHOOK {
	ret, _, _ := procSetWindowsHookExW.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId),
	)
	return HHOOK(ret)
}

// https://docs.microsoft.com/en-us/windows/win32/winmsg/lowlevelkeyboardproc
func LowLevelKeyboardProc(nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procLowLevelKeyboard.Call(
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return LRESULT(ret)
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendinput
func SendInput(cInputs uint, pInputs *INPUT /* LPINPUT */, cbSize int) uint {
	ret, _, _ := procSendInput.Call(
		uintptr(cInputs),
		uintptr(unsafe.Pointer(pInputs)),
		uintptr(cbSize),
	)
	return uint(ret)
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getmessage
func GetMessage(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32) BOOL {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax),
	)
	return BOOL(ret)
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-translatemessage
func TranslateMessage(lpMsg *MSG) BOOL {
	ret, _, _ := procTranslateMessage.Call(uintptr(unsafe.Pointer(lpMsg)))
	return BOOL(ret)
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-dispatchmessage
func DispatchMessage(lpMsg *MSG) LRESULT {
	ret, _, _ := procDispatchMessage.Call(uintptr(unsafe.Pointer(lpMsg)))
	return LRESULT(ret)
}

////////////////////////////////////////////////////////////////////////////////

func (ret LRESULT) keyboardHook(nCode int, wParam WPARAM, lParam LPARAM) {
	kbdstruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
	keyCode := byte(kbdstruct.vkCode)

	if keyCode == VK_CAPITAL {
		// TODO: SendInput VK_LCONTROL
	}
}

func Start() {
	// TODO: Setup keyboardHook, then run MSG loop
}

func main() {
	go Start()
}
