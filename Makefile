.PHONY: build clean deploy all

build: bin/hello.zip bin/prehook.zip

bin/hello.zip: bin/hello/bootstrap
	zip -j bin/hello.zip bin/hello/bootstrap

bin/hello/bootstrap:
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -ldflags="-s -w" -o bin/hello/bootstrap tlsposture/hello

bin/prehook.zip: bin/prehook/bootstrap
	zip -j bin/prehook.zip bin/prehook/bootstrap

bin/prehook/bootstrap:
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -ldflags="-s -w" -o bin/prehook/bootstrap tlsposture/prehook

clean:
	rm -rf ./bin

deploy: build
	npx sls deploy --verbose

all: clean build deploy
