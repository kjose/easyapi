// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package easyapi

import "gitlab.com/kjose/jgmc/api/internal/easyapi/layer"

const (
	SERIALIZER_CONTEXT_KEY_ONE  = "one"
	SERIALIZER_CONTEXT_KEY_LIST = "list"
)

// Serialize a resource (aware of Interface or not)
func Serialize(i interface{}, sc *layer.SerializeGroups) interface{} {
	if is, ok := i.(layer.SerializeAware); ok {
		return is.Serialize(sc)
	}
	return i
}
