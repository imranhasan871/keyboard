package windows

import (
	"syscall"
	"unsafe"
)

var (
	user32Mapper         = syscall.NewLazyDLL("user32.dll")
	procToUnicode        = user32Mapper.NewProc("ToUnicode")
	procGetKeyboardState = user32Mapper.NewProc("GetKeyboardState")
)

// VKCodeToChar converts a Virtual Key code to a character using Windows API
func VKCodeToChar(vkCode uint32) rune {
	var keyState [256]byte

	// Get current keyboard state
	procGetKeyboardState.Call(uintptr(unsafe.Pointer(&keyState[0])))

	var buffer [2]uint16
	scanCode := uint32(0)
	flags := uint32(0)

	ret, _, _ := procToUnicode.Call(
		uintptr(vkCode),
		uintptr(scanCode),
		uintptr(unsafe.Pointer(&keyState[0])),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(2),
		uintptr(flags),
	)

	if ret == 1 {
		return rune(buffer[0])
	}

	return 0
}
