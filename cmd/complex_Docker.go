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

var complexDockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Work with Docker Complexes",
}

var complexDockerGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Docker Complexes",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		complexID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.DockerDockerComplexGet(complexID)
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

var complexDockerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Docker Complex",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		o := c.DockerDockerComplexNew()
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

var complexDockerUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Docker Complex",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		complexID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.DockerDockerComplexGet(complexID)
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
	complexTypes["docker"] = complexTypeEntry{"/api/v1/Docker/", "0.1"}

	complexDockerCreateCmd.Flags().StringVarP(&detailName, "name", "l", "", "Locator of New Docker Complex")
	complexDockerCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New Docker Complex")
	complexDockerCreateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Description of New Docker Complex")
	complexDockerCreateCmd.Flags().IntVarP(&detailBuiltPercentage, "builtperc", "b", 80, "Built Percentage of New Docker Complex\n(Percentage of Built Members at which the complex is considered built)")
	complexDockerCreateCmd.Flags().StringArrayVarP(&detailMembers, "members", "m", []string{}, "Members of the new Docker Complex, specify for each member")

	complexDockerUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Complex with value")
	complexDockerUpdateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Update the Description of Docker Complex with value")
	complexDockerUpdateCmd.Flags().IntVarP(&detailBuiltPercentage, "builtperc", "b", 0, "Update the Built Percentage of Docker Complex with value")
	complexDockerUpdateCmd.Flags().StringArrayVarP(&detailMembers, "members", "m", []string{}, "Update the Members of the Docker Complex, specify for each member")

	complexCmd.AddCommand(complexDockerCmd)
	complexDockerCmd.AddCommand(complexDockerGetCmd)
	complexDockerCmd.AddCommand(complexDockerCreateCmd)
	complexDockerCmd.AddCommand(complexDockerUpdateCmd)
}
