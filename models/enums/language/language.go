package language

import (
	"strings"
)

type Language int

const (
	C Language = iota
	Cpp
	Python
	Java
	Cs
	Rust
	Go
	Perl
	Undefined
)

func (language Language) ToString() string {
	switch language {
	case C:
		return "c"
	case Cpp:
		return "cpp"
	case Python:
		return "python"
	case Java:
		return "java"
	case Cs:
		return "cs"
	case Rust:
		return "rust"
	case Go:
		return "go"
	case Perl:
		return "perl"
	default:
		return "Undefined"
	}
}

func LanguageParser(s string) Language {
	switch strings.ToLower(s) {
	case "c":
		return C
	case "cpp":
		return Cpp
	case "python":
		return Python
	case "java":
		return Java
	case "cs":
		return Cs
	case "rust":
		return Rust
	case "go":
		return Go
	case "perl":
		return Perl
	default:
		return Undefined
	}
}
