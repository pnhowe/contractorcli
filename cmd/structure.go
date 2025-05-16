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

	contractor "github.com/t3kton/contractor_goclient"

	cinp "github.com/cinp/go"
	"github.com/spf13/cobra"
)

func structureArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a Structure Id argument")
	}
	return nil
}

func structureAddressArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a Address Id argument")
	}
	return nil
}

func structureInterfaceArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a Structure Interface Id argument")
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
		ctx := cmd.Context()

		rl := []cinp.Object{}
		vchan, err := contractorClient.BuildingStructureList(ctx, "", map[string]interface{}{})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Id", "Site", "Hostname", "Foundation", "State", "Blueprint", "Created", "Updated"}, "{{.GetURI | extractID}}	{{.Site | extractID}}	{{.Hostname}}	{{.Foundation | extractID}}	{{.State}}	{{.Blueprint | extractID}}	{{.Created}}	{{.Updated}}\n")

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
		ctx := cmd.Context()

		o, err := contractorClient.BuildingStructureGet(ctx, structureID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:            {{.GetURI | extractID}}
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
		ctx := cmd.Context()

		o := contractorClient.BuildingStructureNew()
		o.Hostname = &detailHostname

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			site := r.GetURI()
			o.Site = &site
		}

		if detailBlueprint != "" {
			r, err := contractorClient.BlueprintStructureBluePrintGet(ctx, detailBlueprint)
			if err != nil {
				return err
			}
			blueprint := r.GetURI()
			o.Blueprint = &blueprint
		}

		if detailFoundation != "" {
			r, err := contractorClient.BuildingFoundationGet(ctx, detailFoundation)
			if err != nil {
				return err
			}
			foundation := r.GetURI()
			o.Foundation = &foundation
		}

		o, err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI | extractID}}
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

var structureUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Structure",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o := contractorClient.BuildingStructureNewWithID(structureID)

		if detailHostname != "" {
			o.Hostname = &detailHostname
		}

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			site := r.GetURI()
			o.Site = &site
		}

		if detailBlueprint != "" {
			r, err := contractorClient.BlueprintStructureBluePrintGet(ctx, detailBlueprint)
			if err != nil {
				return err
			}
			blueprint := r.GetURI()
			o.Blueprint = &blueprint
		}

		if detailFoundation != "" {
			r, err := contractorClient.BuildingFoundationGet(ctx, detailFoundation)
			if err != nil {
				return err
			}
			foundation := r.GetURI()
			o.Foundation = &foundation
		}

		o, err = o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI | extractID}}
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

var structureDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Structure",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o, err := contractorClient.BuildingStructureGet(ctx, structureID)
		if err != nil {
			return err
		}
		if err := o.Delete(ctx); err != nil {
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

		ctx := cmd.Context()

		if configSetName != "" {
			r, err := contractorClient.BuildingStructureGet(ctx, structureID)
			if err != nil {
				return err
			}

			o := contractorClient.BuildingStructureNewWithID(structureID)
			o.ConfigValues = r.ConfigValues
			(*o.ConfigValues)[configSetName] = configSetValue

			o, err = o.Update(ctx)
			if err != nil {
				return err
			}
			outputKV(*o.ConfigValues)

		} else if configDeleteName != "" {
			r, err := contractorClient.BuildingStructureGet(ctx, structureID)
			if err != nil {
				return err
			}

			o := contractorClient.BuildingStructureNewWithID(structureID)
			o.ConfigValues = r.ConfigValues
			delete(*o.ConfigValues, configDeleteName)

			o, err = o.Update(ctx)
			if err != nil {
				return err
			}
			outputKV(*o.ConfigValues)

		} else if configFull {
			o := contractorClient.BuildingStructureNewWithID(structureID)
			r, err := o.CallGetConfig(ctx)
			if err != nil {
				return err
			}
			outputKV(r)

		} else {
			o, err := contractorClient.BuildingStructureGet(ctx, structureID)
			if err != nil {
				return err
			}
			outputKV(*o.ConfigValues)
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

		ctx := cmd.Context()

		o, err := contractorClient.BuildingStructureGet(ctx, structureID)
		if err != nil {
			return err
		}

		rl := []cinp.Object{}
		vchan, err := contractorClient.UtilitiesAddressList(ctx, "structure", map[string]interface{}{"structure": o.GetURI()})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Id", "Interface", "Address", "Address Block", "Offset", "Is Primary", "Created", "Updated"}, "{{.GetURI | extractID}}	{{.InterfaceName}}	{{.IPAddress}}	{{.AddressBlock | extractID}}	{{.Offset}}	{{.IsPrimary}}	{{.Updated}}	{{.Created}}\n")

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

		ctx := cmd.Context()

		r, err := contractorClient.BuildingStructureGet(ctx, structureID)
		if err != nil {
			return err
		}

		a, err := contractorClient.UtilitiesAddressBlockGet(ctx, detailAddressBlock)
		if err != nil {
			return err
		}

		addressURI, err := a.CallNextAddress(ctx, strings.Replace(r.GetURI(), "/api/v1/Building/Structure", "/api/v1/Utilities/Networked", 1), detailInterfaceName, detailIsPrimary)
		if err != nil {
			return err
		}

		o, err := contractorClient.UtilitiesAddressGetURI(ctx, addressURI)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI | extractID}}
AddressBlock:  {{.AddressBlock | extractID}}
Offset:        {{.Offset}}
Networked:     {{.Networked | extractID}}
InterfaceName: {{.InterfaceName}}
AliasIndex:    {{.AliasIndex}}
Pointer:       {{.Pointer}}
IsPrimary:     {{.IsPrimary}}
Type:          {{.Type}}
IPAddress:     {{.IPAddress}}
Subnet:        {{.Subnet}}
Netmask:       {{.Netmask}}
Prefix:        {{.Prefix}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)
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

		ctx := cmd.Context()

		s, err := contractorClient.BuildingStructureGet(ctx, structureID)
		if err != nil {
			return err
		}

		a, err := contractorClient.UtilitiesAddressBlockGet(ctx, detailAddressBlock)
		if err != nil {
			return err
		}

		o := contractorClient.UtilitiesAddressNew()
		addressBlock := a.GetURI()
		o.AddressBlock = &addressBlock
		o.Offset = &detailOffset
		networked := strings.Replace(s.GetURI(), "/api/v1/Building/Structure", "/api/v1/Utilities/Networked", 1)
		o.Networked = &networked
		o.InterfaceName = &detailInterfaceName
		o.IsPrimary = &detailIsPrimary

		o, err = o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI | extractID}}
AddressBlock:  {{.AddressBlock | extractID}}
Offset:        {{.Offset}}
Networked:     {{.Networked | extractID}}
InterfaceName: {{.InterfaceName}}
AliasIndex:    {{.AliasIndex}}
Pointer:       {{.Pointer}}
IsPrimary:     {{.IsPrimary}}
Type:          {{.Type}}
IPAddress:     {{.IPAddress}}
Subnet:        {{.Subnet}}
Netmask:       {{.Netmask}}
Prefix:        {{.Prefix}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)
		return nil
	},
}

var structureAddressUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an IP address to structure",
	Args:  structureAddressArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		addressID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o := contractorClient.UtilitiesAddressNewWithID(addressID)

		if detailOffset != 0 {
			o.Offset = &detailOffset
		}

		if detailInterfaceName != "" {
			o.InterfaceName = &detailInterfaceName
		}

		o, err = o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI | extractID}}
AddressBlock:  {{.AddressBlock | extractID}}
Offset:        {{.Offset}}
Networked:     {{.Networked | extractID}}
InterfaceName: {{.InterfaceName}}
AliasIndex:    {{.AliasIndex}}
Pointer:       {{.Pointer}}
IsPrimary:     {{.IsPrimary}}
Type:          {{.Type}}
IPAddress:     {{.IPAddress}}
Subnet:        {{.Subnet}}
Netmask:       {{.Netmask}}
Prefix:        {{.Prefix}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)

		return nil
	},
}

var structureAddressDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete IP address from structure",
	Args:  structureAddressArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		addressID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o, err := contractorClient.UtilitiesAddressGet(ctx, addressID)
		if err != nil {
			return err
		}
		if err := o.Delete(ctx); err != nil {
			return err
		}

		return nil
	},
}

var structureJobCmd = &cobra.Command{
	Use:   "job",
	Short: "Work with Structure Jobs",
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

		ctx := cmd.Context()

		o, err := contractorClient.BuildingStructureGet(ctx, structureID)
		if err != nil {
			return err
		}

		if _, err := o.CallDoCreate(ctx); err != nil {
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

		ctx := cmd.Context()

		o, err := contractorClient.BuildingStructureGet(ctx, structureID)
		if err != nil {
			return err
		}

		if _, err := o.CallDoDestroy(ctx); err != nil {
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
			return errors.New("requires a Structure Id and Utility Job Name Argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		scriptName := args[1]

		ctx := cmd.Context()

		o, err := contractorClient.BuildingStructureGet(ctx, structureID)
		if err != nil {
			return err
		}

		if _, err := o.CallDoJob(ctx, scriptName); err != nil {
			return err
		}
		return nil
	},
}

var structureJobIdCmd = &cobra.Command{
	Use:   "id",
	Short: "Get the Job Id for active job",
	Args:  structureArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		structureID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o, err := contractorClient.BuildingStructureGet(ctx, structureID)
		if err != nil {
			return err
		}
		jobURI, err := o.CallGetJob(ctx)
		if err != nil {
			return err
		}

		if jobURI == "" {
			outputDetail("", "No Job\n")
			return nil
		}

		j, err := contractorClient.ForemanStructureJobGetURI(ctx, jobURI)
		if err != nil {
			return err
		}

		valueMap := map[string]interface{}{"Id": extractID(j.GetURI())}
		outputKV(valueMap)

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

		ctx := cmd.Context()

		o, err := contractorClient.BuildingStructureGet(ctx, structureID)
		if err != nil {
			return err
		}

		rl := []cinp.Object{}
		vchan, err := contractorClient.ForemanJobLogList(ctx, "structure", map[string]interface{}{"structure": o.GetURI()})
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

		ctx := cmd.Context()

		o, err := contractorClient.BuildingStructureGet(ctx, structureID)
		if err != nil {
			return err
		}

		rl := []cinp.Object{}
		vchan, err := contractorClient.UtilitiesAbstractNetworkInterfaceList(ctx, "structure", map[string]interface{}{"structure": o.GetURI()})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}

		sort.Slice(rl, func(i, j int) bool {
			return *(rl[i].(*contractor.UtilitiesAbstractNetworkInterface).ID) < *(rl[j].(*contractor.UtilitiesAbstractNetworkInterface).ID)
		})

		outputList(rl, []string{"Id", "Name", "Type", "Network", "Created", "Update"}, "{{.GetURI | extractID}}	{{.Name}}	{{.Type}}	{{.Network | extractID}}	{{.Created}}	{{.Updated}}\n")

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

		ctx := cmd.Context()

		o, err := contractorClient.UtilitiesAbstractNetworkInterfaceGet(ctx, interfaceID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:               {{.GetURI | extractID}}
Name:             {{.Name}}
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

		ctx := cmd.Context()

		r, err := contractorClient.BuildingStructureGet(ctx, structureID)
		if err != nil {
			return err
		}

		o := contractorClient.UtilitiesAbstractNetworkInterfaceNew()
		structure := r.GetURI()
		o.Structure = &structure
		o.Name = &detailName

		if detailNetwork != 0 {
			r, err := contractorClient.UtilitiesNetworkGet(ctx, detailNetwork)
			if err != nil {
				return err
			}
			network := r.GetURI()
			o.Network = &network
		}

		o, err = o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:               {{.GetURI | extractID}}
Name:             {{.Name}}
Type:             {{.Type}}
Network:          {{.Network | extractID}}
Structure:        {{.Structure | extractID}}
Created:          {{.Created}}
Updated:          {{.Updated}}
`)
		return nil
	},
}

var structureInterfaceUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Structure Interface",
	Args:  structureInterfaceArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		interfaceID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o := contractorClient.UtilitiesAbstractNetworkInterfaceNewWithID(interfaceID)

		if detailName != "" {
			o.Name = &detailName
		}

		if detailNetwork != 0 {
			r, err := contractorClient.UtilitiesNetworkGet(ctx, detailNetwork)
			if err != nil {
				return err
			}
			network := r.GetURI()
			o.Network = &network
		}

		o, err = o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:               {{.GetURI | extractID}}
Name:             {{.Name}}
Type:             {{.Type}}
Network:          {{.Network | extractID}}
Structure:        {{.Structure | extractID}}
Created:          {{.Created}}
Updated:          {{.Updated}}
`)

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

		ctx := cmd.Context()

		o, err := contractorClient.UtilitiesAbstractNetworkInterfaceGet(ctx, interfaceID)
		if err != nil {
			return err
		}
		if err := o.Delete(ctx); err != nil {
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

		ctx := cmd.Context()

		o, err := contractorClient.BuildingStructureGet(ctx, structureID)
		if err != nil {
			return err
		}

		rl := []cinp.Object{}
		vchan, err := contractorClient.UtilitiesAggregatedNetworkInterfaceList(ctx, "structure", map[string]interface{}{"structure": o.GetURI()})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}

		sort.Slice(rl, func(i, j int) bool {
			return *(rl[i].(*contractor.UtilitiesAbstractNetworkInterface).ID) < *(rl[j].(*contractor.UtilitiesAbstractNetworkInterface).ID)
		})

		outputList(rl, []string{"Id", "Name", "Network", "Primary Interface", "Secondary Interface(s)", "Created", "Update"}, "{{.GetURI | extractID}}	{{.Name}}	{{.Network | extractID}}	{{.PrimaryInterface | extractID}}	{{.SecondaryInterfaces | extractIDList}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var structureAggInterfaceGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Structure Aggregated Interface",
	Args:  structureInterfaceArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		interfaceID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.UtilitiesAggregatedNetworkInterfaceGet(ctx, interfaceID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:               {{.GetURI | extractID}}
Name:             {{.Name}}
Network:          {{.Network | extractID}}
Structure:        {{.Structure | extractID}}
Primary:          {{.PrimaryInterface | extractID}}
Secondaries:      {{.SecondaryInterfaces | extractIDList}}
Paramaters:       {{.Paramaters}}
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

		ctx := cmd.Context()

		r, err := contractorClient.BuildingStructureGet(ctx, structureID)
		if err != nil {
			return err
		}

		o := contractorClient.UtilitiesAggregatedNetworkInterfaceNew()
		structure := r.GetURI()
		o.Structure = &structure
		o.Name = &detailName

		if detailNetwork != 0 {
			r, err := contractorClient.UtilitiesNetworkGet(ctx, detailNetwork)
			if err != nil {
				return err
			}
			network := r.GetURI()
			o.Network = &network
		}

		ri := contractorClient.UtilitiesNetworkInterfaceNewWithID(detailPrimary)
		primaryInterface := ri.GetURI()
		o.PrimaryInterface = &primaryInterface

		o.SecondaryInterfaces = &[]string{}
		for _, id := range strings.Split(detailSecondary, ",") {
			i, err := strconv.Atoi(id)
			if err != nil {
				return err
			}
			ri = contractorClient.UtilitiesNetworkInterfaceNewWithID(i)
			(*o.SecondaryInterfaces) = append((*o.SecondaryInterfaces), ri.GetURI())
		}

		o, err = o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:               {{.GetURI | extractID}}
Name:             {{.Name}}
Network:          {{.Network | extractID}}
Structure:        {{.Structure | extractID}}
Primary:          {{.PrimaryInterface | extractID}}
Secondaries:      {{.SecondaryInterfaces | extractIDList}}
Paramaters:       {{.Paramaters}}
Created:          {{.Created}}
Updated:          {{.Updated}}
`)
		return nil
	},
}

var structureAggInterfaceUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Structure Aggregated Interface",
	Args:  structureInterfaceArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		interfaceID := args[0]

		ctx := cmd.Context()

		o := contractorClient.UtilitiesAggregatedNetworkInterfaceNewWithID(interfaceID)

		if detailName != "" {
			o.Name = &detailName
		}

		if detailNetwork != 0 {
			r, err := contractorClient.UtilitiesNetworkGet(ctx, detailNetwork)
			if err != nil {
				return err
			}
			network := r.GetURI()
			o.Network = &network
		}

		if detailPrimary != 0 {
			r := contractorClient.UtilitiesNetworkInterfaceNewWithID(detailPrimary)
			primaryInterface := r.GetURI()
			o.PrimaryInterface = &primaryInterface
		}

		if detailSecondary != "" {
			o.SecondaryInterfaces = &[]string{}
			for id := range strings.Split(detailSecondary, ",") {
				ri := contractorClient.UtilitiesNetworkInterfaceNewWithID(id)
				(*o.SecondaryInterfaces) = append((*o.SecondaryInterfaces), ri.GetURI())
			}
		}

		o, err := o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:               {{.GetURI | extractID}}
Name:             {{.Name}}
Network:          {{.Network | extractID}}
Structure:        {{.Structure | extractID}}
Primary:          {{.PrimaryInterface | extractID}}
Secondaries:      {{.SecondaryInterfaces | extractIDList}}
Paramaters:       {{.Paramaters}}
Created:          {{.Created}}
Updated:          {{.Updated}}
`)

		return nil
	},
}

var structureAggInterfaceDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Structure Aggregated Interface",
	Args:  structureInterfaceArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		interfaceID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.UtilitiesAggregatedNetworkInterfaceGet(ctx, interfaceID)
		if err != nil {
			return err
		}
		if err := o.Delete(ctx); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	structureConfigCmd.Flags().BoolVarP(&configFull, "full", "f", false, "Display the Full/Compiled config")
	structureConfigCmd.Flags().StringVarP(&configSetName, "set-name", "n", "", "Set Config Value Key Name, if set-value is not specified, the value will be set to ''")
	structureConfigCmd.Flags().StringVarP(&configSetValue, "set-value", "v", "", "Set Config Value, ignored if set-name is not specified") // TODO: make a numberic version
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
	structureAggInterfaceCreateCmd.Flags().IntVarP(&detailPrimary, "primary", "p", 0, "Interface name to use as the primary interface")
	structureAggInterfaceCreateCmd.Flags().StringVarP(&detailSecondary, "secondary", "s", "", "Interface names to use as the secondaries, delimited by ','")

	structureAggInterfaceUpdateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Update the Name of the Interface")
	structureAggInterfaceUpdateCmd.Flags().IntVarP(&detailNetwork, "network", "t", 0, "Update Network id the Interface is attached to")
	structureAggInterfaceUpdateCmd.Flags().IntVarP(&detailPrimary, "primary", "p", 0, "Interface name to use as the primary interface")
	structureAggInterfaceUpdateCmd.Flags().StringVarP(&detailSecondary, "secondary", "s", "", "Interface names to use as the secondaries, delimited by ','")

	rootCmd.AddCommand(structureCmd)
	structureCmd.AddCommand(structureListCmd, structureGetCmd, structureCreateCmd, structureUpdateCmd, structureDeleteCmd, structureConfigCmd)

	structureCmd.AddCommand(structureAddressCmd)
	structureAddressCmd.AddCommand(structureAddressListCmd, structureAddressNextCmd, structureAddressAddCmd, structureAddressUpdateCmd, structureAddressDeleteCmd)

	structureCmd.AddCommand(structureJobCmd)
	structureJobCmd.AddCommand(structureJobIdCmd, structureJobDoCreateCmd, structureJobDoDestroyCmd, structureJobDoUtilityCmd)

	structureCmd.AddCommand(structureJobLogCmd)

	structureCmd.AddCommand(structureInterfaceCmd)
	structureInterfaceCmd.AddCommand(structureInterfaceListCmd, structureInterfaceGetCmd, structureInterfaceCreateCmd, structureInterfaceUpdateCmd, structureInterfaceDeleteCmd)

	structureCmd.AddCommand(structureAggInterfaceCmd)
	structureAggInterfaceCmd.AddCommand(structureAggInterfaceListCmd, structureAggInterfaceGetCmd, structureAggInterfaceCreateCmd, structureAggInterfaceUpdateCmd, structureAggInterfaceDeleteCmd)
}
