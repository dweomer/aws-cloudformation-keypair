package main

import (
	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dweomer/aws-cloudformation-keypair/aws/ec2/keypair"
)

func main() {
	lambda.Start(cfn.LambdaWrap(keypair.HandleEvent))
}
