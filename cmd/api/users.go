package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/turut4/social/internal/store"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserByCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}

			return
		}

		ctx = context.WithValue(ctx, "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserByCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value("user").(*store.User)
	return user
}
