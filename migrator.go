package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"plugin"
	"sort"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type Migrator struct{}

var db *sql.DB

func (migrator *Migrator) CreateMigrationFile(fileName string) error {
	/* TODO - VALIDATE THE FILENAME HAS ONLY ALPHABETS AND HYPHENS */
	filePath := fmt.Sprintf(FILE_PATH, TEMPLATE_FOLDER_PATH, MIGRATION_TEMPLATE_NAME, GO_EXT)
	migrationFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	currentTimeStamp := time.Now().Unix()
	fileName = strconv.FormatInt(currentTimeStamp, 10) + "-" + fileName
	filePath = fmt.Sprintf(FILE_PATH, MIGRATION_FOLDER_PATH, fileName, GO_EXT)
	err = ioutil.WriteFile(filePath, migrationFile, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (migrator *Migrator) Prepare(fileName string) error {
	// check if file exists
	filePath := fmt.Sprintf(FILE_PATH, MIGRATION_FOLDER_PATH, fileName, GO_EXT)
	_, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	//check if it has been migrated before
	isMigrated, err := CheckMigrationExistsInDB(fileName)
	if err != nil {
		return err
	}

	if isMigrated {
		errorMessage := "Migration %s already exists in the DB \n"
		return fmt.Errorf(errorMessage, fileName)
	}

	// create migration plugin
	_, err = CreateMigrationPlugin(fileName)
	if err != nil {
		return err
	}

	return nil
}

func (migrator *Migrator) RunMigration(fileName string) error {
	filePath := fmt.Sprintf(FILE_PATH, PLUGIN_FOLDER_PATH, fileName, PLUGIN_EXT)
	migrationPlugin, err := plugin.Open(filePath)
	if err != nil {
		return err
	}

	up, err := migrationPlugin.Lookup("UP")
	if err != nil {
		return err
	}

	err = up.(func(*sql.DB) error)(db)
	if err != nil {
		return err
	}

	err = migrator.SaveMigratedFileInDB(fileName)
	if err != nil {
		return err
	}

	fmt.Printf("%s was successfully migrated \n", fileName)

	return nil
}

func (migrator *Migrator) SaveMigratedFileInDB(fileName string) error {
	fullFileName := fileName + "." + GO_EXT
	query := `INSERT INTO migrations (
		migration_name
	) VALUES ($1) `

	row, err := db.Query(query, fullFileName)
	defer row.Close()

	if err != nil {
		return err
	}

	return nil
}

func CreateMigrationPlugins(migrationNames []string) []string {
	var createdMigrations []string
	for _, migrationName := range migrationNames {
		name, err := CreateMigrationPlugin(migrationName)
		if err != nil {
			fmt.Printf("Create Plugin Failed: %s \n", err)
		} else {
			createdMigrations = append(createdMigrations, name)
		}
	}

	return createdMigrations
}

func CreateMigrationPlugin(migrationName string) (string, error) {
	buildMode := "-buildmode=plugin"
	outputPath := fmt.Sprintf(FILE_PATH, PLUGIN_FOLDER_PATH, migrationName, PLUGIN_EXT)
	outputFlag := "-o=" + outputPath
	filePath := fmt.Sprintf(FILE_PATH, MIGRATION_FOLDER_PATH, migrationName, GO_EXT)
	cmd := exec.Command("go", "build", buildMode, outputFlag, filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return migrationName, err
}

func SetupDatabase() error {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	connString := fmt.Sprintf(CONN_STRING, dbUser, dbPassword, dbHost, dbPort, dbName)
	dbc, err := sql.Open("postgres", connString)
	if err != nil {
		return err
	}
	db = dbc

	query := `CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		migration_name VARCHAR UNIQUE,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW(),
		deleted_at TIMESTAMP DEFAULT NULL
	) `

	_, err = db.Query(query)
	if err != nil {
		return err
	}

	return nil
}

func CheckMigrationExistsInDB(migrationName string) (bool, error) {
	var resultName string
	argument := migrationName + "." + GO_EXT
	query := `SELECT migration_name FROM migrations WHERE migration_name = $1;`
	row, err := db.Query(query, argument)
	if err != nil {
		return false, err
	}
	defer row.Close()

	if !row.Next() {
		return false, nil
	}

	err = row.Scan(&resultName)
	if err != nil {
		return false, err
	}

	if resultName == argument {
		return true, nil
	}

	return false, nil
}

func GetMigrationList() ([]string, error) {
	file, err := os.Open(MIGRATION_FOLDER_PATH)
	migrations, err := file.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	sort.Strings(migrations)

	return migrations, nil
}
