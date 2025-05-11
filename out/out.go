package out

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/twiglab/doggy/hx"
)

type SumArgs struct {
	Table     string
	StartTime int64
	EndTime   int64
}

type SumReply struct {
	Total int64 `json:"total"`
}

type OutServ struct{}

func NewOutServ() *OutServ {
	return &OutServ{}
}

func (h *OutServ) Sum(r *http.Request, args *SumArgs, reply *SumReply) error {
	reply.Total = 8888
	return nil
}

func sum(srv *OutServ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		args := SumArgs{}
		if err := hx.BindAndClose(r, &args); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		reply := SumReply{}

		if err := srv.Sum(r, &args, &reply); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_ = hx.JsonTo(http.StatusOK, reply, w)
	}
}

func OutHandle(srv *OutServ) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.NoCache)

	r.HandleFunc("/sum", sum(srv))
	return r
}
