package windows

import (
	"fmt"
	"google-input-keyboard/internal/core/domain"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")

	procSetWindowsHookExW        = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx           = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx      = user32.NewProc("UnhookWindowsHookEx")
	procGetMessageW              = user32.NewProc("GetMessageW")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procSendMessageW             = user32.NewProc("SendMessageW")
	procGetGUIThreadInfo         = user32.NewProc("GetGUIThreadInfo")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
)

const (
	WH_KEYBOARD_LL = 13
	WM_KEYDOWN     = 0x0100
	WM_KEYUP       = 0x0101
	WM_CHAR        = 0x0102
	WM_SYSKEYDOWN  = 260
	WM_SYSKEYUP    = 261
)

// -- Implementation of InputSimulator --

type WindowsInputSimulator struct{}

func NewWindowsInputSimulator() domain.InputSimulator {
	return &WindowsInputSimulator{}
}

func (w *WindowsInputSimulator) TypeString(text string) error {
	hwnd := getFocusedWindow()
	if hwnd == 0 {
		return fmt.Errorf("no focused window")
	}

	for _, r := range text {
		sendChar(hwnd, r)
		time.Sleep(2 * time.Millisecond)
	}
	return nil
}

func (w *WindowsInputSimulator) DeleteCharacters(count int) error {
	hwnd := getFocusedWindow()
	if hwnd == 0 {
		return fmt.Errorf("no focused window")
	}

	for i := 0; i < count; i++ {
		sendBackspace(hwnd)
		time.Sleep(5 * time.Millisecond)
	}
	return nil
}

// -- Low Level Helpers --

type GUITHREADINFO struct {
	cbSize        uint32
	flags         uint32
	hwndActive    uintptr
	hwndFocus     uintptr
	hwndCapture   uintptr
	hwndMenuOwner uintptr
	hwndMoveSize  uintptr
	hwndCaret     uintptr
	rcCaret       [16]byte
}

func getFocusedWindow() uintptr {
	// Get foreground window
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		return 0
	}

	// Get thread ID
	threadID, _, _ := procGetWindowThreadProcessId.Call(hwnd, 0)

	// Get GUI thread info to find focused control
	var info GUITHREADINFO
	info.cbSize = uint32(unsafe.Sizeof(info))

	ret, _, _ := procGetGUIThreadInfo.Call(threadID, uintptr(unsafe.Pointer(&info)))
	if ret != 0 && info.hwndFocus != 0 {
		return info.hwndFocus
	}

	// Fallback to foreground window
	return hwnd
}

func sendBackspace(hwnd uintptr) {
	VK_BACK := uintptr(0x08)

	// Send WM_KEYDOWN
	procSendMessageW.Call(hwnd, WM_KEYDOWN, VK_BACK, 0)

	// Send WM_KEYUP
	procSendMessageW.Call(hwnd, WM_KEYUP, VK_BACK, 0)
}

func sendChar(hwnd uintptr, r rune) {
	// Send WM_CHAR message directly with the Unicode character
	procSendMessageW.Call(hwnd, WM_CHAR, uintptr(r), 0)
}
