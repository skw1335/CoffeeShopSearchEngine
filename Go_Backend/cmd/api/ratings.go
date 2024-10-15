package main

import (
  "net/http"
  "CoffeeMap/internal/store"
  "strconv"
  "errors"
  "github.com/go-chi/chi/v5"
  "context"
)

type ratingKey string
const ratingCtx ratingKey = "rating"

type CreateRatingPayload struct {
  Coffee   int  `json:"coffee" validate:"required,max=10"`
  Ambiance int  `json"ambiance" validate:"required,max=10"`
  Overall  int  `json:"overall" validate:"required,max=10"`
}
// CreateRating godoc
//
//	@Summary		Creates a rating 
//	@Description	Creates a rating
//	@Tags			rating
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateRatingPayload	true	"Rating payload"
//	@Success		201		{object}	store.Rating
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/ratings [POST]
func (app *application) createRatingHandler(w http.ResponseWriter, r *http.Request) {
  var payload CreateRatingPayload 
  if err := readJSON(w, r, &payload); err != nil {
    app.badRequestResponse(w, r, err)
  }

  if err := Validate.Struct(payload); err != nil {
    app.badRequestResponse(w, r, err)
    return
  }
  rating := &store.Rating{
    Coffee: payload.Coffee,
    Ambiance: payload.Ambiance,
    Overall : payload.Overall,
    UserID: 1, //change after auth
    ShopID: 2,
  }
  
  ctx := r.Context()

  if err := app.store.Ratings.Create(ctx, rating); err != nil {
    app.internalServerError(w, r, err)
    return
  }

  if err := app.jsonResponse(w, http.StatusCreated, rating); err != nil {
    app.internalServerError(w, r, err)
    return
  }
}

// GetRating godoc
//
//	@Summary		Fetches a rating 
//	@Description	Fetches a rating by ID
//	@Tags			ratings
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Rating ID"
//	@Success		200	{object}	store.Rating
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/ratings/{id} [get]
func (app *application) getRatingHandler(w http.ResponseWriter, r *http.Request) {
  rating := getRatingFromCtx(r)

  if err := app.jsonResponse(w, http.StatusOK, rating); err != nil {
    app.internalServerError(w, r, err)
    return
  }
}

// DeleteRating godoc
//
//	@Summary		Deletes a Rating 
//	@Description	Delete a Rating by ID
//	@Tags			ratings
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Rating ID"
//	@Success		204	{object}	string
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/ratings/{id} [delete]
func (app *application) deleteRatingHandler(w http.ResponseWriter, r *http.Request) {
  idParam := chi.URLParam(r, "ratingID")

  id, err := strconv.ParseInt(idParam, 10, 64)
  if err != nil {
    app.internalServerError(w, r, err)
    return
  }

  ctx := r.Context()

  if err := app.store.Ratings.Delete(ctx, id); err != nil {
    switch {
    case errors.Is(err, store.ErrNotFound):
      app.notFoundResponse(w, r, err)
    default:
      app.internalServerError(w, r, err)
    }
    return
  }

  w.WriteHeader(http.StatusNoContent)
}

type UpdateRatingPayload struct {
  Coffee   *int  `json:"coffee" validate:"required,max=10"`
  Ambiance *int  `json"ambiance" validate:"required,max=10"`
  Overall  *int  `json:"overall" validate:"required,max=10"`
}

// UpdateRating godoc
//
//	@Summary		Updates a Rating 
//	@Description	Updates a Rating by ID
//	@Tags			ratings
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Rating ID"
//	@Param			payload	body		UpdateRatingPayload	true	"Rating payload"
//	@Success		200		{object}	store.Rating
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/rating/{id} [patch]
func (app *application) updateRatingHandler (w http.ResponseWriter, r *http.Request) {
  rating := getRatingFromCtx(r)

  var payload UpdateRatingPayload
  if err := readJSON(w, r, &payload); err != nil {
    app.badRequestResponse(w, r, err)
    return
  }

  if err := Validate.Struct(payload); err != nil {
    app.badRequestResponse(w, r, err)
    return
  }
  
  if payload.Coffee != nil {
    rating.Coffee = *payload.Coffee
  }

  if payload.Ambiance != nil {
    rating.Ambiance = *payload.Ambiance
  }

  if payload.Overall != nil {
    rating.Overall = *payload.Ambiance
  }
  
  if err := app.store.Ratings.Update(r.Context(), rating); err != nil {
    app.internalServerError(w, r, err)
  }

  if err := app.jsonResponse(w, http.StatusOK, rating); err != nil {
    app.internalServerError(w, r, err)
    return
  }


}

func (app *application) ratingsContextMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  idParam := chi.URLParam(r, "ratingID")

  if idParam == "" {
            http.Error(w, "ratingID is required", http.StatusBadRequest)
            return
          }
  id, err := strconv.ParseInt(idParam, 10, 64)
  if err != nil {
    app.internalServerError(w, r, err) 
    return
  }

  ctx := r.Context()

  rating, err := app.store.Ratings.GetByID(ctx, id)
  if err != nil {
    switch {
    case errors.Is(err, store.ErrNotFound):
      app.notFoundResponse(w, r, err)
    default:
      app.internalServerError(w, r, err)
      }
    return
    }

    ctx = context.WithValue(ctx, ratingCtx, rating)
    next.ServeHTTP(w, r.WithContext(ctx))
 })

}

func getRatingFromCtx(r *http.Request) *store.Rating {
   rating, _ := r.Context().Value(ratingCtx).(*store.Rating)
   return rating
} 
