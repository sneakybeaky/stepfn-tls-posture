.PHONY: build clean deploy all

build: bin/hello.zip bin/ssllabs.zip

bin/hello.zip: bin/hello/bootstrap
	zip -j bin/hello.zip bin/hello/bootstrap

bin/hello/bootstrap:
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -ldflags="-s -w" -o bin/hello/bootstrap tlsposture/hello

bin/ssllabs.zip: bin/ssllabs/bootstrap
	zip -j bin/ssllabs.zip bin/ssllabs/bootstrap

bin/ssllabs/bootstrap:
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -ldflags="-s -w" -o bin/ssllabs/bootstrap tlsposture/ssllabs

clean:
	rm -rf ./bin

deploy: build
	npx sls deploy --verbose

all: clean build deploy
