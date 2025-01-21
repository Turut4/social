package main

import (
	"net/http"

	"github.com/turut4/social/internal/store"
)

type createCommentPayload struct {
	Content string `json:"content" validate:"required,max=255"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostByCtx(r)
	var payload createCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := &store.Comment{
		UserID:  1, //TODO
		PostID:  post.ID,
		Content: payload.Content,
	}

	if err := app.store.Comments.Create(r.Context(), comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
