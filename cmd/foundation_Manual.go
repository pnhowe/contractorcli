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

var foundationManualCmd = &cobra.Command{
	Use:   "manual",
	Short: "Work with Manual Foundations",
}

var foundationManualGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Manual Foundation",
	Args:  foundationArgCheck,
	Run: func(cmd *cobra.Command, args []string) {
		foundationID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.ManualManualFoundationGet(foundationID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		outputDetail(r, `Locator:        {{.Locator}}
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

var foundationManualCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Manual Foundation",
	Run: func(cmd *cobra.Command, args []string) {
		c := getContractor()
		defer c.Logout()

		o := c.ManualManualFoundationNew()
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

		if err := o.Create(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var foundationManualUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Manual Foundation",
	Args:  foundationArgCheck,
	Run: func(cmd *cobra.Command, args []string) {
		fieldList := []string{}
		foundationID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.ManualManualFoundationGet(foundationID)
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

		if err := o.Update(fieldList); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	fundationTypes["manual"] = typeEntry{"/api/v1/Manual/", "0.1"}

	foundationManualCreateCmd.Flags().StringVarP(&detailLocator, "locator", "l", "", "Locator of New AMT Foundation")
	foundationManualCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New AMT Foundation")
	foundationManualCreateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Blueprint of New AMT Foundation")

	foundationManualUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Foundation with value")
	foundationManualUpdateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Update the Blueprint of Foundation with value")

	foundationCmd.AddCommand(foundationManualCmd)
	foundationManualCmd.AddCommand(foundationManualGetCmd)
	foundationManualCmd.AddCommand(foundationManualCreateCmd)
	foundationManualCmd.AddCommand(foundationManualUpdateCmd)
}
