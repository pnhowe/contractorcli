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

var configSetName, configSetValue, configDeleteName string
var configFull bool
var detailHostname, detailSite, detailBlueprint, detailFoundation string
var jobInfo, jobState, jobCreate, jobDestroy bool
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
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		rl := []cinp.Object{}
		for v := range c.BuildingStructureList("", map[string]interface{}{}) {
			rl = append(rl, v)
		}
		outputList(rl, "Id	Site	Hostname	Foundation	Blueprint	Created	Updated\n", "{{.GetID | extractID}}	{{.Site | extractID}}	{{.Hostname}}	{{.Foundation | extractID}}	{{.Blueprint | extractID}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var structureGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Structure",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
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

		return nil
	},
}

var structureCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Structure",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		o := c.BuildingStructureNew()
		o.Hostname = detailHostname

		if detailSite != "" {
			r, err := c.SiteSiteGet(detailSite)
			if err != nil {
				return err
			}
			o.Site = r.GetID()
		}

		if detailBlueprint != "" {
			r, err := c.BlueprintStructureBluePrintGet(detailBlueprint)
			if err != nil {
				return err
			}
			o.Blueprint = r.GetID()
		}

		if detailFoundation != "" {
			r, err := c.BuildingFoundationGet(detailFoundation)
			if err != nil {
				return err
			}
			o.Foundation = r.GetID()
		}

		if err := o.Create(); err != nil {
			return err
		}

		return nil
	},
}

var structureUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Structure",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		structureID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
		}

		if detailHostname != "" {
			o.Hostname = detailHostname
			fieldList = append(fieldList, "hostname")
		}

		if detailSite != "" {
			r, err := c.SiteSiteGet(detailSite)
			if err != nil {
				return err
			}
			o.Site = r.GetID()
			fieldList = append(fieldList, "site")
		}

		if detailBlueprint != "" {
			r, err := c.BlueprintStructureBluePrintGet(detailBlueprint)
			if err != nil {
				return err
			}
			o.Blueprint = r.GetID()
			fieldList = append(fieldList, "blueprint")
		}

		if detailFoundation != "" {
			r, err := c.BuildingFoundationGet(detailFoundation)
			if err != nil {
				return err
			}
			o.Foundation = r.GetID()
			fieldList = append(fieldList, "foundation")
		}

		if err := o.Update(fieldList); err != nil {
			return err
		}

		return nil
	},
}

var structureDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Structure",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
		}
		if err := r.Delete(); err != nil {
			return err
		}

		return nil
	},
}

var structureConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Work With Structure Config",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID := args[0]
		c := getContractor()
		defer c.Logout()

		if configSetName != "" {
			o, err := c.BuildingStructureGet(structureID)
			if err != nil {
				return err
			}
			o.ConfigValues[configSetName] = configSetValue
			if err := o.Update([]string{"config_values"}); err != nil {
				return err
			}
			outputKV(o.ConfigValues)

		} else if configDeleteName != "" {
			o, err := c.BuildingStructureGet(structureID)
			if err != nil {
				return err
			}
			delete(o.ConfigValues, configDeleteName)
			if err := o.Update([]string{"config_values"}); err != nil {
				return err
			}
			outputKV(o.ConfigValues)

		} else if configFull {
			o := c.BuildingStructureNewWithID(structureID)
			r, err := o.CallGetConfig()
			if err != nil {
				return err
			}
			outputKV(r)

		} else {
			o, err := c.BuildingStructureGet(structureID)
			if err != nil {
				return err
			}
			outputKV(o.ConfigValues)
		}

		return nil
	},
}

var structureJobCmd = &cobra.Command{
	Use:   "job",
	Short: "Work with Structure Jobs",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
		}

		if jobInfo {
			jID, err := o.CallGetJob()
			if err != nil {
				return err
			}
			j, err := c.ForemanStructureJobGet(extractID(jID))
			if err != nil {
				return err
			}
			outputDetail(j, `Site           {{.Site}}
Structure      {{.Structure}}
State:         {{.State}}
Status:        {{.Status}}
Progress:      {{.Progress}}
Message:       {{.Message}}
Script Name:   {{.ScriptName}}
Can Start:     {{.CanStart}}
Updated:       {{.Updated}}
Created:       {{.Created}}
`)
		} else if jobState {
			jID, err := o.CallGetJob()
			if err != nil {
				return err
			}
			j, err := c.ForemanStructureJobGet(extractID(jID))
			if err != nil {
				return err
			}
			vars, err := j.CallJobRunnerVariables()
			if err != nil {
				return err
			}
			state, err := j.CallJobRunnerState()
			if err != nil {
				return err
			}
			outputDetail(map[string]interface{}{"variables": vars, "state": state}, `Variables: {{.variables}}
Script State: {{.state.state}}
Script Line No: {{.state.cur_line}}
-- Script --
{{.state.script}}`)
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
				return err
			}
			outputKV(map[string]interface{}{"job": j})
		}

		return nil
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
	structureJobCmd.Flags().BoolVarP(&jobState, "state", "s", false, "Show Running Job State")
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
