package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /account/create", dynamic.ThenFunc(app.accountGet))
	mux.Handle("POST /account/create", dynamic.ThenFunc(app.accountPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("GET /thread/view/{id}", dynamic.ThenFunc(app.threadView))
	mux.Handle("GET /thread/view/{id}/message/create", dynamic.ThenFunc(app.msgGet))
	mux.Handle("POST /thread/view/{id}/message/create", dynamic.ThenFunc(app.msgPost))
	mux.Handle("GET /thread/{threadId}/message/view/{id}", dynamic.ThenFunc(app.msgView))
	
	protected := dynamic.Append(app.requireAuthentication)
	
	mux.Handle("GET /thread/create", protected.ThenFunc(app.createThreadGet))
	mux.Handle("POST /thread/create", protected.ThenFunc(app.createThreadPost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))
	
	standard := alice.New(app.logRequest, commonHeaders)
	return standard.Then(mux)
}