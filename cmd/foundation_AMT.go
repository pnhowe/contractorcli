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

var foundationAMTCmd = &cobra.Command{
	Use:   "amt",
	Short: "Work with AMT Foundations",
}

var foundationAMTGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get AMT Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]
		ctx := cmd.Context()

		o, err := contractorClient.AmtAMTFoundationGet(ctx, foundationID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
AMT Username:   {{.AmtUsername}}
AMT Password:   {{.AmtPassword}}
AMT Ip Address: {{.AmtIPAddress}}
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

var foundationAMTCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New AMT Foundation",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		o := contractorClient.AmtAMTFoundationNew()
		o.Locator = &detailLocator
		o.AmtUsername = &detailAMTUsername
		o.AmtPassword = &detailAMTPassword
		o.AmtIPAddress = &detailAMTIP

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

		if detailPlot != "" {
			r, err := contractorClient.SurveyPlotGet(ctx, detailPlot)
			if err != nil {
				return err
			}
			o.Plot = cinp.StringAddr(r.GetURI())
		}

		err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
AMT Username:   {{.AmtUsername}}
AMT Password:   {{.AmtPassword}}
AMT Ip Address: {{.AmtIPAddress}}
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

var foundationAMTUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update AMT Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o := contractorClient.AmtAMTFoundationNewWithID(foundationID)

		if detailAMTUsername != "" {
			o.AmtUsername = &detailAMTUsername
		}

		if detailAMTPassword != "" {
			o.AmtPassword = &detailAMTPassword
		}

		if detailAMTIP != "" {
			o.AmtIPAddress = &detailAMTIP
		}

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

		if detailPlot != "" {
			r, err := contractorClient.SurveyPlotGet(ctx, detailPlot)
			if err != nil {
				return err
			}
			o.Plot = cinp.StringAddr(r.GetURI())
		}

		err := o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
AMT Username:   {{.AmtUsername}}
AMT Password:   {{.AmtPassword}}
AMT Ip Address: {{.AmtIPAddress}}
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

func init() {
	fundationTypes["amt"] = foundationTypeEntry{"/api/v1/AMT/", "0.1"}

	foundationAMTCreateCmd.Flags().StringVarP(&detailLocator, "locator", "l", "", "Locator of New AMT Foundation")
	foundationAMTCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New AMT Foundation")
	foundationAMTCreateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Blueprint of New AMT Foundation")
	foundationAMTCreateCmd.Flags().StringVarP(&detailPlot, "plot", "p", "", "Plot of New AMT Foundation")
	foundationAMTCreateCmd.Flags().StringVarP(&detailAMTUsername, "amt-username", "u", "", "AMT Username of New AMT Foundation")
	foundationAMTCreateCmd.Flags().StringVarP(&detailAMTPassword, "amt-password", "a", "", "AMT Password of New AMT Foundation")
	foundationAMTCreateCmd.Flags().StringVarP(&detailAMTIP, "amt-ip", "i", "", "AMT Host Ip Address of New AMT Foundation")

	foundationAMTUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Foundation with value")
	foundationAMTUpdateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Update the Blueprint of Foundation with value")
	foundationAMTUpdateCmd.Flags().StringVarP(&detailPlot, "plot", "p", "", "Update the Plot of the AMT Foundation")
	foundationAMTUpdateCmd.Flags().StringVarP(&detailAMTUsername, "amt-username", "u", "", "Update the AMT Username of the AMT Foundation")
	foundationAMTUpdateCmd.Flags().StringVarP(&detailAMTPassword, "amt-password", "a", "", "Update the AMT Password of the AMT Foundation")
	foundationAMTUpdateCmd.Flags().StringVarP(&detailAMTIP, "amt-ip", "i", "", "Update the AMT Host Ip Address of the AMT Foundation")

	foundationCmd.AddCommand(foundationAMTCmd)
	foundationAMTCmd.AddCommand(foundationAMTGetCmd, foundationAMTCreateCmd, foundationAMTUpdateCmd)
}
