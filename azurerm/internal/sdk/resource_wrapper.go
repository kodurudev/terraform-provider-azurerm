package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
)

type ResourceWrapper struct {
	logger   Logger
	resource Resource
}

func NewResourceWrapper(resource Resource) ResourceWrapper {
	return ResourceWrapper{
		logger:   ConsoleLogger{},
		resource: resource,
	}
}

func (rw *ResourceWrapper) Resource() (*schema.Resource, error) {
	resourceSchema, err := rw.schema()
	if err != nil {
		return nil, fmt.Errorf("building Schema: %+v", err)
	}

	var d = func(duration time.Duration) *time.Duration {
		return &duration
	}

	resource := schema.Resource{
		Schema: *resourceSchema,

		Create: func(d *schema.ResourceData, meta interface{}) error {
			ctx, metaData := rw.runArgs(d, meta)
			err := rw.resource.Create().Func(ctx, metaData)
			if err != nil {
				return err
			}
			return rw.resource.Read().Func(ctx, metaData)
		},

		// looks like these could be reused, easiest if they're not
		Read: func(d *schema.ResourceData, meta interface{}) error {
			ctx, metaData := rw.runArgs(d, meta)
			return rw.resource.Read().Func(ctx, metaData)
		},
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			ctx, metaData := rw.runArgs(d, meta)
			return rw.resource.Delete().Func(ctx, metaData)
		},

		Timeouts: &schema.ResourceTimeout{
			Create: d(rw.resource.Create().Timeout),
			Read:   d(rw.resource.Read().Timeout),
			Delete: d(rw.resource.Delete().Timeout),
		},
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			fn := rw.resource.IDValidationFunc()
			warnings, errors := fn(id, "id")
			if len(warnings) > 0 {
				for _, warning := range warnings {
					rw.logger.Warn(warning)
				}
			}
			if len(errors) > 0 {
				out := ""
				for _, error := range errors {
					out += error.Error()
				}
				return fmt.Errorf(out)
			}

			return err
		}),
	}

	if v, ok := rw.resource.(ResourceWithUpdate); ok {
		resource.Update = func(d *schema.ResourceData, meta interface{}) error {
			ctx, metaData := rw.runArgs(d, meta)
			err := v.Update().Func(ctx, metaData)
			if err != nil {
				return err
			}
			return rw.resource.Read().Func(ctx, metaData)
		}
		resource.Timeouts.Update = d(v.Update().Timeout)
	}

	return &resource, nil
}

func (rw ResourceWrapper) runArgs(d *schema.ResourceData, meta interface{}) (context.Context, ResourceMetaData) {
	ctx := meta.(*clients.Client).StopContext
	client := meta.(*clients.Client)
	metaData := ResourceMetaData{
		Client:       client,
		Logger:       rw.logger,
		ResourceData: d,
	}

	return ctx, metaData
}

func (rw ResourceWrapper) schema() (*map[string]*schema.Schema, error) {
	out := make(map[string]*schema.Schema, 0)

	for k, v := range rw.resource.Arguments() {
		if _, alreadyExists := out[k]; alreadyExists {
			return nil, fmt.Errorf("%q already exists in the schema", k)
		}

		// TODO: if readonly

		out[k] = v
	}

	for k, v := range rw.resource.Attributes() {
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
