package main

import (
	"flag"
	"log"
)

var migrator Migrator

func main() {
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
		createFile(fileName)
	}

	errorMessage := "Operation %s not recognized"
	log.Fatalf(errorMessage, operation)
}

func createFile(fileName string) {
	err := migrator.CreateMigrationFile(fileName)
	if err != nil {
		log.Fatalf("Error occured while creating file: %v", err)

	}
}
