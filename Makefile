.PHONY: build clean deploy

build:
	cd api && env GOOS=linux go build -ldflags="-s -w" -o bin/create lib/create/main.go

clean:
	rm -rf ./api/bin

deploy: clean build
	sls deploy --verbose
