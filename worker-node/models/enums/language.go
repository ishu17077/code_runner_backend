package enums

import "strings"
type Language int



const (
	C Language = iota
	Cpp
	Python
)

func (language Language) ToString() string {
	switch language {
	case C:
		return "c"
	case Cpp:
		return "cpp"
	case Python:
		return "python"
	default:
		return "c"
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
	default:
		return C
	}
}
