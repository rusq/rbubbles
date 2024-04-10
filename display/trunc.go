package display

import "strings"

func Trunc(s string, sz int) string {
	if sz < 1 {
		return s
	}
	firstCR := strings.Index(s, "\n")
	if -1 < firstCR && firstCR < sz {
		return s[:firstCR] + "⏎"
	}
	if len(s) > sz {
		return s[:sz-1] + "…"
	}
	return s
}
