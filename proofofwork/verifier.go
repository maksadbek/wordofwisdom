package proofofwork

type Verifier interface {
	Verify(string) bool
}
