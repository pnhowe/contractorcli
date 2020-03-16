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

var detailIPMIUsername, detailIPMIPassword, detailIPMIIP string

var foundationIPMICmd = &cobra.Command{
	Use:   "ipmi",
	Short: "Work with IPMI Foundations",
}

var foundationIPMIGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get IPMI Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.IpmiIPMIFoundationGet(foundationID)
		if err != nil {
			return err
		}
		outputDetail(r, `Locator:        {{.Locator}}
IPMI Username:   {{.IpmiUsername}}
IPMI Password:   {{.IpmiPassword}}
IPMI Ip Address: {{.IpmiIPAddress}}
Plot:           {{.Plot | extractID}}
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

var foundationIPMICreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New IPMI Foundation",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		o := c.IpmiIPMIFoundationNew()
		o.Locator = detailLocator
		o.IpmiUsername = detailIPMIUsername
		o.IpmiPassword = detailIPMIPassword
		o.IpmiIPAddress = detailIPMIIP

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

		return nil
	},
}

var foundationIPMIUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update IPMI Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		foundationID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.IpmiIPMIFoundationGet(foundationID)
		if err != nil {
			return err
		}

		if detailIPMIUsername != "" {
			o.IpmiUsername = detailIPMIUsername
			fieldList = append(fieldList, "ipmi_username")
		}

		if detailIPMIPassword != "" {
			o.IpmiPassword = detailIPMIPassword
			fieldList = append(fieldList, "ipmi_password")
		}

		if detailIPMIIP != "" {
			o.IpmiIPAddress = detailIPMIIP
			fieldList = append(fieldList, "ipmi_ip_address")
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
	fundationTypes["ipmi"] = foundationTypeEntry{"/api/v1/IPMI/", "0.1"}

	foundationIPMICreateCmd.Flags().StringVarP(&detailLocator, "locator", "l", "", "Locator of New IPMI Foundation")
	foundationIPMICreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New IPMI Foundation")
	foundationIPMICreateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Blueprint of New IPMI Foundation")
	foundationIPMICreateCmd.Flags().StringVarP(&detailPlot, "plot", "p", "", "Plot of New IPMI Foundation")
	foundationIPMICreateCmd.Flags().StringVarP(&detailIPMIUsername, "ipmi-username", "u", "", "IPMI Username of New IPMI Foundation")
	foundationIPMICreateCmd.Flags().StringVarP(&detailIPMIPassword, "ipmi-password", "a", "", "IPMI Password of New IPMI Foundation")
	foundationIPMICreateCmd.Flags().StringVarP(&detailIPMIIP, "ipmi-ip", "i", "", "IPMI Host Ip Address of New IPMI Foundation")

	foundationIPMIUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Foundation with value")
	foundationIPMIUpdateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Update the Blueprint of Foundation with value")
	foundationIPMIUpdateCmd.Flags().StringVarP(&detailPlot, "plot", "p", "", "Update the Plot of the IPMI Foundation")
	foundationIPMIUpdateCmd.Flags().StringVarP(&detailIPMIUsername, "ipmi-username", "u", "", "Update the IPMI Username of the IPMI Foundation")
	foundationIPMIUpdateCmd.Flags().StringVarP(&detailIPMIPassword, "ipmi-password", "a", "", "Update the IPMI Password of the IPMI Foundation")
	foundationIPMIUpdateCmd.Flags().StringVarP(&detailIPMIIP, "ipmi-ip", "i", "", "Update the IPMI Host Ip Address of the IPMI Foundation")

	foundationCmd.AddCommand(foundationIPMICmd)
	foundationIPMICmd.AddCommand(foundationIPMIGetCmd)
	foundationIPMICmd.AddCommand(foundationIPMICreateCmd)
	foundationIPMICmd.AddCommand(foundationIPMIUpdateCmd)
}
