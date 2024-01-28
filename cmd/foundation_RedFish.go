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

var detailRedFishUsername, detailRedFishPassword, detailRedFishIP, detailRedFishSOL string

var foundationRedFishCmd = &cobra.Command{
	Use:   "redfish",
	Short: "Work with RedFish Foundations",
}

var foundationRedFishGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get RedFish Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.RedfishRedFishFoundationGet(foundationID)
		if err != nil {
			return err
		}
		outputDetail(r, `Locator:         {{.Locator}}
RedFish Username:   {{.RedfishUsername}}
RedFish Password:   {{.RedfishPassword}}
RedFish Ip Address: {{.RedfishIPAddress}}
RedFish SOL Port:   {{.RedfishSolPort}}
Plot:            {{.Plot | extractID}}
Type:            {{.Type}}
Site:            {{.Site | extractID}}
Blueprint:       {{.Blueprint | extractID}}
Id Map:          {{.IDMap}}
Class List:      {{.ClassList}}
State:           {{.State}}
Located At:      {{.LocatedAt}}
Built At:        {{.BuiltAt}}
Created:         {{.Created}}
Updated:         {{.Updated}}
`)

		return nil
	},
}

var foundationRedFishCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New RedFish Foundation",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		o := c.RedfishRedFishFoundationNew()
		o.Locator = detailLocator
		o.RedfishUsername = detailRedFishUsername
		o.RedfishPassword = detailRedFishPassword
		o.RedfishIPAddress = detailRedFishIP
		o.RedfishSolPort = detailRedFishSOL

		if detailSite != "" {
			r, err := c.SiteSiteGet(detailSite)
			if err != nil {
				return err
			}
			o.Site = r.GetID()
		}

		if detailBlueprint != "" {
			r, err := c.BlueprintFoundationBluePrintGet(detailBlueprint)
			if err != nil {
				return err
			}
			o.Blueprint = r.GetID()
		}

		if detailPlot != "" {
			r, err := c.SurveyPlotGet(detailPlot)
			if err != nil {
				return err
			}
			o.Plot = r.GetID()
		}

		if err := o.Create(); err != nil {
			return err
		}

		outputKV(map[string]interface{}{"id": o.GetID()})

		return nil
	},
}

var foundationRedFishUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update RedFish Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		foundationID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.RedfishRedFishFoundationGet(foundationID)
		if err != nil {
			return err
		}

		if detailRedFishUsername != "" {
			o.RedfishUsername = detailRedFishUsername
			fieldList = append(fieldList, "redfish_username")
		}

		if detailRedFishPassword != "" {
			o.RedfishPassword = detailRedFishPassword
			fieldList = append(fieldList, "redfish_password")
		}

		if detailRedFishIP != "" {
			o.RedfishIPAddress = detailRedFishIP
			fieldList = append(fieldList, "redfish_ip_address")
		}

		if detailRedFishSOL != "" {
			o.RedfishSolPort = detailRedFishSOL
			fieldList = append(fieldList, "redfish_sol_port")
		}

		if detailSite != "" {
			r, err := c.SiteSiteGet(detailSite)
			if err != nil {
				return err
			}
			o.Site = r.GetID()
			fieldList = append(fieldList, "site")
		}

		if detailBlueprint != "" {
			r, err := c.BlueprintFoundationBluePrintGet(detailBlueprint)
			if err != nil {
				return err
			}
			o.Blueprint = r.GetID()
			fieldList = append(fieldList, "blueprint")
		}

		if detailPlot != "" {
			r, err := c.SurveyPlotGet(detailPlot)
			if err != nil {
				return err
			}
			o.Plot = r.GetID()
			fieldList = append(fieldList, "plot")
		}

		if err := o.Update(fieldList); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	fundationTypes["redfish"] = foundationTypeEntry{"/api/v1/RedFish/", "0.1"}

	foundationRedFishCreateCmd.Flags().StringVarP(&detailLocator, "locator", "l", "", "Locator of New RedFish Foundation")
	foundationRedFishCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New RedFish Foundation")
	foundationRedFishCreateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Blueprint of New RedFish Foundation")
	foundationRedFishCreateCmd.Flags().StringVarP(&detailPlot, "plot", "p", "", "Plot of New RedFish Foundation")
	foundationRedFishCreateCmd.Flags().StringVarP(&detailRedFishUsername, "redfish-username", "u", "", "RedFish Username of New RedFish Foundation")
	foundationRedFishCreateCmd.Flags().StringVarP(&detailRedFishPassword, "redfish-password", "a", "", "RedFish Password of New RedFish Foundation")
	foundationRedFishCreateCmd.Flags().StringVarP(&detailRedFishIP, "redfish-ip", "i", "", "RedFish Host Ip Address of New RedFish Foundation")
	foundationRedFishCreateCmd.Flags().StringVarP(&detailRedFishSOL, "redfish-sol", "o", "", "RedFish SOL Port (console/ttyS1/etc) of New RedFish Foundation")

	foundationRedFishUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Foundation with value")
	foundationRedFishUpdateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Update the Blueprint of Foundation with value")
	foundationRedFishUpdateCmd.Flags().StringVarP(&detailPlot, "plot", "p", "", "Update the Plot of the RedFish Foundation")
	foundationRedFishUpdateCmd.Flags().StringVarP(&detailRedFishUsername, "redfish-username", "u", "", "Update the RedFish Username of the RedFish Foundation")
	foundationRedFishUpdateCmd.Flags().StringVarP(&detailRedFishPassword, "redfish-password", "a", "", "Update the RedFish Password of the RedFish Foundation")
	foundationRedFishUpdateCmd.Flags().StringVarP(&detailRedFishIP, "redfish-ip", "i", "", "Update the RedFish Host Ip Address of the RedFish Foundation")
	foundationRedFishUpdateCmd.Flags().StringVarP(&detailRedFishSOL, "redfish-sol", "o", "", "Update the RedFish SOL Port (console/ttyS1/etc) of the RedFish Foundation")

	foundationCmd.AddCommand(foundationRedFishCmd)
	foundationRedFishCmd.AddCommand(foundationRedFishGetCmd)
	foundationRedFishCmd.AddCommand(foundationRedFishCreateCmd)
	foundationRedFishCmd.AddCommand(foundationRedFishUpdateCmd)
}
