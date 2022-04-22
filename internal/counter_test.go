package internal

import "testing"

func Test_CountWordsInStringr(t *testing.T) {
	text := "example short string"
	count := CountWordsInString(text)
	if count != 3 {
		t.Errorf("expected 3 actual %d", count)
	}
}
