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
	"io"
	"os"

	cinp "github.com/cinp/go"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

func siteArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a Site Id/Name argument")
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
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		rl := []cinp.Object{}
		vchan, err := contractorClient.SiteSiteList(ctx, "", map[string]interface{}{})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Id", "Name", "Description", "Created", "Updated"}, "{{.GetURI | extractID}}	{{.Name}}	{{.Description}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var siteGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Site",
	Args:  siteArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		siteID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.SiteSiteGet(ctx, siteID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:            {{.GetURI | extractID}}
Name:          {{.Name}}
Description:   {{.Description}}
Parent:        {{or .Parent ":<None>:" | extractID}}
Zone:          {{.Zone | extractID}}
Config Values: {{.ConfigValues}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)
		return nil
	},
}

var siteCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Site",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		o := contractorClient.SiteSiteNew()
		o.Name = &detailName
		o.Description = &detailDescription

		if detailParent != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailParent)
			if err != nil {
				return err
			}
			o.Parent = cinp.StringAddr(r.GetURI())
		}

		if detailZone != 0 {
			r, err := contractorClient.DirectoryZoneGet(ctx, detailZone)
			if err != nil {
				return err
			}
			o.Zone = cinp.StringAddr(r.GetURI())
		}

		err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI | extractID}}
Name:          {{.Name}}
Description:   {{.Description}}
Parent:        {{or .Parent ":<None>:" | extractID}}
Zone:          {{.Zone | extractID}}
Config Values: {{.ConfigValues}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)

		return nil
	},
}

var siteUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Site",
	Args:  siteArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		siteID := args[0]

		ctx := cmd.Context()

		o := contractorClient.SiteSiteNewWithID(siteID)

		if detailDescription != "" {
			o.Description = &detailDescription
		}

		if detailParent != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailParent)
			if err != nil {
				return err
			}
			o.Parent = cinp.StringAddr(r.GetURI())
		}

		if detailZone != 0 {
			r, err := contractorClient.DirectoryZoneGet(ctx, detailZone)
			if err != nil {
				return err
			}
			o.Zone = cinp.StringAddr(r.GetURI())
		}

		err := o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI | extractID}}
Name:          {{.Name}}
Description:   {{.Description}}
Parent:        {{or .Parent ":<None>:" | extractID}}
Zone:          {{.Zone | extractID}}
Config Values: {{.ConfigValues}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)

		return nil
	},
}

var siteDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Site",
	Args:  siteArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		siteID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.SiteSiteGet(ctx, siteID)
		if err != nil {
			return err
		}
		if err := o.Delete(ctx); err != nil {
			return err
		}

		return nil
	},
}

var siteConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Work With Site Config",
	Args:  siteArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		siteID := args[0]

		ctx := cmd.Context()

		if configSetName != "" || configFile != "" || configDeleteName != "" {
			r, err := contractorClient.SiteSiteGet(ctx, siteID)
			if err != nil {
				return err
			}

			o := contractorClient.SiteSiteNewWithID(siteID)
			o.ConfigValues = r.ConfigValues

			if configSetName != "" {
				(*o.ConfigValues)[configSetName] = configSetValue

			} else if configFile != "" {
				var reader io.Reader
				if configFile == "-" {
					reader = os.Stdin
				} else {
					f, err := os.Open(configFile)
					if err != nil {
						return err
					}
					defer f.Close()
					reader = f
				}

				var newValues map[string]interface{}
				decoder := toml.NewDecoder(reader)
				err := decoder.Decode(&newValues)
				if err != nil {
					return err
				}

				for k, v := range newValues {
					(*o.ConfigValues)[k] = v
				}

			} else if configDeleteName != "" {
				delete(*o.ConfigValues, configDeleteName)
			}

			err = o.Update(ctx)
			if err != nil {
				return err
			}
			outputKV(*o.ConfigValues)

		} else if configFull {
			o := contractorClient.SiteSiteNewWithID(siteID)
			r, err := o.CallGetConfig(ctx)
			if err != nil {
				return err
			}
			outputKV(r)

		} else {
			o, err := contractorClient.SiteSiteGet(ctx, siteID)
			if err != nil {
				return err
			}
			outputKV(*o.ConfigValues)
		}
		return nil
	},
}

func init() {
	siteConfigCmd.Flags().BoolVarP(&configFull, "full", "f", false, "Display the Full/Compiled config")
	siteConfigCmd.Flags().StringVarP(&configSetName, "set-name", "n", "", "Set Config Value Key Name, if set-value is not specified, the value will be set to ''")
	siteConfigCmd.Flags().StringVarP(&configSetValue, "set-value", "v", "", "Set Config Value, ignored if set-name is not specified")
	siteConfigCmd.Flags().StringVarP(&configDeleteName, "delete", "d", "", "Delete Config Value Key Name")
	siteConfigCmd.Flags().StringVarP(&configFile, "file", "i", "", "Load Values from file in TOML format, this will be merged with the existing config, '-' for reading from stdin")

	siteCreateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Name of New Site")
	siteCreateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Description of New Site")
	siteCreateCmd.Flags().StringVarP(&detailParent, "parent", "p", "", "Parent of New Site")
	siteCreateCmd.Flags().IntVarP(&detailZone, "zone", "z", 0, "Zone of New Site")

	siteUpdateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Update the Description of Site with value")
	siteUpdateCmd.Flags().StringVarP(&detailParent, "parent", "p", "", "Update the Parent of Site with value")
	siteUpdateCmd.Flags().IntVarP(&detailZone, "zone", "z", 0, "Update the Zone of Site with value")

	rootCmd.AddCommand(siteCmd)
	siteCmd.AddCommand(siteListCmd, siteGetCmd, siteCreateCmd, siteUpdateCmd, siteDeleteCmd, siteConfigCmd)
}
