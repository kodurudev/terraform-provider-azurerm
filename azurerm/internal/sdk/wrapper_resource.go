package sdk

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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
	resourceSchema, err := combineSchema(rw.resource.Arguments(), rw.resource.Attributes())
	if err != nil {
		return nil, fmt.Errorf("building Schema: %+v", err)
	}

	var d = func(duration time.Duration) *time.Duration {
		return &duration
	}

	resource := schema.Resource{
		Schema: *resourceSchema,

		Create: func(d *schema.ResourceData, meta interface{}) error {
			ctx, metaData := runArgs(d, meta, rw.logger)
			wrappedCtx, cancel := timeouts.ForCreate(ctx, d)
			defer cancel()
			err := rw.resource.Create().Func(wrappedCtx, metaData)
			if err != nil {
				return err
			}
			// NOTE: whilst this may look like we should use the Read
			// functions timeout here, we're still /technically/ in the
			// Create function so reusing that timeout should be sufficient
			return rw.resource.Read().Func(wrappedCtx, metaData)
		},

		// looks like these could be reused, easiest if they're not
		Read: func(d *schema.ResourceData, meta interface{}) error {
			ctx, metaData := runArgs(d, meta, rw.logger)
			wrappedCtx, cancel := timeouts.ForRead(ctx, d)
			defer cancel()
			return rw.resource.Read().Func(wrappedCtx, metaData)
		},
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			ctx, metaData := runArgs(d, meta, rw.logger)
			wrappedCtx, cancel := timeouts.ForDelete(ctx, d)
			defer cancel()
			return rw.resource.Delete().Func(wrappedCtx, metaData)
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
			ctx, metaData := runArgs(d, meta, rw.logger)
			wrappedCtx, cancel := timeouts.ForUpdate(ctx, d)
			defer cancel()

			err := v.Update().Func(wrappedCtx, metaData)
			if err != nil {
				return err
			}
			// whilst this may look like we should use the Update timeout here
			// we're still "technically" in the update method, so reusing the
			// Update's timeout should be fine
			return rw.resource.Read().Func(wrappedCtx, metaData)
		}
		resource.Timeouts.Update = d(v.Update().Timeout)
	}

	return &resource, nil
}
