package reader

import (
	"testing"
)

func Test_makePromptPrefix(t *testing.T) {
	expected := "\033[96m> \033[0m"
	if prefix := makePromptPrefix("> ", 96); prefix != expected {
		t.Error("prefix is ", prefix, " but should be ", expected)
	}
}