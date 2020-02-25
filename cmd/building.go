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
	"fmt"
	"os"
	"strconv"
	"text/template"

	"github.com/spf13/cobra"
)

var structureCmd = &cobra.Command{
	Use:   "structure",
	Short: "Work with structures",
}

var structureListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Structures",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("List Structures")
	},
}

var structureGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Structure",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a site id argument")
		}
		if _, err := strconv.Atoi(args[0]); err != nil {
			return errors.New("Invalid Id, must be a valid int")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		fmt.Printf("Getting Structure '%s'\n", id)
		c := getContractor()
		r, err := c.BuildingStructureGet(id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		t, err := template.New("structure detail").Parse(`
Hostname:      {{.Hostname}}
Site:          {{.Site}}
Blueprint:     {{.Blueprint}}
Foundation:    {{.Foundation}}
Config UUID:   {{.ConfigUUID}}
Config Values: {{.ConfigValues}}
State:         {{.State}}
Built At:      {{.BuiltAt}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = t.Execute(os.Stdout, r)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var structureGetConfig = &cobra.Command{
	Use:   "getConfig",
	Short: "Get Structure Config",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a structure id argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		siteID := args[0]
		fmt.Printf("Getting Structure '%s'\n", siteID)
		c := getContractor()
		structure := c.BuildingStructureNewWithID(siteID)
		r, err := structure.CallGetConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for k, v := range r {
			fmt.Printf("%s:\t%+v\n", k, v)
		}
	},
}

func init() {
	rootCmd.AddCommand(structureCmd)
	structureCmd.AddCommand(structureListCmd)
	structureCmd.AddCommand(structureGetCmd)
	structureCmd.AddCommand(structureGetConfig)
}
