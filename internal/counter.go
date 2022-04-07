package internal

import "strings"

func CountWordsInString(text string) int {
	return len(strings.Fields(text))
}
