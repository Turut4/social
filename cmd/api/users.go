package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/turut4/social/internal/store"
)

type userKey string

const userCtx userKey = "user"

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserByCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

type followUser struct {
	UserID int64 `json:"user_id"`
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser := getUserByCtx(r)

	//revert back to auth userID from ctx
	var payload followUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Followers.Follow(r.Context(), payload.UserID, followedUser.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}

}

// UnfollowUser gdoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unfollowedUser := getUserByCtx(r)

	//revert back to auth userID from ctx
	var payload followUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Followers.Unfollow(r.Context(), payload.UserID, unfollowedUser.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
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

		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserByCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}

// userActivateHandler godoc
//
//	@Summary		Activate user
//	@Description	Activate user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"Invitation Token"
//	@Success		204		{string}	string "user activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) actvateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
	}
}
