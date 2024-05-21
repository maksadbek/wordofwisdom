package command

import "context"

type CommandKind uint

const (
	CommandUnknown CommandKind = iota
	CommandChallenge
	CommandQuote
)

var CommandNameToKindMapping = map[string]CommandKind{
	"challenge": CommandChallenge,
	"getquote":  CommandQuote,
}

type Command interface {
	Execute(ctx context.Context, args []string) (string, error)
}
