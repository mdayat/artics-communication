.DEFAULT_GOAL := run

.PHONY:fmt vet build run govulncheck staticcheck revive

.SILENT:

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -C cmd/web -o app

run:
	go run cmd/web/main.go

seed:
	docker compose up -d db
	@until docker exec db pg_isready -U postgres; do \
		sleep 1; \
	done
	docker compose run --rm atlas migrate apply -u "postgres://postgres:password@db:5432/postgres?search_path=public&sslmode=disable" --revisions-schema public
	docker compose run --rm seed
	
govulncheck:
	govulncheck ./...

staticcheck:
	staticcheck ./...

revive:
	revive -config revive.toml -formatter friendly ./...