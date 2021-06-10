// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ginh

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type httpError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Implements the Error interface
func (e *httpError) Error() string {
	return fmt.Sprintf("%d - %s - %v", e.Code, e.Message, e.Data)
}

// Create a http error and display it on gin context
func HttpError(c *gin.Context, code int, message string, data interface{}) error {
	if message == "" {
		message = "An error occured"
	}
	e := &httpError{
		Code:    code,
		Message: message,
		Data:    data,
	}
	c.JSON(code, gin.H{"error": e})

	return e
}
