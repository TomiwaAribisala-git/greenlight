package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/TomiwaAribisala-git/greenlight.git/internal/data"
	"github.com/TomiwaAribisala-git/greenlight.git/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		// If we get a ErrDuplicateEmail error, use the v.AddError() method to manually
		// add a message to the validator instance, and then call our
		// failedValidationResponse() helper.
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Add the "movies:read" permission for the new user.
	err = app.models.Permissions.AddForUser(user.ID, "movies:read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		// As there are now multiple pieces of data that we want to pass to our email
		// templates, we create a map to act as a 'holding structure' for the data. This
		// contains the plaintext version of the activation token for the user, along
		// with their ID.
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	// Note that we also change this to send the client a 202 Accepted status code.
	// This status code indicates that the request has been accepted for processing, but
	// the processing has not been completed.
	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	// err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

/*
	go func() {
		// Run a deferred function which uses recover() to catch any panic, and log an
		// error message instead of terminating the application.
		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()

		err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			// Importantly, if there is an error sending the email then we use the
			// app.logger.PrintError() helper to manage it, instead of the
			// app.serverErrorResponse() helper like before.
			app.logger.PrintError(err, nil)
			// app.serverErrorResponse(w, r, err)
			// return
		}
	}()
*/

// one(user) to many(tokens) relationship
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the plaintext activation token from the request body.
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the plaintext token provided by the client.
	v := validator.New()

	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

/*
Hi,
Thanks for signing up for a Greenlight account. We're excited to have you on board!
For future reference, your user ID number is 123.
To activate your Greenlight account please visit h͟t͟t͟p͟s͟:͟/͟/͟e͟x͟a͟m͟p͟l͟e͟.͟c͟o͟m͟/͟u͟s͟e͟r͟s͟/͟a͟c͟t͟i͟v͟a͟t͟e͟and
enter the following code:
--------------------------
Y3QMGX3PJ3WLRL2YRTQGQ6KRHU
--------------------------
Please note that this code will expire in 3 days and can only be used once.
Thanks,
The Greenlight Team

1. your website just provides a form that performs the PUT request for the user, rather than them needing to do it manually using curl or another tool.

2. When creating the link in this email, don’t rely on the Host header from r.Host to construct the URL, as that would be vulnerable to a host header injection attack.
   The URL domain should be either be hard-coded or passed in as a command-line flag when starting the application.
*/
