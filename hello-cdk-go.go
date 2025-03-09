package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
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

	fn.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})

	table.GrantReadWriteData(fn)

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
