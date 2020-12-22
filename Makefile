.PHONY: build
build:
	go build -o auth .application/microservices/auth/cmd
	go build -o vacancy ./application/microservices/vacancy/cmd
	go build -o app ./cmd/api


#.PHONY: migrate
#migrate:
#	migrate -source file://path/to/migrations -database postgres://postgres:mysecretpassword@localhost:5432/testdb up