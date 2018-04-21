package keypair

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const (
	KeyName              = "KeyName"
	KeyFingerprint       = "KeyFingerprint"
	ParameterName        = "ParameterName"
	ParameterPath        = "ParameterPath"
	ParameterKeyID       = "ParameterKeyId"
	ParameterDescription = "ParameterDescription"
	ParameterOverwrite   = "ParameterOverwrite"
)

var (
	DefaultParameterPath      = "/ec2/key-pair"
	DefaultParameterKeyID     = "alias/aws/ssm"
	DefaultParameterOverwrite = false
)

// Resource represents the request/response details of an EC2 key-pair
type Resource struct {
	PhysicalID     string
	KeyName        string `json:"KeyName"`
	ParameterName  string `json:"ParameterName"`
	ParameterKeyID string `json:"ParameterKeyId"`
}

// Source the resource attribute values from the properties map.
func (res *Resource) Source(properties map[string]interface{}) {
	// KeyName
	res.KeyName = StringValue(properties, KeyName, res.PhysicalID)

	// ParameterName
	if path := StringValue(properties, ParameterPath, DefaultParameterPath); path == "" {
		res.ParameterName = res.KeyName
	} else {
		res.ParameterName = fmt.Sprintf("%s/%s", path, res.KeyName)
	}

	// ParameterKeyID
	res.ParameterKeyID = StringValue(properties, ParameterKeyID, DefaultParameterKeyID)
}

// Create the resource.
func (res *Resource) Create(ctx context.Context, evt *cfn.Event) (err error) {
	ses, err := session.NewSession()
	if err != nil {
		return
	}

	out, err := ec2.New(ses).CreateKeyPairWithContext(ctx, &ec2.CreateKeyPairInput{
		KeyName: aws.String(res.KeyName),
	})
	if err != nil {
		return
	}

	parameterDescription := StringValue(evt.ResourceProperties, ParameterDescription, aws.StringValue(out.KeyFingerprint))
	parameterOverwrite := BoolValue(evt.ResourceProperties, ParameterOverwrite, DefaultParameterOverwrite)

	_, err = ssm.New(ses).PutParameterWithContext(ctx, &ssm.PutParameterInput{
		Type:        aws.String(ssm.ParameterTypeSecureString),
		Name:        aws.String(res.ParameterName),
		Value:       out.KeyMaterial,
		Description: aws.String(parameterDescription),
		KeyId:       aws.String(res.ParameterKeyID),
		Overwrite:   aws.Bool(parameterOverwrite),
	})

	return
}

// Update the resource.
func (res *Resource) Update(ctx context.Context, evt *cfn.Event) (err error) {
	if okn, nkn := evt.OldResourceProperties[KeyName], evt.ResourceProperties[KeyName]; reflect.DeepEqual(okn, nkn) {
		log.Printf("%s: %s has not changed, skipping %s", evt.RequestID, KeyName, evt.RequestType)
	} else if !reflect.DeepEqual(evt.ResourceProperties, evt.OldResourceProperties) {
		res.PhysicalID = NewPhysicalResourceID(evt)
		res.KeyName = ""
		res.ParameterName = ""
		res.ParameterKeyID = ""
		res.Source(evt.ResourceProperties)
		err = res.Create(ctx, evt)
	}

	return
}

// Delete the resource.
func (res *Resource) Delete(ctx context.Context, evt *cfn.Event) (err error) {
	ses, err := session.NewSession()
	if err != nil {
		return
	}

	if res.KeyName != "" {
		_, e := ec2.New(ses).DeleteKeyPairWithContext(ctx, &ec2.DeleteKeyPairInput{
			KeyName: aws.String(res.KeyName),
		})
		if e != nil {
			log.Print(e)
		}
	}
	if res.ParameterName != "" {
		_, e := ssm.New(ses).DeleteParameterWithContext(ctx, &ssm.DeleteParameterInput{
			Name: aws.String(res.ParameterName),
		})
		if e != nil {
			log.Print(e)
		}
	}

	return
}

// NewPhysicalResourceIDFunc is the NewPhysicalResourceID function type.
type NewPhysicalResourceIDFunc func(*cfn.Event) string

// NewPhysicalResourceID generates a randomized resource id.
// As this is a var, you may replace it with your own implementation.
var NewPhysicalResourceID NewPhysicalResourceIDFunc = func(evt *cfn.Event) string {
	rns := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	gen := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd := make([]byte, 12)
	for i := range rnd {
		rnd[i] = rns[gen.Intn(len(rns))]
	}
	stack := strings.Split(evt.StackID, "/")[1]
	return fmt.Sprintf("%s-%s-%s", stack, evt.LogicalResourceID, rnd)
}
