// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// This ORM is build on top on gorm https://gorm.io/

package orm

import (
	"log"
	"os"

	"gitlab.com/kjose/jgmc/api/internal/easyapi/db/dao"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init the orm DAO
func Create(driver gorm.Dialector, loggerConfig *logger.Config) (*gorm.DB, error) {
	// config
	config := &gorm.Config{}
	if loggerConfig != nil {
		config.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			*loggerConfig,
		)
	}

	db, err := gorm.Open(driver, config)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Init the orm DAO
func Init(driver gorm.Dialector, isDefault bool, loggerConfig *logger.Config) error {
	db, err := Create(driver, loggerConfig)
	if err != nil {
		return err
	}
	DB = db

	if isDefault {
		dao.InitDefaultDAO(DAO)
	}
	return nil
}
