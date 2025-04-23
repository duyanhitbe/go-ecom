build-api:
	go build -o bin/api cmd/api/main.go
run-api-dev:
	GO_MOD=dev go run cmd/api/main.go
run-api-prod:
	GO_MOD=prod go run cmd/api/main.go
build-cli:
	go build -o bin/cli cmd/cli/main.go
run-cli:
	go run cmd/cli/main.go
sqlc:
	sqlc generate
create-migration:
	migrate create -dir database/migrations -seq -ext sql $(name)
migrate:
	#postgresql://username:password@localhost:5432/dbname?sslmode=disable
	migrate -database $(url) -path database/migrations up
unit-test:
	go test -v ./...
test-coverage:
	rm -f test/coverage.out test/coverage.html
	go test -coverprofile=test/coverage.out ./...
	go tool cover -html=test/coverage.out -o test/coverage.html