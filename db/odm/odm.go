// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// This ODM is build on top on go-bongo https://github.com/go-bongo/bongo

package odm

import (
	"github.com/go-bongo/bongo"
)

type Config struct {
	ConnectionString string
	Database         string
}

var DB *bongo.Connection

func Create(config Config) (*bongo.Connection, error) {
	// config
	c := &bongo.Config{
		ConnectionString: config.ConnectionString,
		Database:         config.Database,
	}

	db, err := bongo.Connect(c)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Init(config Config, isDefault bool) error {
	db, err := Create(config)
	if err != nil {
		return err
	}
	DB = db

	return nil
}
