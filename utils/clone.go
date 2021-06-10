// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utils

import "reflect"

func CloneInterface(i interface{}) interface{} {
	n := reflect.New(reflect.TypeOf(i).Elem())
	oldVal := reflect.ValueOf(i).Elem()
	newVal := n.Elem()
	for i := 0; i < oldVal.NumField(); i++ {
		newValField := newVal.Field(i)
		if newValField.CanSet() {
			newValField.Set(oldVal.Field(i))
		}
	}

	return n.Interface()
}
