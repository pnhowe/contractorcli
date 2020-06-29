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
	"github.com/spf13/cobra"
)

var foundationProxmoxCmd = &cobra.Command{
	Use:   "proxmox",
	Short: "Work with Proxmox Foundations",
}

var foundationProxmoxGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Proxmox Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.ProxmoxProxmoxFoundationGet(foundationID)
		if err != nil {
			return err
		}
		outputDetail(r, `Locator:        {{.Locator}}
Complex:        {{.ProxmoxComplex | extractID}}
VM UUID:        {{.ProxmoxVmID}}
Type:           {{.Type}}
Site:           {{.Site | extractID}}
Blueprint:      {{.Blueprint | extractID}}
Id Map:         {{.IDMap}}
Class List:     {{.ClassList}}
State:          {{.State}}
Located At:     {{.LocatedAt}}
Built At:       {{.BuiltAt}}
Created:        {{.Created}}
Updated:        {{.Updated}}
`)

		return nil
	},
}

var foundationProxmoxCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Proxmox Foundation",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		o := c.ProxmoxProxmoxFoundationNew()
		o.Locator = detailLocator

		if detailSite != "" {
			r, err := c.SiteSiteGet(detailSite)
			if err != nil {
				return err
			}
			o.Site = r.GetID()
		}

		if detailBlueprint != "" {
			r, err := c.BlueprintFoundationBluePrintGet(detailBlueprint)
			if err != nil {
				return err
			}
			o.Blueprint = r.GetID()
		}

		if detailComplex != "" {
			r, err := c.ProxmoxProxmoxComplexGet(detailComplex)
			if err != nil {
				return err
			}
			o.ProxmoxComplex = r.GetID()
		}

		if err := o.Create(); err != nil {
			return err
		}

		outputKV(map[string]interface{}{"id": o.GetID()})

		return nil
	},
}

var foundationProxmoxUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Proxmox Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		foundationID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.ProxmoxProxmoxFoundationGet(foundationID)
		if err != nil {
			return err
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
			r, err := c.BlueprintFoundationBluePrintGet(detailBlueprint)
			if err != nil {
				return err
			}
			o.Blueprint = r.GetID()
			fieldList = append(fieldList, "blueprint")
		}

		if detailComplex != "" {
			r, err := c.ProxmoxProxmoxComplexGet(detailComplex)
			if err != nil {
				return err
			}
			o.ProxmoxComplex = r.GetID()
			fieldList = append(fieldList, "proxmox_complex")
		}

		if err := o.Update(fieldList); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	fundationTypes["proxmox"] = foundationTypeEntry{"/api/v1/Proxmox/", "0.1"}

	foundationProxmoxCreateCmd.Flags().StringVarP(&detailLocator, "locator", "l", "", "Locator of New Proxmox Foundation")
	foundationProxmoxCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New Proxmox Foundation")
	foundationProxmoxCreateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Blueprint of New Proxmox Foundation")
	foundationProxmoxCreateCmd.Flags().StringVarP(&detailComplex, "complex", "x", "", "Plot of New Proxmox Foundation")

	foundationProxmoxUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Foundation with value")
	foundationProxmoxUpdateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Update the Blueprint of Foundation with value")
	foundationProxmoxUpdateCmd.Flags().StringVarP(&detailComplex, "complex", "x", "", "Update the Plot of the Proxmox Foundation")

	foundationCmd.AddCommand(foundationProxmoxCmd)
	foundationProxmoxCmd.AddCommand(foundationProxmoxGetCmd)
	foundationProxmoxCmd.AddCommand(foundationProxmoxCreateCmd)
	foundationProxmoxCmd.AddCommand(foundationProxmoxUpdateCmd)
}
