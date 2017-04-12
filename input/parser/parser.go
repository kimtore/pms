package parser

import "github.com/ambientsound/pms/input"

type Parser interface {
	// Parse the next input token
	Parse(t input.Token) error
}
