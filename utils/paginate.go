package utils

import (
	"log"
	"math"
	"net/http"

	"gorm.io/gorm"
)

type Pagination struct {
	Page         int         `json:"page"`
	PerPage      int         `json:"per_page"`
	Sort         string      `json:"sort"`
	TotalRecords int         `json:"total_records"`
	TotalPages   int         `json:"total_pages"`
	Remaining    int         `json:"remaining"`
	Items        interface{} `json:"items"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPerPage()
}

func (p *Pagination) GetPerPage() int {
	if p.PerPage == 0 {
		p.PerPage = 10
	}
	return p.PerPage
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "created_at asc"
	}
	return p.Sort
}

// Paginate performs pagination on the given database query, based on the page and per_page parameters in the provided HTTP request.
// It returns a PaginatedResult struct containing the paginated items, as well as metadata such as the total number of records, total number of pages, and the number of remaining items.
// The out parameter should be a pointer to a slice that will hold the paginated items.
// If any errors occur during pagination, an error is returned.
func Paginate(r *http.Request, db *gorm.DB, out interface{}, response *APIResponse, associations ...string) (func(db *gorm.DB) *gorm.DB, error) {
	page, perPage := GetPaginationParams(r, 1, 10)

	// calculate the offset based on the current page and number of items per page
	offset := (page - 1) * perPage

	// get the total number of records
	var count int64
	if err := db.Model(out).Count(&count).Error; err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// calculate the total number of pages based on the total number of records and items per page
	totalPages := int(math.Ceil(float64(count) / float64(perPage)))

	// calculate the number of remaining items
	remaining := int(count) - (page * perPage)
	if remaining < 0 {
		remaining = 0
	}

	// Update the pagination struct with the values
	// pagination.Page = page
	// pagination.PerPage = perPage
	// pagination.Sort = pagination.GetSort()
	// pagination.TotalRecords = int(count)
	// pagination.TotalPages = totalPages
	// pagination.Remaining = remaining
	// Update the pagination struct with the values
	response.Meta = map[string]interface{}{
		"page":          page,
		"per_page":      perPage,
		"sort":          "created_at asc",
		"total_records": count,
		"total_pages":   totalPages,
		"remaining":     remaining,
	}
	scopeFunc := func(db *gorm.DB) *gorm.DB {
		// preload specified relationships
		for _, association := range associations {
			db = db.Preload(association)
		}

		// Apply pagination scopes
		return db.Offset(offset).Limit(perPage).Order(response.Meta["sort"])
	}

	return scopeFunc, nil
}
