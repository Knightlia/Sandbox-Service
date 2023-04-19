run:
	go run main.go --debug=true

test:
	go test ./...

coverage:
	go test -cover -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
