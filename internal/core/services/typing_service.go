package services

import (
	"fmt"
	"google-input-keyboard/internal/core/domain"
	"strings"
	"sync"
	"time"
)

type TypingService struct {
	suggestionService domain.SuggestionService
	inputSimulator    domain.InputSimulator

	currentBuffer strings.Builder
	isEnabled     bool
	mu            sync.Mutex
}

func NewTypingService(
	suggestionService domain.SuggestionService,
	inputSimulator domain.InputSimulator,
) *TypingService {
	return &TypingService{
		suggestionService: suggestionService,
		inputSimulator:    inputSimulator,
		isEnabled:         true,
	}
}

// HandleKey processes a key event and returns true if the key should be blocked
func (s *TypingService) HandleKey(keyCode int, char rune, isDown bool) bool {
	if !isDown {
		return false // We only care about KeyDown for logic usually
	}

	if keyCode == domain.KeyF9 {
		s.mu.Lock()
		s.isEnabled = !s.isEnabled
		if s.isEnabled {
			fmt.Println("[ENABLED] Bangla Input is ON")
		} else {
			fmt.Println("[DISABLED] Bangla Input is OFF (typing in English)")
		}
		s.currentBuffer.Reset()
		s.mu.Unlock()
		return true
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isEnabled {
		return false
	}

	// Special keys handling
	if keyCode == domain.KeySpace {
		if s.currentBuffer.Len() > 0 {
			// Buffer exists, user pressed space -> Commit the best suggestion
			s.commitBestSuggestion()
			return true // Block the original Space, we will handle insertion
		}
		return false
	}

	if keyCode == domain.KeyBack {
		if s.currentBuffer.Len() > 0 {
			// Manually remove last char from our buffer
			// We simulate a backspace to the OS if we want visual feedback,
			// OR we just intercept it.
			// For a simple overlay approach, we might want to let the user see what they type
			// in Latin first, OR we might block it.
			// Let's go with: User sees Latin, we delete it all and replace on Space.
			// So we allow Backspace to pass through, but update our internal buffer.
			bufStr := s.currentBuffer.String()
			if len(bufStr) > 0 {
				s.currentBuffer.Reset()
				s.currentBuffer.WriteString(bufStr[:len(bufStr)-1])
			}
			return false // Let system handle the backspace visual
		}
		return false
	}

	// Check if it's a printable character (letter or number)
	if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' {
		// Convert to lowercase for consistency
		if char >= 'A' && char <= 'Z' {
			char = char + 32
		}

		s.currentBuffer.WriteRune(char)

		// Asynchronously log buffer state
		go s.updateSuggestions()

		return false // Let the character appear on screen
	}

	return false
}

func (s *TypingService) commitBestSuggestion() {
	input := s.currentBuffer.String()
	if input == "" {
		return
	}

	// sync fetch for now (blocking on Space press)
	suggestions, err := s.suggestionService.GetSuggestions(input, "bn-t-i0-und")
	if err != nil {
		fmt.Printf("Error fetching suggestions: %v\n", err)
		s.currentBuffer.Reset()
		return
	}

	if len(suggestions) == 0 {
		s.currentBuffer.Reset()
		return
	}

	best := suggestions[0]

	// CRITICAL: We must inject input OUTSIDE the hook callback context
	// Windows blocks SendInput from within hook callbacks
	// So we delay it slightly using a goroutine
	go func() {
		// Small delay to exit hook context
		time.Sleep(50 * time.Millisecond)

		// We need to remove the characters the user just typed from the screen.
		s.inputSimulator.DeleteCharacters(len(input))

		// Type the new word
		s.inputSimulator.TypeString(best + " ") // Append space
	}()

	s.currentBuffer.Reset()
}

func (s *TypingService) updateSuggestions() {
	// In a GUI app, this would update the popup window.
	// For now, we just track the buffer silently
	input := s.currentBuffer.String()
	if input == "" {
		return
	}

	// Future: Display suggestions in a popup UI
}
