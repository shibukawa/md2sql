package md2sql

import (
	"regexp"
	"strings"
)

var indentPattern = regexp.MustCompile("(^[ \t]*)")

var trimIndentCache = make(map[string]string)

func trimIndent(src string, additionalIndent string) string {
	if cached, ok := trimIndentCache[src]; ok {
		return cached
	}
	lines := strings.Split(src, "\n")
	if lines[0] == "" {
		lines = lines[1:]
	}
	matches := indentPattern.FindStringSubmatch(lines[0])
	var b strings.Builder
	for i, line := range lines {
		b.WriteString(additionalIndent + strings.TrimPrefix(line, matches[0]))
		if i != len(lines)-1 {
			b.WriteByte('\n')
		}
	}
	result := b.String()
	trimIndentCache[src] = result
	return result
}
