package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/sethvargo/go-envconfig"
	"net/http"
)

type Handler struct {
	BusName string `env:"BUS_NAME"`
}

type StartScan struct {
	RequestId string `json:"requestId"`
	Host      string `json:"host"`
}

func (h Handler) Handle(ctx context.Context, req events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	fmt.Printf("Got %+v", req)

	ss := StartScan{
		RequestId: req.RequestContext.RequestID,
		Host:      req.QueryStringParameters["host"],
	}

	if ss.Host == "" {
		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Must supply host",
		}, nil
	}

	detail, err := json.Marshal(ss)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Using the Config value, create the EventBridge client
	svc := eventbridge.NewFromConfig(cfg)

	result, err := svc.PutEvents(ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				Detail:       aws.String(string(detail)),
				DetailType:   aws.String("Start Scan"),
				EventBusName: aws.String(h.BusName),
				Source:       aws.String("posture.api"),
			},
		},
		EndpointId: nil,
	})

	if err != nil {
		return nil, err
	}

	fmt.Printf("Result from send is %+v", result)

	return &events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
	}, nil

}

func main() {

	ctx := context.Background()

	var h Handler
	if err := envconfig.Process(ctx, &h); err != nil {
		panic(err)
	}

	lambda.Start(h.Handle)
}
