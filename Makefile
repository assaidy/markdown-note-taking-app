GOOSE_ENV = GOOSE_DRIVER="sqlite3" GOOSE_DBSTRING="./database/notes.db" GOOSE_MIGRATION_DIR="database/migrations"

all: build

run: build 
	@./bin/api-server

build:
	@go build -o ./bin/api-server ./cmd/main.go

clean:
	@rm -rf bin

up:
	$(GOOSE_ENV) goose up

down:
	$(GOOSE_ENV) goose down

reset:
	$(GOOSE_ENV) goose reset

migration:
	$(GOOSE_ENV) goose create -s $(name) sql

