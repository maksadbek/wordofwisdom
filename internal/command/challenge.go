package command

import (
	"context"

	"github.com/maksadbek/wordofwisdom/internal/kvstore"
)

type ChallengeCommand struct {
	randfunc  func() string
	resources *kvstore.Store[string, string]
}

func NewChallangeCommand() Command {
	return &ChallengeCommand{}
}

func (c *ChallengeCommand) Execute(ctx context.Context, args []string) (string, error) {
	res := c.randfunc()

	c.resources.Set(res, "")

	return res, nil
}
