package out

import (
	"net/http"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
)

const SERV_OUT = "out"

func OutHandle(srv *OutServ) http.Handler {
	s := rpc.NewServer()
	s.RegisterCodec(json2.NewCodec(), "application/json")
	s.RegisterService(srv, SERV_OUT)
	return s
}
