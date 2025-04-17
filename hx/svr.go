package hx

import (
	"context"
	"net"
	"net/http"
	"time"
)

type Svr struct {
	hs *http.Server
}

func NewServ() *Svr {
	return &Svr{
		hs: &http.Server{
			IdleTimeout: 10 * time.Second,
		},
	}
}

func (s *Svr) SetRootCtx(rootCtx context.Context) *Svr {
	s.hs.BaseContext = func(_ net.Listener) context.Context { return rootCtx }
	return s
}

func (s *Svr) SetAddr(addr string) *Svr {
	s.hs.Addr = addr
	return s
}

func (s *Svr) SetHandler(h http.Handler) *Svr {
	s.hs.Handler = h
	return s
}

func (s *Svr) Run() error {
	return s.hs.ListenAndServe()
}

func (s *Svr) RunTLS(cert, key string) error {
	return s.hs.ListenAndServeTLS(cert, key)
}
