package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

type Migrator struct{}

var MIGRATION_FOLDER_PATH = "./migrations/%s"
var MIGRATION_TEMPLATE_NAME = "migration-template.go"

func (migrator *Migrator) CreateMigrationFile(fileName string) error {
	migrationFile, err := ioutil.ReadFile(fmt.Sprintf(MIGRATION_FOLDER_PATH, MIGRATION_TEMPLATE_NAME))
	if err != nil {
		return err
	}
	currentTimeStamp := time.Now().Unix()
	fileName = strconv.FormatInt(currentTimeStamp, 10) + "-" + fileName + ".go"
	err = ioutil.WriteFile(fmt.Sprintf(MIGRATION_FOLDER_PATH, fileName), migrationFile, 0644)
	if err != nil {
		return err
	}

	return nil
}
