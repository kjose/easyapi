// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package easyapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/kjose/jgmc/api/internal/easyapi/db/dao"
	"gitlab.com/kjose/jgmc/api/internal/easyapi/event"
	"gitlab.com/kjose/jgmc/api/internal/easyapi/layer"
	"gitlab.com/kjose/jgmc/api/internal/easyapi/security"
	"gitlab.com/kjose/jgmc/api/internal/easyapi/utils"
)

// Gin handler for a POST request
func HandlePost(c *gin.Context, i interface{}) {
	ic := utils.CloneInterface(i) // avoid duplicate variable use
	if err := BindAndValidate(c, ic); err != nil {
		return
	}

	err := event.DispatchEvent(c, event.EVENT_RESOURCE_PRE_CREATE, &event.ResourceActionEvent{
		Resource: ic,
		Action:   event.EVENT_RESOURCE_PRE_CREATE,
	})
	if err != nil {
		return
	}

	RemoveUUIDBindings(ic)

	_, err = dao.GetResourceDAO(ic).Create(ic)
	if err != nil {
		HttpError(c, http.StatusBadRequest, "Creation error", nil)
		return
	}

	err = event.DispatchEvent(c, event.EVENT_RESOURCE_POST_CREATE, &event.ResourceActionEvent{
		Resource: ic,
		Action:   event.EVENT_RESOURCE_POST_CREATE,
	})
	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, NewItem(ic, &security.SerializeGroups{
		Values: []string{security.SERIALIZER_CONTEXT_KEY_ONE},
	}))
}

// Gin handler for a GET request
func HandleGet(c *gin.Context, i interface{}, id string) {
	ic := utils.CloneInterface(i) // avoid duplicate variable use
	_, err := dao.GetResourceDAO(ic).FindById(ic, id)
	if err != nil {
		HttpError(c, http.StatusNotFound, "Not found", nil)
		return
	}

	if err := AppendBindings(ic); err != nil {
		return
	}

	err = event.DispatchEvent(c, event.EVENT_RESOURCE_POST_READ, &event.ResourceActionEvent{
		Resource: ic,
		Action:   event.EVENT_RESOURCE_POST_READ,
	})
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, NewItem(ic, &security.SerializeGroups{
		Values: []string{security.SERIALIZER_CONTEXT_KEY_ONE},
	}))
}

// Gin handler for a LIST request
func HandleList(c *gin.Context, i interface{}) {
	ic := utils.CloneInterface(i) // avoid duplicate variable use

	// Pagination
	var pf *dao.PaginationFilter
	var pQueryName string
	var pc layer.PaginationConfig
	if ipa, ok := i.(layer.PaginationAware); ok {
		pc = ipa.GetPaginationConfig()
		pQueryName = pc.QueryParamName
		pf = pc.GetPaginationFilterFromContext(c)
	}

	// Check filters from
	var ff []dao.FilterFunc
	if iqfa, ok := ic.(layer.QueryFilterAware); ok {
		for key, val := range c.Request.URL.Query() {
			if key == pQueryName {
				continue
			}
			qf := iqfa.GetQueryFilterSet().GetByParam(key)
			if qf == nil {
				HttpError(c, http.StatusNotFound, fmt.Sprintf("Param %s is not a filter", key), nil)
				return
			}
			ff = append(ff, qf.Func(key, val[0], qf.Args))
		}
		// Apply defaults
		for _, f := range iqfa.GetQueryFilterSet() {
			if f.DefaultValue == "" || c.Request.URL.Query().Get(f.UrlParam) != "" {
				continue
			}
			ff = append(ff, f.Func(f.UrlParam, f.DefaultValue, f.Args))
		}
	}

	r, err := dao.GetResourceDAO(ic).FindByFilter(ic, ff, pf)
	if err != nil {
		HttpError(c, http.StatusNotFound, "Get collection request error", nil)
		return
	}

	all := r.All()
	for _, l := range all {
		err = event.DispatchEvent(c, event.EVENT_RESOURCE_POST_READ, &event.ResourceActionEvent{
			Resource: l,
			Action:   event.EVENT_RESOURCE_POST_READ,
		})
		if err != nil {
			return
		}
	}

	collectionItems := NewCollectionItem(all, &security.SerializeGroups{
		Values: []string{security.SERIALIZER_CONTEXT_KEY_LIST},
	})
	collectionItems.Count = len(all)
	collectionItems.Total = r.CountTotal()
	collectionItems.Links = pc.GetLinksFromContext(c, collectionItems.Total)
	c.JSON(http.StatusOK, collectionItems)
}

// Gin handler for a PATCH request
func HandlePatch(c *gin.Context, i interface{}, id string) {
	ic := utils.CloneInterface(i) // avoid duplicate variable use
	_, err := dao.GetResourceDAO(ic).FindById(ic, id)
	clone := utils.CloneInterface(ic)
	if err != nil {
		HttpError(c, http.StatusNotFound, "Not found", nil)
		return
	}

	if err := BindAndValidate(c, ic); err != nil {
		return
	}

	err = event.DispatchEvent(c, event.EVENT_RESOURCE_PRE_UPDATE, &event.ResourceActionEvent{
		Resource: ic,
		Action:   event.EVENT_RESOURCE_PRE_UPDATE,
	})
	if err != nil {
		return
	}

	RemoveUUIDBindings(ic)

	_, err = dao.GetResourceDAO(ic).UpdateFromPrevious(clone, ic)
	if err != nil {
		HttpError(c, http.StatusBadRequest, "Update error", nil)
		return
	}

	err = event.DispatchEvent(c, event.EVENT_RESOURCE_POST_UPDATE, &event.ResourceActionEvent{
		Resource: ic,
		Action:   event.EVENT_RESOURCE_POST_UPDATE,
	})
	if err != nil {
		return
	}

}

// Gin handler for a DELETE request
func HandleDelete(c *gin.Context, i interface{}, id string) {
	ic := utils.CloneInterface(i)

	_, err := dao.GetResourceDAO(ic).FindById(ic, id)
	if err != nil {
		HttpError(c, http.StatusNotFound, "Not found", nil)
		return
	}

	if err := AppendBindings(ic); err != nil {
		return
	}

	err = event.DispatchEvent(c, event.EVENT_RESOURCE_PRE_DELETE, &event.ResourceActionEvent{
		Resource: ic,
		Action:   event.EVENT_RESOURCE_PRE_DELETE,
	})
	if err != nil {
		return
	}

	err = dao.GetResourceDAO(ic).DeleteById(ic, id)
	if err != nil {
		HttpError(c, http.StatusBadRequest, "Delete error", nil)
		return
	}
}

// Shortcut to handle multiple crud requests
func CRUDL(r gin.IRoutes, path string, i interface{}, methods string) {
	if methods == "" {
		methods = "CRUDL"
	}
	if strings.Contains(methods, "C") {
		r.POST(path, func(c *gin.Context) {
			HandlePost(c, i)
		})
	}
	if strings.Contains(methods, "R") {
		r.GET(path+"/:id", func(c *gin.Context) {
			HandleGet(c, i, c.Param("id"))
		})
	}
	if strings.Contains(methods, "U") {
		r.PATCH(path+"/:id", func(c *gin.Context) {
			HandlePatch(c, i, c.Param("id"))
		})
	}
	if strings.Contains(methods, "D") {
		r.DELETE(path+"/:id", func(c *gin.Context) {
			HandleDelete(c, i, c.Param("id"))
		})
	}
	if strings.Contains(methods, "L") {
		r.GET(path, func(c *gin.Context) {
			HandleList(c, i)
		})
	}
}
