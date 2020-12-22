.PHONY: build
build:
	mkdir build
	go build -o build/auth ./application/microservices/auth/cmd
	go build -o build/vacancy ./application/microservices/vacancy/cmd
	go build -o build/app ./cmd/api


#.PHONY: migrate
#migrate:
#	migrate -source file://path/to/migrations -database postgres://postgres:mysecretpassword@localhost:5432/testdb up