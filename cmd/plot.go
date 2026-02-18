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

func plotArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a Plot Id argument")
	}
	return nil
}

var plotCmd = &cobra.Command{
	Use:   "plot",
	Short: "Work with plots",
}

var plotListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Plots",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		rl := []cinp.Object{}
		vchan, err := contractorClient.SurveyPlotList(ctx, "", map[string]interface{}{})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Id", "Name", "Created", "Updated"}, "{{.GetURI | extractID}}	{{.Name}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var plotGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Plot",
	Args:  plotArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		plotID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.SurveyPlotGet(ctx, plotID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:            {{.GetURI | extractID}}
Name:          {{.Name}}
Corners:       {{.Corners}}
Parent:        {{.Parent}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)

		return nil
	},
}

var plotCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Plot",
	RunE: func(cmd *cobra.Command, args []string) error {
		o := contractorClient.SurveyPlotNew()

		ctx := cmd.Context()

		if detailName != "" {
			o.Name = &detailName
		}

		if detailCorners != "" {
			o.Corners = &detailCorners
		}

		if detailParent != "" {
			o.Parent = &detailParent
		}

		err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI | extractID}}
Name:          {{.Name}}
Corners:       {{.Corners}}
Parent:        {{.Parent}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)

		return nil
	},
}

var plotUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Plot",
	Args:  plotArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		plotID := args[0]

		ctx := cmd.Context()

		o := contractorClient.SurveyPlotNewWithID(plotID)

		if detailCorners != "" {
			o.Corners = &detailCorners
		}

		if detailParent != "" {
			o.Parent = &detailParent
		}

		err := o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI | extractID}}
Name:          {{.Name}}
Corners:       {{.Corners}}
Parent:        {{.Parent}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)

		return nil
	},
}

var plotDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Plot",
	Args:  plotArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		plotID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.SurveyPlotGet(ctx, plotID)
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
	plotCreateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Name of New Plot")
	plotCreateCmd.Flags().StringVarP(&detailCorners, "corners", "c", "", "Corners of New Plot")
	plotCreateCmd.Flags().StringVarP(&detailParent, "parent", "p", "", "Parent of New Plot")

	plotUpdateCmd.Flags().StringVarP(&detailCorners, "corners", "s", "", "Update the Corners of the Plot with value")
	plotUpdateCmd.Flags().StringVarP(&detailParent, "parent", "m", "", "Update the Parent of the Plot with the value")

	rootCmd.AddCommand(plotCmd)
	plotCmd.AddCommand(plotListCmd, plotGetCmd, plotCreateCmd, plotUpdateCmd, plotDeleteCmd)

}
