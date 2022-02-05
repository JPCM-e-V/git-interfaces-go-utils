package gitutils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

const FlushPkt = 0
const DelimiterPkt = 1
const ResponseEndPkt = 2

func PktLine(s string) string {
	len_s := len(s)

	if len_s > 65515 {
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
	for _, s := range lines {
		fmt.Fprint(w, PktLine(s))
	}
	fmt.Fprint(w, "0000")
}

func ReadGitProtocol(r io.ReadCloser) ([]string, error) {
	var lines []string

	for {
		b := make([]byte, 4)
		n, err := r.Read(b)
		if err == io.EOF {
			break
		} else if n != 4 {
			return lines, errors.New("hex too short")
		} else {
			num, err0 := strconv.ParseUint(string(b), 16, 16)
			if err0 != nil {
				return lines, errors.New("wrong hex value")
			} else if num >= 4 {
				b_val := make([]byte, num-4)
				n_val, err1 := r.Read(b_val)
				if uint64(n_val) != num-4 || err1 == io.EOF {
					return lines, errors.New("packet too short")
				} else {
					lines = append(lines, "p"+strings.TrimSuffix(string(b_val), "\n"))
				}
			} else if num == 3 {
				return lines, errors.New("hex is 0003")
			} else if num == FlushPkt {
				lines = append(lines, "flush")
			} else if num == DelimiterPkt {
				lines = append(lines, "delimiter")
			} else if num == ResponseEndPkt {
				lines = append(lines, "responseend")
			}
		}
	}

	r.Close()
	return lines, nil
}
