package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/joho/godotenv"
)

var migrator Migrator

func main() {
	godotenv.Load()
	var operation, fileName string
	flag.StringVar(&operation, "o", "", "operation to perform")
	flag.StringVar(&fileName, "n", "", "migration file name")

	flag.Parse()

	if operation == "" {
		log.Fatal("flag 'o':'operation' was not provided")
	}

	if operation == "create" || operation == "c" {
		if fileName == "" {
			log.Fatal("flag 'n':'file name' was not provided")
		}
		CreateFile(fileName)
		return
	}

	if operation == "run" || operation == "r" {
		if fileName != "" {
			RunMigration(fileName)
		} else {
			RunMigrations()
		}
		return
	}

	errorMessage := "Operation %s not recognized"
	log.Fatalf(errorMessage, operation)
}

func CreateFile(fileName string) {
	err := migrator.CreateMigrationFile(fileName)
	if err != nil {
		log.Fatalf("Error occured while creating file: %v", err)

	}
}

func RunMigration(fileName string) {
	// setup database
	err := SetupDatabase()
	if err != nil {
		log.Fatalf("Error occured while setting up database: %v", err)
	}

	err = migrator.Prepare(fileName)
	if err != nil {
		log.Fatalf("Error occured while preparing migration: %v", err)

	}

	err = migrator.RunMigration(fileName)
	if err != nil {
		log.Fatalf("Error occured while running migration: %v", err)
	}
}

func RunMigrations() {
	migrationList, err := GetMigrationList()
	if err != nil {
		log.Fatalf("Error occured while fetching migration list: %v", err)
	}

	if len(migrationList) == 0 {
		log.Fatal("No migration available")
	}

	err = SetupDatabase()
	if err != nil {
		log.Fatalf("Error occured while setting up database: %v", err)
	}

	for _, migration := range migrationList {
		// remove the extension attached to the name
		migrationName := strings.Split(migration, ".")[0]

		err = migrator.Prepare(migrationName)
		if err != nil {
			fmt.Printf("Error occured while preparing migration: %v \n", err)
			continue
		}

		err = migrator.RunMigration(migrationName)
		if err != nil {
			fmt.Printf("Error occured while running migration: %v \n", err)
		}
	}
}
