// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ginh

import (
	"math"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/kjose/jgmc/api/internal/goapi/db/dao"
)

const (
	paginParam = "p"
	nbPerPage  = 20
)

// Interface to implement in a resource to set the configuration of pagination for LIST requests
type PaginationAware interface {
	GetPaginationConfig() PaginationConfig
}

// Pagination config
type PaginationConfig struct {
	QueryParamName string
	NbPerPage      int
}

// Links sent in the request
type Links struct {
	First string `json:"first,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
	Last  string `json:"last,omitempty"`
}

// Create a PaginationConfig with default values
func NewPaginationConfig() PaginationConfig {
	return PaginationConfig{
		QueryParamName: paginParam,
		NbPerPage:      nbPerPage,
	}
}

// Returns pagination filters from a gin context (page number, etc...)
func (pc *PaginationConfig) GetPaginationFilterFromContext(c *gin.Context) *dao.PaginationFilter {
	page := pageNumberFromParam(c.Query(pc.QueryParamName))
	page--

	return &dao.PaginationFilter{
		Limit:  pc.NbPerPage,
		Offset: page * pc.NbPerPage,
	}
}

// Returns links filters from a gin context (next url, prev url...)
func (pc *PaginationConfig) GetLinksFromContext(c *gin.Context, totalCount int) *Links {
	page := pageNumberFromParam(c.Query(pc.QueryParamName))

	l := &Links{}
	pageMax := int(math.Ceil(float64(totalCount) / float64(pc.NbPerPage)))
	l.First = replaceInUrl(c.Request.URL, pc.QueryParamName, 1)
	if page > 1 {
		l.Prev = replaceInUrl(c.Request.URL, pc.QueryParamName, page-1)
	}
	if page < pageMax {
		l.Next = replaceInUrl(c.Request.URL, pc.QueryParamName, page+1)
	}
	l.Last = replaceInUrl(c.Request.URL, pc.QueryParamName, pageMax)
	return l
}

func pageNumberFromParam(p string) int {
	page, err := strconv.Atoi(p)
	if err != nil {
		return 1
	}

	if page < 1 {
		page = 1
	}
	return page
}

func replaceInUrl(url *url.URL, paramName string, newPageNumber int) string {
	q := url.Query()
	q.Set(paramName, strconv.Itoa(newPageNumber))
	url.RawQuery = q.Encode()
	return url.String()
}
