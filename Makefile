.PHONY: postgres migrate run tidy build

postgres:
	docker run --rm -ti -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 postgres:15

migrate:
	# go run cmd/bravo/main.go migrate up

build-local:
	go build -o bravo ./cmd/bravo/main.go

docker: docker-build-bin docker-img docker-run

docker-build-bin:
	GOOS=linux GOARCH=amd64 go build -o bravo ./cmd/bravo/main.go

docker-img:
	docker build -t bravo .

docker-run:
	docker run bravo

run:
	go run cmd/bravo/main.go

tidy:
	go mod tidy
