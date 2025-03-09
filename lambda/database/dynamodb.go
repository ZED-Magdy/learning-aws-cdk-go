package database

import (
	"context"

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
	Name    string `json:"name"`
	Message string `json:"message"`
}

const tableName = "HelloCdkGoTable"

func NewDynamoDBClient() (*DynamoDBClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	return &DynamoDBClient{
		client: dynamodb.NewFromConfig(cfg),
		table:  tableName,
	}, nil
}

func (d *DynamoDBClient) PutItem(item Item) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(d.table),
		Item:      av,
	}

	_, err = d.client.PutItem(context.Background(), input)
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

	if result.Item == nil {
		return nil, nil
	}

	var item Item
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
