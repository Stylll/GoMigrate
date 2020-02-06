# GoMigrate
A database migration tool written in Golang

## Example
1. `./GoMigrate -o=create -n=my-new-migration-file` (If you're using the executable file)
2. `go run ./ -o=create -n=my-new-migration-file` (If you cloned this repo and are running it inside the project folder)

## Setup
1. Download the GoMigrate executable file in this repo and add it in your project folder.
2. If you would like to build the executable file yourself, clone this repo and run `go build -o=GoMigrate` then copy and paste the created executable file in your project folder.
3. Setup the following environment variables in your project. You can reference the sample.env file.
- `DB_USER=your_database_username`
- `DB_PASSWORD=your_database_password`
- `DB_NAME=your_database_name`
- `DB_HOST=your_database_host_name`
- `DB_PORT=your_database_port`
- `MIGRATION_FOLDER_PATH=folder_to_store_all_migrations`
- `MIGRATION_TEMPLATE_NAME=your_migration_template_name` please use the templates/migration-template.go file
- `PLUGIN_FOLDER_PATH=folder_to_store_all_migration_plugins`
- `TEMPLATE_FOLDER_PATH=folder_where_the_migration_template_is_stored`

## Flags
1. -o: (Required) To set the operation to perform. Values that can be parsed are:
- `c or create` to create a new migration file
- `r or run` to run a migration file or run all migration files
- `u or undo` to undo migration file(s)

2. -n: To set the name of the file


# Examples

## To Create A Migration File
1. Run `./GoMigrate -o=create -n=my-new-migration-file`. A new go file will be created inside the migrations folder
2. Open the newly created migration file to edit the queries

## To Run A Migration File
1. Run `./GoMigrate -o=run -n=my-migration-file` to migrate a single file
2. Run `./GoMigrate -o=run` to migrate all files

## To Undo Migration
1. Run `./GoMigrate -o=undo -n=my-migration-file` to undo up to a specific migration file
1. Run `./GoMigrate -o=undo` to undo all migrations

