package pprofhandler

import (
	"net/http"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"net/http/pprof"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	p := context.Get(r, "params").(httprouter.Params)

	switch p.ByName("pprof") {
	case "/cmdline":
		pprof.Cmdline(w, r)
	case "/profile":
		pprof.Profile(w, r)
	case "/symbol":
		pprof.Symbol(w, r)
	default:
		pprof.Index(w, r)
	}
}
