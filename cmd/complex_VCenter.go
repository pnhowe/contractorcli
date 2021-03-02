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
	"strconv"
)

var detailHost, detailDatacenter, detailCluster string

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
		c := getContractor()
		defer c.Logout()

		r, err := c.VcenterVCenterComplexGet(complexID)
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
		c := getContractor()
		defer c.Logout()

		o := c.VcenterVCenterComplexNew()
		o.Name = detailName
		o.Description = detailDescription
		o.BuiltPercentage = detailBuiltPercentage
		o.VcenterUsername = detailUsername
		o.VcenterPassword = detailPassword
		o.VcenterDatacenter = detailDatacenter
		o.VcenterCluster = detailCluster

		if detailSite != "" {
			r, err := c.SiteSiteGet(detailSite)
			if err != nil {
				return err
			}
			o.Site = r.GetID()
		}

		for _, v := range detailMembers {
			s, err := c.BuildingStructureGet(strconv.Itoa(v))
			if err != nil {
				return err
			}
			o.Members = append(o.Members, s.GetID())
		}

		if detailHost != "" {
			r, err := c.BuildingStructureGet(detailHost)
			if err != nil {
				return err
			}
			o.VcenterHost = r.GetID()
		}

		if err := o.Create(); err != nil {
			return err
		}

		outputKV(map[string]interface{}{"id": o.GetID()})

		return nil
	},
}

var complexVCenterUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update VCenter Complex",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		complexID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.VcenterVCenterComplexGet(complexID)
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
			o.VcenterUsername = detailUsername
			fieldList = append(fieldList, "vcenter_username")
		}

		if detailPassword != "" {
			o.VcenterPassword = detailPassword
			fieldList = append(fieldList, "vcenter_password")
		}

		if detailDatacenter != "" {
			o.VcenterDatacenter = detailDatacenter
			fieldList = append(fieldList, "vcenter_datacenter")
		}

		if detailCluster != "" {
			o.VcenterCluster = detailCluster
			fieldList = append(fieldList, "vcenter_cluster")
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
				s, err := c.BuildingStructureGet(strconv.Itoa(v))
				if err != nil {
					return err
				}
				o.Members = append(o.Members, s.GetID())
			}
			fieldList = append(fieldList, "members")
		}

		if detailHost != "" {
			r, err := c.BuildingStructureGet(detailHost)
			if err != nil {
				return err
			}
			o.VcenterHost = r.GetID()
			fieldList = append(fieldList, "vcenter_host")
		}

		if err := o.Update(fieldList); err != nil {
			return err
		}

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
	complexVCenterCreateCmd.Flags().StringVarP(&detailHost, "host", "o", "", "VCenter Host(structure id) of New VCenter Complex")

	complexVCenterUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Complex with value")
	complexVCenterUpdateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Update the Description of VCenter Complex with value")
	complexVCenterUpdateCmd.Flags().IntVarP(&detailBuiltPercentage, "builtperc", "b", 0, "Update the Built Percentage of VCenter Complex with value")
	complexVCenterUpdateCmd.Flags().IntSliceVarP(&detailMembers, "members", "m", []int{}, "Update the Members of the VCenter Complex, specify for each member")
	complexVCenterUpdateCmd.Flags().StringVarP(&detailUsername, "username", "u", "", "Update the VCenter Username of the VCenter Complex with value")
	complexVCenterUpdateCmd.Flags().StringVarP(&detailPassword, "password", "p", "", "Update the VCenter Password of the VCenter Complex with value")
	complexVCenterUpdateCmd.Flags().StringVarP(&detailDatacenter, "datacenter", "a", "", "Update the VCenter DataCenter of New VCenter Complex with value")
	complexVCenterUpdateCmd.Flags().StringVarP(&detailCluster, "cluster", "c", "", "Update the VCenter Cluster of the VCenter Complex with value")
	complexVCenterUpdateCmd.Flags().StringVarP(&detailHost, "host", "o", "", "Update the VCenter Host(structure id) of New VCenter Complex with value")

	complexCmd.AddCommand(complexVCenterCmd)
	complexVCenterCmd.AddCommand(complexVCenterGetCmd)
	complexVCenterCmd.AddCommand(complexVCenterCreateCmd)
	complexVCenterCmd.AddCommand(complexVCenterUpdateCmd)
}
