package hashcash

import (
	"context"
	"testing"
	"time"
)

func GenerateAndVerify(t *testing.T) {
	hashcashInstance := NewHashcash(
		func() time.Time {
			return time.Date(2024, 5, 12, 12, 0, 0, 0, time.UTC)
		},
		RandStringFunc(10),
		20,
		time.Minute*10,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	payload, err := hashcashInstance.Generate(ctx, "resource")
	if err != nil {
		t.Errorf("Error generating hashcash: %v", err)
	}

	if !hashcashInstance.Verify(payload) {
		t.Error("Generated hashcash is not valid")
	}
}

func TestHashcash_GenerateWithTimeout(t *testing.T) {
	hashcashInstance := NewHashcash(
		func() time.Time {
			return time.Date(2024, 5, 12, 12, 0, 0, 0, time.UTC)
		},
		RandStringFunc(10),
		20,
		time.Second*1,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()

	_, err := hashcashInstance.Generate(ctx, "resource")
	if err != ErrTimeout {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}

func VerifyExpiredPayload(t *testing.T) {
	hashcashInstance := NewHashcash(
		func() time.Time {
			return time.Date(2024, 5, 12, 0, 0, 0, 0, time.UTC)
		},
		RandStringFunc(10),
		20,
		time.Minute*10,
	)

	if hashcashInstance.Verify("1:20:1303030600:adam@cypherspace.org::McMybZIhxKXu57jd:ckvi") {
		t.Error("Verification passed for an expired payload")
	}
}
