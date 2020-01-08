package main

import (
	"flag"
	"log"

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
	err := migrator.Prepare(fileName)
	if err != nil {
		log.Fatalf("Error occured while preparing migration: %v", err)

	}

	err = migrator.RunMigration(fileName)
	if err != nil {
		log.Fatalf("Error occured while running migration: %v", err)
	}
}
