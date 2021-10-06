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

func cartographerArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Requires a Id Argument")
	}
	return nil
}

var cartographerCmd = &cobra.Command{
	Use:   "cartographer",
	Short: "Work with cartographers",
}

var cartographerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Cartographers",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		rl := []cinp.Object{}
		for v := range c.SurveyCartographerList("", map[string]interface{}{}) {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Identifier", "Message", "Foundation", "Last Checkin", "Created", "Updated"}, "{{.GetID | extractID}}	{{.Message}}	{{.Foundation | extractID}}	{{.LastCheckin}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var cartographerGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Cartographer",
	Args:  cartographerArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		cartographerID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.SurveyCartographerGet(cartographerID)
		if err != nil {
			return err
		}
		outputDetail(r, `Identifier:    {{.Identifier}}
Message:       {{.Message}}
Foundation:    {{.Foundation | extractID}}
Last Checkin:  {{.Created}}
Info Map:      {{.InfoMap}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)
		return nil
	},
}

var cartographerAssignCmd = &cobra.Command{
	Use:   "assign",
	Short: "Assign Foundation to Cartographer",
	Args:  cartographerArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		cartographerID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.SurveyCartographerGet(cartographerID)
		if err != nil {
			return err
		}

		r, err := c.BuildingFoundationGet(detailFoundation)
		if err != nil {
			return err
		}

		err = o.CallAssign(r.GetID())
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	cartographerAssignCmd.Flags().StringVarP(&detailFoundation, "foundation", "f", "", "Foundation to assign")

	rootCmd.AddCommand(cartographerCmd)
	cartographerCmd.AddCommand(cartographerListCmd)
	cartographerCmd.AddCommand(cartographerGetCmd)
	cartographerCmd.AddCommand(cartographerAssignCmd)
}
