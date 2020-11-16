package sdk

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

type TypedServiceRegistration interface {
	// Name is the name of this Service
	Name() string

	// PackagePath is the relative path to this package
	PackagePath() string

	// SupportedResources returns a list of Data Sources supported by this Service
	SupportedDataSources() []DataSource

	// SupportedResources returns a list of Resources supported by this Service
	SupportedResources() []Resource

	// WebsiteCategories returns a list of categories which can be used for the sidebar
	WebsiteCategories() []string
}

type UntypedServiceRegistration interface {
	// Name is the name of this Service
	Name() string

	// WebsiteCategories returns a list of categories which can be used for the sidebar
	WebsiteCategories() []string

	// SupportedDataSources returns the supported Data Sources supported by this Service
	SupportedDataSources() map[string]*schema.Resource

	// SupportedResources returns the supported Resources supported by this Service
	SupportedResources() map[string]*schema.Resource
}
