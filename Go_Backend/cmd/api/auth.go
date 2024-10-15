package main

import (
  "net/http"
  "encoding/hex"
  "crypto/sha256"
  "github.com/google/uuid"
  "CoffeeMap/internal/store"
)
// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	store.User			"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
type RegisterUserPayload struct {
  Username string `json:"username" validate:"required,max=100"`
  Email    string `json:"email" validate:"required,email,max=255"`
  Password string `json:"password" validate:"required,min=7,max=72"`
}

type UserWithToken struct {
  *store.User
  Token string `json:"token"`
}
// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
  var payload RegisterUserPayload
  if err := readJSON(w, r, &payload); err != nil {
    app.badRequestResponse(w, r, err)
    return
  }

  if err := Validate.Struct(payload); err != nil {
    app.badRequestResponse(w, r, err)
    return
  }

  user := &store.User{
    Username: payload.Username,
    Email: payload.Email,
  }

  if err := user.Password.Set(payload.Password); err != nil {
    app.internalServerError(w, r, err)
    return
  }

  // store the user
  ctx := r.Context()

  plainToken := uuid.New().String() 

  //store
  hash := sha256.Sum256([]byte(plainToken))
  hashToken := hex.EncodeToString(hash[:])

  err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp)
  if err != nil {
    switch err {
    case store.ErrDuplicateEmail:
      app.badRequestResponse(w, r, err)
    case store.ErrDuplicateUsername:
      app.badRequestResponse(w, r, err)
    default:
      app.internalServerError(w, r, err)
    }
    return
  }

  userWithToken := UserWithToken {
    User: user,
    Token: plainToken,
  }
  //send mail
  //


  if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
    app.internalServerError(w, r, err)
  }
}


