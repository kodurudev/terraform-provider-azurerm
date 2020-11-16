package sdk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func combineSchema(arguments map[string]*schema.Schema, attributes map[string]*schema.Schema) (*map[string]*schema.Schema, error) {
	out := make(map[string]*schema.Schema, 0)

	for k, v := range arguments {
		if _, alreadyExists := out[k]; alreadyExists {
			return nil, fmt.Errorf("%q already exists in the schema", k)
		}

		// TODO: if readonly

		out[k] = v
	}

	for k, v := range attributes {
		if _, alreadyExists := out[k]; alreadyExists {
			return nil, fmt.Errorf("%q already exists in the schema", k)
		}

		// TODO: if editable

		// every attribute has to be computed
		v.Computed = true
		out[k] = v
	}

	return &out, nil
}

func runArgs(d *schema.ResourceData, meta interface{}, logger Logger) (context.Context, ResourceMetaData) {
	// NOTE: this is wrapped as a result of this function, so this is "fine" being unwrapped
	stopContext := meta.(*clients.Client).StopContext
	client := meta.(*clients.Client)
	metaData := ResourceMetaData{
		Client:                   client,
		Logger:                   logger,
		ResourceData:             d,
		serializationDebugLogger: NullLogger{},
	}

	return stopContext, metaData
}
