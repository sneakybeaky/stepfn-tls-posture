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

	fns, err := getFunctionNames()
	if err != nil {
		return err
	}

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

	var binaries []string
	for _, fn := range fns {
		// define the application build command
		path := fmt.Sprintf("build/%s/bootstrap", fn)
		//zip := fmt.Sprintf("build/%s.zip", fn)

		golang = golang.WithExec([]string{"go", "build", "-tags", "lambda.norpc", "-ldflags", "-s -w", "-o", path, "tlsposture/functions/" + fn})
		binaries = append(binaries, fn)
	}

	// Package up the binaries into zips
	// get archiver
	archiver := client.Container().From("alpine:3").
		WithMountedDirectory("/src/build", golang.Directory("/src/build")).
		WithWorkdir("/src").
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "add", "zip"})

	for _, binary := range binaries {

		archiver = archiver.WithExec([]string{"zip", "-j", fmt.Sprintf("build/%s.zip", binary), fmt.Sprintf("build/%s/bootstrap", binary)})
		//zip -j bin/ssllabs.zip bin/ssllabs/bootstrap
	}

	// get reference to build output directory in container
	output := archiver.Directory("build")

	// write contents of container build/ directory to the host
	_, err = output.Export(ctx, "build")
	if err != nil {
		return err
	}

	// get `node` image
	node := client.Container().From("node:lts-gallium") // node 16 lts
	// mount cloned repository into `golang` image
	node = node.WithMountedDirectory("/src", src).WithWorkdir("/src")
	stdout, err := node.WithExec([]string{"npm", "install"}).Stdout(ctx)
	// being executed on
	fmt.Printf("npx says: %s\n", stdout)

	return nil
}

func getFunctionNames() ([]string, error) {

	var fns []string

	files, err := os.ReadDir("functions")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			fns = append(fns, file.Name())
		}
	}

	return fns, nil
}
