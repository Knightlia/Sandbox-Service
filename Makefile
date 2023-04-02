run:
	go run main.go

coverage:
	go test -coverprofile coverage.out ./...
	go tool cover -html=coverage.out
