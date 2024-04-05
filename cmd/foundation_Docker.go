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

var foundationDockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Work with Docker Foundations",
}

var foundationDockerGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Docker Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.DockerDockerFoundationGet(ctx, foundationID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
Complex:        {{.DockerComplex | extractID}}
Container Id:   {{.DockerID}}
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

var foundationDockerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Docker Foundation",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		o := contractorClient.DockerDockerFoundationNew()
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
			r, err := contractorClient.DockerDockerComplexGet(ctx, detailComplex)
			if err != nil {
				return err
			}
			(*o.DockerComplex) = r.GetURI()
		}

		o, err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
Complex:        {{.DockerComplex | extractID}}
Container Id:   {{.DockerID}}
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

var foundationDockerUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Docker Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o := contractorClient.DockerDockerFoundationNewWithID(foundationID)

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
			r, err := contractorClient.DockerDockerComplexGet(ctx, detailComplex)
			if err != nil {
				return err
			}
			(*o.DockerComplex) = r.GetURI()
		}

		o, err := o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:        {{.Locator}}
Complex:        {{.DockerComplex | extractID}}
Container Id:   {{.DockerID}}
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
	fundationTypes["docker"] = foundationTypeEntry{"/api/v1/Docker/", "0.1"}

	foundationDockerCreateCmd.Flags().StringVarP(&detailLocator, "locator", "l", "", "Locator of New Docker Foundation")
	foundationDockerCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New Docker Foundation")
	foundationDockerCreateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Blueprint of New Docker Foundation")
	foundationDockerCreateCmd.Flags().StringVarP(&detailComplex, "complex", "x", "", "Plot of New Docker Foundation")

	foundationDockerUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Foundation with value")
	foundationDockerUpdateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Update the Blueprint of Foundation with value")
	foundationDockerUpdateCmd.Flags().StringVarP(&detailComplex, "complex", "x", "", "Update the Plot of the Docker Foundation")

	foundationCmd.AddCommand(foundationDockerCmd)
	foundationDockerCmd.AddCommand(foundationDockerGetCmd, foundationDockerCreateCmd, foundationDockerUpdateCmd)
}
