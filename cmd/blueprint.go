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

func blueprintArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a Blueprint Id/Name argument")
	}
	return nil
}

func scriptArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a Script Id Name argument")
	}
	return nil
}

func pxeArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a PXE Name(Id) Name argument")
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

		ctx := cmd.Context()

		o, err := contractorClient.BlueprintFoundationBluePrintGet(ctx, blueprintID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:                  {{.GetURI | extractID}}
Name:                {{.Name}}
Description:         {{.Description}}
Parents:             {{.ParentList}}
Config Values:       {{.ConfigValues}}
Foundation Types:    {{.FoundationTypeList}}
Validation Template: {{.ValidationTemplate}}
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
		ctx := cmd.Context()

		rl := []cinp.Object{}
		vchan, err := contractorClient.BlueprintFoundationBluePrintList(ctx, "", map[string]interface{}{})
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

var blueprintFoundationCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Foundation Blueprint",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		o := contractorClient.BlueprintFoundationBluePrintNew()
		o.Name = &detailName
		o.Description = &detailDescription
		o.FoundationTypeList = &[]string{detailAddType}

		err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:                  {{.GetURI | extractID}}
Name:                {{.Name}}
Description:         {{.Description}}
Parents:             {{.ParentList}}
Config Values:       {{.ConfigValues}}
Foundation Types:    {{.FoundationTypeList}}
Validation Template: {{.ValidationTemplate}}
Physical Interfaces: {{.PhysicalInterfaceNames}}
Script Map:          {{.ScriptMap}}
Created:             {{.Created}}
Updated:             {{.Updated}}
`)
		return nil
	},
}

var blueprintFoundationUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Foundation Blueprint",
	Args:  blueprintArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]

		ctx := cmd.Context()

		o := contractorClient.BlueprintFoundationBluePrintNewWithID(blueprintID)

		if detailDescription != "" {
			o.Description = &detailDescription
		}

		if detailAddParent != "" {
			p, err := contractorClient.BlueprintFoundationBluePrintGet(ctx, detailAddParent)
			if err != nil {
				return err
			}
			if o.ParentList == nil {
				o.ParentList = &[]string{}
			}
			*o.ParentList = append(*o.ParentList, p.GetURI())
		}

		if detailDeleteParent != "" {
			p, err := contractorClient.BlueprintFoundationBluePrintGet(ctx, detailDeleteParent)
			if err != nil {
				return err
			}
			id := p.GetURI()
			for i := 0; i < len(*o.ParentList); i++ {
				if (*o.ParentList)[i] == id {
					*o.ParentList = append((*o.ParentList)[:i], (*o.ParentList)[i+1:]...)
					break
				}
			}
		}

		if detailAddType != "" {
			foundationTypeList := append(*o.FoundationTypeList, detailAddType)
			o.FoundationTypeList = &foundationTypeList
		}

		if detailDeleteType != "" {
			for i := 0; i < len(*o.FoundationTypeList); i++ {
				if (*o.FoundationTypeList)[i] == detailDeleteType {
					*o.FoundationTypeList = append((*o.FoundationTypeList)[:i], (*o.FoundationTypeList)[i+1:]...)
					break
				}
			}
		}

		if detailAddIfaceName != "" {
			physicalInterfaceNames := append(*o.PhysicalInterfaceNames, detailAddIfaceName)
			o.PhysicalInterfaceNames = &physicalInterfaceNames
		}

		if detailDeleteIfaceName != "" {
			for i := 0; i < len(*o.PhysicalInterfaceNames); i++ {
				if (*o.PhysicalInterfaceNames)[i] == detailDeleteIfaceName {
					*o.PhysicalInterfaceNames = append((*o.PhysicalInterfaceNames)[:i], (*o.PhysicalInterfaceNames)[i+1:]...)
					break
				}
			}
		}

		// will deal with these later
		// Template map[string]interface{} `json:"template"`

		err := o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:                  {{.GetURI | extractID}}
Name:                {{.Name}}
Description:         {{.Description}}
Parents:             {{.ParentList}}
Config Values:       {{.ConfigValues}}
Foundation Types:    {{.FoundationTypeList}}
Validation Template: {{.ValidationTemplate}}
Physical Interfaces: {{.PhysicalInterfaceNames}}
Script Map:          {{.ScriptMap}}
Created:             {{.Created}}
Updated:             {{.Updated}}
`)

		return nil
	},
}

var blueprintFoundationDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Foundation Blueprint",
	Args:  blueprintArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.BlueprintFoundationBluePrintGet(ctx, blueprintID)
		if err != nil {
			return err
		}
		if err := o.Delete(ctx); err != nil {
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

		ctx := cmd.Context()

		if configSetName != "" {
			r, err := contractorClient.BlueprintFoundationBluePrintGet(ctx, blueprintID)
			if err != nil {
				return err
			}

			o := contractorClient.BlueprintFoundationBluePrintNewWithID(blueprintID)
			o.ConfigValues = r.ConfigValues
			(*o.ConfigValues)[configSetName] = configSetValue

			err = o.Update(ctx)
			if err != nil {
				return err
			}
			outputKV(*o.ConfigValues)

		} else if configDeleteName != "" {
			r, err := contractorClient.BlueprintFoundationBluePrintGet(ctx, blueprintID)
			if err != nil {
				return err
			}

			o := contractorClient.BlueprintFoundationBluePrintNewWithID(blueprintID)
			o.ConfigValues = r.ConfigValues
			delete(*o.ConfigValues, configDeleteName)

			err = o.Update(ctx)
			if err != nil {
				return err
			}
			outputKV(*o.ConfigValues)

		} else if configFull {
			o := contractorClient.BlueprintBluePrintNewWithID(blueprintID)
			r, err := o.CallGetConfig(ctx)
			if err != nil {
				return err
			}
			outputKV(r)

		} else {
			o, err := contractorClient.BlueprintFoundationBluePrintGet(ctx, blueprintID)
			if err != nil {
				return err
			}
			outputKV(*o.ConfigValues)
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
			return errors.New("requires a Blueprint Id/Name, Script Id/Nme, and a Link name (ie: create/destroy/etc) Argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]
		scriptID := args[1]
		scriptName := args[2]

		ctx := cmd.Context()

		o, err := contractorClient.BlueprintFoundationBluePrintGet(ctx, blueprintID)
		if err != nil {
			return err
		}

		_, ok := (*o.ScriptMap)[scriptName]
		if ok {
			return fmt.Errorf("blueprint allready has a script linked with that name")
		}

		s, err := contractorClient.BlueprintScriptGet(ctx, scriptID)
		if err != nil {
			return err
		}

		link := contractorClient.BlueprintBluePrintScriptNew()
		link.Name = &scriptName
		blueprint := strings.Replace(o.GetURI(), "/api/v1/BluePrint/FoundationBluePrint", "/api/v1/BluePrint/BluePrint", 1)
		link.Blueprint = &blueprint
		uri := s.GetURI()
		link.Script = &uri
		err = link.Create(ctx)
		if err != nil {
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
			return errors.New("requires a Blueprint Id/Name and a Link name (ie: create/destroy/etc) Argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]
		scriptName := args[1]

		ctx := cmd.Context()

		o, err := contractorClient.BlueprintFoundationBluePrintGet(ctx, blueprintID)
		if err != nil {
			return err
		}

		_, ok := (*o.ScriptMap)[scriptName]
		if !ok {
			return fmt.Errorf("no Script link to Blueprint with that name")
		}

		vchan, err := contractorClient.BlueprintBluePrintScriptList(ctx, "blueprint", map[string]interface{}{"blueprint": strings.Replace(o.GetURI(), "/api/v1/BluePrint/FoundationBluePrint", "/api/v1/BluePrint/BluePrint", 1)})
		if err != nil {
			return err
		}
		for link := range vchan {
			if *link.Name == scriptName {
				if err := link.Delete(ctx); err != nil {
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

		ctx := cmd.Context()

		o, err := contractorClient.BlueprintStructureBluePrintGet(ctx, blueprintID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:                    {{.GetURI | extractID}}
Name:                  {{.Name}}
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
		rl := []cinp.Object{}

		ctx := cmd.Context()

		vchan, err := contractorClient.BlueprintStructureBluePrintList(ctx, "", map[string]interface{}{})
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

var blueprintStructureCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Structure Blueprint",
	RunE: func(cmd *cobra.Command, args []string) error {

		o := contractorClient.BlueprintStructureBluePrintNew()
		o.Name = &detailName
		o.Description = &detailDescription

		ctx := cmd.Context()

		err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:                    {{.GetURI | extractID}}
Name:                  {{.Name}}
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

var blueprintStructureUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Structure Blueprint",
	Args:  blueprintArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]

		ctx := cmd.Context()

		o := contractorClient.BlueprintStructureBluePrintNewWithID(blueprintID)

		if detailDescription != "" {
			o.Description = &detailDescription
		}

		if detailAddParent != "" {
			p, err := contractorClient.BlueprintStructureBluePrintGet(ctx, detailAddParent)
			if err != nil {
				return err
			}
			if o.ParentList == nil {
				o.ParentList = &[]string{}
			}
			*o.ParentList = append(*o.ParentList, p.GetURI())
		}

		if detailDeleteParent != "" {
			p, err := contractorClient.BlueprintStructureBluePrintGet(ctx, detailDeleteParent)
			if err != nil {
				return err
			}
			id := p.GetURI()
			for i := 0; i < len(*o.ParentList); i++ {
				if (*o.ParentList)[i] == id {
					*o.ParentList = append((*o.ParentList)[:i], (*o.ParentList)[i+1:]...)
					break
				}
			}
		}

		if detailAddFoundationBluePrint != "" {
			p, err := contractorClient.BlueprintFoundationBluePrintGet(ctx, detailAddFoundationBluePrint)
			if err != nil {
				return err
			}
			if o.FoundationBlueprintList == nil {
				o.FoundationBlueprintList = &[]string{}
			}
			*o.FoundationBlueprintList = append(*o.FoundationBlueprintList, p.GetURI())
		}

		if detailDeleteFoundationBluePrint != "" {
			p, err := contractorClient.BlueprintFoundationBluePrintGet(ctx, detailDeleteFoundationBluePrint)
			if err != nil {
				return err
			}
			id := p.GetURI()
			for i := 0; i < len(*o.FoundationBlueprintList); i++ {
				if (*o.FoundationBlueprintList)[i] == id {
					*o.FoundationBlueprintList = append((*o.FoundationBlueprintList)[:i], (*o.FoundationBlueprintList)[i+1:]...)
					break
				}
			}
		}

		err := o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:                    {{.GetURI | extractID}}
Name:                  {{.Name}}
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

var blueprintStructureDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Structure Blueprint",
	Args:  blueprintArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.BlueprintStructureBluePrintGet(ctx, blueprintID)
		if err != nil {
			return err
		}
		if err := o.Delete(ctx); err != nil {
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

		ctx := cmd.Context()

		if configSetName != "" {
			r, err := contractorClient.BlueprintStructureBluePrintGet(ctx, blueprintID)
			if err != nil {
				return err
			}

			o := contractorClient.BlueprintStructureBluePrintNewWithID(blueprintID)
			o.ConfigValues = r.ConfigValues
			(*o.ConfigValues)[configSetName] = configSetValue

			err = o.Update(ctx)
			if err != nil {
				return err
			}
			outputKV(*o.ConfigValues)

		} else if configDeleteName != "" {
			r, err := contractorClient.BlueprintStructureBluePrintGet(ctx, blueprintID)
			if err != nil {
				return err
			}

			o := contractorClient.BlueprintStructureBluePrintNewWithID(blueprintID)
			o.ConfigValues = r.ConfigValues
			delete(*o.ConfigValues, configDeleteName)

			err = o.Update(ctx)
			if err != nil {
				return err
			}
			outputKV(*o.ConfigValues)

		} else if configFull {
			o := contractorClient.BlueprintBluePrintNewWithID(blueprintID)
			r, err := o.CallGetConfig(ctx)
			if err != nil {
				return err
			}
			outputKV(r)

		} else {
			o, err := contractorClient.BlueprintStructureBluePrintGet(ctx, blueprintID)
			if err != nil {
				return err
			}
			outputKV(*o.ConfigValues)
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
			return errors.New("requires a Blueprint Id/Name, Script Id/Nme, and a Link name (ie: create/destroy/etc) Argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]
		scriptID := args[1]
		scriptName := args[2]

		ctx := cmd.Context()

		o, err := contractorClient.BlueprintStructureBluePrintGet(ctx, blueprintID)
		if err != nil {
			return err
		}

		_, ok := (*o.ScriptMap)[scriptName]
		if ok {
			return fmt.Errorf("blueprint allready has a script linked with that name")
		}

		s, err := contractorClient.BlueprintScriptGet(ctx, scriptID)
		if err != nil {
			return err
		}

		link := contractorClient.BlueprintBluePrintScriptNew()
		link.Name = &scriptName
		blueprint := strings.Replace(o.GetURI(), "/api/v1/BluePrint/StructureBluePrint", "/api/v1/BluePrint/BluePrint", 1)
		link.Blueprint = &blueprint
		uri := s.GetURI()
		link.Script = &uri

		err = link.Create(ctx)
		if err != nil {
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
			return errors.New("requires a Blueprint Id/Name and a Link name (ie: create/destroy/etc) Argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintID := args[0]
		scriptName := args[1]

		ctx := cmd.Context()

		o, err := contractorClient.BlueprintStructureBluePrintGet(ctx, blueprintID)
		if err != nil {
			return err
		}

		_, ok := (*o.ScriptMap)[scriptName]
		if !ok {
			return fmt.Errorf("no Script link to Blueprint with that name")
		}

		vchan, err := contractorClient.BlueprintBluePrintScriptList(ctx, "blueprint", map[string]interface{}{"blueprint": strings.Replace(o.GetURI(), "/api/v1/BluePrint/StructureBluePrint", "/api/v1/BluePrint/BluePrint", 1)})
		if err != nil {
			return err
		}
		for link := range vchan {
			if *link.Name == scriptName {
				if err := link.Delete(ctx); err != nil {
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

		ctx := cmd.Context()

		o, err := contractorClient.BlueprintScriptGet(ctx, scriptID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:                    {{.GetURI | extractID}}
Name:                  {{.Name}}
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

		ctx := cmd.Context()

		rl := []cinp.Object{}
		vchan, err := contractorClient.BlueprintScriptList(ctx, "", map[string]interface{}{})
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

var scriptCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Script",
	RunE: func(cmd *cobra.Command, args []string) error {

		o := contractorClient.BlueprintScriptNew()
		o.Name = &detailName
		o.Description = &detailDescription
		script := fmt.Sprintf("# %s", detailDescription)
		o.Script = &script

		ctx := cmd.Context()

		err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:                    {{.GetURI | extractID}}
Name:                  {{.Name}}
Description:           {{.Description}}
Created:               {{.Created}}
Updated:               {{.Updated}}
----  Script  ----
{{.Script}}
`)
		return nil
	},
}

var scriptUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Script",
	Args:  scriptArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptID := args[0]

		ctx := cmd.Context()

		o := contractorClient.BlueprintScriptNewWithID(scriptID)

		if detailDescription != "" {
			o.Description = &detailDescription
		}

		err := o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:                    {{.GetURI | extractID}}
Name:                  {{.Name}}
Description:           {{.Description}}
Created:               {{.Created}}
Updated:               {{.Updated}}
----  Script  ----
{{.Script}}
`)

		return nil
	},
}

var scriptDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Script",
	Args:  scriptArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.BlueprintScriptGet(ctx, scriptID)
		if err != nil {
			return err
		}
		if err := o.Delete(ctx); err != nil {
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

		ctx := cmd.Context()

		o := contractorClient.BlueprintScriptNewWithID(scriptID)

		var newScript string
		var err error

		for {
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
				newScript, err = editBuffer(*o.Script)
				if err != nil {
					return err
				}
			}

			o.Script = &newScript
			err := o.Update(ctx)
			if err != nil {
				if scriptFile == "" && strings.HasPrefix(err.Error(), "Invalid Request: 'map[script:[Script is invalid") {
					fmt.Printf("Error parsing the script:%s\n", err.Error())
					fmt.Println("Return to Editor?(Y/N)")
					var b []byte = make([]byte, 1)
					os.Stdin.Read(b)
					if b[0] == 'Y' || b[0] == 'y' {
						continue
					}
					fmt.Println("Aborting")
					break
				}
				return err
			}
			fmt.Println("Changes Saved")

			break
		}

		return nil
	},
}

var pxeCmd = &cobra.Command{
	Use:   "pxe",
	Short: "Work with PXEs",
}

var pxeGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get PXE",
	Args:  pxeArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		pxeID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.BlueprintPXEGet(ctx, pxeID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:                    {{.GetURI | extractID}}
Name:                  {{.Name}}
Created:               {{.Created}}
Updated:               {{.Updated}}
----  Script  ----
{{.BootScript}}

----  Template  ----
{{.Template}}
`)
		return nil
	},
}

var pxeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List PXEs",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := cmd.Context()

		rl := []cinp.Object{}
		vchan, err := contractorClient.BlueprintPXEList(ctx, "", map[string]interface{}{})
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

var pxeEditScriptCmd = &cobra.Command{
	Use:   "editscript",
	Short: "Edit Script of PXE",
	Args:  scriptArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		pxeID := args[0]

		ctx := cmd.Context()

		o := contractorClient.BlueprintPXENewWithID(pxeID)

		var newScript string
		var err error

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
			newScript, err = editBuffer(*o.BootScript)
			if err != nil {
				return err
			}
		}

		o.BootScript = &newScript
		err = o.Update(ctx)
		if err != nil {
			return err
		}
		fmt.Println("Changes Saved")

		return nil
	},
}

var pxeEditTemplateCmd = &cobra.Command{
	Use:   "edittemplate",
	Short: "Edit Template of PXE",
	Args:  scriptArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		pxeID := args[0]

		ctx := cmd.Context()

		o := contractorClient.BlueprintPXENewWithID(pxeID)

		var newTemplate string
		var err error

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
			newTemplate = strings.TrimSpace(string(buf[:len]))

		} else {
			newTemplate, err = editBuffer(*o.Template)
			if err != nil {
				return err
			}
		}

		o.Template = &newTemplate
		err = o.Update(ctx)
		if err != nil {
			return err
		}
		fmt.Println("Changes Saved")

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

	pxeEditScriptCmd.Flags().StringVarP(&scriptFile, "file", "f", "", "File to supply the script, use '-' for stdin or omit for interactive editor")

	pxeEditTemplateCmd.Flags().StringVarP(&scriptFile, "file", "f", "", "File to supply the template, use '-' for stdin or omit for interactive editor")

	rootCmd.AddCommand(blueprintCmd)
	blueprintCmd.AddCommand(blueprintFoundationCmd)
	blueprintFoundationCmd.AddCommand(blueprintFoundationListCmd, blueprintFoundationGetCmd, blueprintFoundationCreateCmd, blueprintFoundationUpdateCmd, blueprintFoundationDeleteCmd, blueprintFoundationConfigCmd)

	blueprintFoundationCmd.AddCommand(blueprintFoundationScriptCmd)
	blueprintFoundationScriptCmd.AddCommand(blueprintFoundationScriptLinkCmd, blueprintFoundationScriptUnlinkCmd)

	blueprintCmd.AddCommand(blueprintStructureCmd)
	blueprintStructureCmd.AddCommand(blueprintStructureListCmd, blueprintStructureGetCmd, blueprintStructureCreateCmd, blueprintStructureUpdateCmd, blueprintStructureDeleteCmd, blueprintStructureConfigCmd)

	blueprintStructureCmd.AddCommand(blueprintStructureScriptCmd)
	blueprintStructureScriptCmd.AddCommand(blueprintStructureScriptLinkCmd, blueprintStructureScriptUnlinkCmd)

	blueprintCmd.AddCommand(scriptCmd)
	scriptCmd.AddCommand(scriptListCmd, scriptGetCmd, scriptCreateCmd, scriptUpdateCmd, scriptDeleteCmd, scriptEditCmd)

	blueprintCmd.AddCommand(pxeCmd)
	pxeCmd.AddCommand(pxeGetCmd, pxeListCmd, pxeEditScriptCmd, pxeEditTemplateCmd)
}
