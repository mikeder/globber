.DEFAULT_GOAL := build

build:
	go build -o ./bin/globber ./cmd/globber

image:
	docker build . -t mikeder/globber:latest

run:
	@docker-compose up --build
