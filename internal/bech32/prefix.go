package bech32

import (
	"fmt"
	"strings"
)

// Decode decodes a Bech32 prefix string. If the string is uppercase, the HRP will be uppercase.
func DecodePrefix(s string) (hrp string, data []byte, bits int, err error) {
	if strings.ToLower(s) != s && strings.ToUpper(s) != s {
		return "", nil, 0, fmt.Errorf("mixed case")
	}
	pos := strings.LastIndex(s, "1")
	if pos < 1 {
		return "", nil, 0, fmt.Errorf("separator '1' at invalid position: pos=%d, len=%d", pos, len(s))
	}
	hrp = s[:pos]
	for p, c := range hrp {
		if c < 33 || c > 126 {
			return "", nil, 0, fmt.Errorf("invalid character human-readable part: s[%d]=%d", p, c)
		}
	}
	s = strings.ToLower(s)[pos+1:]
	for p, c := range s {
		d := strings.IndexRune(charset, c)
		if d == -1 {
			return "", nil, 0, fmt.Errorf("invalid character data part: s[%d]=%c", p, c)
		}
		data = append(data, byte(d))
	}

	bits = 5 * len(data)

	data, err = convertBits(data, 5, 8, true)
	if err != nil {
		return "", nil, 0, err
	}
	return hrp, data, bits, nil
}
