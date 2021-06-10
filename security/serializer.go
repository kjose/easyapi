// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package security

const (
	SERIALIZER_CONTEXT_KEY_ONE  = "one"
	SERIALIZER_CONTEXT_KEY_LIST = "list"
)

// Interface to implement in a resource to support serialization
type SerializeAware interface {
	Serialize(sc *SerializeGroups) interface{}
}

// Use a SerializeGroups to handle the display or not of some resource information based on different context values
type SerializeGroups struct {
	Values []string
}

// Returns true if the SerializeGroups contains a value
func (sc *SerializeGroups) Contains(value string) bool {
	for _, v := range sc.Values {
		if v == value {
			return true
		}
	}
	return false
}

// Create a new SerializeGroups with new value
func (sc *SerializeGroups) WithValue(value string) *SerializeGroups {
	if sc == nil {
		sc = &SerializeGroups{}
	}
	sc.Values = append(sc.Values, value)
	return sc
}

// Serialize a resource (aware of Interface or not)
func Serialize(i interface{}, sc *SerializeGroups) interface{} {
	if is, ok := i.(SerializeAware); ok {
		return is.Serialize(sc)
	}
	return i
}
