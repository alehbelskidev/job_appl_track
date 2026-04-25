.PHONY: sqlc, dev
sqlc:
	sqlc generate

dev:
	go run cmd/app/main.go
