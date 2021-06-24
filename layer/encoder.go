// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package layer

// Interface to implement in a resource (ex: User) to support authentification with strong password
type PasswordEncoderAware interface {
	GetIdentifier() string
	GetPlainPassword() string
	GetEncodedPassword() string
	SetEncodedPassword(pwd string)
}
