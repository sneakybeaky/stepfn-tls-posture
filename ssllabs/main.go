package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(event json.RawMessage) error {

	fmt.Printf("Got %q", event)

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
