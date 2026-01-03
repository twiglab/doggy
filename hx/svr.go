package hx

import (
	"context"
	"net"
	"net/http"
	"time"
)

func NewServer(rootCtx context.Context, addr string, h http.Handler) *http.Server {

	svr := &http.Server{
		Addr:         addr,
		Handler:      h,
		IdleTimeout:  90 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		BaseContext:  func(_ net.Listener) context.Context { return rootCtx },
	}

	return svr
}
