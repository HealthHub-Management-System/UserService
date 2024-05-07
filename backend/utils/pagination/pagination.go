package pagination

import (
	"math"
	"net/url"
	"strconv"

	"gorm.io/gorm"
)

type Pagination struct {
	Limit      int `form:"limit"`
	Page       int `form:"page"`
	TotalRows  int64
	TotalPages int
	Rows       any
}

func (p *Pagination) Parse(query url.Values) {
	if query.Has("page") {
		pageStr := query.Get("page")
		page, _ := strconv.Atoi(pageStr)
		p.Page = page
	} else {
		p.Page = 1
	}

	if query.Has("limit") {
		limitStr := query.Get("limit")
		limit, _ := strconv.Atoi(limitStr)
		p.Limit = limit
	} else {
		p.Limit = 10
	}
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func Paginate(value any, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order("id desc")
	}
}
