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

type typeEntry struct {
	URI        string // Namespace URI, for checking to see if it's loaded and the API version
	APIVersion string
}

var fundationTypes = map[string]typeEntry{}

/*
		func(c *contractor.Contractor) foundation { return c.ManualManualFoundationNew() }, []string{}},
	"amt": typeEntry{"/api/v1/AMT/", "0.1",
		func(c *contractor.Contractor) foundation { return c.AmtAMTFoundationNew() }, []string{}},
	"ipmi": typeEntry{"/api/v1/IPMI/", "0.1",
		func(c *contractor.Contractor) foundation { return c.IpmiIPMIFoundationNew() }, []string{}},
	"vcenter": typeEntry{"/api/v1/VCenter/", "0.1",
		func(c *contractor.Contractor) foundation { return c.VcenterVCenterFoundationNew() }, []string{}},
	"virtualbox": typeEntry{"/api/v1/VirtualBox/", "0.1",
		func(c *contractor.Contractor) foundation { return c.VirtualboxVirtualBoxFoundationNew() }, []string{}},
	"azure": typeEntry{"/api/v1/Azure/", "0.1",
		func(c *contractor.Contractor) foundation { return c.AzureAzureFoundationNew() }, []string{}},
	"docker": typeEntry{"/api/v1/Docker/", "0.1",
		func(c *contractor.Contractor) foundation { return c.DockerDockerFoundationNew() }, []string{}},

*/

var detailLocator, detailPlot string

func foundationArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("Requires a Foundation Id(Locator) Argument")
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
	Run: func(cmd *cobra.Command, args []string) {
		c := getContractor()
		defer c.Logout()

		rl := []cinp.Object{}
		for v := range c.BuildingFoundationList("", map[string]interface{}{}) {
			rl = append(rl, v)
		}
		outputList(rl, "id	Site	Locator	Structure	Blueprint	Created	Updated\n", "{{.GetID | extractID}}	{{.Site | extractID}}	{{.Locator}}	{{.Structure | extractID}}	{{.Blueprint | extractID}}	{{.Created}}	{{.Updated}}\n")
	},
}

var foundationGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Foundation",
	Args:  foundationArgCheck,
	Run: func(cmd *cobra.Command, args []string) {
		foundationID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.BuildingFoundationGet(foundationID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		outputDetail(r, `Locator:       {{.Locator}}
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
	},
}

var foundationTypesCmd = &cobra.Command{
	Use:   "types",
	Short: "List Supported Types",
	Run: func(cmd *cobra.Command, args []string) {
		c := getContractor()
		defer c.Logout()

		typeList := []string{}
		for k, v := range fundationTypes {
			APIVersion, err := c.GetAPIVersion(v.URI)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if APIVersion != v.APIVersion {
				continue // API Version mismatch, print a warning?
			}
			typeList = append(typeList, k)
		}
		outputKV(map[string]interface{}{"type": typeList})
		return
	},
}

var foundationUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Foundation",
	Args:  foundationArgCheck,
	Run: func(cmd *cobra.Command, args []string) {
		fieldList := []string{}
		foundationID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.BuildingFoundationGet(foundationID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
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
			r, err := c.BlueprintFoundationBluePrintGet(detailBlueprint)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			o.Blueprint = r.GetID()
			fieldList = append(fieldList, "blueprint")
		}

		if err := o.Update(fieldList); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var foundationDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Foundation",
	Args:  foundationArgCheck,
	Run: func(cmd *cobra.Command, args []string) {
		foundationID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.BuildingFoundationGet(foundationID)
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

var foundationJobCmd = &cobra.Command{
	Use:   "job",
	Short: "Work with Foundation Jobs",
	Args:  foundationArgCheck,
	Run: func(cmd *cobra.Command, args []string) {
		foundationID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.BuildingFoundationGet(foundationID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if jobInfo {
			jID, err := o.CallGetJob()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			j, err := c.ForemanFoundationJobGet(extractID(jID))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			outputDetail(j, `Site           {{.Site | extractID}}
Foundation     {{.Foundation | extractID}}
State:         {{.State}}
Status:        {{.Status}}
Progress:      {{.Progress}}
Message:       {{.Message}}
Script Name:   {{.ScriptName}}
Can Start:     {{.CanStart}}
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
	foundationJobCmd.Flags().BoolVarP(&jobInfo, "info", "i", false, "Show Running Job Info")
	foundationJobCmd.Flags().BoolVarP(&jobCreate, "do-create", "c", false, "Submit a Create job")
	foundationJobCmd.Flags().BoolVarP(&jobDestroy, "do-destroy", "d", false, "Submit a Destroy job")
	foundationJobCmd.Flags().StringVarP(&jobUtility, "utility", "u", "", "Submit Utility Job")

	rootCmd.AddCommand(foundationCmd)
	foundationCmd.AddCommand(foundationListCmd)
	foundationCmd.AddCommand(foundationGetCmd)
	foundationCmd.AddCommand(foundationTypesCmd)
	foundationCmd.AddCommand(foundationDeleteCmd)
	foundationCmd.AddCommand(foundationJobCmd)
}
