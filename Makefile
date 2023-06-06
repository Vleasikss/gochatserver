build:
	@go build -o bin/out

run: build
	./bin/out

test:
	go test -v ./...