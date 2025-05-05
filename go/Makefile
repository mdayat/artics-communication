.DEFAULT_GOAL := run

.PHONY:fmt vet build run govulncheck staticcheck revive

.SILENT:

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o app

run:
	go run main.go

seed:
	docker compose up -d db
	@until docker exec db pg_isready -U postgres; do \
		sleep 1; \
	done
	atlas migrate apply --env prod -u "postgres://postgres:password@localhost:5432/postgres?search_path=public&sslmode=disable" --revisions-schema public
	go run cmd/seed/main.go

govulncheck:
	govulncheck ./...

staticcheck:
	staticcheck ./...

revive:
	revive -config revive.toml -formatter friendly ./...