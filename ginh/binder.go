// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ginh

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gitlab.com/kjose/jgmc/api/internal/goapi/db/dao"
)

var (
	BinderConfig = &binderConfig{
		// If true the body will be kept in the context key gin.BodyBytesKey
		KeepBody: false,
	}
)

// Interface to implement in a resource to configure bindings
type UUIDBinderInterface interface {
	GetUUIDBindings() []UUIDBinding
}

// Interface to implement in a resource to configure validation
type ValidationAwareInterface interface {
	Validate() []ValidationError
	GetCustomValidationMessages() map[string]string
}

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

// Binder config
type binderConfig struct {
	KeepBody bool
}

// Create a new validation error from a tag and a field
func NewValidationError(tag string, field string, resource interface{}) ValidationError {
	message := fmt.Sprintf("Field %s failed with condition `%s`", field, tag)
	if m, ok := resource.(ValidationAwareInterface); ok {
		if customMessage, ok := m.GetCustomValidationMessages()[fmt.Sprintf("%s:%s", field, tag)]; ok {
			message = customMessage
		}
	}

	return ValidationError{
		Tag:     tag,
		Field:   field,
		Message: message,
	}
}

// Bind and validate recursively a request body to a resource
func BindAndValidate(c *gin.Context, i interface{}) error {
	// request validation
	validationErrors := []ValidationError{}
	var err error
	if BinderConfig.KeepBody {
		err = c.ShouldBindBodyWith(i, binding.JSON)
	} else {
		err = c.ShouldBindJSON(i)
	}
	if err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				validationErrors = append(validationErrors, NewValidationError(e.Tag(), e.Field(), i))
			}
		} else {
			validationErrors = append(validationErrors, NewValidationError("", err.Error(), i))
		}
	}
	if iv, ok := i.(ValidationAwareInterface); ok {
		validationErrors = append(validationErrors, iv.Validate()...)
	}
	if len(validationErrors) > 0 {
		return HttpError(c, http.StatusBadRequest, "Validation errors", validationErrors)
	}

	// uuid bindings
	err = AppendBindings(i)
	if err != nil {
		return HttpError(c, http.StatusNotFound, err.Error(), nil)
	}

	return nil
}

// Append bindings based on GetUUIDBindings resource method
func AppendBindings(item interface{}) error {
	if ib, ok := item.(UUIDBinderInterface); ok {
		for _, b := range ib.GetUUIDBindings() {
			if b.UUID == nil {
				continue
			}

			_, err := dao.GetResourceDAO(b.BindTo).FindById(b.BindTo, b.UUID.String())
			if err != nil {
				return fmt.Errorf("%s not found", b.Name)
			}
			// recursive
			if bb, ok := b.BindTo.(UUIDBinderInterface); ok {
				err := AppendBindings(bb)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Remove bindings of a resource
func RemoveUUIDBindings(item interface{}) {
	if ib, ok := item.(UUIDBinderInterface); ok {
		for _, b := range ib.GetUUIDBindings() {
			if b.UUID == nil {
				continue
			}

			v := reflect.ValueOf(b.BindTo)
			p := v.Elem()
			p.Set(reflect.Zero(p.Type()))
			// recursive
			if bb, ok := b.BindTo.(UUIDBinderInterface); ok {
				RemoveUUIDBindings(bb)
			}
		}
	}
}
