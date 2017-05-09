package commands

type Command interface {
	// Parse the next input token
	Execute(class int, s string) error
}
