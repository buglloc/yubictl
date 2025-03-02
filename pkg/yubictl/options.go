package yubictl

import "time"

const (
	DefaultPingInterval = 5 * time.Second
)

type Option func(*SvcClient)

func WithPingInterval(d time.Duration) Option {
	return func(c *SvcClient) {
		c.pingInterval = d
	}
}
