package httputil

import "strings"

func NormalizeCRLF(raw string) string {
	// Only interpret literal \r\n escapes when the input has no real CR
	// or LF bytes. This supports CLI usage where users type \r\n
	// literally, while protecting body content that legitimately contains
	// the literal sequence (e.g. JSON values with escaped
	// carriage-return/newline) when the input already has real CRLF.
	if !strings.ContainsRune(raw, '\r') && !strings.ContainsRune(raw, '\n') {
		raw = strings.ReplaceAll(raw, `\r\n`, "\r\n")
	}
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\n", "\r\n")
	if !strings.HasSuffix(raw, "\r\n\r\n") {
		if strings.HasSuffix(raw, "\r\n") {
			raw += "\r\n"
		} else {
			raw += "\r\n\r\n"
		}
	}
	return raw
}

func ParseHostHeader(raw string) string {
	for line := range strings.SplitSeq(raw, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "host:") {
			return strings.TrimSpace(line[5:])
		}
	}
	return ""
}
