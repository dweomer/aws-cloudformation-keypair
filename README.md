# Custom::KeyPair

A [custom resource][aws-cfn-custom-resources] Lambda function for creating EC2 key-pairs, written in [Go][golang],
and suitable for [direct][aws-cfn-custom-lambda] invocation by CloudFormation. The private key material is pushed
into an associated SSM Parameter.

## Rationale

Well, at the time of this writing the EC2 key-pair is not a [supported resource type in CloudFormation][aws-cfn-resource-types] and I really wanted it to be.
You see, for demonstration purposes, I'm a big fan of as-self-contained-as-possible [infrastructure definitions][iac-wiki].
I really hate specifying parameters for my templates (everything should have a default).
Because, have you seen the `aws cloudformation` CLI for specifying parameters? **`/me shudders`**

Additionally, for those teams that aren't yet spun up on or are otherwise unable to leverage [Terraform](https://terraform.io) or other
[infrastructure-as-code][iac-book] development tools, this implementation requires no tooling other than the [AWS CLI][aws-cli]
and optionally the [SAM][aws-sam] [CLI][aws-sam-local] for testing.

## Design

### Declarative EC2 Key-Pair With Key Material Dumped into an SSM Parameter (SecureString)

*AS A* developer of infrastructure
*I WANT* to create SSH key-pairs for EC2 instances by declaring such in a CloudFormation template
*SO THAT* when applying said template I am not required to have first created, out of band, EC2 key-pair(s).

#### Input Properties

* optional `KeyName` the EC2 KeyPair name
  * if not specified, this will be generated (see [NewPhysicalResourceID](aws/ec2/keypair/resource.go#L136))

* optional `ParameterPath` the SSM Parameter name prefix
  * if not specified, this will default to `/ec2/key-pair`

* optional `ParameterKeyId` which represents the encryption key used to encipher the private key material
  * if not specified, this will default to `alias/aws/ssm`

* optional `ParameterDescription`
  * if not specified, this will default to value of the key fingerprint

* optional `ParameterOverwrite` determines if a parameter with the existing name with be overwritten with a new version
  * if not specified, this will default to `false`

#### Output Attributes

* [!Ref][aws-cfn-intrinsic-ref] `KeyName`
* [Fn::GetAtt][aws-cfn-intrinsic-getatt] `ParameterName`
* [Fn::GetAtt][aws-cfn-intrinsic-getatt] `ParameterKeyId`

## TODO

* support [indirect invocation via SNS][aws-cfn-custom-sns]
* support alternative methods for handling the private key material, such as:
  * cipher-text as an attribute, suitable for use in an output (NoEcho?)
  * Simple Storage Service (S3)
  * Secrets Manager
  * HTTP PUT

---

[iac-wiki]: <https://en.wikipedia.org/wiki/Infrastructure_as_Code> "Infrastructure as Code"
[iac-book]: <https://info.thoughtworks.com/Infrastructure-as-Code-Kief-Morris.html> "Infrastructure as Code, by Kief Morris"
[aws-sdk-golang]: <https://github.com/aws/aws-sdk-go> "AWS SDK for Go"
[aws-lambda-golang]: <https://github.com/aws/aws-lambda-go> "AWS Lambda for Go"
[aws-cfn]: <https://aws.amazon.com/cloudformation> "AWS CloudFormation"
[aws-cli]: <http://docs.aws.amazon.com/cli/> "AWS CLI"
[aws-sam]: <https://docs.aws.amazon.com/lambda/latest/dg/serverless_app.html> "AWS Serverless Application Model"
[aws-sam-local]: <https://github.com/awslabs/aws-sam-local> "AWS SAM Local"
[aws-cfn-resource-types]: <https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-template-resource-type-ref.html> "AWS Resource Types"
[aws-cfn-custom-resources]: <https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-custom-resources.html> "AWS CloudFormation Custom Resources"
[aws-cfn-custom-lambda]: <https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-custom-resources-lambda.html> "AWS Lambda"
[aws-cfn-custom-sns]: <https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-custom-resources-sns.html> "AWS Simple Notification Service aka SNS"
[aws-cfn-intrinsic-ref]: <https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-ref.html> "Ref"
[aws-cfn-intrinsic-getatt]: <https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html> "Fn::GetAtt"
[aws-resource-property-types-name]: <https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-name.html> "Resource Property Types"
[golang]: <https://golang.org> "The Go Programming Language"
