package windows

import (
	"fmt"
	"google-input-keyboard/internal/core/domain"
	"syscall"
	"unsafe"
)

type WindowsKeyboardListener struct {
	hookHandle uintptr
	callback   func(keyCode int, char rune, isDown bool) bool
}

func NewWindowsKeyboardListener() domain.KeyEventListener {
	return &WindowsKeyboardListener{}
}

// Global variable to hold the callback because the C-callback cannot capture closure state easily
var globalCallback func(keyCode int, char rune, isDown bool) bool

func (w *WindowsKeyboardListener) SetCallback(cb func(keyCode int, char rune, isDown bool) bool) {
	w.callback = cb
	globalCallback = cb
}

func (w *WindowsKeyboardListener) Start() error {
	// Hook must be set in the same thread that pumps messages

	// We pass 0 as hMod because we are hooking in global scope but usually requires DLL for global.
	// However, for WH_KEYBOARD_LL (LowLevel), hMod can be NULL if threadId is 0.
	// Actually, on Windows, for Global Hooks, hMod usually needs to be the handle to the module containing the hook proc.
	// Go compiles to a static binary, so GetModuleHandle(NULL) gives us our own module.

	// NOTE: LowLevel hooks don't need to be in a DLL. They work with context switching.

	hMod, _, _ := syscall.NewLazyDLL("kernel32.dll").NewProc("GetModuleHandleW").Call(0)

	w.hookHandle, _, _ = procSetWindowsHookExW.Call(
		uintptr(WH_KEYBOARD_LL),
		syscall.NewCallback(hookProc),
		hMod,
		0,
	)

	if w.hookHandle == 0 {
		return fmt.Errorf("failed to set hook")
	}

	fmt.Println("Keyboard Hook Installed.")

	var msg struct {
		hwnd    uintptr
		message uint32
		wParam  uintptr
		lParam  uintptr
		time    uint32
		pt      struct{ x, y int32 }
	}

	// Message Loop
	for {
		ret, _, _ := procGetMessageW.Call(
			uintptr(unsafe.Pointer(&msg)),
			0,
			0,
			0,
		)
		if ret == 0 {
			break
		}
		// TranslateMessage & DispatchMessage would go here in a full UI app
	}

	return nil
}

func (w *WindowsKeyboardListener) Stop() {
	if w.hookHandle != 0 {
		procUnhookWindowsHookEx.Call(w.hookHandle)
		w.hookHandle = 0
	}
}

// The Hook Procedure
func hookProc(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode == 0 && globalCallback != nil {
		// Parse KBDLLHOOKSTRUCT
		kbd := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))

		isDown := (wParam == WM_KEYDOWN || wParam == WM_SYSKEYDOWN)

		// Convert VK code to character
		char := VKCodeToChar(kbd.VkCode)

		shouldBlock := globalCallback(int(kbd.VkCode), char, isDown)

		if shouldBlock {
			return 1 // Block the key
		}
	}

	ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
	return ret
}

type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}
