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

var foundationLibVirtCmd = &cobra.Command{
	Use:   "libvirt",
	Short: "Work with LibVirt Foundations",
}

var foundationLibVirtGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get LibVirt Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.LibvirtLibVirtFoundationGet(ctx, foundationID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
Complex:        {{.LibvirtComplex | extractID}}
VM UUID:        {{.LibvirtUUID}}
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

var foundationLibVirtCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New LibVirt Foundation",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		o := contractorClient.LibvirtLibVirtFoundationNew()
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
			r, err := contractorClient.LibvirtLibVirtComplexGet(ctx, detailComplex)
			if err != nil {
				return err
			}
			o.LibvirtComplex = cinp.StringAddr(r.GetURI())
		}

		err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
Complex:        {{.LibvirtComplex | extractID}}
VM UUID:        {{.LibvirtUUID}}
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

var foundationLibVirtUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update LibVirt Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o := contractorClient.LibvirtLibVirtFoundationNewWithID(foundationID)

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
			r, err := contractorClient.LibvirtLibVirtComplexGet(ctx, detailComplex)
			if err != nil {
				return err
			}
			o.LibvirtComplex = cinp.StringAddr(r.GetURI())
		}

		err := o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
Complex:        {{.LibvirtComplex | extractID}}
VM UUID:        {{.LibvirtUUID}}
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
	fundationTypes["libvirt"] = foundationTypeEntry{"/api/v1/LibVirt/", "0.1"}

	foundationLibVirtCreateCmd.Flags().StringVarP(&detailLocator, "locator", "l", "", "Locator of New LibVirt Foundation")
	foundationLibVirtCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New LibVirt Foundation")
	foundationLibVirtCreateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Blueprint of New LibVirt Foundation")
	foundationLibVirtCreateCmd.Flags().StringVarP(&detailComplex, "complex", "x", "", "Plot of New LibVirt Foundation")

	foundationLibVirtUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Foundation with value")
	foundationLibVirtUpdateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Update the Blueprint of Foundation with value")
	foundationLibVirtUpdateCmd.Flags().StringVarP(&detailComplex, "complex", "x", "", "Update the Plot of the LibVirt Foundation")

	foundationCmd.AddCommand(foundationLibVirtCmd)
	foundationLibVirtCmd.AddCommand(foundationLibVirtGetCmd, foundationLibVirtCreateCmd, foundationLibVirtUpdateCmd)
}
