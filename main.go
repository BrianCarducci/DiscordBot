package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	/*
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	*/
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		fmt.Println("Cloudwatch not working.")
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body: "Hello",
	}, nil
}
