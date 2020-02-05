package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
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

	if operation == "undo" || operation == "u" {
		UndoMigration(fileName)

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
	migrationCount := 0
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
			errMessage := err.Error()
			if strings.Contains(errMessage, "already exists in the DB") {
				continue
			}
			fmt.Printf("Error occured while preparing migration: %v \n", err)
			continue
		}

		err = migrator.RunMigration(migrationName)
		if err != nil {
			fmt.Printf("Error occured while running migration: %v \n", err)
			continue
		}

		migrationCount = migrationCount + 1
	}

	fmt.Printf("%d Migration(s) Done \n", migrationCount)
}

func UndoMigration(migrationName string) {
	migrationList, err := GetMigrationList()
	reversedMigrationsCount := 0
	sort.Sort(sort.Reverse(sort.StringSlice(migrationList)))

	if err != nil {
		log.Fatalf("Error occured while fetching migration list: %v", err)
	}

	if len(migrationList) == 0 {
		log.Fatal("No migration available")
	}

	if migrationName != "" {
		fullMigrationName := migrationName + "." + GO_EXT
		index := FindStringIndex(migrationList, fullMigrationName)
		if index != -1 {
			if index+1 < len(migrationList) {
				migrationList = migrationList[:index+1]
			}
		} else {
			log.Fatalf("Migration %s does not exist", migrationName)
		}
	}

	err = SetupDatabase()
	if err != nil {
		log.Fatalf("Error occured while setting up database: %v", err)
	}

	for _, migration := range migrationList {
		// remove the extension attached to the name
		migrationName := strings.Split(migration, ".")[0]

		err = migrator.PrepareUndo(migrationName)
		if err != nil {
			errMessage := err.Error()
			if strings.Contains(errMessage, "doesn't exist in the DB") {
				continue
			}
			fmt.Printf("Error occured while preparing undo migration: %v \n", err)
			continue
		}

		// undo migration
		err = migrator.UndoMigration(migrationName)
		if err != nil {
			fmt.Printf("Error occured while reverting migration: %v \n", err)
			continue
		}

		reversedMigrationsCount = reversedMigrationsCount + 1
	}

	fmt.Printf("%d Migration(s) Reversed \n", reversedMigrationsCount)
}
