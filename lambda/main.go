package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ZED-Magdy/learning-aws-cdk-go/lambda/database"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body map[string]string
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid request body",
		}, nil
	}

	name, ok := body["name"]
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Missing name in request body",
		}, nil
	}

	message, ok := body["message"]
	if !ok {
		message = fmt.Sprintf("Hello, %s!", name)
	}

	// Initialize DynamoDB client
	ddbClient, err := database.NewDynamoDBClient()
	if err != nil {
		log.Printf("Error creating DynamoDB client: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error connecting to database",
		}, nil
	}

	item := database.Item{
		Name:    name,
		Message: message,
	}

	if err := ddbClient.PutItem(item); err != nil {
		log.Printf("Error storing item in DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error storing data",
		}, nil
	}

	response := map[string]string{
		"name":    name,
		"message": message,
		"status":  "stored in database",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error creating response",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonResponse),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}