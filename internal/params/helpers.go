package params

import (
	"go-sober/internal/constants"
	"net/http"
	"strconv"
	"time"
)

// Helper functions
func ParsePaginationParams(r *http.Request) (page, pageSize int) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err = strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = constants.DefaultPageSize
	}
	if pageSize > constants.MaxPageSize {
		pageSize = constants.MaxPageSize
	}

	return page, pageSize
}

func ParseTimeParam(param string) *time.Time {
	if param == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, param)
	if err != nil {
		return nil
	}
	return &t
}

func ParseFloatParam(param string) *float64 {
	if param == "" {
		return nil
	}
	f, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return nil
	}
	return &f
}
