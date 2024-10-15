package main

import (
  "net/http"
  "CoffeeMap/internal/store"
  "strconv"
  "errors"
  "github.com/go-chi/chi/v5"
  "context"
)

type shopKey string
const shopCtx shopKey = "shop"



func (app *application) shopsContextMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  idParam := chi.URLParam(r, "shopID")
  id, err := strconv.ParseInt(idParam, 10, 64)
  if err != nil {
    app.internalServerError(w, r, err) 
    return
  }

  ctx := r.Context()

  shop, err := app.store.Shops.GetByID(ctx, id)
  if err != nil {
    switch {
    case errors.Is(err, store.ErrNotFound):
      app.notFoundResponse(w, r, err)
    default:
      app.internalServerError(w, r, err)
      }
    return
    }

    ctx = context.WithValue(ctx, shopCtx, shop)
    next.ServeHTTP(w, r.WithContext(ctx))
 })

}

// GetShop godoc
//
//	@Summary		Fetches a shop
//	@Description	Fetches a shop by ID
//	@Tags			> shop
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Shop ID"
//	@Success		200	{object}	store.Shop
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/shops/{id} [get]
func (app *application) getShopHandler(w http.ResponseWriter, r *http.Request) {
  shop := getShopFromCtx(r)

  if err := app.jsonResponse(w, http.StatusOK, shop); err != nil {
    app.internalServerError(w, r, err)
    return
  }
}

func getShopFromCtx(r *http.Request) *store.Shop {
   shop, _ := r.Context().Value(shopCtx).(*store.Shop)
   return shop
 }

