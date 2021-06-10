// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ginh

import "gitlab.com/kjose/jgmc/api/internal/goapi/security"

// Create a new item single response
func NewItem(i interface{}, sc *security.SerializeGroups) interface{} {
	return security.Serialize(i, sc)
}

type CollectonItem struct {
	Items []interface{} `json:"items"`
	Count int           `json:"count,omitempty"`
	Total int           `json:"total,omitempty"`
	Links *Links        `json:"_links,omitempty"`
}

// Create a new item collection response
func NewCollectionItem(items []interface{}, sc *security.SerializeGroups) *CollectonItem {
	collection := make([]interface{}, 0)
	for _, i := range items {
		security.Serialize(i, sc)
		collection = append(collection, i)
	}
	return &CollectonItem{
		Items: collection,
	}
}
