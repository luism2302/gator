build:
	@go build -o bin/gator/gator cmd/gator/main.go
run:
	@go run cmd/gator/main.go
test:
	@go test -v ./...