.PHONY: dep run dev build test db-migrate db-rollback db-drop db-fresh db-seed db-fresh-seed

dep:
	go mod tidy

run:
	go run ./cmd/api

dev:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test -v ./...

db-migrate:
	go run ./cmd/api --migrate

db-rollback:
	go run ./cmd/api --rollback

db-drop:
	go run ./cmd/api --drop

db-fresh:
	go run ./cmd/api --fresh

db-seed:
	go run ./cmd/api --seed

db-fresh-seed:
	go run ./cmd/api --fresh --seed
