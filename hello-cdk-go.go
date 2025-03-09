package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type HelloCdkGoStackProps struct {
	awscdk.StackProps
}

func NewHelloCdkGoStack(scope constructs.Construct, id string, props *HelloCdkGoStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	table := awsdynamodb.NewTable(stack, jsii.String("HelloCdkGoTable"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("name"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	fn := awslambda.NewFunction(stack, jsii.String("HelloCdkGoFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("lambda/function.zip"), nil),
		Environment: &map[string]*string{
			"TABLE_NAME": table.TableName(),
		},
	})

	table.GrantReadWriteData(fn)

	api := awsapigateway.NewLambdaRestApi(stack, jsii.String("HelloCdkGoApi"), &awsapigateway.LambdaRestApiProps{
		Handler:     fn,
		Description: jsii.String("API Gateway for Hello CDK Go Lambda"),
		DeployOptions: &awsapigateway.StageOptions{
			StageName: jsii.String("prod"),
		},
	})

	contactResource := api.Root().AddResource(jsii.String("contact"), nil)
	contactResource.AddMethod(jsii.String("POST"), nil, nil)

	awscdk.NewCfnOutput(stack, jsii.String("ApiEndpoint"), &awscdk.CfnOutputProps{
		Value:       api.Url(),
		Description: jsii.String("URL of the API Gateway"),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewHelloCdkGoStack(app, "HelloCdkGoStack", &HelloCdkGoStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}
