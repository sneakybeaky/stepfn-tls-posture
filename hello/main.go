package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"runtime/debug"
)

func HandleRequest(event events.APIGatewayV2HTTPRequest) (string, error) {
	info, _ := debug.ReadBuildInfo()

	_, returnError := event.QueryStringParameters["error"]

	if returnError {
		panic("told to do this")
	}

	return fmt.Sprintf("Hello world !\n%s\n", info), nil
}

func main() {
	lambda.Start(HandleRequest)
}
