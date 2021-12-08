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

var complexProxmoxCmd = &cobra.Command{
	Use:   "proxmox",
	Short: "Work with Proxmox Complexes",
}

var complexProxmoxGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Proxmox Complexes",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		complexID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.ProxmoxProxmoxComplexGet(complexID)
		if err != nil {
			return err
		}
		outputDetail(r, `Name:               {{.Name}}
Description:        {{.Description}}
Type:               {{.Type}}
State:              {{.State}}
Site:               {{.Site | extractID}}
BuiltPercentage:    {{.BuiltPercentage}}
Members:            {{.Members}}
ProxmoxUsername:    {{.ProxmoxUsername}}
ProxmoxPassword:    {{.ProxmoxPassword}}
Created:            {{.Created}}
Updated:            {{.Updated}}
`)

		return nil
	},
}

var complexProxmoxCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Proxmox Complex",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		o := c.ProxmoxProxmoxComplexNew()
		o.Name = detailName
		o.Description = detailDescription
		o.BuiltPercentage = detailBuiltPercentage
		o.ProxmoxUsername = detailUsername
		o.ProxmoxPassword = detailPassword

		if detailSite != "" {
			r, err := c.SiteSiteGet(detailSite)
			if err != nil {
				return err
			}
			o.Site = r.GetID()
		}

		for _, v := range detailMembers {
			s, err := c.BuildingStructureGet(v)
			if err != nil {
				return err
			}
			o.Members = append(o.Members, s.GetID())
		}

		if err := o.Create(); err != nil {
			return err
		}

		outputKV(map[string]interface{}{"id": o.GetID()})

		return nil
	},
}

var complexProxmoxUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Proxmox Complex",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		complexID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.ProxmoxProxmoxComplexGet(complexID)
		if err != nil {
			return err
		}

		if detailDescription != "" {
			o.Description = detailDescription
			fieldList = append(fieldList, "description")
		}

		if detailBuiltPercentage != 0 {
			o.BuiltPercentage = detailBuiltPercentage
			fieldList = append(fieldList, "built_percentage")
		}

		if detailUsername != "" {
			o.ProxmoxUsername = detailUsername
			fieldList = append(fieldList, "proxmox_username")
		}

		if detailPassword != "" {
			o.ProxmoxPassword = detailPassword
			fieldList = append(fieldList, "proxmox_password")
		}

		if detailSite != "" {
			r, err := c.SiteSiteGet(detailSite)
			if err != nil {
				return err
			}
			o.Site = r.GetID()
			fieldList = append(fieldList, "site")
		}

		if len(detailMembers) > 0 {
			for _, v := range detailMembers {
				s, err := c.BuildingStructureGet(v)
				if err != nil {
					return err
				}
				o.Members = append(o.Members, s.GetID())
			}
			fieldList = append(fieldList, "members")
		}

		if err := o.Update(fieldList); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	complexTypes["proxmox"] = complexTypeEntry{"/api/v1/Proxmox/", "0.1"}

	complexProxmoxCreateCmd.Flags().StringVarP(&detailName, "name", "l", "", "Locator of New Proxmox Complex")
	complexProxmoxCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New Proxmox Complex")
	complexProxmoxCreateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Description of New Proxmox Complex")
	complexProxmoxCreateCmd.Flags().IntVarP(&detailBuiltPercentage, "builtperc", "b", 80, "Built Percentage of New Proxmox Complex\n(Percentage of Built Members at which the complex is considered built)")
	complexProxmoxCreateCmd.Flags().IntSliceVarP(&detailMembers, "members", "m", []int{}, "Members of the new Proxmox Complex, specify for each member")
	complexProxmoxCreateCmd.Flags().StringVarP(&detailUsername, "username", "u", "", "Proxmox Username of New Proxmox Complex")
	complexProxmoxCreateCmd.Flags().StringVarP(&detailPassword, "password", "p", "", "Proxmox Password of New Proxmox Complex")

	complexProxmoxUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Complex with value")
	complexProxmoxUpdateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Update the Description of Proxmox Complex with value")
	complexProxmoxUpdateCmd.Flags().IntVarP(&detailBuiltPercentage, "builtperc", "b", 0, "Update the Built Percentage of Proxmox Complex with value")
	complexProxmoxUpdateCmd.Flags().IntSliceVarP(&detailMembers, "members", "m", []int{}, "Update the Members of the Proxmox Complex, specify for each member")
	complexProxmoxUpdateCmd.Flags().StringVarP(&detailUsername, "username", "u", "", "Update the Proxmox Username of the Proxmox Complex with value")
	complexProxmoxUpdateCmd.Flags().StringVarP(&detailPassword, "password", "p", "", "Update the Proxmox Password of the Proxmox Complex with value")

	complexCmd.AddCommand(complexProxmoxCmd)
	complexProxmoxCmd.AddCommand(complexProxmoxGetCmd)
	complexProxmoxCmd.AddCommand(complexProxmoxCreateCmd)
	complexProxmoxCmd.AddCommand(complexProxmoxUpdateCmd)
}
