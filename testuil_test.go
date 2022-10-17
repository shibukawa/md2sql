package md2sql

import (
	"regexp"
	"strings"
	"testing"
)

var stripSpacePattern = regexp.MustCompile("(^[ \t]*)")

func TrimIndent(t *testing.T, src string, replace ...string) string {
	t.Helper()
	lines := strings.Split(src, "\n")
	if lines[0] == "" {
		lines = lines[1:]
	}
	matches := stripSpacePattern.FindStringSubmatch(lines[0])
	var b strings.Builder
	for i, line := range lines {
		b.WriteString(strings.TrimPrefix(line, matches[0]))
		if i != len(lines)-1 {
			b.WriteByte('\n')
		}
	}
	result := strings.TrimRight(b.String(), "\n")
	if len(replace) > 1 {
		result = strings.ReplaceAll(result, replace[0], replace[1])
	}
	return result
}
