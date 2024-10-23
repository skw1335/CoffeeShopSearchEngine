package main

import (
  "net/http"
  "CoffeeMap/internal/store"
  "strconv"
  "errors"
  "github.com/go-chi/chi/v5"
)
type userKey string

const userCtx userKey = "user"

//  GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
  userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
  if err != nil {
    app.badRequestResponse(w, r, err)
    return
  }
  
  ctx := r.Context()
  
  user, err := app.store.Users.GetByID(ctx, userID) 

  if err != nil {
    switch {
      case errors.Is(err, store.ErrNotFound):
        app.badRequestResponse(w, r, err)
        return
      default: 
        app.internalServerError(w, r, err)
    }
  }

  if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
    app.internalServerError(w, r, err)
  }
}
// ActivateUser godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
  token := chi.URLParam(r, "token")

  err := app.store.Users.Activate(r.Context(), token)
  if err != nil {
    switch err{
    case store.ErrNotFound:
      app.notFoundResponse(w, r, err)
    default:
      app.internalServerError(w, r, err)
    }
    return
  }

  if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
    app.internalServerError(w, r, err)
    }
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
