package controller

import (
	"net/http"

	"github.com/krzysztofSkolimowski/webapplication/service/shared/view"
)

func AboutGET(w http.ResponseWriter, r *http.Request) {
	v := view.New(r)
	v.Name = "about/about"
	v.Render(w)
}
