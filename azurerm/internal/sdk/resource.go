package sdk

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

// TODO: Docs

type DataSource interface {
	Arguments() map[string]*schema.Schema
	Attributes() map[string]*schema.Schema
	ResourceType() string
	Read() ResourceFunc
}

type Resource interface {
	Arguments() map[string]*schema.Schema
	Attributes() map[string]*schema.Schema

	ResourceType() string

	Create() ResourceFunc
	Read() ResourceFunc
	Delete() ResourceFunc
	IDValidationFunc() schema.SchemaValidateFunc
}

// TODO: resource with state migration
// TODO: a generic state migration for updating ID's

type ResourceWithUpdate interface {
	Update() ResourceFunc
}

type ResourceRunFunc func(ctx context.Context, metadata ResourceMetaData) error

type ResourceFunc struct {
	Func ResourceRunFunc

	// Timeout is the default timeout, which can be overridden by users
	// for this method - in-turn used for the Azure API
	Timeout time.Duration
}

type ResourceMetaData struct {
	Client       *clients.Client
	Logger       Logger
	ResourceData *schema.ResourceData
}
