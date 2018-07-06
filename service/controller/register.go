package controller

import (
	"log"
	"net/http"


	"github.com/josephspurrier/csrfbanana"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/session"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/view"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/recaptcha"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/passhash"
	"github.com/krzysztofSkolimowski/webapplication/service/model"
)

func RegisterGET(w http.ResponseWriter, r *http.Request) {
	s := session.Instance(r)

	v := view.New(r)
	v.Name = "register/register"
	v.Vars["token"] = csrfbanana.Token(w, r, s)
	view.Repopulate([]string{"first_name", "last_name", "email"}, r.Form, v.Vars)
	v.Render(w)
}

func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	s := session.Instance(r)

	if s.Values["register_attempt"] != nil && s.Values["register_attempt"].(int) >= 5 {
		log.Println("Brute force register prevented")
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	if ok, missing := view.Validate(r, []string{"first_name", "last_name", "email", "password"}); !ok {
		s.AddFlash(view.Flash{"Field missing: " + missing, view.FlashError})
		s.Save(r, w)
		RegisterGET(w, r)
		return
	}

	if !recaptcha.Verified(r) {
		s.AddFlash(view.Flash{"reCAPTCHA invalid!", view.FlashError})
		s.Save(r, w)
		RegisterGET(w, r)
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	password, err := passhash.HashString(r.FormValue("password"))

	if err != nil {
		log.Println(err)
		s.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		s.Save(r, w)
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	_, err = model.UserByEmail(email)

	if err == model.ErrNoResult {
		ex := model.UserCreate(firstName, lastName, email, password)
		if ex != nil {
			log.Println(ex)
			s.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
			s.Save(r, w)
		} else {
			s.AddFlash(view.Flash{"Account created successfully for: " + email, view.FlashSuccess})
			s.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
	} else if err != nil {
		log.Println(err)
		s.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		s.Save(r, w)
	} else {
		s.AddFlash(view.Flash{"Account already exists for: " + email, view.FlashError})
		s.Save(r, w)
	}

	RegisterGET(w, r)
}
