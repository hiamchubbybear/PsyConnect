package stringprep

import "fmt"



type Error struct {
	Msg  string
	Rune rune
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (rune: '\\u%04x')", e.Msg, e.Rune)
}
