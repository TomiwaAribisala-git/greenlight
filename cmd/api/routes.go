package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// func (app *application) routes() *httprouter.Router {}
func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/movies", app.listMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler) // v1 API Versioning
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	return app.recoverPanic(app.rateLimit(router))

	// return app.recoverPanic(router)
	// return router
}

// case of spinning up additional goroutines from within your handlers and there is any
// chance of a panic, you must make sure that you recover any panics from within those goroutines too.
