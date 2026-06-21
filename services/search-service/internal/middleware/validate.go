package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Tejas-Panchal/sravas/services/search-service/internal/model"
)

// ParseSearchRequest extracts and validates search parameters from query string
func ParseSearchRequest(r *http.Request) *model.SearchRequest {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	category := r.URL.Query().Get("category")
	sort := r.URL.Query().Get("sort")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	if sort == "" {
		sort = "relevance"
	}

	return &model.SearchRequest{
		Query:    q,
		Category: category,
		Sort:     sort,
		Page:     page,
		PageSize: pageSize,
	}
}
