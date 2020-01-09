# GoMigrate
A database migration tool written in Golang

## Example
go run ./ -o=create -n=my-new-migration-file

## Flags
1. -o: (Required) To set the operation to perform. Values that can be parsed are:
- `c or create` to create a new migration file
- `r or run` to run a migration file or run all migration files

2. -n: To set the name of the file


# Examples

## To Create A Migration File
1. Run `go run ./ -o=create -n=my-new-migration-file`. A new go file will be created inside the migrations folder
2. Open the newly created migration file to edit the query

## To Run A Migration File
1. Run `go run ./ -o=run -n=my-migration-file` to migrate a single file
2. Run `go run ./ -o=run` to migrate all files

