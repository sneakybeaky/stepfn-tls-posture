package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
	"time"

	"dagger.io/dagger"
)

func main() {

	ctx := context.Background()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()

	archives, err := build(ctx, client.Pipeline("build"))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// write contents of container build/ directory to the host
	//for _, archive := range archives {
	//	ok, err := archive.file.Export(ctx, fmt.Sprintf("build/%s.zip", archive.functionName))
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//
	//	if !ok {
	//		fmt.Println("Failed to export build output to host")
	//		os.Exit(1)
	//	}
	//
	//}

	if err := deploy(ctx, client.Pipeline("Deploy"), archives); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func deploy(ctx context.Context, client *dagger.Client, archives []archive) error {

	// get `node` image
	node := client.Container().From("node:lts-gallium").
		WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
		WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
		WithEnvVariable("AWS_SESSION_TOKEN", os.Getenv("AWS_SESSION_TOKEN")).
		WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION")).
		WithMountedDirectory("/src", client.Host().Directory(".")).
		WithWorkdir("/src")

	for _, a := range archives {
		node = node.WithFile(fmt.Sprintf("build/%s.zip", a.functionName), a.file)
	}

	exitCode, err := node.WithExec([]string{"npm", "install"}).WithExec([]string{"npx", "sls", "deploy", "--verbose"}).ExitCode(ctx)
	// being executed on
	fmt.Printf("npx says: %d\n", exitCode)

	return err

}

type archive struct {
	functionName string
	file         *dagger.File
}

func build(ctx context.Context, client *dagger.Client) ([]archive, error) {
	fmt.Println("Building with Dagger")

	gocache := client.CacheVolume(time.Now().Weekday().String())
	g, ctx := errgroup.WithContext(ctx)

	fnnames, err := getFunctionNames()

	if err != nil {
		return nil, err
	}

	zips := make([]archive, len(fnnames))

	golang := client.Container().
		From("golang:latest").
		WithMountedCache("/cache", gocache).
		WithEnvVariable("GOMODCACHE", "/cache").
		WithEnvVariable("GOOS", "linux").
		WithEnvVariable("GOARCH", "amd64").
		WithMountedDirectory("/src", client.Host().Directory(".")).
		WithWorkdir("/src")

	packager := client.Container().From("alpine:3").
		WithWorkdir("/src").
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "add", "zip"})

	for i, fn := range fnnames {

		f := fn
		id := i
		g.Go(func() error {

			path := fmt.Sprintf("build/%s/bootstrap", f)

			builder := golang.
				WithExec([]string{"go", "mod", "download"}).
				WithExec([]string{"go", "build", "-tags", "lambda.norpc", "-ldflags", "-s -w", "-o", path, "tlsposture/functions/" + f})

			output := builder.Directory("build")

			archiver := packager.
				WithMountedDirectory("/src/build", output).
				WithExec([]string{"zip", "-j", fmt.Sprintf("build/%s.zip", f), fmt.Sprintf("build/%s/bootstrap", f)})

			zips[id] = archive{
				functionName: f,
				file:         archiver.File(fmt.Sprintf("build/%s.zip", f)),
			}

			return nil
		})

	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return zips, err

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
