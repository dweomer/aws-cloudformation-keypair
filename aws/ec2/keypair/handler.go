package keypair

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/cfn"
)

// HandleEvent ...
func HandleEvent(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
	enc := json.NewEncoder(os.Stdout)

	if e := enc.Encode(event); e != nil {
		log.Printf("event encoding error: %v", e)
	}

	resource := &Resource{
		PhysicalID: event.PhysicalResourceID,
	}
	if resource.PhysicalID == "" {
		resource.PhysicalID = NewPhysicalResourceID(&event)
	}
	resource.Source(event.ResourceProperties)

	switch event.RequestType {
	case cfn.RequestCreate:
		err = resource.Create(ctx, &event)
	case cfn.RequestUpdate:
		err = resource.Update(ctx, &event)
	case cfn.RequestDelete:
		err = resource.Delete(ctx, &event)
	}

	physicalResourceID, data = resource.PhysicalID, map[string]interface{}{
		KeyName:        resource.KeyName,
		ParameterName:  resource.ParameterName,
		ParameterKeyID: resource.ParameterKeyID,
	}

	if e := enc.Encode(data); e != nil {
		log.Printf("event encoding error: %v", e)
	}

	return
}
