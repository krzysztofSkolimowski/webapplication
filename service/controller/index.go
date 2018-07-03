package controller

import (
	"net/http"

	"github.com/krzysztofSkolimowski/webapplication/service/shared/session"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/view"
)

func IndexGET(w http.ResponseWriter, r *http.Request) {
	session := session.Instance(r)

	if session.Values["id"] != nil {
		v := view.New(r)
		v.Name = "index/auth"
		v.Vars["first_name"] = session.Values["first_name"]
		v.Render(w)
	} else {
		v := view.New(r)
		v.Name = "index/anon"
		v.Render(w)
		return
	}
}
