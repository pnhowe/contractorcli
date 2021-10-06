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
	"sort"

	cinp "github.com/cinp/go"
	"github.com/spf13/cobra"
)

var detailSubnet, detailReason, detailPXE string
var detailPrefix, detailGatewayOffset, detailOffset int

func addressblockArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Requires a AddressBlock Id Argument")
	}
	return nil
}

var addressblockCmd = &cobra.Command{
	Use:   "addressblock",
	Short: "Work with AddressBlocks",
}

var addressblockListCmd = &cobra.Command{
	Use:   "list",
	Short: "List AddressBlocks",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		rl := []cinp.Object{}
		for v := range c.UtilitiesAddressBlockList("", map[string]interface{}{}) {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Id", "Name", "Site", "SubNet", "Prefix", "Created", "Updated"}, "{{.GetID | extractID}}	{{.Name}}	{{.Site | extractID}}	{{.Subnet}}	{{.Prefix}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var addressblockGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get AddressBlock",
	Args:  addressblockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		addressblockID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.UtilitiesAddressBlockGet(addressblockID)
		if err != nil {
			return err
		}
		outputDetail(r, `Name:          {{.Name}}
Site:          {{.Site | extractID}}
Subnet:        {{.Subnet}}
Prefix:        {{.Prefix}} ({{.Netmask}})
Gateway Offset:{{.GatewayOffset}} ({{.Gateway}})
Max Addresse:  {{.MaxAddress}}
Size:          {{.Size}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)
		return nil
	},
}

var addressblockCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New AddressBlock",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		o := c.UtilitiesAddressBlockNew()
		o.Name = detailName
		o.Subnet = detailSubnet
		o.Prefix = detailPrefix
		if detailGatewayOffset != 0 {
			o.GatewayOffset = detailGatewayOffset
		}

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

var addressblockUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update AddressBlock",
	Args:  addressblockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		addressblockID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.UtilitiesAddressBlockGet(addressblockID)
		if err != nil {
			return err
		}

		if detailName != "" {
			o.Name = detailName
			fieldList = append(fieldList, "name")
		}

		if detailSubnet != "" {
			o.Subnet = detailSubnet
			fieldList = append(fieldList, "subnet")
		}

		if detailPrefix != 0 {
			o.Prefix = detailPrefix
			fieldList = append(fieldList, "prefix")
		}

		if detailGatewayOffset != 0 {
			o.GatewayOffset = detailGatewayOffset
			fieldList = append(fieldList, "gateway_offset")
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

var addressblockDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete AddressBlock",
	Args:  addressblockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		addressblockID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.UtilitiesAddressBlockGet(addressblockID)
		if err != nil {
			return err
		}
		if err := r.Delete(); err != nil {
			return err
		}

		return nil
	},
}

var addressblockUsageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Display the usage for an AddressBlock",
	Args:  addressblockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		addressblockID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.UtilitiesAddressBlockGet(addressblockID)
		if err != nil {
			return err
		}

		u, err := r.CallUsage()
		if err != nil {
			return err
		}

		outputDetail(u, `Total:         {{.total}}
Static:        {{.static}}
Reserved:      {{.reserved}}
Dynammic:      {{.dynamic}}
`)

		return nil
	},
}

var addressblockAllocationCmd = &cobra.Command{
	Use:   "allocation",
	Short: "Display the usage for an AddressBlock",
	Args:  addressblockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		addressblockID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.UtilitiesAddressBlockGet(addressblockID)
		if err != nil {
			return err
		}

		rl := []cinp.Object{}
		for v := range c.UtilitiesAddressList("address_block", map[string]interface{}{"address_block": r.GetID()}) {
			rl = append(rl, v)
		}
		for v := range c.UtilitiesReservedAddressList("address_block", map[string]interface{}{"address_block": r.GetID()}) {
			rl = append(rl, v)
		}
		for v := range c.UtilitiesDynamicAddressList("address_block", map[string]interface{}{"address_block": r.GetID()}) {
			rl = append(rl, v)
		}

		sort.Slice(rl, func(i, j int) bool {
			return (*rl[i].AsMap(false))["offset"].(int) < (*rl[j].AsMap(false))["offset"].(int)
		})

		outputList(rl, []string{"Id", "Offset", "Ip Address", "Type"}, "{{.GetID | extractID}}	{{.Offset}}	{{.IPAddress}}	{{.Type}}\n")

		return nil
	},
}

var addressblockReserveCmd = &cobra.Command{
	Use:   "reserve",
	Short: "Reserve and ip/offset in an Address Block",
	Args:  addressblockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		addressblockID := args[0]
		c := getContractor()
		defer c.Logout()

		if detailOffset == 0 {
			return fmt.Errorf("Offset required")
		}

		if detailReason == "" {
			return fmt.Errorf("Reason Required")
		}

		r, err := c.UtilitiesAddressBlockGet(addressblockID)
		if err != nil {
			return err
		}

		o := c.UtilitiesReservedAddressNew()
		o.AddressBlock = r.GetID()
		o.Offset = detailOffset
		o.Reason = detailReason

		if err := o.Create(); err != nil {
			return err
		}

		outputKV(map[string]interface{}{"id": o.GetID()})

		return nil
	},
}

var addressblockDeReserveCmd = &cobra.Command{
	Use:   "dereserve",
	Short: "Dereserve a Reserved Ip From an Address Block",
	Args:  addressblockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		addressblockID := args[0]
		c := getContractor()
		defer c.Logout()

		if detailOffset == 0 {
			return fmt.Errorf("Offset required")
		}

		r, err := c.UtilitiesAddressBlockGet(addressblockID)
		if err != nil {
			return err
		}
		for v := range c.UtilitiesReservedAddressList("address_block", map[string]interface{}{"address_block": r.GetID()}) {
			if v.Offset == detailOffset {
				v.Delete()
				return nil
			}
		}

		return fmt.Errorf("Offset not found")
	},
}

var addressblockDynamicCmd = &cobra.Command{
	Use:   "dynamic",
	Short: "Assign an ip/offset in an Address Block as Dynamic",
	Args:  addressblockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		addressblockID := args[0]
		c := getContractor()
		defer c.Logout()

		if detailOffset == 0 {
			return fmt.Errorf("Offset required")
		}

		r, err := c.UtilitiesAddressBlockGet(addressblockID)
		if err != nil {
			return err
		}

		o := c.UtilitiesDynamicAddressNew()
		o.AddressBlock = r.GetID()
		o.Offset = detailOffset
		if detailPXE != "" {
			r2, err := c.BlueprintPXEGet(detailPXE)
			if err != nil {
				return err
			}
			o.Pxe = r2.GetID()
		}

		if err := o.Create(); err != nil {
			return err
		}

		outputKV(map[string]interface{}{"id": o.GetID()})

		return nil
	},
}

var addressblockDeDynamicCmd = &cobra.Command{
	Use:   "dedynamic",
	Short: "De-assign an ip/offset in an Address Block as Dynamic",
	Args:  addressblockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		addressblockID := args[0]
		c := getContractor()
		defer c.Logout()

		if detailOffset == 0 {
			return fmt.Errorf("Offset required")
		}

		r, err := c.UtilitiesAddressBlockGet(addressblockID)
		if err != nil {
			return err
		}
		for v := range c.UtilitiesDynamicAddressList("address_block", map[string]interface{}{"address_block": r.GetID()}) {
			if v.Offset == detailOffset {
				v.Delete()
				return nil
			}
		}

		return fmt.Errorf("Offset not found")
	},
}

func init() {
	addressblockCreateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Name of the New AddressBlock")
	addressblockCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of the New AddressBlock")
	addressblockCreateCmd.Flags().StringVarP(&detailSubnet, "subnet", "u", "", "Subnet of the New AddressBlock")
	addressblockCreateCmd.Flags().IntVarP(&detailPrefix, "prefix", "p", 0, "Prefix of the New AddressBlock")
	addressblockCreateCmd.Flags().IntVarP(&detailGatewayOffset, "gateway", "g", 0, "Gateway Offset of the New AddressBlock")

	addressblockUpdateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Update the Name of the AddressBlock")
	addressblockUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of the AddressBlock")
	addressblockUpdateCmd.Flags().StringVarP(&detailSubnet, "subnet", "u", "", "Update the Subnet of the AddressBlock")
	addressblockUpdateCmd.Flags().IntVarP(&detailPrefix, "prefix", "p", 0, "Update the Prefix of the AddressBlock")
	addressblockUpdateCmd.Flags().IntVarP(&detailGatewayOffset, "gateway", "g", 0, "Update the Gateway Offset of the AddressBlock")

	addressblockReserveCmd.Flags().IntVarP(&detailOffset, "offset", "o", 0, "Offset for the New Reservation")
	addressblockReserveCmd.Flags().StringVarP(&detailReason, "reason", "r", "", "Reason for the New Reservation")

	addressblockDeReserveCmd.Flags().IntVarP(&detailOffset, "offset", "o", 0, "Offset of the Reservation to Remove")

	addressblockDynamicCmd.Flags().IntVarP(&detailOffset, "offset", "o", 0, "Offset for the New Dynamic Ip")
	addressblockDynamicCmd.Flags().StringVarP(&detailPXE, "pxe", "p", "", "PXE for the New Dynamic Ip")

	addressblockDeDynamicCmd.Flags().IntVarP(&detailOffset, "offset", "o", 0, "Offset of the Dynamic Ip to Remove")

	rootCmd.AddCommand(addressblockCmd)
	addressblockCmd.AddCommand(addressblockListCmd)
	addressblockCmd.AddCommand(addressblockGetCmd)
	addressblockCmd.AddCommand(addressblockCreateCmd)
	addressblockCmd.AddCommand(addressblockUpdateCmd)
	addressblockCmd.AddCommand(addressblockDeleteCmd)
	addressblockCmd.AddCommand(addressblockUsageCmd)
	addressblockCmd.AddCommand(addressblockAllocationCmd)
	addressblockCmd.AddCommand(addressblockReserveCmd)
	addressblockCmd.AddCommand(addressblockDeReserveCmd)
	addressblockCmd.AddCommand(addressblockDynamicCmd)
	addressblockCmd.AddCommand(addressblockDeDynamicCmd)
}
