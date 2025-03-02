package touchctl

import (
	"fmt"

	"github.com/buglloc/yubictl/pkg/toucher"
)

type TouchCtl struct {
	tKind    ToucherKind
	touch    toucher.Toucher
	yubikeys map[uint32]uint8
}

func NewTouchCtl(opts ...Option) (*TouchCtl, error) {
	t := &TouchCtl{}
	for _, opt := range opts {
		opt(t)
	}

	return t, t.init()
}

func (t *TouchCtl) init() error {
	switch t.tKind {
	case ToucherKindNone:
		t.touch = toucher.NewNopToucher()

	case ToucherKindHw:
		tDev, err := toucher.FirstDevice()
		if err != nil {
			return fmt.Errorf("find toucher device: %w", err)
		}

		t.touch = toucher.NewHwToucher(tDev)

	default:
		return fmt.Errorf("unknown toucher kind %s", t.tKind)
	}

	return nil
}

func (t *TouchCtl) Touch(serial uint32) error {
	pin, ok := t.yubikeys[serial]
	if !ok {
		return fmt.Errorf("pin doesn't configured for yubikey with serial %d", serial)
	}

	return t.touch.Touch(pin)
}
