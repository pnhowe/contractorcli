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

var foundationPacketCmd = &cobra.Command{
	Use:   "packet",
	Short: "Work with Packet/Equinix Metal Foundations",
}

var foundationPacketGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Packet Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		foundationID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.PacketPacketFoundationGet(foundationID)
		if err != nil {
			return err
		}
		outputDetail(r, `Locator:        {{.Locator}}
Complex:        {{.PacketComplex | extractID}}
Server UUID:    {{.PacketUUID}}
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

var foundationPacketCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Packet Foundation",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		o := c.PacketPacketFoundationNew()
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
			r, err := c.PacketPacketComplexGet(detailComplex)
			if err != nil {
				return err
			}
			o.PacketComplex = r.GetID()
		}

		if err := o.Create(); err != nil {
			return err
		}

		outputKV(map[string]interface{}{"id": o.GetID()})

		return nil
	},
}

var foundationPacketUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Packet Foundation",
	Args:  foundationArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		foundationID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.PacketPacketFoundationGet(foundationID)
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
			r, err := c.PacketPacketComplexGet(detailComplex)
			if err != nil {
				return err
			}
			o.PacketComplex = r.GetID()
			fieldList = append(fieldList, "packet_complex")
		}

		if err := o.Update(fieldList); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	fundationTypes["packet"] = foundationTypeEntry{"/api/v1/Packet/", "0.1"}

	foundationPacketCreateCmd.Flags().StringVarP(&detailLocator, "locator", "l", "", "Locator of New Packet Foundation")
	foundationPacketCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New Packet Foundation")
	foundationPacketCreateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Blueprint of New Packet Foundation")
	foundationPacketCreateCmd.Flags().StringVarP(&detailComplex, "complex", "x", "", "Plot of New Packet Foundation")

	foundationPacketUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of Foundation with value")
	foundationPacketUpdateCmd.Flags().StringVarP(&detailBlueprint, "blueprint", "b", "", "Update the Blueprint of Foundation with value")
	foundationPacketUpdateCmd.Flags().StringVarP(&detailComplex, "complex", "x", "", "Update the Plot of the Packet Foundation")

	foundationCmd.AddCommand(foundationPacketCmd)
	foundationPacketCmd.AddCommand(foundationPacketGetCmd)
	foundationPacketCmd.AddCommand(foundationPacketCreateCmd)
	foundationPacketCmd.AddCommand(foundationPacketUpdateCmd)
}
