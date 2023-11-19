.PHONY: web_server \
		docker-up docker-down

module := $(shell go list -m)
out_dir := build/out
all: web_server

web_server: 
	GOOS=linux GOARCH=amd64 go build -o $(out_dir) $(module)/cmd/$@

docker-up:
	docker compose -f build/docker-compose.yaml up -d

docker-down:
	docker compose -f build/docker-compose.yaml down