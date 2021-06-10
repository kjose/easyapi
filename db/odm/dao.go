// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package odm

import (
	"reflect"
	"strings"

	"github.com/go-bongo/bongo"
	"gitlab.com/kjose/jgmc/api/internal/goapi/db/dao"
	"gitlab.com/kjose/jgmc/api/internal/goapi/utils"
	"gopkg.in/mgo.v2/bson"
)

var (
	DAO *nosqlDAO
)

type nosqlDAO struct {
	IdentifierKey string
}

type daoResult struct {
	r dao.S
}

func (r *daoResult) Get() dao.S {
	return r.r
}

type daoResults struct {
	r          dao.SS
	totalCount int
}

type statement struct {
	Collection *bongo.Collection
	Filters    bson.M
}

func (r *daoResults) All() dao.SS {
	return r.r
}

func (r *daoResults) CountTotal() int {
	return int(r.totalCount)
}

func NewNosqlDAO(identifierKey string) *nosqlDAO {
	return &nosqlDAO{
		IdentifierKey: identifierKey,
	}
}

func (n *nosqlDAO) FindByFilter(dest interface{}, ff []dao.FilterFunc, pf *dao.PaginationFilter) (dao.DAOResultsInterface, error) {
	st := &statement{
		Collection: DB.Collection(getCollectionName(dest)),
		Filters:    bson.M{},
	}
	stCtx := utils.NewContext().With("s", st)
	for _, f := range ff {
		stCtx = f(stCtx)
	}

	results := st.Collection.Find(st.Filters)
	count, err := st.Collection.Find(st.Filters).Query.Count()
	if err != nil {
		return nil, err
	}

	if pf != nil {
		results.Query.Sort("_id")
		results.Query.Skip(pf.Offset)
		results.Query.Limit(pf.Limit)
	}

	var list []interface{}
	for results.Next(dest) {
		list = append(list, utils.CloneInterface(dest))
	}

	return &daoResults{
		r:          list,
		totalCount: count,
	}, nil
}

func (n *nosqlDAO) FindBy(dest interface{}, params map[string]string, pf *dao.PaginationFilter) (dao.DAOResultsInterface, error) {
	var ff []dao.FilterFunc
	for k, p := range params {
		ff = append(ff, ApplyExactFilter(k, p, nil))
	}
	return n.FindByFilter(dest, ff, pf)
}

func (n *nosqlDAO) FindById(dest interface{}, id string) (dao.DAOResultInterface, error) {
	err := DB.Collection(getCollectionName(dest)).FindById(bson.ObjectIdHex(id), dest)
	if err != nil {
		return nil, err
	}

	return &daoResult{
		r: dest,
	}, nil
}

func (n *nosqlDAO) UpdateFromPrevious(from interface{}, to interface{}) (dao.DAOResultInterface, error) {
	err := DB.Collection(getCollectionName(to)).Save(to.(bongo.Document))
	if err != nil {
		return nil, err
	}

	return &daoResult{
		r: to,
	}, nil
}

func (n *nosqlDAO) Create(resource interface{}) (dao.DAOResultInterface, error) {
	err := DB.Collection(getCollectionName(resource)).Save(resource.(bongo.Document))
	if err != nil {
		return nil, err
	}

	return &daoResult{
		r: resource,
	}, nil
}

func (n *nosqlDAO) DeleteById(resource interface{}, id string) error {
	err := DB.Collection(getCollectionName(resource)).DeleteDocument(resource.(bongo.Document))
	if err != nil {
		return err
	}
	return nil
}

func getCollectionName(resource interface{}) string {
	split := strings.Split(reflect.TypeOf(resource).String(), ".")
	collection := split[len(split)-1] + "s"
	return strings.ToLower(collection)
}

// Query filter of type EXACT MATCH
func ApplyExactFilter(param string, value string, args interface{}) dao.FilterFunc {
	return func(s *utils.Context) *utils.Context {
		filters := s.Get("s").(*statement).Filters
		filters[param] = value
		return s
	}
}

func init() {
	DAO = NewNosqlDAO("id")
}
