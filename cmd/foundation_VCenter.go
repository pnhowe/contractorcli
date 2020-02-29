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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var foundationVCenterCmd = &cobra.Command{
	Use:   "virtualbox",
	Short: "Work with VCenter Foundations",
}

var foundationVCenterGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get VCenter Foundation",
	Args:  foundationArgCheck,
	Run: func(cmd *cobra.Command, args []string) {
		foundationID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.VcenterVCenterFoundationGet(foundationID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		outputDetail(r, `Locator:        {{.Locator}}
Complex:        {{.VcenterComplex | extractID}}
VM UUID:        {{.VcenterUUID}}
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
	},
}

var foundationVCenterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New VCenter Foundation",
	Run: func(cmd *cobra.Command, args []string) {
		c := getContractor()
		defer c.Logout()

		o := c.VcenterVCenterFoundationNew()
		o.Locator = detailLocator

		if detailSite != "" {
			r, err := c.SiteSiteGet(detailSite)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			o.Site = r.GetID()
		}

		if detailBlueprint != "" {
			r, err := c.BlueprintFoundationBluePrintGet(detailBlueprint)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			o.Blueprint = r.GetID()
		}

		if detailComplex != "" {
			r, err := c.VcenterVCenterFoundationGet(detailComplex)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			o.VcenterComplex = r.GetID()
		}

		if err := o.Create(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var foundationVCenterUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update VCenter Foundation",
	Args:  foundationArgCheck,
	Run: func(cmd *cobra.Command, args []string) {
		fieldList := []string{}
		foundationID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.VcenterVCenterFoundationGet(foundationID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if detailSite != "" {
			r, err := c.SiteSiteGet(detailSite)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			o.Site = r.GetID()
			fieldList = append(fieldList, "site")
		}

		if detailBlueprint != "" {
			r, err := c.BlueprintFoundationBluePrintGet(detailBlueprint)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			o.Blueprint = r.GetID()
			fieldList = append(fieldList, "blueprint")
		}

		if detailComplex != "" {
			r, err := c.VcenterVCenterComplexGet(detailComplex)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			o.VcenterComplex = r.GetID()
			fieldList = append(fieldList, "vcenter_complex")
		}

		if err := o.Update(fieldList); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	fundationTypes["vcenter"] = typeEntry{"/api/v1/VCenter/", "0.1"}

	foundationVCenterCreateCmd.Flags().StringVarP(&detailLocator, "locator", "l", "", "Locator of New VCenter Foundation")
	foundationVCenterCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New VCenter Foundation")
	foundationVCenterCreateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Blueprint of New VCenter Foundation")
	foundationVCenterCreateCmd.Flags().StringVarP(&detailComplex, "complex", "x", "", "Plot of New VCenter Foundation")

	foundationVCenterUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Foundation with value")
	foundationVCenterUpdateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Update the Blueprint of Foundation with value")
	foundationVCenterUpdateCmd.Flags().StringVarP(&detailComplex, "complex", "x", "", "Update the Plot of the VCenter Foundation")

	foundationCmd.AddCommand(foundationVCenterCmd)
	foundationVCenterCmd.AddCommand(foundationVCenterGetCmd)
	foundationVCenterCmd.AddCommand(foundationVCenterCreateCmd)
	foundationVCenterCmd.AddCommand(foundationVCenterUpdateCmd)
}
