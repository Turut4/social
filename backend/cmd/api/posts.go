package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/turut4/social/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=200"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type postKey string

const postCtx postKey = "post"

// CreatePost gdoc
//
//	@Summary		Create a new post
//	@Description	Create a new post with title, content, and optional tags
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreatePostPayload	true	"Post payload"
//	@Success		201		{object}	store.Post			"Created post"
//	@Failure		400		{object}	error				"Invalid input"
//	@Failure		500		{object}	error				"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := getUserFromContext(r)
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  user.ID,
	}

	ctx := r.Context()
	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetPost gdoc
//
//	@Summary		Get a post by ID
//	@Description	Retrieve a post and its comments by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int			true	"Post ID"
//	@Success		200		{object}	store.Post	"Post details with comments"
//	@Failure		404		{object}	error		"Post not found"
//	@Failure		500		{object}	error		"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

type updatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

// UpdatePost gdoc
//
//	@Summary		Update a post
//	@Description	Update the title or content of a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int					true	"Post ID"
//	@Param			request	body		updatePostPayload	true	"Update payload"
//	@Success		200		{object}	store.Post			"Updated post"
//	@Failure		400		{object}	error				"Invalid input"
//	@Failure		404		{object}	error				"Post not found"
//	@Failure		500		{object}	error				"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)
	var payload updatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		switch err {
		case store.ErrNotFound:
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

// DeletePost gdoc
//
//	@Summary		Delete a post
//	@Description	Delete a post by its ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int		true	"Post ID"
//	@Success		204		{string}	string	"Post deleted"
//	@Failure		404		{object}	error	"Post not found"
//	@Failure		500		{object}	error	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.store.Posts.Delete(ctx, id); err != nil {
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

// CreateComment gdoc
//
//	@Summary		Add a comment to a post
//	@Description	Add a comment to a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int						true	"Post ID"
//	@Param			request	body		createCommentPayload	true	"Comment payload"
//	@Success		201		{object}	store.Comment			"Created comment"
//	@Failure		400		{object}	error					"Invalid input"
//	@Failure		500		{object}	error					"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID}/comments [post]
func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		post, err := app.store.Posts.GetByID(ctx, postID)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromContext(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}

type createCommentPayload struct {
	Content string `json:"content" validate:"required,max=255"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)
	var payload createCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := getUserFromContext(r)
	comment := &store.Comment{
		UserID:  user.ID,
		PostID:  post.ID,
		Content: payload.Content,
	}

	if err := app.store.Comments.Create(r.Context(), comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
	}
}
