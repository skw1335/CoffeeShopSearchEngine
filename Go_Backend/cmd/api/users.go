package main

import (
  "net/http"
  "CoffeeMap/internal/store"
  "time"
)

type CreateUserPayload struct {
  ID        int64     `json:"id"`
  Username  string    `json:"username"`
  FirstName string    `json:"first_name"`
  LastName  string    `json:"last_name"`
  Email     string    `json:"email"`
  Password  string    `json:"-"`
  CreatedAt time.Time `json:"created_at"`
}


func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
  var payload CreateUserPayload 
  if err := readJSON(w, r, &payload); err != nil {
    writeJSONError(w, http.StatusBadRequest, err.Error())
  }

  user := &store.User{
    Username:  payload.Username,
    FirstName: payload.FirstName,
    LastName:  payload.LastName,
    Email:     payload.Email,
    Password:  payload.Password,
    ID:        1,
  }
  
  ctx := r.Context()

  if err := app.store.Users.Create(ctx, user); err != nil {
    writeJSONError(w, http.StatusInternalServerError, err.Error())
    return
  }

  if err := writeJSON(w, http.StatusCreated, user); err != nil {
    writeJSONError(w, http.StatusInternalServerError, err.Error())
    return
  }
}
