// Package tabcomplete provides tab-complete functionality for Command.
package tabcomplete

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/commands"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/parser"
	"github.com/ambientsound/pms/utils"
)

// TabComplete provides tab-complete functionality. It takes a string, which is
// fed to the Command. The Command returns a tab-complete list, and TabComplete
// manages a cyclic list of possible tabcomplete candidates.
type TabComplete struct {
	api      api.API
	base     string   // the portion of the text that is static
	cursor   int      // where in the index the cursor is
	items    []string // an index of possible tabcomplete candidates
	original string   // the last token in the stream
	source   string   // string that was used to initialize TabComplete
}

// New returns TabComplete.
func New(source string, a api.API) *TabComplete {
	return &TabComplete{
		api:    a,
		source: source,
	}
}

// Len returns the number of tab completion candidates.
func (t *TabComplete) Len() int {
	return len(t.items)
}

// Active returns true if the tab completion is initialized successfully and
// has more than one candidate.
func (t *TabComplete) Active() bool {
	return t.Len() > 0
}

// Scan cycles through tabcomplete entries and returns the next full string.
// If no tabcomplete candidates can be found, Scan returns an empty string along with an error.
func (t *TabComplete) Scan() (string, error) {

	if !t.Active() {
		if err := t.init(); err != nil {
			return "", err
		}
	}

	return t.next(), nil
}

// next returns the next tabcomplete item in the list.
func (t *TabComplete) next() string {
	// Wrap around
	if t.cursor >= t.Len() {
		t.cursor = 0
	}

	// Return next in line
	s := t.base + t.items[t.cursor]
	t.cursor++

	return s
}

// init initializes the tab completion. The source string is fed to the
// command, and its TabComplete function is invoked to return a list of strings.
func (t *TabComplete) init() error {

	// Set up the input token stream for the parser
	reader := strings.NewReader(t.source)
	scanner := lexer.NewScanner(reader)
	parser := parser.New(scanner)

	// Find the verb
	tok, verb := parser.ScanIgnoreWhitespace()
	if tok != lexer.TokenIdentifier {
		return fmt.Errorf("Tab completing verb '%s', but this is not an identifier", verb)
	}

	// Instantiate the Command registered with this verb
	cmd := commands.New(verb, t.api)
	if cmd == nil {

		// No command. If the next token is anything other than END, return an error.
		tok, _ := parser.Scan()
		if tok != lexer.TokenEnd {
			return fmt.Errorf("No command '%s', cannot autocomplete", verb)
		}

		// Otherwise, try command tabcompletion.
		items := utils.TokenFilter(verb, commands.Keys())
		if len(items) == 0 {
			return fmt.Errorf("No tab complete candidates")
		}
		t.set([]string{}, items, false)

		return nil
	}

	// Parse the remaining text
	cmd.Parse(scanner)

	// Concatenate scanned tokens, except the last one
	tokens := cmd.Scanned()
	if len(tokens) < 2 {
		return fmt.Errorf("Not enough data to tab complete")
	}

	// Get tabcomplete items
	items := cmd.TabComplete()
	if len(items) == 0 {
		return fmt.Errorf("No tab complete candidates")
	}

	// Copy all parsed tokens into the final completion string, except the end token
	lastToken := lexer.TokenEnd
	stringTokens := make([]string, 1)
	stringTokens[0] = verb
	for i := range tokens {
		if tokens[i].Tok == lexer.TokenEnd {
			break
		}
		lastToken = tokens[i].Tok
		stringTokens = append(stringTokens, tokens[i].Lit)
	}

	// Initialize autocomplete
	whitespaceEnds := lastToken == lexer.TokenWhitespace
	t.set(stringTokens, items, whitespaceEnds)

	return nil
}

// set stores the tabcomplete configuration.
func (t *TabComplete) set(stringTokens []string, items []string, whitespaceEnds bool) {
	t.items = items
	t.cursor = t.Len()
	if whitespaceEnds {
		t.base = strings.Join(stringTokens, "")
	} else if len(stringTokens) > 0 {
		maxLen := len(stringTokens) - 1
		t.original = stringTokens[maxLen]
		t.base = strings.Join(stringTokens[:maxLen], "")
	}
}
