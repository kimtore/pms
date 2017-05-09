package lexer_test

import (
	"strings"
	"testing"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/stretchr/testify/assert"
)

type result struct {
	class int
	str   string
}

var lexluthorTests = []struct {
	input    string
	expected []result
}{
	{
		`a normal sentence`,
		[]result{
			{class: lexer.TokenIdentifier, str: `a`},
			{class: lexer.TokenWhitespace, str: ` `},
			{class: lexer.TokenIdentifier, str: `normal`},
			{class: lexer.TokenWhitespace, str: ` `},
			{class: lexer.TokenIdentifier, str: `sentence`},
			{class: lexer.TokenEnd, str: ``},
		},
	},
	{
		`some "quoted text" here`,
		[]result{
			{class: lexer.TokenIdentifier, str: `some`},
			{class: lexer.TokenWhitespace, str: ` `},
			{class: lexer.TokenIdentifier, str: `quoted text`},
			{class: lexer.TokenWhitespace, str: ` `},
			{class: lexer.TokenIdentifier, str: `here`},
			{class: lexer.TokenEnd, str: ``},
		},
	},
	{
		`;|${}# comment ;|;`,
		[]result{
			{class: lexer.TokenStop, str: `;`},
			{class: lexer.TokenSeparator, str: `|`},
			{class: lexer.TokenVariable, str: `$`},
			{class: lexer.TokenOpen, str: `{`},
			{class: lexer.TokenClose, str: `}`},
			{class: lexer.TokenComment, str: `# comment ;|;`},
			{class: lexer.TokenEnd, str: ``},
		},
	},
	{
		`$"quoted variable" ok`,
		[]result{
			{class: lexer.TokenVariable, str: `$`},
			{class: lexer.TokenIdentifier, str: `quoted variable`},
			{class: lexer.TokenWhitespace, str: ` `},
			{class: lexer.TokenIdentifier, str: `ok`},
			{class: lexer.TokenEnd, str: ``},
		},
	},
	{
		`\v\e\ \r\y "quo\"\\ted $pecial" \$pec\|al`,
		[]result{
			{class: lexer.TokenIdentifier, str: `ve ry`},
			{class: lexer.TokenWhitespace, str: ` `},
			{class: lexer.TokenIdentifier, str: `quo"\ted $pecial`},
			{class: lexer.TokenWhitespace, str: ` `},
			{class: lexer.TokenIdentifier, str: `$pec|al`},
			{class: lexer.TokenEnd, str: ``},
		},
	},
}

func TestLexluthor(t *testing.T) {

	for n, test := range lexluthorTests {

		index := 0
		reader := strings.NewReader(test.input)
		scanner := lexer.NewScanner(reader)

		t.Logf("### Test %d: '%s'", n+1, test.input)

		for {
			class, str := scanner.Scan()

			if index == len(test.expected) {
				if class == lexer.TokenEnd {
					break
				}
				t.Fatalf("Tokenizer generated too many tokens!")
			}

			t.Logf("Token %d: class='%d', literal='%s'", index, class, str)

			check := test.expected[index]

			assert.Equal(t, check.class, class,
				"Token class for token %d is wrong; expected %d but got %d", index, check.class, class)
			assert.Equal(t, check.str, str,
				"String check against token %d failed; expected '%s' but got '%s'", index, check.str, str)

			index++
		}
	}
}
