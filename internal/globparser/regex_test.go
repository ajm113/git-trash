package globparser_test

import (
	"testing"

	"github.com/ajm113/git-trash/internal/globparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegex_Matching(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		input   string
		want    bool
	}{
		// Literal matches
		{name: "exact literal match", pattern: "foo.txt", input: "foo.txt", want: true},
		{name: "literal no match", pattern: "foo.txt", input: "bar.txt", want: false},
		{name: "empty pattern matches empty input", pattern: "", input: "", want: true},
		{name: "empty pattern does not match non-empty input", pattern: "", input: "x", want: false},

		// Anchoring: the match must cover the full string
		{name: "anchored prefix only does not match", pattern: "foo", input: "foobar", want: false},
		{name: "anchored suffix only does not match", pattern: "bar", input: "foobar", want: false},

		// * wildcard
		{name: "star matches anything", pattern: "*", input: "anything at all", want: true},
		{name: "star matches empty", pattern: "*", input: "", want: true},
		{name: "star prefix", pattern: "*.txt", input: "file.txt", want: true},
		{name: "star prefix wrong ext", pattern: "*.txt", input: "file.md", want: false},
		{name: "star suffix", pattern: "foo*", input: "foobar", want: true},
		{name: "star in middle", pattern: "foo*bar", input: "fooXYZbar", want: true},
		{name: "star in middle matches empty", pattern: "foo*bar", input: "foobar", want: true},
		{name: "star in middle no match", pattern: "foo*bar", input: "fooXYZbaz", want: false},
		{name: "multiple stars", pattern: "*foo*", input: "xxfooyy", want: true},

		// ? wildcard
		{name: "question matches single char", pattern: "f?o", input: "foo", want: true},
		{name: "question matches another single char", pattern: "f?o", input: "fxo", want: true},
		{name: "question does not match zero chars", pattern: "f?o", input: "fo", want: false},
		{name: "question does not match two chars", pattern: "f?o", input: "fxxo", want: false},
		{name: "multiple questions", pattern: "??", input: "ab", want: true},
		{name: "multiple questions wrong length", pattern: "??", input: "abc", want: false},

		// Combined wildcards
		{name: "star and question combined", pattern: "?oo*", input: "foobar", want: true},
		{name: "star and question combined no match", pattern: "?oo*", input: "fxobar", want: false},

		// Regex special characters must be treated literally
		{name: "dot is literal", pattern: "a.b", input: "a.b", want: true},
		{name: "dot does not match arbitrary char", pattern: "a.b", input: "axb", want: false},
		{name: "plus is literal", pattern: "a+b", input: "a+b", want: true},
		{name: "parens are literal", pattern: "(foo)", input: "(foo)", want: true},
		{name: "brackets are literal", pattern: "[abc]", input: "[abc]", want: true},
		{name: "brackets not a char class", pattern: "[abc]", input: "a", want: false},
		{name: "backslash is literal", pattern: `a\b`, input: `a\b`, want: true},
		{name: "dollar is literal", pattern: "a$b", input: "a$b", want: true},
		{name: "caret is literal", pattern: "a^b", input: "a^b", want: true},
		{name: "pipe is literal", pattern: "a|b", input: "a|b", want: true},

		// Special chars mixed with wildcards
		{name: "escaped dot with star", pattern: "*.go", input: "main.go", want: true},
		{name: "escaped dot with star rejects substring", pattern: "*.go", input: "maingo", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			re, err := globparser.Regex(tc.pattern)
			require.NoError(t, err)
			require.NotNil(t, re)

			assert.Equal(t, tc.want, re.MatchString(tc.input),
				"pattern %q against input %q", tc.pattern, tc.input)
		})
	}
}

func TestRegex_ReturnsCompiledPattern(t *testing.T) {
	re, err := globparser.Regex("foo*")
	require.NoError(t, err)
	require.NotNil(t, re)

	// The pattern should be anchored on both ends.
	assert.Equal(t, "^foo.*$", re.String())
}
