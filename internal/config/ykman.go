package config

import (
	"time"

	"github.com/buglloc/yubictl/internal/ykman"
)

type YkManCfg struct {
	LockTTL time.Duration `yaml:"lock_ttl"`
}

func (r *Runtime) NewYkMan() (*ykman.YkMan, error) {
	yk := ykman.NewYkMan(
		ykman.WithLockTTL(r.cfg.YkMan.LockTTL),
	)
	return yk, yk.ReloadDevices()
}
