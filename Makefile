build:
	docker-compose build

run:
	docker-compose up

test:
	go test -v ./...