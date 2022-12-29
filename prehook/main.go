package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy"
	"os"
)

type CodeDeployEvent struct {
	DeploymentId                  string `json:"deploymentId"`
	LifecycleEventHookExecutionId string `json:"lifecycleEventHookExecutionId"`
}

func HandleRequest(ctx context.Context, event CodeDeployEvent) error {

	fmt.Printf("Prehook called for '%+v'\n", event)
	fmt.Printf("I'm checking %q\n", os.Getenv("NEW_VERSION"))
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := codedeploy.NewFromConfig(cfg)

	//d, err := client.GetDeployment(ctx, &codedeploy.GetDeploymentInput{DeploymentId: &event.DeploymentId})
	//
	//if err != nil {
	//	return err
	//}
	//
	//gs, err := client.ListDeploymentGroups(ctx, &codedeploy.ListDeploymentGroupsInput{
	//	ApplicationName: d.DeploymentInfo.ApplicationName,
	//})
	//
	//if err != nil {
	//	return err
	//}
	//
	//for _, group := range gs.DeploymentGroups {
	//	g, err := client.GetDeploymentGroup(ctx, &codedeploy.GetDeploymentGroupInput{
	//		ApplicationName:     d.DeploymentInfo.ApplicationName,
	//		DeploymentGroupName: aws.String(group),
	//	})
	//
	//	if err != nil {
	//		return err
	//	}
	//
	//	//fmt.Printf("Appspec content is %q ", g.DeploymentGroupInfo.TargetRevision.AppSpecContent.Content)
	//
	//}

	params := &codedeploy.PutLifecycleEventHookExecutionStatusInput{
		DeploymentId:                  &event.DeploymentId,
		LifecycleEventHookExecutionId: &event.LifecycleEventHookExecutionId,
		Status:                        "Succeeded",
	}

	fmt.Printf("Reporting status '%+v'\n", params)

	_, err = client.PutLifecycleEventHookExecutionStatus(ctx, params)

	return err
}

func main() {
	lambda.Start(HandleRequest)
}
