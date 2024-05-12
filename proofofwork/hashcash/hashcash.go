package hashcash

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	hashFormat  = "X-Hashcash: %v:%v:%v:%v:%v:%v:%v"
	hashVersion = 1

	timeFormatFull    = "060102150405"
	timeFormatMinutes = "0601021504"
	timeFormatHours   = "06010215"
	timeFormatDate    = "060102"
)

var ErrInvalidPayload = errors.New("invalid payload")
var ErrTimeout = errors.New("timeout")

type Hashcash struct {
	timeFunc func() time.Time
	prngFunc func() string

	bits   int
	period time.Duration
}

type Layout struct {
	Version  int
	Bits     int
	Datetime time.Time
	Resource string
	Ext      string
	Rand     string
	Counter  int
}

type timeFunc func() time.Time
type prngFunc func() string

func NewHashcash(tfn timeFunc, prng prngFunc, bits int, period time.Duration) *Hashcash {
	return &Hashcash{
		timeFunc: tfn,
		prngFunc: prng,
		bits:     bits,
		period:   period,
	}
}

func (h *Hashcash) Verify(payload string) bool {
	layout, err := h.parse(payload)
	if err != nil {
		return false
	}

	if layout.Datetime.Before(h.timeFunc().Add(-h.period)) {
		return false
	}

	sha := sha1.New()
	io.WriteString(sha, payload)
	hash := sha.Sum(nil)

	// zero bits may not be multiply of bytes(8bits)
	zbits := layout.Bits

	zbytes := int(zbits / 8)
	for i := range zbytes {
		if hash[i] != 0x00 {
			return false
		}
	}

	zrembits := int(zbits % 8)
	if zrembits > 0 {
		rembyte := hash[zbytes]
		if rembyte>>(8-zrembits) != 0 {
			return false
		}
	}

	return true
}

func (h *Hashcash) Generate(ctx context.Context, resource string) (string, error) {
	now := h.timeFunc()

	for cnt := 1; ; cnt++ {
		select {
		case <-ctx.Done():
			return "", ErrTimeout
		default:
		}

		payload := fmt.Sprintf(
			hashFormat,
			hashVersion,
			h.bits,
			now.Format(timeFormatFull),
			resource,
			"", // ext is empty for version 1
			base64.StdEncoding.EncodeToString([]byte(h.prngFunc())),
			base64.StdEncoding.EncodeToString(binary.AppendVarint(nil, int64(cnt))),
		)

		if h.Verify(payload) {
			return payload, nil
		}
	}

}

func (h *Hashcash) parse(payload string) (*Layout, error) {
	signature := "X-Hashcash: "

	chunks := strings.Split(payload[len(signature):], ":")
	if len(chunks) != 7 {
		return nil, ErrInvalidPayload
	}

	version, err := strconv.Atoi(chunks[0])
	if err != nil {
		return nil, err
	}

	bits, err := strconv.Atoi(chunks[1])
	if err != nil {
		return nil, err
	}

	// SHA1 hash size is 160 bytes, which is equal to 1280 bits
	// required number of zero bits cannot be more than that.
	if bits > 1280 {
		return nil, errors.New("number of zero bits exceeds sha1 hash size")
	}

	var datetimeLayout string

	switch len(chunks[2]) {
	case 12:
		datetimeLayout = timeFormatFull
	case 10:
		datetimeLayout = timeFormatMinutes
	case 8:
		datetimeLayout = timeFormatHours
	case 6:
		datetimeLayout = timeFormatDate
	default:
		return nil, fmt.Errorf("invalid datetime: %v", chunks[2])
	}

	datetime, err := time.Parse(datetimeLayout, string(chunks[2]))
	if err != nil {
		return nil, err
	}

	resource := chunks[3]

	ext := chunks[4]

	rand, err := base64.StdEncoding.DecodeString(chunks[5])
	if err != nil {
		return nil, fmt.Errorf("failed to decode random string: %w", err)
	}

	counterBytes, err := base64.StdEncoding.DecodeString(chunks[6])
	if err != nil {
		return nil, fmt.Errorf("failed to decode counter: %w", err)
	}

	cnt, err := binary.ReadVarint(bytes.NewReader(counterBytes))
	if err != nil {
		return nil, fmt.Errorf("failed decode counter: %w", err)
	}

	return &Layout{
		Version:  version,
		Bits:     bits,
		Datetime: datetime,
		Resource: resource,
		Ext:      ext,
		Rand:     string(rand),
		Counter:  int(cnt),
	}, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func RandStringFunc(n int) prngFunc {
	return func() string {
		return RandString(n)
	}
}
