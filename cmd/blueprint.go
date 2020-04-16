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

func blueprintArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("Requires a Blueprint Id/Name Argument")
	}
	return nil
}

var blueprintCmd = &cobra.Command{
	Use:   "blueprint",
	Short: "Work with blueprints",
}

var blueprintFoundationCmd = &cobra.Command{
	Use:   "foundation",
	Short: "Work with Foundation Blueprints",
}

var blueprintFoundationGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Foundation Blueprint",
	Args:  blueprintArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.BlueprintFoundationBluePrintGet(blueprintID)
		if err != nil {
			return err
		}
		outputDetail(r, `Name:                {{.Name}}
Description:         {{.Description}}
Parents:             {{.ParentList}}
Config Values:       {{.ConfigValues}}
Foundation Types:    {{.FoundationTypeList}}
Template:            {{.Template}}
Physical Interfaces: {{.PhysicalInterfaceNames}}
Scripts:             {{.Scripts}}
Created:             {{.Created}}
Updated:             {{.Updated}}
`)
		return nil
	},
}

var blueprintFoundationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Foundation Blueprints",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		rl := []cinp.Object{}
		for v := range c.BlueprintFoundationBluePrintList("", map[string]interface{}{}) {
			rl = append(rl, v)
		}
		outputList(rl, "Id	Name	Description	Created	Updated\n", "{{.GetID | extractID}}	{{.Name}}	{{.Description}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var blueprintFoundationCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Foundation Blueprint",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		o := c.BlueprintFoundationBluePrintNew()
		o.Name = detailName
		o.Description = detailDescription

		// will deal with these later
		// FoundationTypeList []string `json:"foundation_type_list"`
		// Template map[string]interface{} `json:"template"`
		// PhysicalInterfaceNames []string `json:"physical_interface_names"`
		// Scripts []string `json:"scripts"`
		// ParentList []string `json:"parent_list"`

		if err := o.Create(); err != nil {
			return err
		}

		outputKV(map[string]interface{}{"id": o.GetID()})

		return nil
	},
}

var blueprintFoundationUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Foundation Blueprint",
	Args:  blueprintArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		blueprintID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.BlueprintBluePrintGet(blueprintID)
		if err != nil {
			return err
		}
		if detailDescription != "" {
			o.Description = detailDescription
			fieldList = append(fieldList, "description")
		}

		// will deal with these later
		// FoundationTypeList []string `json:"foundation_type_list"`
		// Template map[string]interface{} `json:"template"`
		// PhysicalInterfaceNames []string `json:"physical_interface_names"`
		// Scripts []string `json:"scripts"`
		// ParentList []string `json:"parent_list"`

		if err := o.Update(fieldList); err != nil {
			return err
		}

		return nil
	},
}

var blueprintFoundationDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Foundation Blueprint",
	Args:  blueprintArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.BlueprintFoundationBluePrintGet(blueprintID)
		if err != nil {
			return err
		}
		if err := r.Delete(); err != nil {
			return err
		}

		return nil
	},
}

var blueprintFoundationConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Work With Foundation Blueprint Config",
	Args:  blueprintArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]
		c := getContractor()
		defer c.Logout()

		if configSetName != "" {
			o, err := c.BlueprintFoundationBluePrintGet(blueprintID)
			if err != nil {
				return err
			}
			o.ConfigValues[configSetName] = configSetValue
			if err := o.Update([]string{"config_values"}); err != nil {
				return err
			}
			outputKV(o.ConfigValues)

		} else if configDeleteName != "" {
			o, err := c.BlueprintBluePrintGet(blueprintID)
			if err != nil {
				return err
			}
			delete(o.ConfigValues, configDeleteName)
			if err := o.Update([]string{"config_values"}); err != nil {
				return err
			}
			outputKV(o.ConfigValues)

		} else if configFull {
			o := c.BlueprintBluePrintNewWithID(blueprintID)
			r, err := o.CallGetConfig()
			if err != nil {
				return err
			}
			outputKV(r)

		} else {
			o, err := c.BlueprintBluePrintGet(blueprintID)
			if err != nil {
				return err
			}
			outputKV(o.ConfigValues)
		}
		return nil
	},
}

var blueprintStructureCmd = &cobra.Command{
	Use:   "structure",
	Short: "Work with Structure Blueprints",
}

var blueprintStructureGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Structure Blueprint",
	Args:  blueprintArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.BlueprintFoundationBluePrintGet(blueprintID)
		if err != nil {
			return err
		}
		outputDetail(r, `Name:                  {{.Name}}
Description:           {{.Description}}
Parents:               {{.ParentList}}
Config Values:         {{.ConfigValues}}
Foundation BluePrints: {{.FoundationBlueprintList}}
Scripts:               {{.Scripts}}
Created:               {{.Created}}
Updated:               {{.Updated}}
`)
		return nil
	},
}

var blueprintStructureListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Structure Blueprints",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		rl := []cinp.Object{}
		for v := range c.BlueprintStructureBluePrintList("", map[string]interface{}{}) {
			rl = append(rl, v)
		}
		outputList(rl, "Id	Name	Description	Created	Updated\n", "{{.GetID | extractID}}	{{.Name}}	{{.Description}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var blueprintStructureCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Structure Blueprint",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		o := c.BlueprintStructureBluePrintNew()
		o.Name = detailName
		o.Description = detailDescription

		// will deal with these later
		// Scripts []string `json:"scripts"`
		// ParentList []string `json:"parent_list"`
		// FoundationBlueprintList []string `json:"foundation_blueprint_list"`

		if err := o.Create(); err != nil {
			return err
		}

		outputKV(map[string]interface{}{"id": o.GetID()})

		return nil
	},
}

var blueprintStructureUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Structure Blueprint",
	Args:  blueprintArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		blueprintID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.BlueprintBluePrintGet(blueprintID)
		if err != nil {
			return err
		}
		if detailDescription != "" {
			o.Description = detailDescription
			fieldList = append(fieldList, "description")
		}

		// will deal with these later
		// Scripts []string `json:"scripts"`
		// ParentList []string `json:"parent_list"`
		// FoundationBlueprintList []string `json:"foundation_blueprint_list"`

		if err := o.Update(fieldList); err != nil {
			return err
		}

		return nil
	},
}

var blueprintStructureDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Structure Blueprint",
	Args:  blueprintArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.BlueprintStructureBluePrintGet(blueprintID)
		if err != nil {
			return err
		}
		if err := r.Delete(); err != nil {
			return err
		}

		return nil
	},
}

var blueprintStructureConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Work With Structure Blueprint Config",
	Args:  blueprintArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]
		c := getContractor()
		defer c.Logout()

		if configSetName != "" {
			o, err := c.BlueprintStructureBluePrintGet(blueprintID)
			if err != nil {
				return err
			}
			o.ConfigValues[configSetName] = configSetValue
			if err := o.Update([]string{"config_values"}); err != nil {
				return err
			}
			outputKV(o.ConfigValues)

		} else if configDeleteName != "" {
			o, err := c.BlueprintBluePrintGet(blueprintID)
			if err != nil {
				return err
			}
			delete(o.ConfigValues, configDeleteName)
			if err := o.Update([]string{"config_values"}); err != nil {
				return err
			}
			outputKV(o.ConfigValues)

		} else if configFull {
			o := c.BlueprintBluePrintNewWithID(blueprintID)
			r, err := o.CallGetConfig()
			if err != nil {
				return err
			}
			outputKV(r)

		} else {
			o, err := c.BlueprintBluePrintGet(blueprintID)
			if err != nil {
				return err
			}
			outputKV(o.ConfigValues)
		}
		return nil
	},
}

func init() {
	blueprintFoundationConfigCmd.Flags().BoolVarP(&configFull, "full", "f", false, "Display the Full/Compiled config")
	blueprintFoundationConfigCmd.Flags().StringVarP(&configSetName, "set-name", "n", "", "Set Config Value Key Name, if set-value is not specified, the value will be set to ''")
	blueprintFoundationConfigCmd.Flags().StringVarP(&configSetValue, "set-value", "v", "", "Set Config Value, ignored if set-name is not specified")
	blueprintFoundationConfigCmd.Flags().StringVarP(&configDeleteName, "delete", "d", "", "Delete Config Value Key Name")

	blueprintStructureConfigCmd.Flags().BoolVarP(&configFull, "full", "f", false, "Display the Full/Compiled config")
	blueprintStructureConfigCmd.Flags().StringVarP(&configSetName, "set-name", "n", "", "Set Config Value Key Name, if set-value is not specified, the value will be set to ''")
	blueprintStructureConfigCmd.Flags().StringVarP(&configSetValue, "set-value", "v", "", "Set Config Value, ignored if set-name is not specified")
	blueprintStructureConfigCmd.Flags().StringVarP(&configDeleteName, "delete", "d", "", "Delete Config Value Key Name")

	blueprintFoundationCreateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Name of New Foundation Blueprint")
	blueprintFoundationCreateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Description of New Foundation Blueprint")

	blueprintFoundationUpdateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Update the Description of Foundation Blueprint with value")

	blueprintStructureCreateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Name of New Foundation Blueprint")
	blueprintStructureCreateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Description of New Structure Blueprint")

	blueprintStructureUpdateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Update the Description of Structure Blueprint with value")

	rootCmd.AddCommand(blueprintCmd)
	blueprintCmd.AddCommand(blueprintFoundationCmd)
	blueprintFoundationCmd.AddCommand(blueprintFoundationListCmd)
	blueprintFoundationCmd.AddCommand(blueprintFoundationGetCmd)
	blueprintFoundationCmd.AddCommand(blueprintFoundationCreateCmd)
	blueprintFoundationCmd.AddCommand(blueprintFoundationUpdateCmd)
	blueprintFoundationCmd.AddCommand(blueprintFoundationDeleteCmd)
	blueprintFoundationCmd.AddCommand(blueprintFoundationConfigCmd)

	blueprintCmd.AddCommand(blueprintStructureCmd)
	blueprintStructureCmd.AddCommand(blueprintStructureListCmd)
	blueprintStructureCmd.AddCommand(blueprintStructureGetCmd)
	blueprintStructureCmd.AddCommand(blueprintStructureCreateCmd)
	blueprintStructureCmd.AddCommand(blueprintStructureUpdateCmd)
	blueprintStructureCmd.AddCommand(blueprintStructureDeleteCmd)
	blueprintStructureCmd.AddCommand(blueprintStructureConfigCmd)
}
