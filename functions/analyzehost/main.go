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
	ID     string `json:"id"`
}

func HandleRequest(event Scan) (json.RawMessage, error) {

	fmt.Printf("Got %q", event)

	switch event.Status {
	case "START":
		r, err := http.Get(fmt.Sprintf("https://api.ssllabs.com/api/v3/analyze?host=%s&startNew=on&all=done", event.Host))
		if err != nil {
			return []byte{}, err
		}
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		raw := json.RawMessage{}
		err = json.Unmarshal(body, &raw)
		if err != nil {
			return nil, err
		}

		return raw, err
	default:
		r, err := http.Get(fmt.Sprintf("https://api.ssllabs.com/api/v3/analyze?host=%s", event.Host))
		if err != nil {
			return []byte{}, err
		}
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)

		raw := json.RawMessage{}
		err = json.Unmarshal(body, &raw)
		if err != nil {
			return nil, err
		}

		return raw, err
	}

}

func main() {
	lambda.Start(HandleRequest)
}
