// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/kjose/jgmc/api/internal/goapi/ginh"
	"gitlab.com/kjose/jgmc/api/internal/goapi/security"
)

// Middleware to check the token sent in the header
// It needs TOKEN_COOKIE_NAME env var to know the cookie where it is registered
func SecurityTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tkn, err := c.Cookie(os.Getenv("TOKEN_COOKIE_NAME"))
		if err != nil {
			tkn = c.GetHeader("Authorization")
			if tkn == "" {
				ginh.HttpError(c, http.StatusUnauthorized, "Authorization token is required", nil)
				c.Abort()
				return
			}
			tkn = strings.Replace(tkn, "Bearer ", "", 1)
		}

		tknData, err := security.ParseToken(tkn)
		if err != nil {
			ginh.HttpError(c, http.StatusUnauthorized, "Authorization token is invalid", nil)
			c.Abort()
			return
		}

		c.Set(ginh.CONTEXT_KEY_TOKEN, tknData)

		c.Next()
	}
}
