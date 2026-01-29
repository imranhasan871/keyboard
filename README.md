# Google Input Tools Keyboard for Windows

A native Windows keyboard application that brings **Google Input Tools** transliteration to your entire system. Type phonetically in English and get Bangla (or other languages) everywhere!

## 🎯 Features

- ✅ **System-wide support** - Works in ANY application (Notepad, Browser, Word, VS Code, etc.)
- ✅ **Google Input Tools API** - Uses the same API as Google's web input tools
- ✅ **Real-time transliteration** - Type "amar" → Press Space → Get "আমার"
- ✅ **Toggle ON/OFF** - Press F9 to enable/disable anytime
- ✅ **Clean Architecture** - Built with SOLID principles and proper separation of concerns
- ✅ **Lightweight** - No heavy dependencies, just pure Go

## 🚀 Quick Start

### Run the Application

```powershell
.\keyboard_app.exe
```

### Usage

1. **Start typing** in any application (e.g., Notepad)
2. **Type phonetically** in English: `amar`, `bangla`, `tumi`, etc.
3. **Press Space** - The text will be replaced with Bangla: `আমার`, `বাংলা`, `তুমি`
4. **Press F9** - Toggle the input tool ON/OFF when you want to type in English

### Example

```
Type: amar [Space]
Result: আমার 

Type: ami banglay gan gai [Space after each word]
Result: আমি বাংলায় গান গাই
```

## 🏗️ Architecture

This project follows **Clean Architecture** principles:

```
keyboard/
├── cmd/app/                    # Application entry point
│   └── main.go                 # Dependency injection & wiring
├── internal/
│   ├── core/
│   │   ├── domain/             # Interfaces & domain models (DIP)
│   │   │   └── interfaces.go
│   │   └── services/           # Business logic (Use Cases)
│   │       └── typing_service.go
│   └── infrastructure/
│       ├── google/             # Google API Gateway
│       │   └── google_gateway.go
│       └── windows/            # Windows-specific implementations
│           ├── char_mapper.go
│           ├── input_simulator.go
│           └── keyboard_listener.go
├── go.mod
└── keyboard_app.exe            # Compiled binary
```

### SOLID Principles Applied

- **Single Responsibility** - Each service has one clear purpose
- **Open/Closed** - Easy to extend with new languages or input methods
- **Liskov Substitution** - All implementations follow their interfaces
- **Interface Segregation** - Small, focused interfaces
- **Dependency Inversion** - Core logic depends on abstractions, not implementations

## 🔧 Building from Source

```powershell
# Build
go build -o keyboard_app.exe cmd/app/main.go

# Run
.\keyboard_app.exe
```

## 🌍 Supported Languages

Currently configured for **Bangla** (`bn-t-i0-und`), but you can easily add:

- Hindi: `hi-t-i0-und`
- Tamil: `ta-t-i0-und`
- Telugu: `te-t-i0-und`
- Gujarati: `gu-t-i0-und`
- Kannada: `kn-t-i0-und`
- Malayalam: `ml-t-i0-und`
- Marathi: `mr-t-i0-und`
- Punjabi: `pa-t-i0-und`
- Urdu: `ur-t-i0-und`

To change language, edit `internal/core/services/typing_service.go` line 111:
```go
suggestions, err := s.suggestionService.GetSuggestions(input, "hi-t-i0-und") // Hindi
```

## ⚠️ Important Notes

### Antivirus Warning
Since this app uses **low-level keyboard hooks**, your antivirus might flag it as a keylogger. This is a **false positive** - the app only monitors your typing to provide transliteration. You may need to:
- Add an exception in Windows Defender
- Whitelist the executable

### Administrator Privileges
Some applications (like Task Manager or apps running as Admin) may not receive the simulated input. This is a Windows security feature.

## 🎮 Controls

| Key | Action |
|-----|--------|
| **F9** | Toggle input tool ON/OFF |
| **Ctrl+C** (in terminal) | Stop the application |
| **Space** | Commit transliteration |
| **Backspace** | Delete characters (updates buffer) |

## 🛠️ How It Works

1. **Keyboard Hook** - Intercepts all keyboard events globally
2. **Character Mapping** - Converts Virtual Key codes to actual characters
3. **Buffer Tracking** - Builds up words as you type
4. **API Call** - Fetches suggestions from Google Input Tools when you press Space
5. **Text Replacement** - Deletes the English text and types the Bangla equivalent

## 📝 Future Enhancements

- [ ] GUI popup showing suggestions before Space press
- [ ] Multiple suggestion selection (Tab to cycle)
- [ ] Custom dictionary support
- [ ] Auto-learn from user corrections
- [ ] System tray icon with settings
- [ ] Multi-language switching (Ctrl+Shift to cycle)
- [ ] Installer package

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Ways to Contribute
- Report bugs and issues
- Suggest new features
- Add support for more languages
- Improve documentation
- Submit pull requests

## License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

### What this means:
- Free to use, modify, and distribute
- Can be used in commercial projects
- No warranty provided
- Must include license notice

## Acknowledgments

- **Google Input Tools API** - For providing the transliteration service
- **Windows API** - For low-level keyboard hooks and input simulation
- **Go Community** - For excellent tools and libraries

## Disclaimer

This is an independent project and is not officially affiliated with or endorsed by Google. Google Input Tools is a trademark of Google LLC.

---

**Made with Go and Clean Architecture principles**

**Star this repo if you find it useful!**
