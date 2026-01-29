package domain

// SuggestionService defines the contract for fetching suggestions
type SuggestionService interface {
	GetSuggestions(text string, lang string) ([]string, error)
}

// KeyEventListener defines how we receive input events
type KeyEventListener interface {
	Start() error
	Stop()
	// SetCallback registers the function to call when a key is pressed.
	// It should return true if the key event should be consumed/blocked, false to let it pass through.
	SetCallback(func(keyCode int, char rune, isDown bool) bool)
}

// InputSimulator defines how we inject text back to the system
type InputSimulator interface {
	TypeString(text string) error
	DeleteCharacters(count int) error
}

// KeyMap Defines specific key codes we care about in a platform-agnostic way
const (
	KeySpace = 32
	KeyEnter = 13
	KeyBack  = 8
	KeyEsc   = 27
	KeyF9    = 120
)
