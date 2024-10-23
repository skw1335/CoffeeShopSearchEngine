package main

import (
  "net/http"
  "CoffeeMap/internal/store"
  "strconv"
  "errors"
  "github.com/go-chi/chi/v5"
  "context"
)

type commentKey string
const commentCtx commentKey = "comment"

type CreateCommentPayload struct {
  Title   string  `json:"title" validate:"required,max=100"`
  Content string  `json:"content" validate:"required,max=1000"`
}
// CreateComment godoc
//
//	@Summary		Creates a comment 
//	@Description	Creates a comment 
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateCommentPayload	true	"Comment payload"
//	@Success		201		{object}	store.Comment
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/comments [POST]
func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
  var payload CreateCommentPayload 
  if err := readJSON(w, r, &payload); err != nil {
    app.badRequestResponse(w, r, err)
  }

  if err := Validate.Struct(payload); err != nil {
    app.badRequestResponse(w, r, err)
    return
  }
  comment := &store.Comment{
    Title: payload.Title,
    Content: payload.Content,
    UserID: 1, //change after auth
    ShopID: 2,
  }
  
  ctx := r.Context()

  if err := app.store.Comments.Create(ctx, comment); err != nil {
    app.internalServerError(w, r, err)
    return
  }

  if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
    app.internalServerError(w, r, err)
    return
  }
}

// GetComment godoc
//
//	@Summary		Fetches a comment 
//	@Description	Fetches a comment by ID
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Comment ID"
//	@Success		200	{object}	store.Comment
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/comments/{id} [GET]
func (app *application) getCommentHandler(w http.ResponseWriter, r *http.Request) {
  comment := getCommentFromContext(r)

  if err := app.jsonResponse(w, http.StatusOK, comment); err != nil {
    app.internalServerError(w, r, err)
    return
  }
}

// DeleteComment godoc
//
//	@Summary		Deletes a Comment 
//	@Description	Delete a Comment by ID
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Comment ID"
//	@Success		204	{object}	string
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/comments/{id} [DELETE]
func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
  idParam := chi.URLParam(r, "commentID")

  id, err := strconv.ParseInt(idParam, 10, 64)
  if err != nil {
    app.internalServerError(w, r, err)
    return
  }

  ctx := r.Context()

  if err := app.store.Comments.Delete(ctx, id); err != nil {
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

type UpdateCommentPayload struct {
  Title     *string    `json:"title" validate:"omitempty,max=100"`
  Content   *string    `json:"content" validate:"omitempty,max=1000"`
}

// UpdateComment godoc
//
//	@Summary		Updates a Comment 
//	@Description	Updates a Comment by ID
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Comment ID"
//	@Param			payload	body		UpdateCommentPayload	true	"Comment payload"
//	@Success		200		{object}	store.Comment
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/comments/{id} [patch]
func (app *application) updateCommentHandler (w http.ResponseWriter, r *http.Request) {
  comment := getCommentFromContext(r)

  var payload UpdateCommentPayload
  if err := readJSON(w, r, &payload); err != nil {
    app.badRequestResponse(w, r, err)
    return
  }

  if err := Validate.Struct(payload); err != nil {
    app.badRequestResponse(w, r, err)
    return
  }
  
  if payload.Content != nil {
    comment.Content = *payload.Content
  }

  if payload.Title != nil {
    comment.Title = *payload.Title
  }
  
  if err := app.store.Comments.Update(r.Context(), comment); err != nil {
    app.internalServerError(w, r, err)
  }

  if err := app.jsonResponse(w, http.StatusOK, comment); err != nil {
    app.internalServerError(w, r, err)
    return
  }


}

func (app *application) commentsContextMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  idParam := chi.URLParam(r, "commentID")

  if idParam == "" {
            http.Error(w, "commentID is required", http.StatusBadRequest)
            return
          }
  id, err := strconv.ParseInt(idParam, 10, 64)
  if err != nil {
    app.internalServerError(w, r, err) 
    return
  }

  ctx := r.Context()

  comment, err := app.store.Comments.GetByID(ctx, id)
  if err != nil {
    switch {
    case errors.Is(err, store.ErrNotFound):
      app.notFoundResponse(w, r, err)
    default:
      app.internalServerError(w, r, err)
      }
    return
    }

    ctx = context.WithValue(ctx, commentCtx, comment)
    next.ServeHTTP(w, r.WithContext(ctx))
 })

}

func getCommentFromContext(r *http.Request) *store.Comment {
   comment, _ := r.Context().Value(commentCtx).(*store.Comment)
   return comment
 }


func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
   type envelope struct {
     Data any `json:"data"`
   }
   return writeJSON(w, status, &envelope{Data: data})
 }
