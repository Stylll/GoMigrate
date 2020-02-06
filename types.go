package main

type Migrator struct{}

type Config struct {
	MIGRATION_FOLDER_PATH   string
	MIGRATION_TEMPLATE_NAME string
	PLUGIN_FOLDER_PATH      string
	TEMPLATE_FOLDER_PATH    string
}
