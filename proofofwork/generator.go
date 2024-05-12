package proofofwork

import "context"

type Generator interface {
	Generate(context.Context, string) (string, error)
}
