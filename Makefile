start:
	go run main.go

build:
	go build

build-windows:
	GOARCH=amd64 GOOS=windows go build

build-linux:
	GOARCH=amd64 GOOS=linux go build
