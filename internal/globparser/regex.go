package globparser

import (
	"regexp"
	"strings"
)

// GlobToRegex safely converts a glob pattern to a compiled regex.
// Supports * (any chars) and ? (single char) wildcards.
func Regex(pattern string) (*regexp.Regexp, error) {
	var sb strings.Builder

	// Anchor the match to the full string
	sb.WriteString("^")

	for _, ch := range pattern {
		switch ch {
		case '*':
			sb.WriteString(".*") // * matches anything
		case '?':
			sb.WriteString(".") // ? matches one character
		default:
			// Escape everything else so regex special chars are treated literally
			sb.WriteString(regexp.QuoteMeta(string(ch)))
		}
	}

	sb.WriteString("$")

	return regexp.Compile(sb.String())
}
