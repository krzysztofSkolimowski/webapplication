package controller

import (
	"fmt"
	"log"
	"net/http"


	"github.com/gorilla/sessions"
	"github.com/josephspurrier/csrfbanana"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/session"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/view"
	"github.com/krzysztofSkolimowski/webapplication/service/model"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/passhash"
)

const (
	sessionLoginAttempt = "login_attempt"
)

func loginAttempt(s *sessions.Session) {
	if s.Values[sessionLoginAttempt] == nil {
		s.Values[sessionLoginAttempt] = 1
	} else {
		s.Values[sessionLoginAttempt] = s.Values[sessionLoginAttempt].(int) + 1
	}
}

func LoginGET(w http.ResponseWriter, r *http.Request) {
	s := session.Instance(r)

	v := view.New(r)
	v.Name = "login/login"
	v.Vars["token"] = csrfbanana.Token(w, r, s)
	view.Repopulate([]string{"email"}, r.Form, v.Vars)
	v.Render(w)
}

func LoginPOST(w http.ResponseWriter, r *http.Request) {
	s := session.Instance(r)


	if s.Values[sessionLoginAttempt] != nil && s.Values[sessionLoginAttempt].(int) >= 5 {
		log.Println("Brute force login prevented")
		s.AddFlash(view.Flash{"Sorry, no brute force :-)", view.FlashNotice})
		s.Save(r, w)
		LoginGET(w, r)
		return
	}

	if validate, missingField := view.Validate(r, []string{"email", "password"}); !validate {
		s.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		s.Save(r, w)
		LoginGET(w, r)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	result, err := model.UserByEmail(email)

	if err == model.ErrNoResult {
		loginAttempt(s)
		s.AddFlash(view.Flash{"Pass is incorrect - Attempt: " + fmt.Sprintf("%v", s.Values[sessionLoginAttempt]), view.FlashWarning})
		s.Save(r, w)
	} else if err != nil {
		log.Println(err)
		s.AddFlash(view.Flash{"There was an error. Please try again later.", view.FlashError})
		s.Save(r, w)
	} else if passhash.MatchString(result.Password, password) {
		if result.StatusID != 1 {
			s.AddFlash(view.Flash{"Account is inactive so login is disabled.", view.FlashNotice})
			s.Save(r, w)
		} else {
			session.Empty(s)
			s.AddFlash(view.Flash{"Ok!", view.FlashSuccess})
			s.Values["id"] = result.UserID()
			s.Values["email"] = email
			s.Values["first_name"] = result.FirstName
			s.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	} else {
		loginAttempt(s)
		s.AddFlash(view.Flash{"Pass is incorrect - Attempt: " + fmt.Sprintf("%v", s.Values[sessionLoginAttempt]), view.FlashWarning})
		s.Save(r, w)
	}

	LoginGET(w, r)
}

func LogoutGET(w http.ResponseWriter, r *http.Request) {
	s := session.Instance(r)

	if s.Values["id"] != nil {
		session.Empty(s)
		s.AddFlash(view.Flash{"byee, hope to see you soon!", view.FlashNotice})
		s.Save(r, w)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
