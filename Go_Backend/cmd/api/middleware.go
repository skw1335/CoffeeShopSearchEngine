package main

import (
	"fmt"
	"net/http"
	"strings"
	"strconv"
	"github.com/golang-jwt/jwt/v5"
	"context"
	"CoffeeMap/internal/store"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
				app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
				return
			}
		parts := strings.Split(authHeader, " ") 
			if len(parts) != 2 || parts[0] != "Bearer" {
				app.unauthorizedErrorResponse(w, r, fmt.Errorf("malformed authorization Header"))
				return
			}

			token := parts[1]
			jwtToken, err := app.authenticator.ValidateToken(token)
			if err != nil {
				app.unauthorizedErrorResponse(w, r, err)
				return
			}

			claims := jwtToken.Claims.(jwt.MapClaims)

			userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
			if err != nil {
				app.unauthorizedErrorResponse(w, r, err)
				return

			}

			ctx := r.Context()

			user, err := app.store.Users.GetByID(r.Context(), userID)
			if err != nil {
				app.unauthorizedErrorResponse(w, r, err)
				return
			}

			ctx = context.WithValue(ctx, userCtx, user)
			next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		user := getUserFromContext(r)
		post := getCommentFromContext(r)

		// check if it is the users post
		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		//role precedence check
	allowed, err := app.checkRolePrecedence(r.Context(), user, requiredRole) 
	if err != nil {
		app.internalServerError(w,r,err)
	} 
	
	if !allowed {
		app.forbiddenResponse(w,r)
	}
	next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil

}

