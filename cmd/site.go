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

	cinp "github.com/cinp/go"
	"github.com/spf13/cobra"
	//"github.com/t3kton/contractorcli/contractor"
)

var configSetName, configSetValue, configDeleteName string
var configFull bool

func siteArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires a site id argument")
	}
	return nil
}

var siteCmd = &cobra.Command{
	Use:   "site",
	Short: "Work with sites",
}

var siteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Sites",
	Run: func(cmd *cobra.Command, args []string) {
		c := getContractor()
		//objList := map[string]contractor.SiteSite{}
		//err := c.Cinp().GetMulti("/api/v1/Site/Site:site1:", &objList)
		rl := []cinp.Object{}
		for v := range c.SiteSiteList("", map[string]interface{}{}) {
			rl = append(rl, v)
		}
		outputList(rl, "Id	Name	Description	Created	Updated\n", "{{.GetID}}	{{.Name}}	{{.Description}}	{{.Created}}	{{.Updated}}\n")
	},
}

var siteGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Site",
	Args:  siteArgCheck,
	Run: func(cmd *cobra.Command, args []string) {
		siteID := args[0]
		c := getContractor()
		r, err := c.SiteSiteGet(siteID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		outputDetail(r, `Site:          {{.Name}}
Description:   {{.Description}}
Parent:        {{.Parent}}
Zone:          {{.Zone}}
Config Values: {{.ConfigValues}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)
	},
}

var siteConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Work With Site Config",
	Args:  siteArgCheck,
	Run: func(cmd *cobra.Command, args []string) {
		siteID := args[0]
		c := getContractor()
		if configSetName != "" {
			o, err := c.SiteSiteGet(siteID)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			o.ConfigValues[configSetName] = configSetValue
			err = o.Update([]string{"config_values"})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			outputKV(o.ConfigValues)
		} else if configDeleteName != "" {
			o, err := c.SiteSiteGet(siteID)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			delete(o.ConfigValues, configDeleteName)
			err = o.Update([]string{"config_values"})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			outputKV(o.ConfigValues)
		} else if configFull {
			o := c.SiteSiteNewWithID(siteID)
			r, err := o.CallGetConfig()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			outputKV(r)
		} else {
			o, err := c.SiteSiteGet(siteID)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			outputKV(o.ConfigValues)
		}
	},
}

func init() {
	siteConfigCmd.Flags().BoolVarP(&configFull, "full", "f", false, "Display the Full/Compiled config")
	siteConfigCmd.Flags().StringVarP(&configSetName, "set-name", "n", "", "Set Config Value Key Name, if set-value is not specified, the value will be set to ''")
	siteConfigCmd.Flags().StringVarP(&configSetValue, "set-value", "v", "", "Set Config Value, ignored if set-name is not specified")
	siteConfigCmd.Flags().StringVarP(&configDeleteName, "delete", "d", "", "Delete Config Value Key Name")
	rootCmd.AddCommand(siteCmd)
	siteCmd.AddCommand(siteListCmd)
	siteCmd.AddCommand(siteGetCmd)
	siteCmd.AddCommand(siteConfigCmd)
}
