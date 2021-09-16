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

var detailAuthToken, detailProject, detailFacility string

var complexPacketCmd = &cobra.Command{
	Use:   "packet",
	Short: "Work with Packet/Equinix Metal Complexes",
}

var complexPacketGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Packet Complexes",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		complexID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.PacketPacketComplexGet(complexID)
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
Packet Auth Token:  {{.PacketAuthToken}}
Packet Project:     {{.PacketProject}}
Packet Facility:    {{.PacketFacility}}
Created:            {{.Created}}
Updated:            {{.Updated}}
`)

		return nil
	},
}

var complexPacketCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Packet Complex",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		o := c.PacketPacketComplexNew()
		o.Name = detailName
		o.Description = detailDescription
		o.PacketAuthToken = detailAuthToken
		o.PacketProject = detailProject
		o.PacketFacility = detailFacility

		if detailSite != "" {
			r, err := c.SiteSiteGet(detailSite)
			if err != nil {
				return err
			}
			o.Site = r.GetID()
		}

		if err := o.Create(); err != nil {
			return err
		}

		outputKV(map[string]interface{}{"id": o.GetID()})

		return nil
	},
}

var complexPacketUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Packet Complex",
	Args:  complexArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		complexID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.PacketPacketComplexGet(complexID)
		if err != nil {
			return err
		}

		if detailDescription != "" {
			o.Description = detailDescription
			fieldList = append(fieldList, "description")
		}

		if detailAuthToken != "" {
			o.PacketAuthToken = detailAuthToken
			fieldList = append(fieldList, "packet_auth_token")
		}

		if detailProject != "" {
			o.PacketProject = detailProject
			fieldList = append(fieldList, "packet_project")
		}

		if detailFacility != "" {
			o.PacketFacility = detailFacility
			fieldList = append(fieldList, "packet_facility")
		}

		if detailSite != "" {
			r, err := c.SiteSiteGet(detailSite)
			if err != nil {
				return err
			}
			o.Site = r.GetID()
			fieldList = append(fieldList, "site")
		}

		if err := o.Update(fieldList); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	complexTypes["packet"] = complexTypeEntry{"/api/v1/Packet/", "0.1"}

	complexPacketCreateCmd.Flags().StringVarP(&detailName, "name", "l", "", "Locator of New Packet Complex")
	complexPacketCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New Packet Complex")
	complexPacketCreateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Description of New Packet Complex")
	complexPacketCreateCmd.Flags().StringVarP(&detailAuthToken, "token", "t", "", "Packet Auth Token of New Packet Complex")
	complexPacketCreateCmd.Flags().StringVarP(&detailProject, "project", "p", "", "Packet Project of New Packet Complex")
	complexPacketCreateCmd.Flags().StringVarP(&detailFacility, "facility", "f", "", "Packet Facility of New Packet Complex")

	complexPacketUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Complex with value")
	complexPacketUpdateCmd.Flags().StringVarP(&detailDescription, "description", "d", "", "Update the Description of Packet Complex with value")
	complexPacketUpdateCmd.Flags().StringVarP(&detailAuthToken, "token", "t", "", "Update the Packet Auth Token of the Packet Complex with value")
	complexPacketUpdateCmd.Flags().StringVarP(&detailProject, "project", "p", "", "Update the Packet Project of New Packet Complex with value")
	complexPacketUpdateCmd.Flags().StringVarP(&detailFacility, "facility", "f", "", "Update the Packet Facility of the Packet Complex with value")

	complexCmd.AddCommand(complexPacketCmd)
	complexPacketCmd.AddCommand(complexPacketGetCmd)
	complexPacketCmd.AddCommand(complexPacketCreateCmd)
	complexPacketCmd.AddCommand(complexPacketUpdateCmd)
}
