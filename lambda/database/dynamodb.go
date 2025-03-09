package database

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClient struct {
	client *dynamodb.Client
	table  string
}

type Item struct {
	Name    string `json:"name" dynamodbav:"name"`
	Message string `json:"message" dynamodbav:"message"`
}

const defaultTableName = "HelloCdkGoTable"

func NewDynamoDBClient() (*DynamoDBClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if (err != nil) {
		return nil, err
	}

	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		tableName = defaultTableName
	}

	return &DynamoDBClient{
		client: dynamodb.NewFromConfig(cfg),
		table:  tableName,
	}, nil
}

func (d *DynamoDBClient) PutItem(item Item) error {
	av := make(map[string]types.AttributeValue, 2)
	
	if _, ok := av["name"]; !ok {
		av["name"] = &types.AttributeValueMemberS{Value: item.Name}
	}

	if _, ok := av["message"]; !ok {
		av["message"] = &types.AttributeValueMemberS{Value: item.Message}
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(d.table),
		Item:      av,
	}

	_, err := d.client.PutItem(context.Background(), input)
	return err
}

func (d *DynamoDBClient) GetItem(name string) (*Item, error) {
	key := map[string]types.AttributeValue{
		"name": &types.AttributeValueMemberS{Value: name},
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(d.table),
		Key:       key,
	}

	result, err := d.client.GetItem(context.Background(), input)
	if err != nil {
		return nil, err
	}

	item := new(Item)
	if err := attributevalue.UnmarshalMap(result.Item, item); err != nil {
		return nil, err
	}

	return item, nil
}
