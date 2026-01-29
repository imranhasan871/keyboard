# Contributing to Google Input Keyboard

Thank you for your interest in contributing! This project welcomes contributions from everyone.

## How to Contribute

### Reporting Bugs

If you find a bug, please open an issue with:
- A clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Your OS version and Go version

### Suggesting Features

Feature requests are welcome! Please open an issue describing:
- The feature you'd like to see
- Why it would be useful
- How it might work

### Code Contributions

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/your-feature-name`
3. **Make your changes** following the code style
4. **Test your changes** thoroughly
5. **Commit with clear messages**: `git commit -m "Add: feature description"`
6. **Push to your fork**: `git push origin feature/your-feature-name`
7. **Open a Pull Request**

### Code Style

- Follow Go best practices and conventions
- Maintain Clean Architecture principles
- Add comments for complex logic
- Keep functions small and focused
- Write descriptive variable names

### Architecture Guidelines

This project follows **Clean Architecture**:

```
internal/
├── core/
│   ├── domain/      # Interfaces (no dependencies)
│   └── services/    # Business logic (depends only on domain)
└── infrastructure/  # External implementations (depends on domain)
```

**Rules:**
- Domain layer has NO external dependencies
- Services depend only on domain interfaces
- Infrastructure implements domain interfaces
- Dependency flow: Infrastructure → Services → Domain

### Testing

- Test your changes on Windows 10/11
- Verify it works in multiple applications (Notepad, Browser, Word)
- Check that F9 toggle works correctly
- Ensure no memory leaks for long-running sessions

### Adding New Languages

To add support for a new language:

1. Find the language code from [Google Input Tools](https://www.google.com/inputtools/)
2. Update `typing_service.go` line 111:
   ```go
   suggestions, err := s.suggestionService.GetSuggestions(input, "hi-t-i0-und") // Hindi
   ```
3. Update documentation

### Pull Request Process

1. Ensure your code builds: `go build cmd/app/main.go`
2. Update README.md if needed
3. Describe your changes clearly in the PR
4. Link any related issues
5. Wait for review and address feedback

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Help others learn

## Questions?

Feel free to open an issue for any questions about contributing!

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
