package touchctl

import (
	"encoding"
	"fmt"
	"strings"
)

var _ encoding.TextUnmarshaler = (*ToucherKind)(nil)
var _ encoding.TextMarshaler = (*ToucherKind)(nil)

type ToucherKind string

const (
	ToucherKindNone ToucherKind = ""
	ToucherKindHw   ToucherKind = "hw"
)

func (k *ToucherKind) UnmarshalText(data []byte) error {
	switch strings.ToLower(string(data)) {
	case "", "none":
		*k = ToucherKindNone
	case "hw":
		*k = ToucherKindHw
	default:
		return fmt.Errorf("invalid upstream kind: %s", string(data))
	}
	return nil
}

func (k ToucherKind) MarshalText() ([]byte, error) {
	return []byte(k), nil
}
