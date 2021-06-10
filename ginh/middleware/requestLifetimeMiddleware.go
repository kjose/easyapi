// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/kjose/jgmc/api/internal/goapi/event"
)

// Middleware to add in gin configuration to support request listeners
func RequestLifetimeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		event.DispatchEvent(c, event.EVENT_REQUEST_START, &event.RequestEvent{
			Context: c,
		})

		c.Next()

		go func() {
			event.DispatchEvent(c, event.EVENT_REQUEST_TERMINATE, &event.RequestEvent{
				Context: c,
			})
		}()
	}
}
