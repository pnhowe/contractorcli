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

var foundationProxmoxCmd = &cobra.Command{
	Use:   "proxmox",
	Short: "Work with Proxmox Foundations",
}

var foundationProxmoxGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Proxmox Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.ProxmoxProxmoxFoundationGet(ctx, foundationID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
Complex:        {{.ProxmoxComplex | extractID}}
VM ID:          {{.ProxmoxID}}
Type:           {{.Type}}
Site:           {{.Site | extractID}}
Blueprint:      {{.Blueprint | extractID}}
Id Map:         {{.IDMap}}
Class List:     {{.ClassList}}
State:          {{.State}}
Located At:     {{.LocatedAt}}
Built At:       {{.BuiltAt}}
Created:        {{.Created}}
Updated:        {{.Updated}}
`)

		return nil
	},
}

var foundationProxmoxCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Proxmox Foundation",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		o := contractorClient.ProxmoxProxmoxFoundationNew()
		o.Locator = &detailLocator

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			o.Site = cinp.StringAddr(r.GetURI())
		}

		if detailBlueprint != "" {
			r, err := contractorClient.BlueprintFoundationBluePrintGet(ctx, detailBlueprint)
			if err != nil {
				return err
			}
			o.Blueprint = cinp.StringAddr(r.GetURI())
		}

		if detailComplex != "" {
			r, err := contractorClient.ProxmoxProxmoxComplexGet(ctx, detailComplex)
			if err != nil {
				return err
			}
			o.ProxmoxComplex = cinp.StringAddr(r.GetURI())
		}

		err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
Complex:        {{.ProxmoxComplex | extractID}}
VM ID:          {{.ProxmoxID}}
Type:           {{.Type}}
Site:           {{.Site | extractID}}
Blueprint:      {{.Blueprint | extractID}}
Id Map:         {{.IDMap}}
Class List:     {{.ClassList}}
State:          {{.State}}
Located At:     {{.LocatedAt}}
Built At:       {{.BuiltAt}}
Created:        {{.Created}}
Updated:        {{.Updated}}
`)
		return nil
	},
}

var foundationProxmoxUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Proxmox Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o := contractorClient.ProxmoxProxmoxFoundationNewWithID(foundationID)

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			o.Site = cinp.StringAddr(r.GetURI())
		}

		if detailBlueprint != "" {
			r, err := contractorClient.BlueprintFoundationBluePrintGet(ctx, detailBlueprint)
			if err != nil {
				return err
			}
			o.Blueprint = cinp.StringAddr(r.GetURI())
		}

		if detailComplex != "" {
			r, err := contractorClient.ProxmoxProxmoxComplexGet(ctx, detailComplex)
			if err != nil {
				return err
			}
			o.ProxmoxComplex = cinp.StringAddr(r.GetURI())
		}

		err := o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
Complex:        {{.ProxmoxComplex | extractID}}
VM ID:          {{.ProxmoxID}}
Type:           {{.Type}}
Site:           {{.Site | extractID}}
Blueprint:      {{.Blueprint | extractID}}
Id Map:         {{.IDMap}}
Class List:     {{.ClassList}}
State:          {{.State}}
Located At:     {{.LocatedAt}}
Built At:       {{.BuiltAt}}
Created:        {{.Created}}
Updated:        {{.Updated}}
`)

		return nil
	},
}

func init() {
	fundationTypes["proxmox"] = foundationTypeEntry{"/api/v1/Proxmox/", "0.1"}

	foundationProxmoxCreateCmd.Flags().StringVarP(&detailLocator, "locator", "l", "", "Locator of New Proxmox Foundation")
	foundationProxmoxCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New Proxmox Foundation")
	foundationProxmoxCreateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Blueprint of New Proxmox Foundation")
	foundationProxmoxCreateCmd.Flags().StringVarP(&detailComplex, "complex", "x", "", "Plot of New Proxmox Foundation")

	foundationProxmoxUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Foundation with value")
	foundationProxmoxUpdateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Update the Blueprint of Foundation with value")
	foundationProxmoxUpdateCmd.Flags().StringVarP(&detailComplex, "complex", "x", "", "Update the Plot of the Proxmox Foundation")

	foundationCmd.AddCommand(foundationProxmoxCmd)
	foundationProxmoxCmd.AddCommand(foundationProxmoxGetCmd, foundationProxmoxCreateCmd, foundationProxmoxUpdateCmd)
}
