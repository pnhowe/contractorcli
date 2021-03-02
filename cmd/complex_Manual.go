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

var complexManualCmd = &cobra.Command{
	Use:   "manual",
	Short: "Work with Manual Complexes",
}

var complexManualGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Manual Complexes",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		complexID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.ManualManualComplexGet(complexID)
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
Created:            {{.Created}}
Updated:            {{.Updated}}
`)

		return nil
	},
}

var complexManualCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Manual Complex",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		o := c.ManualManualComplexNew()
		o.Name = detailName
		o.Description = detailDescription
		o.BuiltPercentage = detailBuiltPercentage

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

		if err := o.Create(); err != nil {
			return err
		}

		outputKV(map[string]interface{}{"id": o.GetID()})

		return nil
	},
}

var complexManualUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Manual Complex",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		complexID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.ManualManualComplexGet(complexID)
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

		if err := o.Update(fieldList); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	complexTypes["manual"] = complexTypeEntry{"/api/v1/Manual/", "0.1"}

	complexManualCreateCmd.Flags().StringVarP(&detailName, "name", "l", "", "Locator of New Manual Complex")
	complexManualCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New Manual Complex")
	complexManualCreateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Description of New Manual Complex")
	complexManualCreateCmd.Flags().IntVarP(&detailBuiltPercentage, "builtperc", "b", 80, "Built Percentage of New Manual Complex\n(Percentage of Built Members at which the complex is considered built)")
	complexManualCreateCmd.Flags().IntSliceVarP(&detailMembers, "members", "m", []int{}, "Members of the new Manual Complex, specify for each member")

	complexManualUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Complex with value")
	complexManualUpdateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Update the Description of Manual Complex with value")
	complexManualUpdateCmd.Flags().IntVarP(&detailBuiltPercentage, "builtperc", "b", 0, "Update the Built Percentage of Manual Complex with value")
	complexManualUpdateCmd.Flags().IntSliceVarP(&detailMembers, "members", "m", []int{}, "Update the Members of the Manual Complex, specify for each member")

	complexCmd.AddCommand(complexManualCmd)
	complexManualCmd.AddCommand(complexManualGetCmd)
	complexManualCmd.AddCommand(complexManualCreateCmd)
	complexManualCmd.AddCommand(complexManualUpdateCmd)
}
