package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"io"
	"net/http"
)

type Scan struct {
	Host   string `json:"host"`
	Status string `json:"status"`
	ID     string `json:"requestId"`
}

func HandleRequest(event Scan) (map[string]any, error) {

	fmt.Printf("Got %q", event)

	switch event.Status {
	case "START":
		r, err := http.Get(fmt.Sprintf("https://api.ssllabs.com/api/v3/analyze?host=%s&startNew=on&all=done", event.Host))
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		var raw = map[string]any{}
		err = json.Unmarshal(body, &raw)
		if err != nil {
			return nil, err
		}

		raw["requestId"] = event.ID
		return raw, err
	default:
		r, err := http.Get(fmt.Sprintf("https://api.ssllabs.com/api/v3/analyze?host=%s", event.Host))
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)

		var raw = map[string]any{}
		err = json.Unmarshal(body, &raw)
		if err != nil {
			return nil, err
		}

		raw["requestId"] = event.ID
		return raw, err
	}

}

func main() {
	lambda.Start(HandleRequest)
}
