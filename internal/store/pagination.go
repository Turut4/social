package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = time.RFC3339

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=0,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"omitempty,max=100"`
	Since  *string   `json:"since" validate:"omitempty"`
	Until  *string   `json:"until" validate:"omitempty"`
}

func NewPaginatedQuery() PaginatedFeedQuery {
	return PaginatedFeedQuery{
		Limit: 20,
		Sort:  "desc",
		Tags:  make([]string, 0),
	}
}

func (fq PaginatedFeedQuery) ParseFromRequest(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	if limit := qs.Get("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, err
		}

		fq.Limit = l
	}

	if offset := qs.Get("offset"); offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, err
		}

		fq.Offset = o
	}

	if sort := qs.Get("sort"); sort != "" {
		fq.Sort = sort
	}

	if tags := qs.Get("tags"); tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	if search := qs.Get("search"); search != "" {
		fq.Search = search
	}

	if since := qs.Get("since"); since != "" {
		s := parseTime(since)
		fq.Since = &s
	}

	if until := qs.Get("until"); until != "" {
		u := parseTime(until)
		fq.Since = &u
	}

	return fq, nil
}

func parseTime(s string) string {
	t, err := time.Parse(dateFormat, s)
	if err != nil {
		return ""
	}

	return t.Format(dateFormat)
}
