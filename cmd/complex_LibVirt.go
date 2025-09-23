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
	cinp "github.com/cinp/go"
	"github.com/spf13/cobra"
)

var complexLibVirtCmd = &cobra.Command{
	Use:   "libvirt",
	Short: "Work with LibVirt Complexes",
}

var complexLibVirtGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get LibVirt Complexes",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		complexID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.LibvirtLibVirtComplexGet(ctx, complexID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:                 {{.GetURI | extractID}}
Name:               {{.Name}}
Description:        {{.Description}}
Type:               {{.Type}}
State:              {{.State}}
Site:               {{.Site | extractID}}
Member:             {{.Members}}
Created:            {{.Created}}
Updated:            {{.Updated}}
`)

		return nil
	},
}

var complexLibVirtCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New LibVirt Complex",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		o := contractorClient.LibvirtLibVirtComplexNew()
		o.BuiltPercentage = nil
		o.Name = &detailName
		o.Description = &detailDescription

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			o.Site = cinp.StringAddr(r.GetURI())
		}

		if detailMember != 0 {

			s, err := contractorClient.BuildingStructureGet(ctx, detailMember)
			if err != nil {
				return err
			}
			o.Members = &[]string{s.GetURI()}
		}

		err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:                 {{.GetURI | extractID}}
Name:               {{.Name}}
Description:        {{.Description}}
Type:               {{.Type}}
State:              {{.State}}
Site:               {{.Site | extractID}}
Member:             {{.Members}}
Created:            {{.Created}}
Updated:            {{.Updated}}
`)

		return nil
	},
}

var complexLibVirtUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update LibVirt Complex",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		complexID := args[0]

		ctx := cmd.Context()

		o := contractorClient.LibvirtLibVirtComplexNewWithID(complexID)

		if detailDescription != "" {
			o.Description = &detailDescription
		}

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			o.Site = cinp.StringAddr(r.GetURI())
		}

		if detailMember != 0 {
			s, err := contractorClient.BuildingStructureGet(ctx, detailMember)
			if err != nil {
				return err
			}
			o.Members = &[]string{s.GetURI()}
		}

		err := o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:                 {{.GetURI | extractID}}
Name:               {{.Name}}
Description:        {{.Description}}
Type:               {{.Type}}
State:              {{.State}}
Site:               {{.Site | extractID}}
Member:             {{.Members}}
Created:            {{.Created}}
Updated:            {{.Updated}}
`)

		return nil
	},
}

func init() {
	complexTypes["libvirt"] = complexTypeEntry{"/api/v1/LibVirt/", "0.1"}

	complexLibVirtCreateCmd.Flags().StringVarP(&detailName, "name", "l", "", "Locator of New LibVirt Complex")
	complexLibVirtCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New LibVirt Complex")
	complexLibVirtCreateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Description of New LibVirt Complex")
	complexLibVirtCreateCmd.Flags().IntVarP(&detailMember, "member", "m", 0, "Members of the new LibVirt Complex")

	complexLibVirtUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Complex with value")
	complexLibVirtUpdateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Update the Description of LibVirt Complex with value")
	complexLibVirtUpdateCmd.Flags().IntVarP(&detailMember, "member", "m", 0, "Update the Member of the LibVirt Complex")

	complexCmd.AddCommand(complexLibVirtCmd)
	complexLibVirtCmd.AddCommand(complexLibVirtGetCmd, complexLibVirtCreateCmd, complexLibVirtUpdateCmd)
}
