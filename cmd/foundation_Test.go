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

var foundationTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Work with Test Foundations",
}

var foundationTestGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Test Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.TestTestFoundationGet(ctx, foundationID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
Delay Variance  {{.TestDelayVariance}}
Fail Likelihood {{.TestFailLikelihood}}
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

var foundationTestCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Test Foundation",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		o := contractorClient.TestTestFoundationNew()
		o.Locator = &detailLocator

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			site := r.GetURI()
			o.Site = &site
		}

		if detailBlueprint != "" {
			r, err := contractorClient.BlueprintFoundationBluePrintGet(ctx, detailBlueprint)
			if err != nil {
				return err
			}
			blueprint := r.GetURI()
			o.Blueprint = &blueprint
		}

		if detailDelayVariance != -1 {
			o.TestDelayVariance = &detailDelayVariance
		}

		if detailFailLikelihood != -1 {
			o.TestFailLikelihood = &detailFailLikelihood
		}

		o, err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
Delay Variance  {{.TestDelayVariance}}
Fail Likelihood {{.TestFailLikelihood}}
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

var foundationTestUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Test Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {

		foundationID := args[0]

		ctx := cmd.Context()

		o := contractorClient.TestTestFoundationNewWithID(foundationID)

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			site := r.GetURI()
			o.Site = &site
		}

		if detailBlueprint != "" {
			r, err := contractorClient.BlueprintFoundationBluePrintGet(ctx, detailBlueprint)
			if err != nil {
				return err
			}
			blueprint := r.GetURI()
			o.Blueprint = &blueprint
		}

		if detailDelayVariance != -1 {
			o.TestDelayVariance = &detailDelayVariance
		}

		if detailFailLikelihood != -1 {
			o.TestFailLikelihood = &detailFailLikelihood
		}

		o, err := o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
Delay Variance  {{.TestDelayVariance}}
Fail Likelihood {{.TestFailLikelihood}}
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
	fundationTypes["test"] = foundationTypeEntry{"/api/v1/Test/", "0.1"}

	foundationTestCreateCmd.Flags().StringVarP(&detailLocator, "locator", "l", "", "Locator of New Test Foundation")
	foundationTestCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New Test Foundation")
	foundationTestCreateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Blueprint of New Test Foundation")
	foundationTestCreateCmd.Flags().IntVarP(&detailDelayVariance, "delay", "d", -1, "The Variance of operations, in seconds")
	foundationTestCreateCmd.Flags().IntVarP(&detailFailLikelihood, "fail", "f", -1, "Likelyhood of failures per 1000, per 'run' execution (ie: the longer the delay the more likely), if greater than 9, Unrecoverable errors my happen, 0 to disable")

	foundationTestUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Foundation with value")
	foundationTestUpdateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Update the Blueprint of Foundation with value")
	foundationTestUpdateCmd.Flags().IntVarP(&detailDelayVariance, "delay", "d", -1, "The Variance of operations, in seconds")
	foundationTestUpdateCmd.Flags().IntVarP(&detailFailLikelihood, "fail", "f", -1, "Likelyhood of failures per 1000, per 'run' execution (ie: the longer the delay the more likely), if greater than 9, Unrecoverable errors my happen, 0 to disable")

	foundationCmd.AddCommand(foundationTestCmd)
	foundationTestCmd.AddCommand(foundationTestGetCmd, foundationTestCreateCmd, foundationTestUpdateCmd)
}
