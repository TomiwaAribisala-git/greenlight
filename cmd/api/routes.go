package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// func (app *application) routes() *httprouter.Router {}
func (app *application) routes() http.Handler {

	router := httprouter.New()

	// Routing errors: Any error messages that our own API handlers send will now be well-formed JSON responses
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// router.HandlerFunc(http.MethodGet, "/v1/movies", app.listMoviesHandler)
	// router.HandlerFunc(http.MethodGet, "/v1/movies", app.requireActivatedUser(app.listMoviesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission("movies:read", app.listMoviesHandler))

	// router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	// router.HandlerFunc(http.MethodPost, "/v1/movies", app.requireActivatedUser(app.createMovieHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission("movies:write", app.createMovieHandler))

	// router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
	// router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requireActivatedUser(app.showMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission("movies:read", app.showMovieHandler))

	// router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateMovieHandler)
	// router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requireActivatedUser(app.updateMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission("movies:write", app.updateMovieHandler))

	// router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)
	// router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requireActivatedUser(app.deleteMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovieHandler))

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	// Using PUT because nothing in our application state (i.e. database) changes after that first request
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createLoginAuthenticationTokenHandler)

	// return app.recoverPanic(router)
	// return router

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))

}

// case of spinning up additional goroutines from within your handlers and there is any
// chance of a panic, you must make sure that you recover any panics from within those goroutines too.
