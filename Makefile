.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./userland-names/userland-names
	
build:
	GOOS=linux GOARCH=amd64 go build -o userland-names/userland-names ./userland-names
