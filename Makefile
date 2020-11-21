CMD := worker

build:
	go build -o $(CMD) ./cmd/main.go

run:
	make build
	./$(CMD)

docker:
	docker build --tag $(CMD) -f ./build/Dockerfile .
