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
	cinp "github.com/cinp/go"
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

		ctx := cmd.Context()

		o, err := contractorClient.ProxmoxProxmoxComplexGet(ctx, complexID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:                 {{.GetURI | extractID}}
Name:               {{.Name}}
Description:        {{.Description}}
Type:               {{.Type}}
State:              {{.State}}
Site:               {{.Site | extractID}}
BuiltPercentage:    {{.BuiltPercentage}}
Members:            {{.Members}}
Proxmox Host:       {{.ProxmoxHost}}
Proxmox Username:   {{.ProxmoxUsername}}
Proxmox Password:   {{.ProxmoxPassword}}
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
		ctx := cmd.Context()

		o := contractorClient.ProxmoxProxmoxComplexNew()
		o.Name = &detailName
		o.Description = &detailDescription
		o.BuiltPercentage = &detailBuiltPercentage
		o.ProxmoxUsername = &detailUsername
		o.ProxmoxPassword = &detailPassword

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			o.Site = cinp.StringAddr(r.GetURI())
		}

		o.Members = &[]string{}
		for _, v := range detailMembers {
			s, err := contractorClient.BuildingStructureGet(ctx, v)
			if err != nil {
				return err
			}
			*o.Members = append(*o.Members, s.GetURI())
		}

		if detailHost != 0 {
			r, err := contractorClient.BuildingStructureGet(ctx, detailHost)
			if err != nil {
				return err
			}
			o.ProxmoxHost = cinp.StringAddr(r.GetURI())
		}

		err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:                 {{.GetURI | extractID}}
Name:               {{.Name}}
Description:        {{.Description}}
Type:               {{.Type}}
State:              {{.State}}
Site:               {{.Site | extractID}}
BuiltPercentage:    {{.BuiltPercentage}}
Members:            {{.Members}}
Proxmox Host:       {{.ProxmoxHost}}
Proxmox Username:   {{.ProxmoxUsername}}
Proxmox Password:   {{.ProxmoxPassword}}
Created:            {{.Created}}
Updated:            {{.Updated}}
`)
		return nil
	},
}

var complexProxmoxUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Proxmox Complex",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		complexID := args[0]

		ctx := cmd.Context()

		o := contractorClient.ProxmoxProxmoxComplexNewWithID(complexID)

		if detailDescription != "" {
			o.Description = &detailDescription
		}

		if detailBuiltPercentage != 0 {
			o.BuiltPercentage = &detailBuiltPercentage
		}

		if detailUsername != "" {
			o.ProxmoxUsername = &detailUsername
		}

		if detailPassword != "" {
			o.ProxmoxPassword = &detailPassword
		}

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			o.Site = cinp.StringAddr(r.GetURI())
		}

		if len(detailMembers) > 0 {
			if o.Members == nil {
				o.Members = &[]string{}
			}
			for _, v := range detailMembers {
				s, err := contractorClient.BuildingStructureGet(ctx, v)
				if err != nil {
					return err
				}
				*o.Members = append(*o.Members, s.GetURI())
			}
		}

		if detailHost != 0 {
			r, err := contractorClient.BuildingStructureGet(ctx, detailHost)
			if err != nil {
				return err
			}
			o.ProxmoxHost = cinp.StringAddr(r.GetURI())
		}

		err := o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:                 {{.GetURI | extractID}}
Name:               {{.Name}}
Description:        {{.Description}}
Type:               {{.Type}}
State:              {{.State}}
Site:               {{.Site | extractID}}
BuiltPercentage:    {{.BuiltPercentage}}
Members:            {{.Members}}
Proxmox Host:       {{.ProxmoxHost}}
Proxmox Username:   {{.ProxmoxUsername}}
Proxmox Password:   {{.ProxmoxPassword}}
Created:            {{.Created}}
Updated:            {{.Updated}}
`)

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
	complexProxmoxCreateCmd.Flags().IntVarP(&detailHost, "host", "o", 0, "Proxmox Host(structure id) of New Proxmox Complex")

	complexProxmoxUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Complex with value")
	complexProxmoxUpdateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Update the Description of Proxmox Complex with value")
	complexProxmoxUpdateCmd.Flags().IntVarP(&detailBuiltPercentage, "builtperc", "b", 0, "Update the Built Percentage of Proxmox Complex with value")
	complexProxmoxUpdateCmd.Flags().IntSliceVarP(&detailMembers, "members", "m", []int{}, "Update the Members of the Proxmox Complex, specify for each member")
	complexProxmoxUpdateCmd.Flags().StringVarP(&detailUsername, "username", "u", "", "Update the Proxmox Username of the Proxmox Complex with value")
	complexProxmoxUpdateCmd.Flags().StringVarP(&detailPassword, "password", "p", "", "Update the Proxmox Password of the Proxmox Complex with value")
	complexProxmoxUpdateCmd.Flags().IntVarP(&detailHost, "host", "o", 0, "Update the Proxmox Host(structure id) of New Proxmox Complex with value")

	complexCmd.AddCommand(complexProxmoxCmd)
	complexProxmoxCmd.AddCommand(complexProxmoxGetCmd, complexProxmoxCreateCmd, complexProxmoxUpdateCmd)
}
