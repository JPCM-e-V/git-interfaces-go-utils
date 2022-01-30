package gitutils

import (
	"fmt"
	"net/http"
	"unicode"
)

func PktLine(s string) string {
	len_s := len(s)

	if len_s > 65516 {
		return PktLine("ERR To long response.")
	}

	for i := 0; i < len_s; i++ {
		if s[i] > unicode.MaxASCII {
			return PktLine("ERR Non ASCII character found.")
		}
	}
	length := len_s + 5
	return fmt.Sprintf("%04x%s\n", length, s)
}

func WriteGitProtocol(w http.ResponseWriter, lines []string) {
	for _, line := range lines {
		fmt.Fprint(w, PktLine(line))
	}
	fmt.Fprint(w, "0000")
}
