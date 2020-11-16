package example

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/sdk"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Example"
}

// PackagePath is the relative path to this package
func (r Registration) PackagePath() string {
	return "./azurerm/internal/services/example"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Example",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() []sdk.Resource {
	return []sdk.Resource{
		ExampleResource{},
	}
}
