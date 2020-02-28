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
)

var configSetName, configSetValue, configDeleteName string
var configFull bool
var detailHostname, detailSite, detailBlueprint, detailFoundation string
var jobInfo, jobCreate, jobDestroy bool
var jobUtility string

func structureArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("Requires a Structure Id Argument")
	}
	return nil
}

var structureCmd = &cobra.Command{
	Use:   "structure",
	Short: "Work with Structures",
}

var structureListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Structures",
	Run: func(cmd *cobra.Command, args []string) {
		c := getContractor()
		defer c.Logout()

		rl := []cinp.Object{}
		for v := range c.BuildingStructureList("", map[string]interface{}{}) {
			rl = append(rl, v)
		}
		outputList(rl, "id	Site	Hostname	Foundation	Blueprint	Created	Updated\n", "{{.GetID | extractID}}	{{.Site | extractID}}	{{.Hostname}}	{{.Blueprint | extractID}}	{{.Foundation | extractID}}	{{.Created}}	{{.Updated}}\n")
	},
}

var structureGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Structure",
	Args:  structureArgCheck,
	Run: func(cmd *cobra.Command, args []string) {
		structureID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.BuildingStructureGet(structureID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		outputDetail(r, `
Hostname:      {{.Hostname}}
Site:          {{.Site | extractID}}
Blueprint:     {{.Blueprint | extractID}}
Foundation:    {{.Foundation | extractID}}
Config UUID:   {{.ConfigUUID}}
Config Values: {{.ConfigValues}}
State:         {{.State}}
Built At:      {{.BuiltAt}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)
	},
}

var structureCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Structure",
	Run: func(cmd *cobra.Command, args []string) {
		c := getContractor()
		defer c.Logout()

		o := c.BuildingStructureNew()
		o.Hostname = detailHostname

		if detailSite != "" {
			r, err := c.SiteSiteGet(detailSite)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			o.Site = r.GetID()
		}

		if detailBlueprint != "" {
			r, err := c.BlueprintStructureBluePrintGet(detailBlueprint)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			o.Blueprint = r.GetID()
		}

		if detailFoundation != "" {
			r, err := c.BuildingFoundationGet(detailFoundation)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			o.Foundation = r.GetID()
		}

		if err := o.Create(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var structureUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Structure",
	Args:  structureArgCheck,
	Run: func(cmd *cobra.Command, args []string) {
		fieldList := []string{}
		structureID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.BuildingStructureGet(structureID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if detailHostname != "" {
			o.Hostname = detailHostname
			fieldList = append(fieldList, "hostname")
		}

		if detailSite != "" {
			r, err := c.SiteSiteGet(detailSite)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			o.Site = r.GetID()
			fieldList = append(fieldList, "site")
		}

		if detailBlueprint != "" {
			r, err := c.BlueprintStructureBluePrintGet(detailBlueprint)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			o.Blueprint = r.GetID()
			fieldList = append(fieldList, "blueprint")
		}

		if detailFoundation != "" {
			r, err := c.BuildingFoundationGet(detailFoundation)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			o.Foundation = r.GetID()
			fieldList = append(fieldList, "foundation")
		}

		if err := o.Update(fieldList); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var structureDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Structure",
	Args:  structureArgCheck,
	Run: func(cmd *cobra.Command, args []string) {
		structureID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.BuildingStructureGet(structureID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := r.Delete(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var structureConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Work With Structure Config",
	Args:  structureArgCheck,
	Run: func(cmd *cobra.Command, args []string) {
		structureID := args[0]
		c := getContractor()
		defer c.Logout()

		if configSetName != "" {
			o, err := c.BuildingStructureGet(structureID)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			o.ConfigValues[configSetName] = configSetValue
			if err := o.Update([]string{"config_values"}); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			outputKV(o.ConfigValues)

		} else if configDeleteName != "" {
			o, err := c.BuildingStructureGet(structureID)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			delete(o.ConfigValues, configDeleteName)
			if err := o.Update([]string{"config_values"}); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			outputKV(o.ConfigValues)

		} else if configFull {
			o := c.BuildingStructureNewWithID(structureID)
			r, err := o.CallGetConfig()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			outputKV(r)

		} else {
			o, err := c.BuildingStructureGet(structureID)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			outputKV(o.ConfigValues)
		}
	},
}

var structureJobCmd = &cobra.Command{
	Use:   "job",
	Short: "Work with Structure Jobs",
	Args:  structureArgCheck,
	Run: func(cmd *cobra.Command, args []string) {
		structureID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.BuildingStructureGet(structureID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if jobInfo {
			j, err := o.CallGetJob()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			outputDetail(j, `Site           {{.Site}}
Structure      {{.Structure}}
State:         {{.State}}
Status:        {{.Status}}
Message:       {{.Message}}
Script Name:   {{.ScriptName}}
Script Runner: {{.ScriptRunner}}
Updated:       {{.Updated}}
Created:       {{.Created}}
`)
		} else {
			var j int
			if jobCreate {
				j, err = o.CallDoCreate()
			} else if jobDestroy {
				j, err = o.CallDoDestroy()
			} else if jobUtility != "" {
				j, err = o.CallDoJob(jobUtility)
			}
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			outputKV(map[string]interface{}{"job": j})
		}
	},
}

func init() {
	structureConfigCmd.Flags().BoolVarP(&configFull, "full", "f", false, "Display the Full/Compiled config")
	structureConfigCmd.Flags().StringVarP(&configSetName, "set-name", "n", "", "Set Config Value Key Name, if set-value is not specified, the value will be set to ''")
	structureConfigCmd.Flags().StringVarP(&configSetValue, "set-value", "v", "", "Set Config Value, ignored if set-name is not specified")
	structureConfigCmd.Flags().StringVarP(&configDeleteName, "delete", "d", "", "Delete Config Value Key Name")

	structureCreateCmd.Flags().StringVarP(&detailHostname, "hostname", "o", "", "Hostname of New Structure")
	structureCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New Structure")
	structureCreateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Blueprint of New Structure")
	structureCreateCmd.Flags().StringVarP(&detailFoundation, "foundation", "f", "", "Foundation of New Structure")

	structureUpdateCmd.Flags().StringVarP(&detailHostname, "hostname", "o", "", "Update the Hostname of Structure with value")
	structureUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Structure with value")
	structureUpdateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Update the Blueprint of Structure with value")
	structureUpdateCmd.Flags().StringVarP(&detailFoundation, "foundation", "f", "", "Update the Foundation of Structure with value")

	structureJobCmd.Flags().BoolVarP(&jobInfo, "info", "i", false, "Show Running Job Info")
	structureJobCmd.Flags().BoolVarP(&jobCreate, "do-create", "c", false, "Submit a Create job")
	structureJobCmd.Flags().BoolVarP(&jobDestroy, "do-destroy", "d", false, "Submit a Destroy job")
	structureJobCmd.Flags().StringVarP(&jobUtility, "utility", "u", "", "Submit Utility Job")

	rootCmd.AddCommand(structureCmd)
	structureCmd.AddCommand(structureListCmd)
	structureCmd.AddCommand(structureGetCmd)
	structureCmd.AddCommand(structureCreateCmd)
	structureCmd.AddCommand(structureUpdateCmd)
	structureCmd.AddCommand(structureDeleteCmd)
	structureCmd.AddCommand(structureConfigCmd)
	structureCmd.AddCommand(structureJobCmd)
}
