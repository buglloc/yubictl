package ykman

import (
	"fmt"
	"sync"
	"time"

	"github.com/buglloc/fidoctl"
)

type Yubikey struct {
	dev        fidoctl.Device
	serial     uint32
	client     string
	mu         sync.Mutex
	lastAccess time.Time
}

func NewYubikey(dev fidoctl.Device) (*Yubikey, error) {
	serial, err := dev.Serial()
	if err != nil {
		return nil, fmt.Errorf("get Yubikey serial: %w", err)
	}

	return &Yubikey{
		dev:    dev,
		serial: serial,
	}, nil
}

func (y *Yubikey) IsFree() bool {
	y.mu.Lock()
	defer y.mu.Unlock()

	return y.client == ""
}

func (y *Yubikey) Acquire(clientID string) error {
	y.mu.Lock()
	defer y.mu.Unlock()

	y.client = clientID
	y.lastAccess = time.Now()

	return nil
}

func (y *Yubikey) Reboot() error {
	y.mu.Lock()
	defer y.mu.Unlock()

	if err := y.dev.Reboot(); err != nil {
		return err
	}

	return y.Ping()
}

func (y *Yubikey) Ping() error {
	y.mu.Lock()
	defer y.mu.Unlock()

	y.lastAccess = time.Now()

	return nil
}

func (y *Yubikey) Release() error {
	y.mu.Lock()
	defer y.mu.Unlock()

	y.client = ""

	return nil
}

func (y *Yubikey) Serial() uint32 {
	return y.serial
}

func (y *Yubikey) Path() string {
	return y.dev.Path()
}

func (y *Yubikey) String() string {
	return fmt.Sprintf("Yubikey #%d", y.serial)
}
