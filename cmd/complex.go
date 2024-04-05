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
	"errors"

	cinp "github.com/cinp/go"
	"github.com/spf13/cobra"
)

type complexTypeEntry struct {
	URI        string // Namespace URI, for checking to see if it's loaded and the API version
	APIVersion string
}

var complexTypes = map[string]complexTypeEntry{}

func complexArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a Complex Id(Name) argument")
	}
	return nil
}

var complexCmd = &cobra.Command{
	Use:   "complex",
	Short: "Work with Complexs",
}

var complexListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Complexs",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		rl := []cinp.Object{}
		vchan, err := contractorClient.BuildingComplexList(ctx, "", map[string]interface{}{})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Id", "Site", "Name", "State", "Type", "Created", "Updated"}, "{{.GetURI | extractID}}	{{.Site | extractID}}	{{.Name}}	{{.State}}	{{.Type}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var complexGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Complex",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		complexID := args[0]
		ctx := cmd.Context()

		o, err := contractorClient.BuildingComplexGet(ctx, complexID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:            {{.GetURI | extractID}}
Name:          {{.Name}}
Description:   {{.Description}}
Site:          {{.Site | extractID}}
Type:          {{.Type}}
State:         {{.State}}
Members:       {{.Members}}
Built %:       {{.BuiltPercentage}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)

		return nil
	},
}

var complexTypesCmd = &cobra.Command{
	Use:   "types",
	Short: "List Supported Types",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		typeList := []string{}
		for k, v := range complexTypes {
			APIVersion, err := contractorClient.GetAPIVersion(ctx, v.URI)
			if err != nil {
				// return err
				continue // TODO: really should only do this if it is a 404
			}
			if APIVersion != v.APIVersion {
				continue // API Version mismatch, print a warning?
			}
			typeList = append(typeList, k)
		}
		outputKV(map[string]interface{}{"type": typeList})

		return nil
	},
}

var complexDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Complex",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		complexID := args[0]
		ctx := cmd.Context()

		o, err := contractorClient.BuildingComplexGet(ctx, complexID)
		if err != nil {
			return err
		}
		if err := o.Delete(ctx); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(complexCmd)
	complexCmd.AddCommand(complexListCmd, complexGetCmd, complexTypesCmd, complexDeleteCmd)
}
