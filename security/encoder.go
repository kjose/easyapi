// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package security

import (
	"golang.org/x/crypto/bcrypt"
)

// Interface to implement in a resource (ex: User) to support authentification with strong password
type PasswordEncoderAware interface {
	GetIdentifier() string
	GetPlainPassword() string
	GetEncodedPassword() string
	SetEncodedPassword(pwd string)
}

// Strongly encode a password based on the resource ID
func EncodePassword(p PasswordEncoderAware) error {
	pwdWithSalt := p.GetIdentifier() + p.GetPlainPassword()
	pwd, err := bcrypt.GenerateFromPassword([]byte(pwdWithSalt), 14)
	if err != nil {
		return err
	}

	p.SetEncodedPassword(string(pwd))
	return nil
}

// Check the strong password
func CheckPassword(password string, p PasswordEncoderAware) bool {
	hash := p.GetEncodedPassword()
	pwdWithSalt := p.GetIdentifier() + password
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwdWithSalt))
	return err == nil
}
