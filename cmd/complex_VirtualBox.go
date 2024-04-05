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

var complexVirtualBoxCmd = &cobra.Command{
	Use:   "virtualbox",
	Short: "Work with VirtualBox Complexes",
}

var complexVirtualBoxGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get VirtualBox Complexes",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		complexID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.VirtualboxVirtualBoxComplexGet(ctx, complexID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:                 {{.GetURI | extractID}}
Name:               {{.Name}}
Description:        {{.Description}}
Type:               {{.Type}}
State:              {{.State}}
Site:               {{.Site | extractID}}
BuiltPercentage:    {{.BuiltPercentage}}
Members:            {{.Members}}
VirtualboxUsername: {{.VirtualboxUsername}}
VirtualboxPassword: {{.VirtualboxPassword}}
Created:            {{.Created}}
Updated:            {{.Updated}}
`)

		return nil
	},
}

var complexVirtualBoxCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New VirtualBox Complex",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		o := contractorClient.VirtualboxVirtualBoxComplexNew()
		o.Name = &detailName
		o.Description = &detailDescription
		o.BuiltPercentage = &detailBuiltPercentage
		o.VirtualboxUsername = &detailUsername
		o.VirtualboxPassword = &detailPassword

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			(*o.Site) = r.GetURI()
		}

		for _, v := range detailMembers {
			s, err := contractorClient.BuildingStructureGet(ctx, v)
			if err != nil {
				return err
			}
			(*o.Members) = append((*o.Members), s.GetURI())
		}

		o, err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:                 {{.GetURI | extractID}}
Name:               {{.Name}}
Description:        {{.Description}}
Type:               {{.Type}}
State:              {{.State}}
Site:               {{.Site | extractID}}
BuiltPercentage:    {{.BuiltPercentage}}
Members:            {{.Members}}
VirtualboxUsername: {{.VirtualboxUsername}}
VirtualboxPassword: {{.VirtualboxPassword}}
Created:            {{.Created}}
Updated:            {{.Updated}}
`)

		return nil
	},
}

var complexVirtualBoxUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update VirtualBox Complex",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		complexID := args[0]

		ctx := cmd.Context()

		o := contractorClient.VirtualboxVirtualBoxComplexNewWithID(complexID)

		if detailDescription != "" {
			o.Description = &detailDescription
		}

		if detailBuiltPercentage != 0 {
			o.BuiltPercentage = &detailBuiltPercentage
		}

		if detailUsername != "" {
			o.VirtualboxUsername = &detailUsername
		}

		if detailPassword != "" {
			o.VirtualboxPassword = &detailPassword
		}

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			(*o.Site) = r.GetURI()
		}

		if len(detailMembers) > 0 {
			for _, v := range detailMembers {
				s, err := contractorClient.BuildingStructureGet(ctx, v)
				if err != nil {
					return err
				}
				(*o.Members) = append((*o.Members), s.GetURI())
			}
		}

		o, err := o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:                 {{.GetURI | extractID}}
Name:               {{.Name}}
Description:        {{.Description}}
Type:               {{.Type}}
State:              {{.State}}
Site:               {{.Site | extractID}}
BuiltPercentage:    {{.BuiltPercentage}}
Members:            {{.Members}}
VirtualboxUsername: {{.VirtualboxUsername}}
VirtualboxPassword: {{.VirtualboxPassword}}
Created:            {{.Created}}
Updated:            {{.Updated}}
`)

		return nil
	},
}

func init() {
	complexTypes["virtualbox"] = complexTypeEntry{"/api/v1/VirtualBox/", "0.1"}

	complexVirtualBoxCreateCmd.Flags().StringVarP(&detailName, "name", "l", "", "Locator of New VirtualBox Complex")
	complexVirtualBoxCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New VirtualBox Complex")
	complexVirtualBoxCreateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Description of New VirtualBox Complex")
	complexVirtualBoxCreateCmd.Flags().IntVarP(&detailBuiltPercentage, "builtperc", "b", 80, "Built Percentage of New VirtualBox Complex\n(Percentage of Built Members at which the complex is considered built)")
	complexVirtualBoxCreateCmd.Flags().IntSliceVarP(&detailMembers, "members", "m", []int{}, "Members of the new VirtualBox Complex, specify for each member")
	complexVirtualBoxCreateCmd.Flags().StringVarP(&detailUsername, "username", "u", "", "VirtualBox Username of New VirtualBox Complex")
	complexVirtualBoxCreateCmd.Flags().StringVarP(&detailPassword, "password", "p", "", "VirtualBox Password of New VirtualBox Complex")

	complexVirtualBoxUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Complex with value")
	complexVirtualBoxUpdateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Update the Description of VirtualBox Complex with value")
	complexVirtualBoxUpdateCmd.Flags().IntVarP(&detailBuiltPercentage, "builtperc", "b", 0, "Update the Built Percentage of VirtualBox Complex with value")
	complexVirtualBoxUpdateCmd.Flags().IntSliceVarP(&detailMembers, "members", "m", []int{}, "Update the Members of the VirtualBox Complex, specify for each member")
	complexVirtualBoxUpdateCmd.Flags().StringVarP(&detailUsername, "username", "u", "", "Update the VirtualBox Username of the VirtualBox Complex with value")
	complexVirtualBoxUpdateCmd.Flags().StringVarP(&detailPassword, "password", "p", "", "Update the VirtualBox Password of the VirtualBox Complex with value")

	complexCmd.AddCommand(complexVirtualBoxCmd)
	complexVirtualBoxCmd.AddCommand(complexVirtualBoxGetCmd, complexVirtualBoxCreateCmd, complexVirtualBoxUpdateCmd)
}
