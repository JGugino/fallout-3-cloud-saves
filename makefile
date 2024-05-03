default:
	@go run main.go

init:
	@go run main.go init

upload:
	@go run main.go upload

sync:
	@go run main.go sync

build:
	@go build -o ./bin/fcs