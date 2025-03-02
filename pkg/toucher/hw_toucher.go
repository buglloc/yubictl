package toucher

import (
	"fmt"
	"sync"
)

var _ Toucher

type HwToucher struct {
	mu  sync.Mutex
	dev *Device
}

func NewHwToucher(d *Device) *HwToucher {
	return &HwToucher{
		dev: d,
	}
}

func (t *HwToucher) Touch(pin uint8) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	err := t.dev.Open()
	if err != nil {
		return fmt.Errorf("open toucher: %w", err)
	}
	defer func() {
		_ = t.dev.Close()
	}()

	n, err := t.dev.Write([]byte{'t', pin})
	if err != nil {
		return fmt.Errorf("write command: %w", err)
	}

	if n != 2 {
		return fmt.Errorf("write command: expected to write 2, got %d", n)
	}

	return nil
}
