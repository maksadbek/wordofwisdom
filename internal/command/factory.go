package command

import (
	"errors"
	"strings"
	"sync"

	"github.com/maksadbek/wordofwisdom/hashcash"
	"github.com/maksadbek/wordofwisdom/internal/kvstore"
	"github.com/maksadbek/wordofwisdom/pkg/random"
)

var errCommandNotFound = errors.New("command not found")

type Factory struct {
	commands map[CommandKind]Command

	lock sync.RWMutex
}

func NewFactory(hc *hashcash.Hashcash, quotesdb *kvstore.Store[string, string], resourcesDB *kvstore.Store[string, string]) *Factory {
	commands := map[CommandKind]Command{
		CommandQuote: &QueryCommand{
			verifier:    hc,
			quotesdb:    quotesdb,
			resourcesdb: resourcesDB,
		},
		CommandChallenge: &ChallengeCommand{
			resources: resourcesDB,
			randfunc:  func() string { return random.Alnum(8) },
		},
	}

	return &Factory{
		commands: commands,
	}
}

func (f *Factory) Find(cmd string) (Command, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	kind, ok := CommandNameToKindMapping[strings.ToLower(cmd)]
	if !ok {
		return nil, errCommandNotFound
	}

	return f.commands[kind], nil
}
