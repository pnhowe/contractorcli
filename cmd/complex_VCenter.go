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

var complexVCenterCmd = &cobra.Command{
	Use:   "vcenter",
	Short: "Work with VCenter Complexes",
}

var complexVCenterGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get VCenter Complexes",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		complexID := args[0]

		ctx := cmd.Context()

		o, err := contractorClient.VcenterVCenterComplexGet(ctx, complexID)
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
VCenter Host:       {{.VcenterHost}}
VCenter Username:   {{.VcenterUsername}}
VCenter Password:   {{.VcenterPassword}}
VCenter DataCenter: {{.VcenterDatacenter}}
VCenter Cluster:    {{.VcenterCluster}}
Created:            {{.Created}}
Updated:            {{.Updated}}
`)

		return nil
	},
}

var complexVCenterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New VCenter Complex",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		o := contractorClient.VcenterVCenterComplexNew()
		o.Name = &detailName
		o.Description = &detailDescription
		o.BuiltPercentage = &detailBuiltPercentage
		o.VcenterUsername = &detailUsername
		o.VcenterPassword = &detailPassword
		o.VcenterDatacenter = &detailDatacenter
		o.VcenterCluster = &detailCluster

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			o.Site = cinp.StringAddr(r.GetURI())
		}

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
			o.VcenterHost = cinp.StringAddr(r.GetURI())
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
VCenter Host:       {{.VcenterHost}}
VCenter Username:   {{.VcenterUsername}}
VCenter Password:   {{.VcenterPassword}}
VCenter DataCenter: {{.VcenterDatacenter}}
VCenter Cluster:    {{.VcenterCluster}}
Created:            {{.Created}}
Updated:            {{.Updated}}
`)
		return nil
	},
}

var complexVCenterUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update VCenter Complex",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		complexID := args[0]

		ctx := cmd.Context()

		o := contractorClient.VcenterVCenterComplexNewWithID(complexID)

		if detailDescription != "" {
			o.Description = &detailDescription
		}

		if detailBuiltPercentage != 0 {
			o.BuiltPercentage = &detailBuiltPercentage
		}

		if detailUsername != "" {
			o.VcenterUsername = &detailUsername
		}

		if detailPassword != "" {
			o.VcenterPassword = &detailPassword
		}

		if detailDatacenter != "" {
			o.VcenterDatacenter = &detailDatacenter
		}

		if detailCluster != "" {
			o.VcenterCluster = &detailCluster
		}

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			o.Site = cinp.StringAddr(r.GetURI())
		}

		if len(detailMembers) > 0 {
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
			o.VcenterHost = cinp.StringAddr(r.GetURI())
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
VCenter Host:       {{.VcenterHost}}
VCenter Username:   {{.VcenterUsername}}
VCenter Password:   {{.VcenterPassword}}
VCenter DataCenter: {{.VcenterDatacenter}}
VCenter Cluster:    {{.VcenterCluster}}
Created:            {{.Created}}
Updated:            {{.Updated}}
`)

		return nil
	},
}

func init() {
	complexTypes["vcenter"] = complexTypeEntry{"/api/v1/VCenter/", "0.1"}

	complexVCenterCreateCmd.Flags().StringVarP(&detailName, "name", "l", "", "Locator of New VCenter Complex")
	complexVCenterCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New VCenter Complex")
	complexVCenterCreateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Description of New VCenter Complex")
	complexVCenterCreateCmd.Flags().IntVarP(&detailBuiltPercentage, "builtperc", "b", 80, "Built Percentage of New VCenter Complex\n(Percentage of Built Members at which the complex is considered built)")
	complexVCenterCreateCmd.Flags().IntSliceVarP(&detailMembers, "members", "m", []int{}, "Members of the new VCenter Complex, specify for each member")
	complexVCenterCreateCmd.Flags().StringVarP(&detailUsername, "username", "u", "", "VCenter Username of New VCenter Complex")
	complexVCenterCreateCmd.Flags().StringVarP(&detailPassword, "password", "p", "", "VCenter Password of New VCenter Complex")
	complexVCenterCreateCmd.Flags().StringVarP(&detailDatacenter, "datacenter", "a", "", "VCenter DataCenter of New VCenter Complex")
	complexVCenterCreateCmd.Flags().StringVarP(&detailCluster, "cluster", "c", "", "VCenter Cluster of New VCenter Complex")
	complexVCenterCreateCmd.Flags().IntVarP(&detailHost, "host", "o", 0, "VCenter Host(structure id) of New VCenter Complex")

	complexVCenterUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Complex with value")
	complexVCenterUpdateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Update the Description of VCenter Complex with value")
	complexVCenterUpdateCmd.Flags().IntVarP(&detailBuiltPercentage, "builtperc", "b", 0, "Update the Built Percentage of VCenter Complex with value")
	complexVCenterUpdateCmd.Flags().IntSliceVarP(&detailMembers, "members", "m", []int{}, "Update the Members of the VCenter Complex, specify for each member")
	complexVCenterUpdateCmd.Flags().StringVarP(&detailUsername, "username", "u", "", "Update the VCenter Username of the VCenter Complex with value")
	complexVCenterUpdateCmd.Flags().StringVarP(&detailPassword, "password", "p", "", "Update the VCenter Password of the VCenter Complex with value")
	complexVCenterUpdateCmd.Flags().StringVarP(&detailDatacenter, "datacenter", "a", "", "Update the VCenter DataCenter of New VCenter Complex with value")
	complexVCenterUpdateCmd.Flags().StringVarP(&detailCluster, "cluster", "c", "", "Update the VCenter Cluster of the VCenter Complex with value")
	complexVCenterUpdateCmd.Flags().IntVarP(&detailHost, "host", "o", 0, "Update the VCenter Host(structure id) of New VCenter Complex with value")

	complexCmd.AddCommand(complexVCenterCmd)
	complexVCenterCmd.AddCommand(complexVCenterGetCmd, complexVCenterCreateCmd, complexVCenterUpdateCmd)
}
