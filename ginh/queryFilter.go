// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ginh

import "gitlab.com/kjose/jgmc/api/internal/goapi/db/dao"

// Function type to implement to create custom filters handlers in a DAO
type QueryFilterFunc func(string, string, interface{}) dao.FilterFunc

// Type QueryFilterSet
type QueryFilterSet []QueryFilter

// Interface to implement in a resource to support filters
type QueryFilterAware interface {
	GetQueryFilterSet() QueryFilterSet
}

// QueryFilter in a QueryFilterSet
type QueryFilter struct {
	UrlParam     string
	Func         QueryFilterFunc
	Args         interface{}
	DefaultValue string
}

// Returns the queryfilter of url param
func (qfs QueryFilterSet) GetByParam(param string) *QueryFilter {
	for _, i := range qfs {
		if i.UrlParam == param {
			return &i
		}
	}
	return nil
}
