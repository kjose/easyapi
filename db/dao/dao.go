// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dao

import "gitlab.com/kjose/jgmc/api/internal/easyapi/utils"

// Default DAO of application, it is used when a resource has no custom DAO configured
var defaultDAO DAOInterface

// DAOInterface is is the interface to implement if you need to have a custom database manager, by default you can use the dao in orm folder.
type DAOInterface interface {
	// FindByFilter is the global function to get multiple results from the database, with custom filters and pagination.
	FindByFilter(dest interface{}, ff []FilterFunc, pf *PaginationFilter) (DAOResultsInterface, error)
	// FindBy is the shortcut function of FindByFilter with simple filters
	FindBy(dest interface{}, params map[string]string, pf *PaginationFilter) (DAOResultsInterface, error)
	// FindById is the function to get a single result by its ID
	FindById(dest interface{}, id string) (DAOResultInterface, error)
	// UpdateFromPrevious is the function to update a resource from a previous resource
	UpdateFromPrevious(from interface{}, to interface{}) (DAOResultInterface, error)
	// Create is the function to create a resource
	Create(resource interface{}) (DAOResultInterface, error)
	// DeleteById is the function to delete a resource by its id
	DeleteById(resource interface{}, id string) error
}

// Interface of a single results
type S interface{}

// Interface of a multiple results
type SS []interface{}

// Interface to get a DAOResultInterface
type DAOResultInterface interface {
	Get() S
}

// Interface to get a DAOResultsInterface
type DAOResultsInterface interface {
	All() SS
	CountTotal() int
}

// Interface to implement in the resource when you want to use a custom DAO
type GetDAOInterface interface {
	GetDAO() DAOInterface
}

// Init the default DAO of application
func InitDefaultDAO(dao DAOInterface) {
	defaultDAO = dao
}

// Get the DAO of a resource, or the default one if no one is configured
func GetResourceDAO(resource interface{}) DAOInterface {
	if rdo, ok := resource.(GetDAOInterface); ok {
		return rdo.GetDAO()
	}
	if defaultDAO == nil {
		panic("orm not initialized, default dao is missing")
	}
	return defaultDAO
}

// Type function to create to handle filters, it starts from a Statement and returns the same updated Statement
type FilterFunc func(s *utils.Context) *utils.Context

// Pagination filter is a model to use to handle pagination results
type PaginationFilter struct {
	Limit  int
	Offset int
}
