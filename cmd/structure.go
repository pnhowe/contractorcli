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
	"sort"
	"strconv"
	"strings"

	contractor "github.com/t3kton/contractor_client/go"

	cinp "github.com/cinp/go"
	"github.com/spf13/cobra"
)

var configSetName, configSetValue, configDeleteName string
var configFull, detailIsPrimary bool
var detailHostname, detailSite, detailBlueprint, detailFoundation, detailInterfaceName string
var detailPrimary, detailSecondary string

func structureArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Requires a Structure Id Argument")
	}
	return nil
}

func structureAddressArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Requires a Address Id Argument")
	}
	return nil
}

func structureInterfaceArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Requires a Structure Interface Id Argument")
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
		vchan, err := c.BuildingStructureList("", map[string]interface{}{})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Id", "Site", "Hostname", "Foundation", "Blueprint", "Created", "Updated"}, "{{.GetID | extractID}}	{{.Site | extractID}}	{{.Hostname}}	{{.Foundation | extractID}}	{{.Blueprint | extractID}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var structureGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Structure",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		c := getContractor()
		defer c.Logout()

		r, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
		}
		outputDetail(r, `Hostname:      {{.Hostname}}
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

		outputKV(map[string]interface{}{"id": o.GetID()})

		return nil
	},
}

var structureUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Structure",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
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
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
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
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
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

var structureAddressCmd = &cobra.Command{
	Use:   "address",
	Short: "Work with Structure Ip Addresses",
}

var structureAddressListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Ip Addresses attached to a structure",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		c := getContractor()
		defer c.Logout()

		o, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
		}

		rl := []cinp.Object{}
		vchan, err := c.UtilitiesAddressList("structure", map[string]interface{}{"structure": o.GetID()})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Id", "Interface", "Address", "Address Block", "Offset", "Is Primary", "Created", "Updated"}, "{{.GetID | extractID}}	{{.InterfaceName}}	{{.IPAddress}}	{{.AddressBlock | extractID}}	{{.Offset}}	{{.IsPrimary}}	{{.Updated}}	{{.Created}}\n")

		return nil
	},
}

var structureAddressNextCmd = &cobra.Command{
	Use:   "next",
	Short: "Assign Next available IP address in Address Block to structure",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		c := getContractor()
		defer c.Logout()

		o, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
		}

		a, err := c.UtilitiesAddressBlockGet(detailAddressBlock)
		if err != nil {
			return err
		}

		addrID, err := a.CallNextAddress(strings.Replace(o.GetID(), "/api/v1/Building/Structure", "/api/v1/Utilities/Networked", 1), detailInterfaceName, detailIsPrimary)
		if err != nil {
			return err
		}

		outputKV(map[string]interface{}{"id": addrID})

		return nil
	},
}

var structureAddressAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an IP address to structure",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		c := getContractor()
		defer c.Logout()

		s, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
		}

		a, err := c.UtilitiesAddressBlockGet(detailAddressBlock)
		if err != nil {
			return err
		}

		o := c.UtilitiesAddressNew()
		o.AddressBlock = a.GetID()
		o.Offset = detailOffset
		o.Networked = strings.Replace(s.GetID(), "/api/v1/Building/Structure", "/api/v1/Utilities/Networked", 1)
		o.InterfaceName = detailInterfaceName
		o.IsPrimary = detailIsPrimary

		if err := o.Create(); err != nil {
			return err
		}

		outputKV(map[string]interface{}{"id": o.GetID()})

		return nil
	},
}

var structureAddressUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an IP address to structure",
	Args:  structureAddressArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		addressID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		c := getContractor()
		defer c.Logout()

		o, err := c.UtilitiesAddressGet(addressID)
		if err != nil {
			return err
		}

		if detailOffset != 0 {
			o.Offset = detailOffset
			fieldList = append(fieldList, "offset")
		}

		if detailInterfaceName != "" {
			o.InterfaceName = detailInterfaceName
			fieldList = append(fieldList, "interface_name")
		}

		if err := o.Update(fieldList); err != nil {
			return err
		}

		return nil
	},
}

var structureJobCmd = &cobra.Command{
	Use:   "job",
	Short: "Work with Structure Jobs",
}

var structureJobInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show Structure Job Info",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		c := getContractor()
		defer c.Logout()

		o, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
		}

		jID, err := o.CallGetJob()
		if err != nil {
			return err
		}
		jIDi, err := extractIDInt(jID)
		if err != nil {
			return err
		}
		j, err := c.ForemanStructureJobGet(jIDi)
		if err != nil {
			return err
		}

		outputDetail(j, `Job:           {{.GetID | extractID}}
Site:          {{.Site}}
Structure:     {{.Structure | extractID}}
State:         {{.State}}
Status:        {{.Status}}
Progress:      {{.Progress}}
Message:       {{.Message}}
Script Name:   {{.ScriptName}}
Can Start:     {{.CanStart}}
Updated:       {{.Updated}}
Created:       {{.Created}}
`)
		return nil
	},
}

var structureJobStateCmd = &cobra.Command{
	Use:   "state",
	Short: "Show Structure Job State",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		c := getContractor()
		defer c.Logout()

		o, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
		}
		jID, err := o.CallGetJob()
		if err != nil {
			return err
		}
		jIDi, err := extractIDInt(jID)
		if err != nil {
			return err
		}
		j, err := c.ForemanStructureJobGet(jIDi)
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
{{.state.script}}
`)
		return nil
	},
}

var structureJobDoCreateCmd = &cobra.Command{
	Use:   "do-create",
	Short: "Start Create Job for Structure",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		c := getContractor()
		defer c.Logout()

		o, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
		}

		if _, err := o.CallDoCreate(); err != nil {
			return err
		}

		return nil
	},
}

var structureJobDoDestroyCmd = &cobra.Command{
	Use:   "do-destroy",
	Short: "Start Destroy Job for Structure",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		c := getContractor()
		defer c.Logout()

		o, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
		}

		if _, err := o.CallDoDestroy(); err != nil {
			return err
		}
		return nil
	},
}

var structureJobDoUtilityCmd = &cobra.Command{
	Use:   "do-utility",
	Short: "Start Utility Job for Structure",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Requires a Structure Id and Utility Job Name Argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		scriptName := args[1]
		c := getContractor()
		defer c.Logout()

		o, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
		}

		if _, err := o.CallDoJob(scriptName); err != nil {
			return err
		}
		return nil
	},
}

var structureJobLogCmd = &cobra.Command{
	Use:   "joblog",
	Short: "Job Log List for a Structure",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		c := getContractor()
		defer c.Logout()

		o, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
		}

		rl := []cinp.Object{}
		vchan, err := c.ForemanJobLogList("structure", map[string]interface{}{"structure": o.GetID()})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Script Name", "Created By", "Started At", "Finished At", "Canceled By", "Cancled At", "Created", "Updated"}, "{{.ScriptName}}	{{.Creator}}	{{.StartedAt}}	{{.FinishedAt}}	{{.CanceledBy}}	{{.CanceledAt}}	{{.Updated}}	{{.Created}}\n")

		return nil
	},
}

var structureInterfaceCmd = &cobra.Command{
	Use:   "interface",
	Short: "Work with Structure Interfaces",
}

var structureInterfaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Interfaces attached to a structure",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		c := getContractor()
		defer c.Logout()

		r, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
		}

		rl := []cinp.Object{}
		vchan, err := c.UtilitiesAbstractNetworkInterfaceList("structure", map[string]interface{}{"structure": r.GetID()})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}

		sort.Slice(rl, func(i, j int) bool {
			return rl[i].(*contractor.UtilitiesAbstractNetworkInterface).ID < rl[i].(*contractor.UtilitiesAbstractNetworkInterface).ID
		})

		outputList(rl, []string{"Id", "Name", "Type", "Network", "Created", "Update"}, "{{.GetID | extractID}}	{{.Name}}	{{.Type}}	{{.Network | extractID}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var structureInterfaceGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Structure Interface",
	Args:  structureInterfaceArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		interfaceID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		c := getContractor()
		defer c.Logout()

		r, err := c.UtilitiesAbstractNetworkInterfaceGet(interfaceID)
		if err != nil {
			return err
		}
		outputDetail(r, `Name:             {{.Name}}
Type:             {{.Type}}
Network:          {{.Network | extractID}}
Structure:        {{.Structure | extractID}}
Created:          {{.Created}}
Updated:          {{.Updated}}
`)
		return nil
	},
}

var structureInterfaceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Structure Interface",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		c := getContractor()
		defer c.Logout()

		r, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
		}

		o := c.UtilitiesAbstractNetworkInterfaceNew()
		o.Structure = r.GetID()
		o.Name = detailName

		if detailNetwork != 0 {
			r, err := c.UtilitiesNetworkGet(detailNetwork)
			if err != nil {
				return err
			}
			o.Network = r.GetID()
		}

		if err := o.Create(); err != nil {
			return err
		}

		outputKV(map[string]interface{}{"id": o.GetID()})

		return nil
	},
}

var structureInterfaceUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Structure Interface",
	Args:  structureInterfaceArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		interfaceID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		c := getContractor()
		defer c.Logout()

		o, err := c.UtilitiesAbstractNetworkInterfaceGet(interfaceID)
		if err != nil {
			return err
		}

		if detailName != "" {
			o.Name = detailName
			fieldList = append(fieldList, "name")
		}

		if detailNetwork != 0 {
			r, err := c.UtilitiesNetworkGet(detailNetwork)
			if err != nil {
				return err
			}
			o.Network = r.GetID()
			fieldList = append(fieldList, "network")
		}

		if err := o.Update(fieldList); err != nil {
			return err
		}

		return nil
	},
}

var structureInterfaceDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Structure Interface",
	Args:  structureInterfaceArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		interfaceID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		c := getContractor()
		defer c.Logout()

		r, err := c.UtilitiesAbstractNetworkInterfaceGet(interfaceID)
		if err != nil {
			return err
		}
		if err := r.Delete(); err != nil {
			return err
		}

		return nil
	},
}

var structureAggInterfaceCmd = &cobra.Command{
	Use:   "agginterface",
	Short: "Work with Structure Aggregated (bonded) Interfaces",
}

var structureAggInterfaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Aggregated Interfaces attached to a structure",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		c := getContractor()
		defer c.Logout()

		r, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
		}

		rl := []cinp.Object{}
		vchan, err := c.UtilitiesAggregatedNetworkInterfaceList("structure", map[string]interface{}{"structure": r.GetID()})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}

		sort.Slice(rl, func(i, j int) bool {
			return rl[i].(*contractor.UtilitiesAbstractNetworkInterface).ID < rl[i].(*contractor.UtilitiesAbstractNetworkInterface).ID
		})

		outputList(rl, []string{"Id", "Name", "Network", "Primary Interface", "Created", "Update"}, "{{.GetID | extractID}}	{{.Name}}	{{.PrimaryInterface}}	{{.Network | extractID}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var structureAggInterfaceGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Structure Aggregated Interface",
	Args:  structureInterfaceArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		interfaceID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.UtilitiesAggregatedNetworkInterfaceGet(interfaceID)
		if err != nil {
			return err
		}
		outputDetail(r, `Name:             {{.Name}}
Network:          {{.Network | extractID}}
Structure:        {{.Structure | extractID}}
Primary:          {{.PrimaryInterface}}
Secondaries:      {{.SecondaryInterfaces}}
Created:          {{.Created}}
Updated:          {{.Updated}}
`)
		return nil
	},
}

var structureAggInterfaceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Structure Aggregated Interface",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		c := getContractor()
		defer c.Logout()

		r, err := c.BuildingStructureGet(structureID)
		if err != nil {
			return err
		}

		o := c.UtilitiesAggregatedNetworkInterfaceNew()
		o.Structure = r.GetID()
		o.Name = detailName

		if detailNetwork != 0 {
			r, err := c.UtilitiesNetworkGet(detailNetwork)
			if err != nil {
				return err
			}
			o.Network = r.GetID()
		}

		o.PrimaryInterface = detailPrimary
		o.SecondaryInterfaces = strings.Split(detailSecondary, ",")

		if err := o.Create(); err != nil {
			return err
		}

		outputKV(map[string]interface{}{"id": o.GetID()})

		return nil
	},
}

var structureAggInterfaceUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Structure Aggregated Interface",
	Args:  structureInterfaceArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		interfaceID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.UtilitiesAggregatedNetworkInterfaceGet(interfaceID)
		if err != nil {
			return err
		}

		if detailName != "" {
			o.Name = detailName
			fieldList = append(fieldList, "name")
		}

		if detailNetwork != 0 {
			r, err := c.UtilitiesNetworkGet(detailNetwork)
			if err != nil {
				return err
			}
			o.Network = r.GetID()
			fieldList = append(fieldList, "network")
		}

		if detailPrimary != "" {
			o.PrimaryInterface = detailPrimary
			fieldList = append(fieldList, "primary_interface")
		}

		if detailSecondary != "" {
			o.SecondaryInterfaces = strings.Split(detailSecondary, ",")
			fieldList = append(fieldList, "secondary_interfaces")
		}

		if err := o.Update(fieldList); err != nil {
			return err
		}

		return nil
	},
}

var structureAggInterfaceDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Structure Aggregated Interface",
	Args:  structureInterfaceArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		interfaceID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.UtilitiesAggregatedNetworkInterfaceGet(interfaceID)
		if err != nil {
			return err
		}
		if err := r.Delete(); err != nil {
			return err
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

	structureAddressNextCmd.Flags().IntVarP(&detailAddressBlock, "addressblock", "a", 0, "Address Block to get an IP From")
	structureAddressNextCmd.Flags().StringVarP(&detailInterfaceName, "interfacename", "n", "", "Name of the Interface to assigne the IP To")
	structureAddressNextCmd.Flags().BoolVarP(&detailIsPrimary, "primary", "p", false, "If this is the primary IP Address")

	structureAddressAddCmd.Flags().IntVarP(&detailAddressBlock, "addressblock", "a", 0, "Address Block to get an IP From")
	structureAddressAddCmd.Flags().StringVarP(&detailInterfaceName, "interfacename", "n", "", "Name of the Interface to assigne the IP To")
	structureAddressAddCmd.Flags().IntVarP(&detailOffset, "offset", "o", 0, "Offset inside the Address Block to use")
	structureAddressAddCmd.Flags().BoolVarP(&detailIsPrimary, "primary", "p", false, "If this is the primary IP Address")

	structureAddressUpdateCmd.Flags().StringVarP(&detailInterfaceName, "interfacename", "n", "", "Name of the Interface to assigne the IP To")
	structureAddressUpdateCmd.Flags().IntVarP(&detailOffset, "offset", "o", 0, "Offset inside the Address Block to use")

	structureInterfaceCreateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Name of the new Interface")
	structureInterfaceCreateCmd.Flags().IntVarP(&detailNetwork, "network", "t", 0, "Network id to attach the new Interface to")

	structureInterfaceUpdateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Update the Name of the Interface")
	structureInterfaceUpdateCmd.Flags().IntVarP(&detailNetwork, "network", "t", 0, "Update Network id the Interface is attached to")

	structureAggInterfaceCreateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Name of the new Interface")
	structureAggInterfaceCreateCmd.Flags().IntVarP(&detailNetwork, "network", "t", 0, "Network id to attach the new Interface to")
	structureAggInterfaceCreateCmd.Flags().StringVarP(&detailPrimary, "primary", "p", "", "Interface name to use as the primary interface")
	structureAggInterfaceCreateCmd.Flags().StringVarP(&detailSecondary, "secondary", "s", "", "Interface names to use as the secondaries, delimited by ','")

	structureAggInterfaceUpdateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Update the Name of the Interface")
	structureAggInterfaceUpdateCmd.Flags().IntVarP(&detailNetwork, "network", "t", 0, "Update Network id the Interface is attached to")
	structureAggInterfaceUpdateCmd.Flags().StringVarP(&detailPrimary, "primary", "p", "", "Interface name to use as the primary interface")
	structureAggInterfaceUpdateCmd.Flags().StringVarP(&detailSecondary, "secondary", "s", "", "Interface names to use as the secondaries, delimited by ','")

	rootCmd.AddCommand(structureCmd)
	structureCmd.AddCommand(structureListCmd)
	structureCmd.AddCommand(structureGetCmd)
	structureCmd.AddCommand(structureCreateCmd)
	structureCmd.AddCommand(structureUpdateCmd)
	structureCmd.AddCommand(structureDeleteCmd)
	structureCmd.AddCommand(structureConfigCmd)

	structureCmd.AddCommand(structureAddressCmd)
	structureAddressCmd.AddCommand(structureAddressListCmd)
	structureAddressCmd.AddCommand(structureAddressNextCmd)
	structureAddressCmd.AddCommand(structureAddressAddCmd)
	structureAddressCmd.AddCommand(structureAddressUpdateCmd)

	structureCmd.AddCommand(structureJobCmd)
	structureJobCmd.AddCommand(structureJobInfoCmd)
	structureJobCmd.AddCommand(structureJobStateCmd)
	structureJobCmd.AddCommand(structureJobDoCreateCmd)
	structureJobCmd.AddCommand(structureJobDoDestroyCmd)
	structureJobCmd.AddCommand(structureJobDoUtilityCmd)

	structureCmd.AddCommand(structureJobLogCmd)

	structureCmd.AddCommand(structureInterfaceCmd)
	structureInterfaceCmd.AddCommand(structureInterfaceListCmd)
	structureInterfaceCmd.AddCommand(structureInterfaceGetCmd)
	structureInterfaceCmd.AddCommand(structureInterfaceCreateCmd)
	structureInterfaceCmd.AddCommand(structureInterfaceUpdateCmd)
	structureInterfaceCmd.AddCommand(structureInterfaceDeleteCmd)

	structureCmd.AddCommand(structureAggInterfaceCmd)
	structureAggInterfaceCmd.AddCommand(structureAggInterfaceListCmd)
	structureAggInterfaceCmd.AddCommand(structureAggInterfaceGetCmd)
	structureAggInterfaceCmd.AddCommand(structureAggInterfaceCreateCmd)
	structureAggInterfaceCmd.AddCommand(structureAggInterfaceUpdateCmd)
	structureAggInterfaceCmd.AddCommand(structureAggInterfaceDeleteCmd)
}
