run:
	go run ./cmd/server

test:
	go test ./internal/calculator/... -v

build:
	go build -o pack-calculator ./cmd/server

docker-build:
	docker build -t pack-calculator .

docker-run:
	docker run --rm -p 8080:8080 pack-calculator
