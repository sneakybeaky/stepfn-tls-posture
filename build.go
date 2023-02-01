package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	if err := build(context.Background()); err != nil {
		fmt.Println(err)
	}
}

func build(ctx context.Context) error {
	fmt.Println("Building with Dagger")

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// get reference to the local project
	src := client.Host().Directory(".")

	// get `golang` image
	golang := client.Container().From("golang:latest")

	// mount cloned repository into `golang` image
	golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/src")

	// set build variables
	golang = golang.WithEnvVariable("GOOS", "linux")
	golang = golang.WithEnvVariable("GOARCH", "amd64")

	// define the application build command
	path := "build/"
	golang = golang.WithExec([]string{"go", "build", "-tags", "lambda.norpc", "-ldflags", "-s -w", "-o", path, "tlsposture/functions/ssllabs"})
	//	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -ldflags="-s -w" -o bin/ssllabs/bootstrap tlsposture/functions/ssllabs

	// get reference to build output directory in container
	output := golang.Directory(path)

	// write contents of container build/ directory to the host
	_, err = output.Export(ctx, path)
	if err != nil {
		return err
	}

	return nil
}
