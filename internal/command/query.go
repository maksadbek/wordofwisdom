package command

import (
	"context"
	"errors"

	"github.com/maksadbek/wordofwisdom/hashcash"
	"github.com/maksadbek/wordofwisdom/internal/kvstore"
)

type QueryCommand struct {
	verifier    *hashcash.Hashcash
	quotesdb    *kvstore.Store[string, string]
	resourcesdb *kvstore.Store[string, string]
}

func (c *QueryCommand) Execute(ctx context.Context, args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("invalid input")
	}

	payload := args[0]

	if resource, ok := c.verifier.Verify(payload); ok {
		_, ok := c.resourcesdb.Del(resource)
		if !ok {
			return "", errors.New("invalid resource")
		}

		_, quote, _ := c.quotesdb.Rand()
		return quote, nil
	}

	return "", errors.New("invalid token")
}
