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
	"strconv"

	contractor "github.com/t3kton/contractor_client/go/autogen"

	cinp "github.com/cinp/go"
	"github.com/spf13/cobra"
)

func addressblockArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires an AddressBlock Id argument")
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
		ctx := cmd.Context()

		rl := []cinp.Object{}
		vchan, err := contractorClient.UtilitiesAddressBlockList(ctx, "", map[string]interface{}{})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Id", "Name", "Site", "SubNet", "Prefix", "Created", "Updated"}, "{{.GetURI | extractID}}	{{.Name}}	{{.Site | extractID}}	{{.Subnet}}	{{.Prefix}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var addressblockGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get AddressBlock",
	Args:  addressblockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		addressblockID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o, err := contractorClient.UtilitiesAddressBlockGet(ctx, addressblockID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:            {{.GetURI | extractID}}
Name:          {{.Name}}
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
		o := contractorClient.UtilitiesAddressBlockNew()
		o.Name = &detailName
		o.Subnet = &detailSubnet
		o.Prefix = &detailPrefix
		if detailGatewayOffset != 0 {
			o.GatewayOffset = &detailGatewayOffset
		}

		ctx := cmd.Context()

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			(*o.Site) = r.GetURI()
		}

		o, err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI  | extractID}}
Name:          {{.Name}}
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

var addressblockUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update AddressBlock",
	Args:  addressblockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		addressblockID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		ctx := cmd.Context()

		o := contractorClient.UtilitiesAddressBlockNewWithID(addressblockID)

		if detailName != "" {
			o.Name = &detailName
		}

		if detailSubnet != "" {
			o.Subnet = &detailSubnet
		}

		if detailPrefix != 0 {
			o.Prefix = &detailPrefix
		}

		if detailGatewayOffset != 0 {
			o.GatewayOffset = &detailGatewayOffset
		}

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			(*o.Site) = r.GetURI()
		}

		o, err = o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:               {{.GetURI | extractID}}
Name:             {{.Name}}
Site:             {{.Site | extractID}}
Subnet:           {{.Subnet}}
Prefix:           {{.Prefix}} ({{.Netmask}})
Gateway Offset:   {{.GatewayOffset}} ({{.Gateway}})
Max Addresse:     {{.MaxAddress}}
Size:             {{.Size}}
Created:          {{.Created}}
Updated:          {{.Updated}}
`)
		return nil
	},
}

var addressblockDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete AddressBlock",
	Args:  addressblockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		addressblockID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		ctx := cmd.Context()

		o, err := contractorClient.UtilitiesAddressBlockGet(ctx, addressblockID)
		if err != nil {
			return err
		}
		if err := o.Delete(ctx); err != nil {
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
		addressblockID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		ctx := cmd.Context()

		o, err := contractorClient.UtilitiesAddressBlockGet(ctx, addressblockID)
		if err != nil {
			return err
		}

		u, err := o.CallUsage(ctx)
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
		addressblockID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		ctx := cmd.Context()

		o, err := contractorClient.UtilitiesAddressBlockGet(ctx, addressblockID)
		if err != nil {
			return err
		}

		rl := []cinp.Object{}
		vchan, err := contractorClient.UtilitiesAddressList(ctx, "address_block", map[string]interface{}{"address_block": o.GetURI()})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		vchan2, err := contractorClient.UtilitiesReservedAddressList(ctx, "address_block", map[string]interface{}{"address_block": o.GetURI()})
		if err != nil {
			return err
		}
		for v := range vchan2 {
			rl = append(rl, v)
		}
		vchan3, err := contractorClient.UtilitiesDynamicAddressList(ctx, "address_block", map[string]interface{}{"address_block": o.GetURI()})
		if err != nil {
			return err
		}
		for v := range vchan3 {
			rl = append(rl, v)
		}

		sort.Slice(rl, func(i, j int) bool {
			var offsetI, offsetJ int

			switch a := rl[i].(type) {
			case *contractor.UtilitiesDynamicAddress:
				offsetI = *a.Offset
			case *contractor.UtilitiesReservedAddress:
				offsetI = *a.Offset
			case *contractor.UtilitiesAddress:
				offsetI = *a.Offset
			default:
				offsetI = 0
			}

			switch a := rl[j].(type) {
			case *contractor.UtilitiesDynamicAddress:
				offsetJ = *a.Offset
			case *contractor.UtilitiesReservedAddress:
				offsetJ = *a.Offset
			case *contractor.UtilitiesAddress:
				offsetJ = *a.Offset
			default:
				offsetJ = 0
			}
			return offsetI < offsetJ
		})

		outputList(rl, []string{"Id", "Offset", "Ip Address", "Type"}, "{{.GetURI | extractID}}	{{.Offset}}	{{.IPAddress}}	{{.Type}}\n")

		return nil
	},
}

var addressblockReserveCmd = &cobra.Command{
	Use:   "reserve",
	Short: "Reserve and ip/offset in an Address Block",
	Args:  addressblockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		addressblockID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		if detailOffset == 0 {
			return fmt.Errorf("offset required")
		}

		if detailReason == "" {
			return fmt.Errorf("reason Required")
		}

		ctx := cmd.Context()

		r, err := contractorClient.UtilitiesAddressBlockGet(ctx, addressblockID)
		if err != nil {
			return err
		}

		o := contractorClient.UtilitiesReservedAddressNew()
		(*o.AddressBlock) = r.GetURI()
		o.Offset = &detailOffset
		o.Reason = &detailReason

		o, err = o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:             {{.GetURI | extractID}}
AddressBlock:   {{.AddressBlock | extractID}}
Offset:         {{.Offset}}
Reason:         {{.Reason}}
Updated:        {{.Updated}}
Created:        {{.Created}}
`)

		return nil
	},
}

var addressblockDeReserveCmd = &cobra.Command{
	Use:   "dereserve",
	Short: "Dereserve a Reserved Ip From an Address Block",
	Args:  addressblockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		addressblockID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		if detailOffset == 0 {
			return fmt.Errorf("offset required")
		}

		ctx := cmd.Context()

		o, err := contractorClient.UtilitiesAddressBlockGet(ctx, addressblockID)
		if err != nil {
			return err
		}
		vchan, err := contractorClient.UtilitiesReservedAddressList(ctx, "address_block", map[string]interface{}{"address_block": o.GetURI()})
		if err != nil {
			return err
		}
		for v := range vchan {
			if *v.Offset == detailOffset {
				v.Delete(ctx)
				return nil
			}
		}

		return fmt.Errorf("offset not found")
	},
}

var addressblockDynamicCmd = &cobra.Command{
	Use:   "dynamic",
	Short: "Assign an ip/offset in an Address Block as Dynamic",
	Args:  addressblockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		addressblockID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		if detailOffset == 0 {
			return fmt.Errorf("offset required")
		}

		ctx := cmd.Context()

		r, err := contractorClient.UtilitiesAddressBlockGet(ctx, addressblockID)
		if err != nil {
			return err
		}

		o := contractorClient.UtilitiesDynamicAddressNew()
		(*o.AddressBlock) = r.GetURI()
		o.Offset = &detailOffset
		if detailPXE != "" {
			r2, err := contractorClient.BlueprintPXEGet(ctx, detailPXE)
			if err != nil {
				return err
			}
			(*o.Pxe) = r2.GetURI()
		}

		o, err = o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:           {{.GetURI | extractID}}
AddressBlock: {{.AddressBlock | extractID}}
Offset:       {{.Offset}}
Updated:      {{.Updated}}
Created:      {{.Created}}
`)

		return nil
	},
}

var addressblockDeDynamicCmd = &cobra.Command{
	Use:   "dedynamic",
	Short: "De-assign an ip/offset in an Address Block as Dynamic",
	Args:  addressblockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		addressblockID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		if detailOffset == 0 {
			return fmt.Errorf("offset required")
		}

		ctx := cmd.Context()

		o, err := contractorClient.UtilitiesAddressBlockGet(ctx, addressblockID)
		if err != nil {
			return err
		}
		vchan, err := contractorClient.UtilitiesDynamicAddressList(ctx, "address_block", map[string]interface{}{"address_block": o.GetURI()})
		if err != nil {
			return err
		}
		for v := range vchan {
			if *v.Offset == detailOffset {
				v.Delete(ctx)
				return nil
			}
		}

		return fmt.Errorf("offset not found")
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
	addressblockCmd.AddCommand(addressblockListCmd, addressblockGetCmd, addressblockCreateCmd, addressblockUpdateCmd, addressblockDeleteCmd, addressblockUsageCmd, addressblockAllocationCmd, addressblockReserveCmd, addressblockDeReserveCmd, addressblockDynamicCmd, addressblockDeDynamicCmd)
}
