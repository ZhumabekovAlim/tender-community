package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders, makeResponseJSON)

	dynamicMiddleware := alice.New()

	mux := pat.New()

	mux.Post("/clients/signup", dynamicMiddleware.ThenFunc(app.userHandler.SignUp))
	mux.Get("/clients", standardMiddleware.ThenFunc(app.userHandler.GetAllUsers))
	return standardMiddleware.Then(mux)
}
