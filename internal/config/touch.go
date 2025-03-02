package config

import (
	"fmt"

	"github.com/buglloc/yubictl/internal/touchctl"
)

type TouchYubikey struct {
	Serial uint32 `koanf:"serial"`
	Pin    uint8  `koanf:"pin"`
}

type TouchCfg struct {
	Kind     touchctl.ToucherKind `koanf:"kind"`
	Yubikeys []TouchYubikey       `koanf:"yubikeys"`
}

func (r *Runtime) NewTouchCtl() (*touchctl.TouchCtl, error) {
	yubiMap := make(map[uint32]uint8, len(r.cfg.Touch.Yubikeys))
	for _, yk := range r.cfg.Touch.Yubikeys {
		if _, exists := yubiMap[yk.Serial]; exists {
			return nil, fmt.Errorf("yubikey with serial %v already exists", yk.Serial)
		}

		yubiMap[yk.Serial] = yk.Pin
	}

	return touchctl.NewTouchCtl(
		touchctl.WithToucherKind(r.cfg.Touch.Kind),
		touchctl.WithYubikeys(yubiMap),
	)
}
