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

var foundationVirtualBoxCmd = &cobra.Command{
	Use:   "virtualbox",
	Short: "Work with VirtualBox Foundations",
}

var foundationVirtualBoxGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get VirtualBox Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.VirtualboxVirtualBoxFoundationGet(ctx, foundationID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
Complex:        {{.VirtualboxComplex | extractID}}
VM UUID:        {{.VirtualboxUUID}}
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

var foundationVirtualBoxCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New VirtualBox Foundation",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		o := contractorClient.VirtualboxVirtualBoxFoundationNew()
		o.Locator = &detailLocator

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			(*o.Site) = r.GetURI()
		}

		if detailBlueprint != "" {
			r, err := contractorClient.BlueprintFoundationBluePrintGet(ctx, detailBlueprint)
			if err != nil {
				return err
			}
			(*o.Blueprint) = r.GetURI()
		}

		if detailComplex != "" {
			r, err := contractorClient.VirtualboxVirtualBoxComplexGet(ctx, detailComplex)
			if err != nil {
				return err
			}
			(*o.VirtualboxComplex) = r.GetURI()
		}

		o, err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
Complex:        {{.VirtualboxComplex | extractID}}
VM UUID:        {{.VirtualboxUUID}}
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

var foundationVirtualBoxUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update VirtualBox Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o := contractorClient.VirtualboxVirtualBoxFoundationNewWithID(foundationID)

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			(*o.Site) = r.GetURI()
		}

		if detailBlueprint != "" {
			r, err := contractorClient.BlueprintFoundationBluePrintGet(ctx, detailBlueprint)
			if err != nil {
				return err
			}
			(*o.Blueprint) = r.GetURI()
		}

		if detailComplex != "" {
			r, err := contractorClient.VirtualboxVirtualBoxComplexGet(ctx, detailComplex)
			if err != nil {
				return err
			}
			(*o.VirtualboxComplex) = r.GetURI()
		}

		o, err := o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
Complex:        {{.VirtualboxComplex | extractID}}
VM UUID:        {{.VirtualboxUUID}}
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
	fundationTypes["virtualbox"] = foundationTypeEntry{"/api/v1/VirtualBox/", "0.1"}

	foundationVirtualBoxCreateCmd.Flags().StringVarP(&detailLocator, "locator", "l", "", "Locator of New VirtualBox Foundation")
	foundationVirtualBoxCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New VirtualBox Foundation")
	foundationVirtualBoxCreateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Blueprint of New VirtualBox Foundation")
	foundationVirtualBoxCreateCmd.Flags().StringVarP(&detailComplex, "complex", "x", "", "Plot of New VirtualBox Foundation")

	foundationVirtualBoxUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Foundation with value")
	foundationVirtualBoxUpdateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Update the Blueprint of Foundation with value")
	foundationVirtualBoxUpdateCmd.Flags().StringVarP(&detailComplex, "complex", "x", "", "Update the Plot of the VirtualBox Foundation")

	foundationCmd.AddCommand(foundationVirtualBoxCmd)
	foundationVirtualBoxCmd.AddCommand(foundationVirtualBoxGetCmd, foundationVirtualBoxCreateCmd, foundationVirtualBoxUpdateCmd)
}
