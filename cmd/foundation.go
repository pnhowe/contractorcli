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
	"sort"
	"strconv"

	contractor "github.com/t3kton/contractor_goclient"

	cinp "github.com/cinp/go"
	"github.com/spf13/cobra"
)

type foundationTypeEntry struct {
	URI        string // Namespace URI, for checking to see if it's loaded and the API version
	APIVersion string
}

var fundationTypes = map[string]foundationTypeEntry{}

func foundationArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a Foundation Id(Locator) argument")
	}
	return nil
}

func foundationInterfaceArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a Structure Interface Id argument")
	}
	return nil
}

var foundationCmd = &cobra.Command{
	Use:   "foundation",
	Short: "Work with Foundations",
}

var foundationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Foundations",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		rl := []cinp.Object{}
		vchan, err := contractorClient.BuildingFoundationList(ctx, "", map[string]interface{}{})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Id", "Site", "Locator", "Structure", "State", "Blueprint", "Created", "Updated"}, "{{.GetURI | extractID}}	{{.Site | extractID}}	{{.Locator}}	{{.Structure | extractID}}	{{.State}}	{{.Blueprint | extractID}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var foundationGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.BuildingFoundationGet(ctx, foundationID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:             {{.GetURI | extractID}}
Locator:       {{.Locator}}
Type:          {{.Type}}
Site:          {{.Site | extractID}}
Blueprint:     {{.Blueprint | extractID}}
Structure:     {{.Structure | extractID}}
Id Map:        {{.IDMap}}
Class List:    {{.ClassList}}
State:         {{.State}}
Located At:    {{.LocatedAt}}
Built At:      {{.BuiltAt}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)
		return nil
	},
}

var foundationTypesCmd = &cobra.Command{
	Use:   "types",
	Short: "List Supported Types",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		typeList := []string{}
		for k, v := range fundationTypes {
			APIVersion, err := contractorClient.GetAPIVersion(ctx, v.URI)
			if err != nil {
				// return err
				continue // TODO: really should only do this if it is a 404
			}
			if APIVersion != v.APIVersion {
				continue // API Version mismatch, print a warning?
			}
			typeList = append(typeList, k)
		}
		outputKV(map[string]interface{}{"type": typeList})

		return nil
	},
}

var foundationDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.BuildingFoundationGet(ctx, foundationID)
		if err != nil {
			return err
		}
		if err := o.Delete(ctx); err != nil {
			return err
		}

		return nil
	},
}

var foundationInterfaceCmd = &cobra.Command{
	Use:   "interface",
	Short: "Work with Foundation Interfaces",
}

var foundationInterfaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Interfaces attached to a foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.BuildingFoundationGet(ctx, foundationID)
		if err != nil {
			return err
		}

		rl := []cinp.Object{}
		vchan, err := contractorClient.UtilitiesRealNetworkInterfaceList(ctx, "foundation", map[string]interface{}{"foundation": o.GetURI()})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}

		sort.Slice(rl, func(i, j int) bool {
			return *(rl[i].(*contractor.UtilitiesRealNetworkInterface).ID) < *(rl[j].(*contractor.UtilitiesRealNetworkInterface).ID)
		})

		outputList(rl, []string{"Id", "Name", "Physical Location", "MAC", "Is Provisioning", "Network", "Link Name", "PXE", "Created", "Update"}, "{{.GetURI | extractID}}	{{.Name}}	{{.PhysicalLocation}}	{{.Mac}}	{{.IsProvisioning}}	{{.Network | extractID}}	{{.LinkName}}	{{.Pxe| extractID}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var foundationInterfaceGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Foundation Interface",
	Args:  foundationInterfaceArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		interfaceID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o, err := contractorClient.UtilitiesRealNetworkInterfaceGet(ctx, interfaceID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:               {{.GetURI | extractID}}
Name:             {{.Name}}
Type:             {{.Type}}
Network:          {{.Network | extractID}}
Foundation:       {{.Foundation | extractID}}
Mac:              {{.Mac}}
IsProvisioning:   {{.IsProvisioning}}
PhysicalLocation: {{.PhysicalLocation}}
LinkName:         {{.LinkName}}
Pxe:              {{.Pxe}}
Created:          {{.Created}}
Updated:          {{.Updated}}
`)
		return nil
	},
}

var foundationInterfaceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Foundation Interface",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		r, err := contractorClient.BuildingFoundationGet(ctx, foundationID)
		if err != nil {
			return err
		}

		o := contractorClient.UtilitiesRealNetworkInterfaceNew()
		(*o.Foundation) = r.GetURI()
		o.Name = &detailName
		o.PhysicalLocation = &detailPhysicalLocation
		o.IsProvisioning = &detailIsProvisioning
		o.LinkName = &detailLinkName
		o.Mac = &detailMac

		if detailNetwork != 0 {
			r, err := contractorClient.UtilitiesNetworkGet(ctx, detailNetwork)
			if err != nil {
				return err
			}
			(*o.Network) = r.GetURI()
		}

		o, err = o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:               {{.GetURI | extractID}}
Name:             {{.Name}}
Type:             {{.Type}}
Network:          {{.Network | extractID}}
Foundation:       {{.Foundation | extractID}}
Mac:              {{.Mac}}
IsProvisioning:   {{.IsProvisioning}}
PhysicalLocation: {{.PhysicalLocation}}
LinkName:         {{.LinkName}}
Pxe:              {{.Pxe}}
Created:          {{.Created}}
Updated:          {{.Updated}}
`)

		return nil
	},
}

var foundationInterfaceUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Foundation Interface",
	Args:  foundationInterfaceArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		interfaceID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o := contractorClient.UtilitiesRealNetworkInterfaceNewWithID(interfaceID)

		if detailName != "" {
			o.Name = &detailName
		}

		if detailPhysicalLocation != "" {
			o.PhysicalLocation = &detailPhysicalLocation
		}

		// if detailIsProvisioning != "" {
		// 	o.IsProvisioning = detailIsProvisioning
		// }

		if detailLinkName != "" {
			o.LinkName = &detailLinkName
		}

		if detailMac != "" {
			o.Mac = &detailMac
		}

		if detailNetwork != 0 {
			r, err := contractorClient.UtilitiesNetworkGet(ctx, detailNetwork)
			if err != nil {
				return err
			}
			(*o.Network) = r.GetURI()
		}

		o, err = o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:               {{.GetURI | extractID}}
Name:             {{.Name}}
Type:             {{.Type}}
Network:          {{.Network | extractID}}
Foundation:       {{.Foundation | extractID}}
Mac:              {{.Mac}}
IsProvisioning:   {{.IsProvisioning}}
PhysicalLocation: {{.PhysicalLocation}}
LinkName:         {{.LinkName}}
Pxe:              {{.Pxe}}
Created:          {{.Created}}
Updated:          {{.Updated}}
`)

		return nil
	},
}

var foundationInterfacePXECmd = &cobra.Command{
	Use:   "pxe",
	Short: "SetFoundation Interface PXE",
	Args:  foundationInterfaceArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		interfaceID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o := contractorClient.UtilitiesRealNetworkInterfaceNewWithID(interfaceID)
		(*o.Pxe) = fmt.Sprintf("/api/v1/BluePrint/PXE:%s:", detailPxeName)

		_, err = o.Update(ctx)
		if err != nil {
			return err
		}

		return nil
	},
}

var foundationInterfaceDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Foundation Interface",
	Args:  foundationInterfaceArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		interfaceID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o, err := contractorClient.UtilitiesRealNetworkInterfaceGet(ctx, interfaceID)
		if err != nil {
			return err
		}
		if err := o.Delete(ctx); err != nil {
			return err
		}

		return nil
	},
}

var foundationJobCmd = &cobra.Command{
	Use:   "job",
	Short: "Work with Foundation Jobs",
}

var foundationJobInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show Foundation Job Info",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.BuildingFoundationGet(ctx, foundationID)
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

		outputDetail(j, `Job:           {{.GetURI | extractID}}
Site:          {{.Site}}
Foundation:    {{.Foundation | extractID}}
State:         {{.State}}
Status:        {{.Status}}
Message:       {{.Message}}
Script Name:   {{.ScriptName}}
Can Start:     {{.CanStart}}
Updated:       {{.Updated}}
Created:       {{.Created}}
`)
		return nil
	},
}

var foundationJobStateCmd = &cobra.Command{
	Use:   "state",
	Short: "Show Foundation Job State",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.BuildingFoundationGet(ctx, foundationID)
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

		vars, err := j.CallJobRunnerVariables(ctx)
		if err != nil {
			return err
		}

		state, err := j.CallJobRunnerState(ctx)
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

var foundationJobDoCreateCmd = &cobra.Command{
	Use:   "do-create",
	Short: "Start Create Job for Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.BuildingFoundationGet(ctx, foundationID)
		if err != nil {
			return err
		}

		if _, err := o.CallDoCreate(ctx); err != nil {
			return err
		}

		return nil
	},
}

var foundationJobDoDestroyCmd = &cobra.Command{
	Use:   "do-destroy",
	Short: "Start Destroy Job for Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.BuildingFoundationGet(ctx, foundationID)
		if err != nil {
			return err
		}

		if _, err := o.CallDoDestroy(ctx); err != nil {
			return err
		}
		return nil
	},
}

var foundationJobDoUtilityCmd = &cobra.Command{
	Use:   "do-utility",
	Short: "Start Utility Job for Foundation",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("requires a Foundation Id(Locator) and Utility Job Name argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]
		scriptName := args[1]

		ctx := cmd.Context()

		o, err := contractorClient.BuildingFoundationGet(ctx, foundationID)
		if err != nil {
			return err
		}

		if _, err := o.CallDoJob(ctx, scriptName); err != nil {
			return err
		}
		return nil
	},
}

var foundationJobLogCmd = &cobra.Command{
	Use:   "joblog",
	Short: "Job Log List for a Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.BuildingFoundationGet(ctx, foundationID)
		if err != nil {
			return err
		}

		rl := []cinp.Object{}
		vchan, err := contractorClient.ForemanJobLogList(ctx, "foundation", map[string]interface{}{"foundation": o.GetURI()})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Script Name", "Created By", "Started At", "Finished At", "Cancled By", "Cancled At", "Created", "Updated"}, "{{.ScriptName}}	{{.Creator}}	{{.StartedAt}}	{{.FinishedAt}}	{{.CanceledBy}}	{{.CanceledAt}}	{{.Updated}}	{{.Created}}\n")

		return nil
	},
}

var foundationBootToCmd = &cobra.Command{
	Use:   "bootto",
	Short: "Set Default Interface PXE",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.BuildingFoundationGet(ctx, foundationID)
		if err != nil {
			return err
		}

		vchan, err := contractorClient.UtilitiesRealNetworkInterfaceList(ctx, "foundation", map[string]interface{}{"foundation": o.GetURI()})
		if err != nil {
			return err
		}
		for v := range vchan {
			if *v.IsProvisioning {
				v = contractorClient.UtilitiesRealNetworkInterfaceNewWithID(*v.ID)
				(*v.Pxe) = fmt.Sprintf("/api/v1/BluePrint/PXE:%s:", detailPxeName)

				_, err := v.Update(ctx)
				if err != nil {
					return err
				}
				break
			}

		}

		return nil
	},
}

func init() {
	foundationInterfaceCreateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Name of the new Interface")
	foundationInterfaceCreateCmd.Flags().StringVarP(&detailPhysicalLocation, "physical", "y", "", "Physical Location of the new Interface")
	foundationInterfaceCreateCmd.Flags().IntVarP(&detailNetwork, "network", "t", 0, "Network id to attach the new Interface to")
	foundationInterfaceCreateCmd.Flags().BoolVarP(&detailIsProvisioning, "provisioning", "p", false, "New Interface is the provisioning interface")
	foundationInterfaceCreateCmd.Flags().StringVarP(&detailLinkName, "linkname", "l", "", "Link Name of the new Interface")
	foundationInterfaceCreateCmd.Flags().StringVarP(&detailMac, "mac", "m", "", "MAC of the new Interface")

	foundationInterfaceUpdateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Update the Name of the Interface")
	foundationInterfaceUpdateCmd.Flags().StringVarP(&detailPhysicalLocation, "physical", "y", "", "Update the Physical Location of the Interface")
	foundationInterfaceUpdateCmd.Flags().IntVarP(&detailNetwork, "network", "t", 0, "Update Network id the Interface is attached to")
	//foundationInterfaceUpdateCmd.Flags().BoolVarP(&detailIsProvisioning, "provisioning", "p", false, "New Interface is the provisioning interface")
	foundationInterfaceUpdateCmd.Flags().StringVarP(&detailLinkName, "linkname", "l", "", "Update the Link Name of the Interface")
	foundationInterfaceUpdateCmd.Flags().StringVarP(&detailMac, "mac", "m", "", "Update the MAC of the Interface")

	foundationInterfacePXECmd.Flags().StringVarP(&detailPxeName, "name", "n", "normal-boot", "Update the Name PXE to set")

	foundationBootToCmd.Flags().StringVarP(&detailPxeName, "name", "n", "normal-boot", "PXE to boot to")

	rootCmd.AddCommand(foundationCmd)
	foundationCmd.AddCommand(foundationListCmd, foundationGetCmd, foundationTypesCmd, foundationDeleteCmd, foundationBootToCmd)

	foundationCmd.AddCommand(foundationInterfaceCmd)
	foundationInterfaceCmd.AddCommand(foundationInterfaceListCmd, foundationInterfaceGetCmd, foundationInterfaceCreateCmd, foundationInterfaceUpdateCmd, foundationInterfaceDeleteCmd, foundationInterfacePXECmd)

	foundationCmd.AddCommand(foundationJobCmd)
	foundationJobCmd.AddCommand(foundationJobInfoCmd, foundationJobStateCmd, foundationJobDoCreateCmd, foundationJobDoDestroyCmd, foundationJobDoUtilityCmd)

	foundationCmd.AddCommand(foundationJobLogCmd)
}
