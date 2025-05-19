MIGRATE_EXE=C:\ProgramData\chocolatey\bin\migrate.exe
MIGRATE_DB=postgres://postgres:postgres@localhost:5432/todo?sslmode=disable

migrate-up:
	cmd.exe /C "$(MIGRATE_EXE) -path ./migrations -database $(MIGRATE_DB) up"

migrate-down:
	cmd.exe /C "$(MIGRATE_EXE) -path ./migrations -database $(MIGRATE_DB) down 1"

migrate-create:
	cmd.exe /C "$(MIGRATE_EXE) create -ext sql -dir migrations -seq $(name)"
