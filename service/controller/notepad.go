package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/session"
	"github.com/krzysztofSkolimowski/webapplication/service/model"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/view"
)

func NotepadReadGET(w http.ResponseWriter, r *http.Request) {
	s := session.Instance(r)

	userID := fmt.Sprintf("%s", s.Values["id"])

	notes, err := model.NotesByUserID(userID)
	if err != nil {
		log.Println(err)
		notes = []model.Note{}
	}

	v := view.New(r)
	v.Name = "notepad/read"
	v.Vars["first_name"] = s.Values["first_name"]
	v.Vars["notes"] = notes
	v.Render(w)
}

func NotepadCreateGET(w http.ResponseWriter, r *http.Request) {
	s := session.Instance(r)

	v := view.New(r)
	v.Name = "notepad/create"
	v.Vars["token"] = csrfbanana.Token(w, r, s)
	v.Render(w)
}

func NotepadCreatePOST(w http.ResponseWriter, r *http.Request) {
	s := session.Instance(r)

	if validate, missingField := view.Validate(r, []string{"note"}); !validate {
		s.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		s.Save(r, w)
		NotepadCreateGET(w, r)
		return
	}

	n := r.FormValue("note")

	userID := fmt.Sprintf("%s", s.Values["id"])

	err := model.NoteCreate(n, userID)
	if err != nil {
		log.Println(err)
		s.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		s.Save(r, w)
	} else {
		s.AddFlash(view.Flash{"Note added!", view.FlashSuccess})
		s.Save(r, w)
		http.Redirect(w, r, "/notepad", http.StatusFound)
		return
	}

	NotepadCreateGET(w, r)
}

func NotepadUpdateGET(w http.ResponseWriter, r *http.Request) {
	s := session.Instance(r)

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	noteID := params.ByName("id")

	userID := fmt.Sprintf("%s", s.Values["id"])

	n, err := model.NoteByID(userID, noteID)
	if err != nil { // If the n doesn't exist
		log.Println(err)
		s.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		s.Save(r, w)
		http.Redirect(w, r, "/notepad", http.StatusFound)
		return
	}

	v := view.New(r)
	v.Name = "notepad/update"
	v.Vars["token"] = csrfbanana.Token(w, r, s)
	v.Vars["n"] = n.Content
	v.Render(w)
}

func NotepadUpdatePOST(w http.ResponseWriter, r *http.Request) {
	s := session.Instance(r)

	if ok, missing := view.Validate(r, []string{"note"}); !ok {
		s.AddFlash(view.Flash{"Field missing: " + missing, view.FlashError})
		s.Save(r, w)
		NotepadUpdateGET(w, r)
		return
	}

	content := r.FormValue("note")

	usrID := fmt.Sprintf("%s", s.Values["id"])

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	noteID := params.ByName("id")

	err := model.NoteUpdate(content, usrID, noteID)
	if err != nil {
		log.Println(err)
		s.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		s.Save(r, w)
	} else {
		s.AddFlash(view.Flash{"Note updated!", view.FlashSuccess})
		s.Save(r, w)
		http.Redirect(w, r, "/notepad", http.StatusFound)
		return
	}

	NotepadUpdateGET(w, r)
}

func NotepadDeleteGET(w http.ResponseWriter, r *http.Request) {
	s := session.Instance(r)

	usrID := fmt.Sprintf("%s", s.Values["id"])

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	noteID := params.ByName("id")

	err := model.NoteDelete(usrID, noteID)
	if err != nil {
		log.Println(err)
		s.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		s.Save(r, w)
	} else {
		s.AddFlash(view.Flash{"Note deleted!", view.FlashSuccess})
		s.Save(r, w)
	}

	http.Redirect(w, r, "/notepad", http.StatusFound)
	return
}
