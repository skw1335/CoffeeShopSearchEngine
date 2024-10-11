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


func (app *application) getCommentHandler(w http.ResponseWriter, r *http.Request) {
  comment := getCommentFromCtx(r)

  if err := app.jsonResponse(w, http.StatusOK, comment); err != nil {
    app.internalServerError(w, r, err)
    return
  }
}

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

func (app *application) updateCommentHandler (w http.ResponseWriter, r *http.Request) {
  comment := getCommentFromCtx(r)

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

func getCommentFromCtx(r *http.Request) *store.Comment {
   comment, _ := r.Context().Value(commentCtx).(*store.Comment)
   return comment
 }


func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
   type envelope struct {
     Data any `json:"data"`
   }
   return writeJSON(w, status, &envelope{Data: data})
 }
