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
	"strings"

	cinp "github.com/cinp/go"
	"github.com/spf13/cobra"
)

var scriptFile string
var detailAddParent string
var detailDeleteParent string
var detailAddFoundationBluePrint string
var detailDeleteFoundationBluePrint string
var detailAddType string
var detailDeleteType string
var detailAddIfaceName string
var detailDeleteIfaceName string

func blueprintArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Requires a Blueprint Id/Name Argument")
	}
	return nil
}

func scriptArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Requires a Script Id Name Argument")
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
Script Map:          {{.ScriptMap}}
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
		o.FoundationTypeList = []string{detailAddType}

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

		o, err := c.BlueprintFoundationBluePrintGet(blueprintID)
		if err != nil {
			return err
		}
		if detailDescription != "" {
			o.Description = detailDescription
			fieldList = append(fieldList, "description")
		}

		if detailAddParent != "" {
			p, err := c.BlueprintStructureBluePrintGet(detailAddParent)
			if err != nil {
				return err
			}
			o.ParentList = append(o.ParentList, p.GetID())
			fieldList = append(fieldList, "parent_list")
		}

		if detailDeleteParent != "" {
			p, err := c.BlueprintStructureBluePrintGet(detailDeleteParent)
			if err != nil {
				return err
			}
			id := p.GetID()
			for i := 0; i < len(o.ParentList); i++ {
				if o.ParentList[i] == id {
					o.ParentList = append(o.ParentList[:i], o.ParentList[i+1:]...)
					break
				}
			}
			fieldList = append(fieldList, "parent_list")
		}

		if detailAddType != "" {
			o.FoundationTypeList = append(o.FoundationTypeList, detailAddType)
			fieldList = append(fieldList, "foundation_type_list")
		}

		if detailDeleteType != "" {
			for i := 0; i < len(o.FoundationTypeList); i++ {
				if o.FoundationTypeList[i] == detailDeleteType {
					o.FoundationTypeList = append(o.FoundationTypeList[:i], o.FoundationTypeList[i+1:]...)
					break
				}
			}
			fieldList = append(fieldList, "foundation_type_list")
		}

		if detailAddIfaceName != "" {
			o.PhysicalInterfaceNames = append(o.PhysicalInterfaceNames, detailAddIfaceName)
			fieldList = append(fieldList, "physical_interface_names")
		}

		if detailDeleteIfaceName != "" {
			for i := 0; i < len(o.PhysicalInterfaceNames); i++ {
				if o.PhysicalInterfaceNames[i] == detailDeleteIfaceName {
					o.PhysicalInterfaceNames = append(o.PhysicalInterfaceNames[:i], o.PhysicalInterfaceNames[i+1:]...)
					break
				}
			}
			fieldList = append(fieldList, "physical_interface_names")
		}

		// will deal with these later
		// Template map[string]interface{} `json:"template"`

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

var blueprintFoundationScriptCmd = &cobra.Command{
	Use:   "script",
	Short: "Work with Foundation Blueprint Scripts",
}

var blueprintFoundationScriptLinkCmd = &cobra.Command{
	Use:   "link",
	Short: "Link Foundation Blueprint to Script",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 3 {
			return errors.New("Requires a Blueprint Id/Name, Script Id/Nme, and a Link name (ie: create/destroy/etc) Argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]
		scriptID := args[1]
		scriptName := args[2]
		c := getContractor()
		defer c.Logout()

		r, err := c.BlueprintFoundationBluePrintGet(blueprintID)
		if err != nil {
			return err
		}

		_, ok := r.ScriptMap[scriptName]
		if ok {
			return fmt.Errorf("Blueprint allready has a script linked with that name")
		}

		s, err := c.BlueprintScriptGet(scriptID)
		if err != nil {
			return err
		}

		link := c.BlueprintBluePrintScriptNew()
		link.Name = scriptName
		link.Blueprint = strings.Replace(r.GetID(), "/api/v1/BluePrint/FoundationBluePrint", "/api/v1/BluePrint/BluePrint", 1)
		link.Script = s.GetID()
		if err := link.Create(); err != nil {
			return err
		}

		return nil
	},
}

var blueprintFoundationScriptUnlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "UnLink Script from Foundation Blueprint",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Requires a Blueprint Id/Name and a Link name (ie: create/destroy/etc) Argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]
		scriptName := args[1]
		c := getContractor()
		defer c.Logout()

		r, err := c.BlueprintFoundationBluePrintGet(blueprintID)
		if err != nil {
			return err
		}

		_, ok := r.ScriptMap[scriptName]
		if !ok {
			return fmt.Errorf("No Script link to Blueprint with that name")
		}

		for link := range c.BlueprintBluePrintScriptList("blueprint", map[string]interface{}{"blueprint": strings.Replace(r.GetID(), "/api/v1/BluePrint/FoundationBluePrint", "/api/v1/BluePrint/BluePrint", 1)}) {
			if link.Name == scriptName {
				if err := link.Delete(); err != nil {
					return err
				}
			}
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

		r, err := c.BlueprintStructureBluePrintGet(blueprintID)
		if err != nil {
			return err
		}
		outputDetail(r, `Name:                  {{.Name}}
Description:           {{.Description}}
Parents:               {{.ParentList}}
Config Values:         {{.ConfigValues}}
Foundation BluePrints: {{.FoundationBlueprintList}}
Script Map:            {{.ScriptMap}}
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

		o, err := c.BlueprintStructureBluePrintGet(blueprintID)
		if err != nil {
			return err
		}
		if detailDescription != "" {
			o.Description = detailDescription
			fieldList = append(fieldList, "description")
		}

		if detailAddParent != "" {
			p, err := c.BlueprintStructureBluePrintGet(detailAddParent)
			if err != nil {
				return err
			}
			o.ParentList = append(o.ParentList, p.GetID())
			fieldList = append(fieldList, "parent_list")
		}

		if detailDeleteParent != "" {
			p, err := c.BlueprintStructureBluePrintGet(detailDeleteParent)
			if err != nil {
				return err
			}
			id := p.GetID()
			for i := 0; i < len(o.ParentList); i++ {
				if o.ParentList[i] == id {
					o.ParentList = append(o.ParentList[:i], o.ParentList[i+1:]...)
					break
				}
			}
			fieldList = append(fieldList, "parent_list")
		}

		if detailAddFoundationBluePrint != "" {
			p, err := c.BlueprintFoundationBluePrintGet(detailAddFoundationBluePrint)
			if err != nil {
				return err
			}
			o.FoundationBlueprintList = append(o.FoundationBlueprintList, p.GetID())
			fieldList = append(fieldList, "foundation_blueprint_list")
		}

		if detailDeleteFoundationBluePrint != "" {
			p, err := c.BlueprintFoundationBluePrintGet(detailDeleteFoundationBluePrint)
			if err != nil {
				return err
			}
			id := p.GetID()
			for i := 0; i < len(o.FoundationBlueprintList); i++ {
				if o.FoundationBlueprintList[i] == id {
					o.FoundationBlueprintList = append(o.FoundationBlueprintList[:i], o.FoundationBlueprintList[i+1:]...)
					break
				}
			}
			fieldList = append(fieldList, "foundation_blueprint_list")
		}

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

var blueprintStructureScriptCmd = &cobra.Command{
	Use:   "script",
	Short: "Work with Structure Blueprint Scripts",
}

var blueprintStructureScriptLinkCmd = &cobra.Command{
	Use:   "link",
	Short: "Link Structure Blueprint to Script",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 3 {
			return errors.New("Requires a Blueprint Id/Name, Script Id/Nme, and a Link name (ie: create/destroy/etc) Argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]
		scriptID := args[1]
		scriptName := args[2]
		c := getContractor()
		defer c.Logout()

		r, err := c.BlueprintStructureBluePrintGet(blueprintID)
		if err != nil {
			return err
		}

		_, ok := r.ScriptMap[scriptName]
		if ok {
			return fmt.Errorf("Blueprint allready has a script linked with that name")
		}

		s, err := c.BlueprintScriptGet(scriptID)
		if err != nil {
			return err
		}

		link := c.BlueprintBluePrintScriptNew()
		link.Name = scriptName
		link.Blueprint = strings.Replace(r.GetID(), "/api/v1/BluePrint/StructureBluePrint", "/api/v1/BluePrint/BluePrint", 1)
		link.Script = s.GetID()
		if err := link.Create(); err != nil {
			return err
		}

		return nil
	},
}

var blueprintStructureScriptUnlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "UnLink Script from Structure Blueprint",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Requires a Blueprint Id/Name and a Link name (ie: create/destroy/etc) Argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]
		scriptName := args[1]
		c := getContractor()
		defer c.Logout()

		r, err := c.BlueprintStructureBluePrintGet(blueprintID)
		if err != nil {
			return err
		}

		_, ok := r.ScriptMap[scriptName]
		if !ok {
			return fmt.Errorf("No Script link to Blueprint with that name")
		}

		for link := range c.BlueprintBluePrintScriptList("blueprint", map[string]interface{}{"blueprint": strings.Replace(r.GetID(), "/api/v1/BluePrint/StructureBluePrint", "/api/v1/BluePrint/BluePrint", 1)}) {
			if link.Name == scriptName {
				if err := link.Delete(); err != nil {
					return err
				}
			}
		}

		return nil
	},
}

var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "Work with Blueprint Scripts",
}

var scriptGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Script",
	Args:  scriptArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.BlueprintScriptGet(scriptID)
		if err != nil {
			return err
		}
		outputDetail(r, `Name:                  {{.Name}}
Description:           {{.Description}}
Created:               {{.Created}}
Updated:               {{.Updated}}
----  Script  ----
{{.Script}}
`)
		return nil
	},
}

var scriptListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Scriptss",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		rl := []cinp.Object{}
		for v := range c.BlueprintScriptList("", map[string]interface{}{}) {
			rl = append(rl, v)
		}
		outputList(rl, "Id	Name	Description	Created	Updated\n", "{{.GetID | extractID}}	{{.Name}}	{{.Description}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var scriptCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Script",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		o := c.BlueprintScriptNew()
		o.Name = detailName
		o.Description = detailDescription
		o.Script = fmt.Sprintf("# %s", detailDescription)

		if err := o.Create(); err != nil {
			return err
		}

		outputKV(map[string]interface{}{"id": o.GetID()})

		return nil
	},
}

var scriptUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Script",
	Args:  scriptArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		scriptID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.BlueprintScriptGet(scriptID)
		if err != nil {
			return err
		}
		if detailDescription != "" {
			o.Description = detailDescription
			fieldList = append(fieldList, "description")
		}

		if err := o.Update(fieldList); err != nil {
			return err
		}

		return nil
	},
}

var scriptDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Script",
	Args:  scriptArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.BlueprintScriptGet(scriptID)
		if err != nil {
			return err
		}
		if err := r.Delete(); err != nil {
			return err
		}

		return nil
	},
}

var scriptEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit Script of Blueprint Script",
	Args:  scriptArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.BlueprintScriptGet(scriptID)
		if err != nil {
			return err
		}

		var newScript string

		if scriptFile != "" {
			var source *os.File
			if scriptFile == "-" {
				source = os.Stdin
			} else {
				source, err = os.Open(scriptFile)
				if err != nil {
					return err
				}
			}
			buf := make([]byte, 4096*1024)
			len, err := source.Read(buf)
			if err != nil {
				return err
			}
			newScript = strings.TrimSpace(string(buf[:len]))

		} else {
			newScript, err = editBuffer(r.Script)
			if err != nil {
				return err
			}
		}

		if newScript != r.Script {
			r.Script = newScript
			if err := r.Update([]string{"script"}); err != nil {
				return err
			}
			fmt.Println("Changes Saved")
		} else {
			fmt.Println("No Change Detected")
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
	blueprintFoundationCreateCmd.Flags().StringVarP(&detailAddType, "type", "t", "", "Type of New Foundation Blueprint")

	blueprintFoundationUpdateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Update the Description of Foundation Blueprint with value")
	blueprintFoundationUpdateCmd.Flags().StringVarP(&detailAddParent, "add-parent", "p", "", "Add Parent to Foundation Blueprint")
	blueprintFoundationUpdateCmd.Flags().StringVarP(&detailDeleteParent, "delete-parent", "q", "", "Remove Parent from Foundation Blueprint")
	blueprintFoundationUpdateCmd.Flags().StringVarP(&detailAddType, "add-type", "t", "", "Add Type to Foundation Blueprint")
	blueprintFoundationUpdateCmd.Flags().StringVarP(&detailDeleteType, "delete-type", "u", "", "Remove Type from Foundation Blueprint")
	blueprintFoundationUpdateCmd.Flags().StringVarP(&detailAddIfaceName, "add-iface-name", "i", "", "Add Physical Interface Name to Foundation Blueprint")
	blueprintFoundationUpdateCmd.Flags().StringVarP(&detailDeleteIfaceName, "delete-iface-name", "k", "", "Remove Physical Interface Name from Foundation Blueprint")

	blueprintStructureCreateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Name of New Structure Blueprint")
	blueprintStructureCreateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Description of New Structure Blueprint")

	blueprintStructureUpdateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Update the Description of Structure Blueprint with value")
	blueprintStructureUpdateCmd.Flags().StringVarP(&detailAddParent, "add-parent", "p", "", "Add Parent to Structure Blueprint")
	blueprintStructureUpdateCmd.Flags().StringVarP(&detailDeleteParent, "delete-parent", "q", "", "Remove Parent from Structure Blueprint")
	blueprintStructureUpdateCmd.Flags().StringVarP(&detailAddFoundationBluePrint, "add-foundation-blueprint", "f", "", "Add Foundation Blueprint to Structure Blueprint")
	blueprintStructureUpdateCmd.Flags().StringVarP(&detailDeleteFoundationBluePrint, "delete-foundation-blueprint", "g", "", "Remove Foundation Blueprint from Structure Blueprint")

	scriptCreateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Name of New Scriptt")
	scriptCreateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Description of New Script")

	scriptUpdateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Update the Description of Script with value")

	scriptEditCmd.Flags().StringVarP(&scriptFile, "file", "f", "", "File to supply the script, use '-' for stdin or omit for interactive editor")

	rootCmd.AddCommand(blueprintCmd)
	blueprintCmd.AddCommand(blueprintFoundationCmd)
	blueprintFoundationCmd.AddCommand(blueprintFoundationListCmd)
	blueprintFoundationCmd.AddCommand(blueprintFoundationGetCmd)
	blueprintFoundationCmd.AddCommand(blueprintFoundationCreateCmd)
	blueprintFoundationCmd.AddCommand(blueprintFoundationUpdateCmd)
	blueprintFoundationCmd.AddCommand(blueprintFoundationDeleteCmd)
	blueprintFoundationCmd.AddCommand(blueprintFoundationConfigCmd)

	blueprintFoundationCmd.AddCommand(blueprintFoundationScriptCmd)
	blueprintFoundationScriptCmd.AddCommand(blueprintFoundationScriptLinkCmd)
	blueprintFoundationScriptCmd.AddCommand(blueprintFoundationScriptUnlinkCmd)

	blueprintCmd.AddCommand(blueprintStructureCmd)
	blueprintStructureCmd.AddCommand(blueprintStructureListCmd)
	blueprintStructureCmd.AddCommand(blueprintStructureGetCmd)
	blueprintStructureCmd.AddCommand(blueprintStructureCreateCmd)
	blueprintStructureCmd.AddCommand(blueprintStructureUpdateCmd)
	blueprintStructureCmd.AddCommand(blueprintStructureDeleteCmd)
	blueprintStructureCmd.AddCommand(blueprintStructureConfigCmd)

	blueprintStructureCmd.AddCommand(blueprintStructureScriptCmd)
	blueprintStructureScriptCmd.AddCommand(blueprintStructureScriptLinkCmd)
	blueprintStructureScriptCmd.AddCommand(blueprintStructureScriptUnlinkCmd)

	blueprintCmd.AddCommand(scriptCmd)
	scriptCmd.AddCommand(scriptListCmd)
	scriptCmd.AddCommand(scriptGetCmd)
	scriptCmd.AddCommand(scriptCreateCmd)
	scriptCmd.AddCommand(scriptUpdateCmd)
	scriptCmd.AddCommand(scriptDeleteCmd)
	scriptCmd.AddCommand(scriptEditCmd)

}
