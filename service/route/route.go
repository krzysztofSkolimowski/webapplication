package route

import (
	"net/http"
	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/krzysztofSkolimowski/webapplication/service/route/middleware/pprofhandler"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/session"
	"github.com/krzysztofSkolimowski/webapplication/service/route/middleware/logrequest"
	"github.com/krzysztofSkolimowski/webapplication/service/controller"
	"github.com/krzysztofSkolimowski/webapplication/service/route/middleware/acl"
	wrapper "github.com/krzysztofSkolimowski/webapplication/service/route/middleware/httprouterwrapper"
)

func Load() http.Handler {
	return middleware(routes())
}

func LoadHTTPS() http.Handler {
	return middleware(routes())
}

func LoadHTTP() http.Handler {
	return middleware(routes())
	//return http.HandlerFunc(redirectToHTTPS)
}

func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://"+req.Host, http.StatusMovedPermanently)
}

func routes() *httprouter.Router {
	r := httprouter.New()

	r.NotFound = alice.
		New().
		ThenFunc(controller.Error404)

	r.GET("/static/*filepath", wrapper.Handler(alice.
		New().
		ThenFunc(controller.Static)))

	// Home page
	r.GET("/", wrapper.Handler(alice.
		New().
		ThenFunc(controller.IndexGET)))

	// Login
	r.GET("/login", wrapper.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.LoginGET)))
	r.POST("/login", wrapper.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.LoginPOST)))
	r.GET("/logout", wrapper.Handler(alice.
		New().
		ThenFunc(controller.LogoutGET)))

	// Register
	r.GET("/register", wrapper.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.RegisterGET)))
	r.POST("/register", wrapper.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.RegisterPOST)))

	// About
	r.GET("/about", wrapper.Handler(alice.
		New().
		ThenFunc(controller.AboutGET)))

	// Notepad
	r.GET("/notepad", wrapper.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.NotepadReadGET)))
	r.GET("/notepad/create", wrapper.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.NotepadCreateGET)))
	r.POST("/notepad/create", wrapper.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.NotepadCreatePOST)))
	r.GET("/notepad/update/:id", wrapper.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.NotepadUpdateGET)))
	r.POST("/notepad/update/:id", wrapper.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.NotepadUpdatePOST)))
	r.GET("/notepad/delete/:id", wrapper.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.NotepadDeleteGET)))

	r.GET("/debug/pprof/*pprof", wrapper.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(pprofhandler.Handler)))

	return r
}

func middleware(handler http.Handler) http.Handler {
	cs := csrfbanana.New(handler, session.Store, session.Name)
	cs.FailureHandler(http.HandlerFunc(controller.InvalidToken))
	cs.ClearAfterUsage(true)
	cs.ExcludeRegexPaths([]string{"/static(.*)"})
	csrfbanana.TokenLength = 32
	csrfbanana.TokenName = "token"
	csrfbanana.SingleToken = false
	handler = cs
	handler = logrequest.Handler(handler)
	handler = context.ClearHandler(handler)
	return handler
}
