.PHONY: build
build:
	go build -o app ./application


#.PHONY: migrate
#migrate:
#	migrate -source file://path/to/migrations -database postgres://postgres:mysecretpassword@localhost:5432/testdb up