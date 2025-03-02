package httpd

import (
	"github.com/buglloc/yubictl/internal/touchctl"
	"github.com/buglloc/yubictl/internal/ykman"
)

type Option func(server *Server)

func WithAddr(addr string) Option {
	return func(s *Server) {
		s.addr = addr
	}
}

func WithTouchCtl(t *touchctl.TouchCtl) Option {
	return func(s *Server) {
		s.touch = t
	}
}

func WithYkMan(yk *ykman.YkMan) Option {
	return func(s *Server) {
		s.yk = yk
	}
}
