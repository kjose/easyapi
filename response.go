// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package easyapi

import (
	"gitlab.com/kjose/jgmc/api/internal/easyapi/layer"
)

// Create a new item single response
func NewItem(i interface{}, sc *layer.SerializeGroups) interface{} {
	return Serialize(i, sc)
}

type CollectonItem struct {
	Items []interface{} `json:"items"`
	Count int           `json:"count,omitempty"`
	Total int           `json:"total,omitempty"`
	Links *layer.Links  `json:"_links,omitempty"`
}

// Create a new item collection response
func NewCollectionItem(items []interface{}, sc *layer.SerializeGroups) *CollectonItem {
	collection := make([]interface{}, 0)
	for _, i := range items {
		Serialize(i, sc)
		collection = append(collection, i)
	}
	return &CollectonItem{
		Items: collection,
	}
}
