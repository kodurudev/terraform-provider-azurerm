package sdk

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"

func (rmd ResourceMetaData) SetID(formatter resourceid.Formatter) {
	subscriptionId := rmd.Client.Account.SubscriptionId
	rmd.ResourceData.SetId(formatter.ID(subscriptionId))
}
