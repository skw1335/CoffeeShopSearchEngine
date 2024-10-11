package main

import (
  "net/http"
  "CoffeeMap/internal/store"
  "strconv"
  "errors"
  "github.com/go-chi/chi/v5"
)



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

