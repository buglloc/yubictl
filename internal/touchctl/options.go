package touchctl

type Option func(*TouchCtl)

func WithToucherKind(kind ToucherKind) Option {
	return func(t *TouchCtl) {
		t.tKind = kind
	}
}

func WithYubikeys(yubikeys map[uint32]uint8) Option {
	return func(t *TouchCtl) {
		t.yubikeys = yubikeys
	}
}
