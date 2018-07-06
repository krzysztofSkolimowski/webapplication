package session

import (
	"github.com/gorilla/sessions"
	"net/http"
)

var (
	Store *sessions.CookieStore
	Name  string
)

type Session struct {
	Options   sessions.Options `json:"Options"`
	Name      string           `json:"Name"`
	SecretKey string           `json:"SecretKey"`
}

func Configure(s Session) {
	Store = sessions.NewCookieStore([]byte(s.SecretKey))
	Store.Options = &s.Options
	Name = s.Name
}

func Instance(r *http.Request) *sessions.Session {
	session, _ := Store.Get(r, Name)
	return session
}

func Empty(s *sessions.Session) {
	for k := range s.Values {
		delete(s.Values, k)
	}
}
