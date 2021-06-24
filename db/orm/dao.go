// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package orm

import (
	"strings"

	"gitlab.com/kjose/jgmc/api/internal/easyapi/db/dao"
	"gitlab.com/kjose/jgmc/api/internal/easyapi/utils"
	"gorm.io/gorm"
)

var (
	DAO *relationalDAO
	UOW *unitOfWork
)

// relationalDAO implements DAOInterface and allow to query on relational databases
type relationalDAO struct {
	IdentifierKey string
}

// Unit of work, the working resource
type unitOfWork struct {
	From interface{}
	To   interface{}
}

func NewRelationalDAO(identifierKey string) *relationalDAO {
	return &relationalDAO{
		IdentifierKey: identifierKey,
	}
}

func (rdao *relationalDAO) FindByFilter(dest interface{}, ff []dao.FilterFunc, pf *dao.PaginationFilter) (dao.DAOResultsInterface, error) {
	st := DB.Model(dest)
	stCtx := utils.NewContext().With("c", st)
	for _, f := range ff {
		stCtx = f(stCtx)
	}

	// result
	ret := &relationalDAOResults{}

	// Pagination
	if pf != nil {
		st.Count(&ret.totalCount)
		st = st.Limit(pf.Limit).Offset(pf.Offset)
	}

	r, err := st.Rows()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var list []interface{}
	for r.Next() {
		DB.ScanRows(r, dest)
		list = append(list, utils.CloneInterface(dest))
	}

	ret.r = list
	return ret, nil
}

func (rdao *relationalDAO) FindBy(dest interface{}, params map[string]string, pf *dao.PaginationFilter) (dao.DAOResultsInterface, error) {
	var ff []dao.FilterFunc
	for k, p := range params {
		ff = append(ff, ApplyExactFilter(k, p, nil))
	}
	return rdao.FindByFilter(dest, ff, pf)
}

func (rdao *relationalDAO) FindById(dest interface{}, id string) (dao.DAOResultInterface, error) {
	r := DB.First(dest, rdao.IdentifierKey+" = ?", id)
	if r.Error != nil {
		return nil, r.Error
	}
	ret := &relationalDAOResult{
		r: dest,
	}
	return ret, nil
}

func (rdao *relationalDAO) UpdateFromPrevious(from interface{}, to interface{}) (dao.DAOResultInterface, error) {
	UOW = &unitOfWork{
		From: from,
		To:   to,
	}
	r := DB.Model(from).Updates(to)
	if r.Error != nil {
		return nil, r.Error
	}
	ret := &relationalDAOResult{
		r: to,
	}
	return ret, nil
}

func (rdao *relationalDAO) Create(resource interface{}) (dao.DAOResultInterface, error) {
	r := DB.Create(resource)
	if r.Error != nil {
		return nil, r.Error
	}
	ret := &relationalDAOResult{
		r: resource,
	}
	return ret, nil
}

func (rdao *relationalDAO) DeleteById(resource interface{}, id string) error {
	r := DB.Where(rdao.IdentifierKey+" = ?", id).Delete(resource)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

type relationalDAOResult struct {
	r dao.S
}

func (r *relationalDAOResult) Get() dao.S {
	return r.r
}

type relationalDAOResults struct {
	r          dao.SS
	totalCount int64
}

func (r *relationalDAOResults) All() dao.SS {
	return r.r
}

func (r *relationalDAOResults) CountTotal() int {
	return int(r.totalCount)
}

// Query filter of type EXACT MATCH
func ApplyExactFilter(param string, value string, args interface{}) dao.FilterFunc {
	return func(s *utils.Context) *utils.Context {
		s.Get("c").(*gorm.DB).Where(param+" = ?", value)
		return s
	}
}

// Query filter of type LIKE MATCH
func ApplyLikeFilter(param string, value string, args interface{}) dao.FilterFunc {
	return func(s *utils.Context) *utils.Context {
		s.Get("c").(*gorm.DB).Where(param+" LIKE ?", "%"+value+"%")
		return s
	}
}

// Query filter to support ordering in query results
func ApplyOrderFilter(param string, value string, args interface{}) dao.FilterFunc {
	allowedFilters := args.([]string)
	return func(s *utils.Context) *utils.Context {
		isdesc := ""
		if strings.Index(value, "-") == 0 {
			isdesc = " DESC"
		}
		replacer := strings.NewReplacer("+", "", "-", "")
		val := strings.TrimSpace(replacer.Replace(value))

		for _, f := range allowedFilters {
			if f == val {
				s.Get("c").(*gorm.DB).Order("`" + val + "`" + isdesc)
			}
		}
		return s
	}
}

func init() {
	DAO = NewRelationalDAO("id")
}
