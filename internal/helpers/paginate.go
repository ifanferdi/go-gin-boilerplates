package helpers

import (
	"math"
	"strings"
)

// PaginateParam adalah parameter input untuk fungsi Paginate.
type PaginateParam[T any] struct {
	Page  int
	Limit int
	Total int64
	Data  []T
}

// PaginateResult adalah struktur response pagination standar (sama seperti PaginateResult di POS-nextjs-express).
type PaginateResult[T any] struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
	Data       []T   `json:"data"`
}

// Paginate menghitung total halaman dan mengembalikan PaginateResult[T].
func Paginate[T any](param PaginateParam[T]) PaginateResult[T] {
	page := param.Page
	if page < 1 {
		page = 1
	}

	limit := param.Limit
	if limit == 0 {
		limit = 10
	}

	totalPages := 1
	if limit > 0 {
		totalPages = int(math.Ceil(float64(param.Total) / float64(limit)))
		if totalPages == 0 {
			totalPages = 1
		}
	}

	data := param.Data
	if data == nil {
		data = make([]T, 0)
	}

	return PaginateResult[T]{
		Page:       page,
		Limit:      limit,
		Total:      param.Total,
		TotalPages: totalPages,
		Data:       data,
	}
}

// ParseOrderBy memvalidasi dan memformat parameter sorting (contoh: "id:desc" atau "name:asc" -> "id DESC").
func ParseOrderBy(orderBy string, allowedColumns map[string]string, defaultOrder string) string {
	if strings.TrimSpace(orderBy) == "" {
		return defaultOrder
	}

	parts := strings.Split(orderBy, ":")
	rawCol := strings.TrimSpace(parts[0])
	dbCol, allowed := allowedColumns[rawCol]
	if !allowed {
		return defaultOrder
	}

	direction := "ASC"
	if len(parts) > 1 && strings.EqualFold(strings.TrimSpace(parts[1]), "desc") {
		direction = "DESC"
	}

	return dbCol + " " + direction
}
