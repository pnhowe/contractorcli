package cmd

/*
Copyright © 2020 Peter Howe <pnhowe@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"github.com/spf13/cobra"
)

var detailSubscriptionID, detailLocation, detailResourceGroup, detailClientID, detailTenantID string

var complexAzureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Work with Azure Complexes",
}

var complexAzureGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Azure Complexes",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		complexID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.AzureAzureComplexGet(complexID)
		if err != nil {
			return err
		}
		outputDetail(r, `Name:                {{.Name}}
Description:         {{.Description}}
Type:                {{.Type}}
State:               {{.State}}
Site:                {{.Site | extractID}}
BuiltPercentage:     {{.BuiltPercentage}}
Members:             {{.Members}}
AzureSubscriptionID: {{.AzureSubscriptionID}}
AzureLocation:       {{.AzureLocation}}
AzureResourceGroup:  {{.AzureResourceGroup}}
AzureClientID:       {{.AzureClientID}}
AzurePassword:       {{.AzurePassword}}
AzureTenantID:       {{.AzureTenantID}}
Created:             {{.Created}}
Updated:             {{.Updated}}
`)

		return nil
	},
}

var complexAzureCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Azure Complex",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		o := c.AzureAzureComplexNew()
		o.Name = detailName
		o.Description = detailDescription
		o.BuiltPercentage = detailBuiltPercentage
		o.AzureSubscriptionID = detailSubscriptionID
		o.AzureLocation = detailLocation
		o.AzureResourceGroup = detailResourceGroup
		o.AzureClientID = detailClientID
		o.AzurePassword = detailPassword
		o.AzureTenantID = detailTenantID

		if detailSite != "" {
			r, err := c.SiteSiteGet(detailSite)
			if err != nil {
				return err
			}
			o.Site = r.GetID()
		}

		for _, v := range detailMembers {
			s, err := c.BuildingStructureGet(v)
			if err != nil {
				return err
			}
			o.Members = append(o.Members, s.GetID())
		}

		if err := o.Create(); err != nil {
			return err
		}

		return nil
	},
}

var complexAzureUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Azure Complex",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		complexID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.AzureAzureComplexGet(complexID)
		if err != nil {
			return err
		}

		if detailDescription != "" {
			o.Description = detailDescription
			fieldList = append(fieldList, "description")
		}

		if detailBuiltPercentage != 0 {
			o.BuiltPercentage = detailBuiltPercentage
			fieldList = append(fieldList, "built_percentage")
		}

		if detailSubscriptionID != "" {
			o.AzureSubscriptionID = detailSubscriptionID
			fieldList = append(fieldList, "azure_subscription_id")
		}

		if detailLocation != "" {
			o.AzureLocation = detailLocation
			fieldList = append(fieldList, "azure_location")
		}

		if detailResourceGroup != "" {
			o.AzureResourceGroup = detailResourceGroup
			fieldList = append(fieldList, "azure_resource_group")
		}

		if detailClientID != "" {
			o.AzureClientID = detailClientID
			fieldList = append(fieldList, "azure_client_id")
		}

		if detailPassword != "" {
			o.AzurePassword = detailPassword
			fieldList = append(fieldList, "azure_password")
		}

		if detailTenantID != "" {
			o.AzureTenantID = detailTenantID
			fieldList = append(fieldList, "azure_tenant_id")
		}

		if detailSite != "" {
			r, err := c.SiteSiteGet(detailSite)
			if err != nil {
				return err
			}
			o.Site = r.GetID()
			fieldList = append(fieldList, "site")
		}

		if len(detailMembers) > 0 {
			for _, v := range detailMembers {
				s, err := c.BuildingStructureGet(v)
				if err != nil {
					return err
				}
				o.Members = append(o.Members, s.GetID())
			}
			fieldList = append(fieldList, "members")
		}

		if err := o.Update(fieldList); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	complexTypes["azure"] = complexTypeEntry{"/api/v1/Azure/", "0.1"}

	complexAzureCreateCmd.Flags().StringVarP(&detailName, "name", "l", "", "Locator of New Azure Complex")
	complexAzureCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New Azure Complex")
	complexAzureCreateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Description of New Azure Complex")
	complexAzureCreateCmd.Flags().IntVarP(&detailBuiltPercentage, "builtperc", "b", 80, "Built Percentage of New Azure Complex\n(Percentage of Built Members at which the complex is considered built)")
	complexAzureCreateCmd.Flags().StringArrayVarP(&detailMembers, "members", "m", []string{}, "Members of the new Azure Complex, specify for each member")
	complexAzureCreateCmd.Flags().StringVarP(&detailUsername, "username", "u", "", "Azure Username of New Azure Complex")
	complexAzureCreateCmd.Flags().StringVarP(&detailPassword, "password", "p", "", "Azure Password of New Azure Complex")

	complexAzureUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Complex with value")
	complexAzureUpdateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Update the Description of Azure Complex with value")
	complexAzureUpdateCmd.Flags().IntVarP(&detailBuiltPercentage, "builtperc", "b", 0, "Update the Built Percentage of Azure Complex with value")
	complexAzureUpdateCmd.Flags().StringArrayVarP(&detailMembers, "members", "m", []string{}, "Update the Members of the Azure Complex, specify for each member")
	complexAzureUpdateCmd.Flags().StringVarP(&detailUsername, "username", "u", "", "Update the Azure Username of the Azure Complex with value")
	complexAzureUpdateCmd.Flags().StringVarP(&detailPassword, "password", "p", "", "Update the Azure Password of the Azure Complex with value")

	complexCmd.AddCommand(complexAzureCmd)
	complexAzureCmd.AddCommand(complexAzureGetCmd)
	complexAzureCmd.AddCommand(complexAzureCreateCmd)
	complexAzureCmd.AddCommand(complexAzureUpdateCmd)
}
