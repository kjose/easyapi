// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/kjose/jgmc/api/internal/easyapi/event"
)

// Middleware to add in gin configuration to add event listeners
func EventListenersMiddleware(els []event.EventListener) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, el := range els {
			event.RegisterEventListener(el)
		}

		c.Next()

		go event.ResetEventListeners()
	}
}
