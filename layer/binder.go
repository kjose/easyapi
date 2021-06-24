// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package layer

import "github.com/google/uuid"

// Object to return in resources to configure bindings of the resource
type UUIDBinding struct {
	UUID   *uuid.UUID
	BindTo interface{}
	Name   string
}

// Validation error
type ValidationError struct {
	Tag     string `json:"tag"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Interface to implement in a resource to configure bindings
type UUIDBinderInterface interface {
	GetUUIDBindings() []UUIDBinding
}

// Interface to implement in a resource to configure validation
type ValidationAwareInterface interface {
	Validate() []ValidationError
	GetCustomValidationMessages() map[string]string
}
