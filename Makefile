.DEFAULT_GOAL := build

build:
	go build -o ./bin/admin   ./cmd/admin
	go build -o ./bin/globber ./cmd/globber

image:
	docker build . -t mikeder/globber:latest --platform=linux/amd64

run:
	@docker-compose up --build
