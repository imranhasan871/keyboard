package main

import (
	"fmt"
	"google-input-keyboard/internal/core/services"
	"google-input-keyboard/internal/infrastructure/google"
	"google-input-keyboard/internal/infrastructure/windows"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Clear screen and show banner
	fmt.Print("\033[H\033[2J")
	fmt.Println("===============================================================")
	fmt.Println("     Google Input Tools - Bangla Keyboard for Windows")
	fmt.Println("===============================================================")
	fmt.Println()

	// Initialize services
	suggestionService := google.NewGoogleInputGateway()
	inputSim := windows.NewWindowsInputSimulator()
	keyListener := windows.NewWindowsKeyboardListener()
	typingService := services.NewTypingService(suggestionService, inputSim)
	keyListener.SetCallback(typingService.HandleKey)

	// Start keyboard hook
	go func() {
		if err := keyListener.Start(); err != nil {
			fmt.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Println("[OK] Keyboard hook installed successfully")
	fmt.Println()
	fmt.Println("How to use:")
	fmt.Println("   - Type in English (phonetically): amar, bangla, tumi")
	fmt.Println("   - Press SPACE to convert to Bangla")
	fmt.Println("   - Press F9 to toggle ON/OFF")
	fmt.Println()
	fmt.Println("Status: ENABLED")
	fmt.Println()
	fmt.Println("Press Ctrl+C to exit...")
	fmt.Println("---------------------------------------------------------------")

	// Wait for termination
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	fmt.Println()
	fmt.Println("Shutting down... Goodbye!")
	keyListener.Stop()
}
