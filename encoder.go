// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package easyapi

import (
	"gitlab.com/kjose/jgmc/api/internal/easyapi/layer"
	"golang.org/x/crypto/bcrypt"
)

// Strongly encode a password based on the resource ID
func EncodePassword(p layer.PasswordEncoderAware) error {
	pwdWithSalt := p.GetIdentifier() + p.GetPlainPassword()
	pwd, err := bcrypt.GenerateFromPassword([]byte(pwdWithSalt), 14)
	if err != nil {
		return err
	}

	p.SetEncodedPassword(string(pwd))
	return nil
}

// Check the strong password
func CheckPassword(password string, p layer.PasswordEncoderAware) bool {
	hash := p.GetEncodedPassword()
	pwdWithSalt := p.GetIdentifier() + password
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwdWithSalt))
	return err == nil
}
