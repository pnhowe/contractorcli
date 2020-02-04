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
package cmd

import (
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

var siteCmd = &cobra.Command{
	Use:   "site",
	Short: "Work with sites",
}

var siteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Sites",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("List Sites")
	},
}

var siteGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Site",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a site id argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		siteID := args[0]
		fmt.Printf("Getting Site %s\n", siteID)
		c := getContractor()
		r, err := c.Get(fmt.Sprintf("/api/v1/Site/Site:%s:", siteID))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		t, err := template.New("site detail").Parse(`
Site:          {{.name}}
Description:   {{.description}}
Parent:        {{.parent}}
Zone:          {{.zone}}
Config Values: {{.config_values}}
Created:       {{.created}}
Updated:       {{.updated}}
`)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		r["zone"], err = c.ExtractIds([]string{r["zone"].(string)})[0]
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

var siteGetConfig = &cobra.Command{
	Use:   "get",
	Short: "Get Site Config",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a site id argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		siteID := args[0]
		fmt.Printf("Getting Site %s\n", siteID)
		c := getContractor()
		r, err := c.Call(fmt.Sprintf("/api/v1/Site/Site:%s:(getConfig)", siteID), nil)
		fmt.Println(err)
		fmt.Printf("%+v\n", r)

	},
}

func init() {
	rootCmd.AddCommand(siteCmd)
	siteCmd.AddCommand(siteListCmd)
	siteCmd.AddCommand(siteGetCmd)
	siteCmd.AddCommand(siteGetConfig)
}
