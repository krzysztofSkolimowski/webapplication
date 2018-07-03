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
	sess := session.Instance(r)

	userID := fmt.Sprintf("%s", sess.Values["id"])

	notes, err := model.NotesByUserID(userID)
	if err != nil {
		log.Println(err)
		notes = []model.Note{}
	}

	v := view.New(r)
	v.Name = "notepad/read"
	v.Vars["first_name"] = sess.Values["first_name"]
	v.Vars["notes"] = notes
	v.Render(w)
}

func NotepadCreateGET(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	v := view.New(r)
	v.Name = "notepad/create"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
}

func NotepadCreatePOST(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	if validate, missingField := view.Validate(r, []string{"note"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		NotepadCreateGET(w, r)
		return
	}

	content := r.FormValue("note")

	userID := fmt.Sprintf("%s", sess.Values["id"])

	err := model.NoteCreate(content, userID)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Note added!", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/notepad", http.StatusFound)
		return
	}

	NotepadCreateGET(w, r)
}

func NotepadUpdateGET(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	noteID := params.ByName("id")

	userID := fmt.Sprintf("%s", sess.Values["id"])

	note, err := model.NoteByID(userID, noteID)
	if err != nil { // If the note doesn't exist
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/notepad", http.StatusFound)
		return
	}

	v := view.New(r)
	v.Name = "notepad/update"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["note"] = note.Content
	v.Render(w)
}

func NotepadUpdatePOST(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	if validate, missingField := view.Validate(r, []string{"note"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		NotepadUpdateGET(w, r)
		return
	}

	content := r.FormValue("note")

	userID := fmt.Sprintf("%s", sess.Values["id"])

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	noteID := params.ByName("id")

	err := model.NoteUpdate(content, userID, noteID)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Note updated!", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/notepad", http.StatusFound)
		return
	}

	NotepadUpdateGET(w, r)
}

func NotepadDeleteGET(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	userID := fmt.Sprintf("%s", sess.Values["id"])

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	noteID := params.ByName("id")

	err := model.NoteDelete(userID, noteID)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Note deleted!", view.FlashSuccess})
		sess.Save(r, w)
	}

	http.Redirect(w, r, "/notepad", http.StatusFound)
	return
}
