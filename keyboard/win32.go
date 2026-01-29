package keyboard

import (
	"syscall"
	"unsafe"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")

	procSetWindowsHookExW   = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procGetMessageW         = user32.NewProc("GetMessageW")
	procSendInput           = user32.NewProc("SendInput")
)

const (
	WH_KEYBOARD_LL = 13
	WM_KEYDOWN     = 256
	WM_KEYUP       = 257
	WM_SYSKEYDOWN  = 260
	WM_SYSKEYUP    = 261
)

type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

// LowLevelKeyboardProc is the callback type
type LowLevelKeyboardProc func(nCode int, wparam uintptr, lparam uintptr) uintptr

func SetWindowsHookEx(idHook int, lpfn LowLevelKeyboardProc, hMod uintptr, dwThreadId uint32) uintptr {
	ret, _, _ := procSetWindowsHookExW.Call(
		uintptr(idHook),
		syscall.NewCallback(lpfn),
		hMod,
		uintptr(dwThreadId),
	)
	return ret
}

func CallNextHookEx(hhk uintptr, nCode int, wparam uintptr, lparam uintptr) uintptr {
	ret, _, _ := procCallNextHookEx.Call(hhk, uintptr(nCode), wparam, lparam)
	return ret
}

func UnhookWindowsHookEx(hhk uintptr) bool {
	ret, _, _ := procUnhookWindowsHookEx.Call(hhk)
	return ret != 0
}

func GetMessage(msg *uintptr, hwnd uintptr, msgFilterMin uint32, msgFilterMax uint32) int {
	ret, _, _ := procGetMessageW.Call(
		uintptr(unsafe.Pointer(msg)),
		hwnd,
		uintptr(msgFilterMin),
		uintptr(msgFilterMax),
	)
	return int(ret)
}

// Input Types for SendInput

const (
	INPUT_MOUSE    = 0
	INPUT_KEYBOARD = 1
	INPUT_HARDWAR  = 2

	KEYEVENTF_KEYUP   = 0x0002
	KEYEVENTF_UNICODE = 0x0004
)

type INPUT struct {
	Type uint32
	Ki   KEYBDINPUT
	// MOUSEINPUT and HARDWAREINPUT would overlap here in a C union, but we only need Keyboard
	// We add padding to match C union size if necessary (40 bytes on 64bit usually)
	// For simplicity in Go, we might need a specific struct layout or byte array if we were being very strict,
	// but usually constructing it carefully works.
	Padding [8]byte
}

type KEYBDINPUT struct {
	WVk         uint16
	WScan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

func SendInput(nInputs uint32, pInputs *INPUT, cbSize int) uint32 {
	ret, _, _ := procSendInput.Call(
		uintptr(nInputs),
		uintptr(unsafe.Pointer(pInputs)),
		uintptr(cbSize),
	)
	return uint32(ret)
}

func SendUnicodeChar(r rune) {
	// Press
	i1 := INPUT{Type: INPUT_KEYBOARD}
	i1.Ki.WScan = uint16(r)
	i1.Ki.DwFlags = KEYEVENTF_UNICODE

	// Release
	i2 := INPUT{Type: INPUT_KEYBOARD}
	i2.Ki.WScan = uint16(r)
	i2.Ki.DwFlags = KEYEVENTF_UNICODE | KEYEVENTF_KEYUP

	// Send
	// We can send them in batch
	inputs := []INPUT{i1, i2}
	SendInput(2, &inputs[0], int(unsafe.Sizeof(i1)))
}

func SendBackspace() {
	// VK_BACK = 0x08
	i1 := INPUT{Type: INPUT_KEYBOARD}
	i1.Ki.WVk = 0x08

	i2 := INPUT{Type: INPUT_KEYBOARD}
	i2.Ki.WVk = 0x08
	i2.Ki.DwFlags = KEYEVENTF_KEYUP

	inputs := []INPUT{i1, i2}
	SendInput(2, &inputs[0], int(unsafe.Sizeof(i1)))
}
