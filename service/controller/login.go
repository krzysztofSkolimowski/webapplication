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
	sessLoginAttempt = "login_attempt"
)

func loginAttempt(sess *sessions.Session) {
	if sess.Values[sessLoginAttempt] == nil {
		sess.Values[sessLoginAttempt] = 1
	} else {
		sess.Values[sessLoginAttempt] = sess.Values[sessLoginAttempt].(int) + 1
	}
}

func LoginGET(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	v := view.New(r)
	v.Name = "login/login"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	view.Repopulate([]string{"email"}, r.Form, v.Vars)
	v.Render(w)
}

func LoginPOST(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	//TODO: uncomment to enable brute force protection
	//if sess.Values[sessLoginAttempt] != nil && sess.Values[sessLoginAttempt].(int) >= 5 {
	//	log.Println("Brute force login prevented")
	//	sess.AddFlash(view.Flash{"Sorry, no brute force :-)", view.FlashNotice})
	//	sess.Save(r, w)
	//	LoginGET(w, r)
	//	return
	//}

	if validate, missingField := view.Validate(r, []string{"email", "password"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	result, err := model.UserByEmail(email)

	if err == model.ErrNoResult {
		loginAttempt(sess)
		sess.AddFlash(view.Flash{"Password is incorrect - Attempt: " + fmt.Sprintf("%v", sess.Values[sessLoginAttempt]), view.FlashWarning})
		sess.Save(r, w)
	} else if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"There was an error. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else if passhash.MatchString(result.Password, password) {
		if result.StatusID != 1 {
			sess.AddFlash(view.Flash{"Account is inactive so login is disabled.", view.FlashNotice})
			sess.Save(r, w)
		} else {
			session.Empty(sess)
			sess.AddFlash(view.Flash{"Ok!", view.FlashSuccess})
			sess.Values["id"] = result.UserID()
			sess.Values["email"] = email
			sess.Values["first_name"] = result.FirstName
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	} else {
		loginAttempt(sess)
		sess.AddFlash(view.Flash{"Password is incorrect - Attempt: " + fmt.Sprintf("%v", sess.Values[sessLoginAttempt]), view.FlashWarning})
		sess.Save(r, w)
	}

	LoginGET(w, r)
}

func LogoutGET(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	if sess.Values["id"] != nil {
		session.Empty(sess)
		sess.AddFlash(view.Flash{"byee, hope to see you soon!", view.FlashNotice})
		sess.Save(r, w)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
