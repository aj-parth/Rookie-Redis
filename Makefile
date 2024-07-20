run: build
	@./bin/rookie-redis-bin

build:
	@go build -o bin/rookie-redis-bin .