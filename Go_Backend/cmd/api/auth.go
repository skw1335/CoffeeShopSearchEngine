package main

import (
  "fmt"
  "log"
  "net/http"
  "encoding/hex"
  "crypto/sha256"
  "time"
	
   "github.com/golang-jwt/jwt/v5"
   "github.com/google/uuid"
  "CoffeeMap/internal/store"
  "CoffeeMap/internal/mailer"

)

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
    Role : store.Role{
	name: "user",
	},
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
  activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)

  isProdEnv := app.config.env == "production"

  vars := struct {
    Username string
    ActivationURL string
  } {
    Username: user.Username,
    ActivationURL: activationURL,
  }

  //send mail
  status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
  if err != nil {
    log.Printf("error sending welcome email")
 
    if err := app.store.Users.Delete(ctx, user.ID); err != nil {
      log.Printf("error deleting user!")
    }

    app.internalServerError(w, r, err)
    return
  } 

  log.Printf("Email sent with status code %v", status)


  if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
    app.internalServerError(w, r, err)
  }
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255` 
	Password string `json:"password" validate:"required,min=7,max=255"`
}

// Handler godoc
//
//	@Summary		Creates a token
//	@Description	Creates a token for a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User credentials"
//	@Success		200		{string}	string					"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler (w http.ResponseWriter, r *http.Request) {
	//parse payload credentials
	var payload CreateUserTokenPayload
  if err := readJSON(w, r, &payload); err != nil {
    app.badRequestResponse(w, r, err)
    return
  }

  if err := Validate.Struct(payload); err != nil {
    app.badRequestResponse(w, r, err)
    return
  }
	//retrieve user (check is user exists)
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.unauthorizedErrorResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	// generate the token -> add claims 
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}
	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
    app.internalServerError(w, r, err)
  }
}
