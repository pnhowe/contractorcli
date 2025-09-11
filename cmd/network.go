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
	"strconv"

	cinp "github.com/cinp/go"
	"github.com/spf13/cobra"
)

func networkArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a Network Id argument")
	}
	return nil
}

func networkAddressBlockArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a Network AddressBlock Link Id argument")
	}
	return nil
}

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Work with networks",
}

var networkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Networks",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		rl := []cinp.Object{}
		vchan, err := contractorClient.UtilitiesNetworkList(ctx, "", map[string]interface{}{})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Id", "Name", "Site", "Created", "Updated"}, "{{.GetURI | extractID}}	{{.Name}}	{{.Site | extractID}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var networkGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Network",
	Args:  networkArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		networkID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o, err := contractorClient.UtilitiesNetworkGet(ctx, networkID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:            {{.GetURI | extractID}}
Name:          {{.Name}}
Site:          {{.Site | extractID}}
MTU:           {{.Mtu}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)
		rl := []cinp.Object{}
		vchan, err := contractorClient.UtilitiesNetworkAddressBlockList(ctx, "network", map[string]interface{}{"network": o.GetURI()})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		outputList(rl, []string{"link id", "Address Block", "vlan id", "Created", "Update"}, "{{.GetURI | extractID}}	{{.AddressBlock | extractID}}	{{.Vlan}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var networkCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Network",
	RunE: func(cmd *cobra.Command, args []string) error {
		o := contractorClient.UtilitiesNetworkNew()
		o.Name = &detailName
		o.Mtu = &detailMTU

		ctx := cmd.Context()

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			o.Site = cinp.StringAddr(r.GetURI())
		}

		err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI | extractID}}
Name:          {{.Name}}
Site:          {{.Site | extractID}}
MTU:           {{.Mtu}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)

		return nil
	},
}

var networkUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Network",
	Args:  networkArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		networkID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o := contractorClient.UtilitiesNetworkNewWithID(networkID)

		if detailName != "" {
			o.Name = &detailName
		}

		if detailSite != "" {
			r, err := contractorClient.SiteSiteGet(ctx, detailSite)
			if err != nil {
				return err
			}
			o.Site = cinp.StringAddr(r.GetURI())
		}

		if detailMTU != 0 {
			o.Mtu = &detailMTU
		}

		err = o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI | extractID}}
Name:          {{.Name}}
Site:          {{.Site | extractID}}
MTU:           {{.Mtu}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)

		return nil
	},
}

var networkDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Network",
	Args:  networkArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		networkID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o, err := contractorClient.UtilitiesNetworkGet(ctx, networkID)
		if err != nil {
			return err
		}
		if err := o.Delete(ctx); err != nil {
			return err
		}

		return nil
	},
}

var networkAddressBlockCmd = &cobra.Command{
	Use:   "link",
	Short: "Work with Address Block Links",
}

var networkAddressBlockCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Network AddressBlock Link",
	Args:  networkArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		networkID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		r, err := contractorClient.UtilitiesNetworkGet(ctx, networkID)
		if err != nil {
			return err
		}

		o := contractorClient.UtilitiesNetworkAddressBlockNew()
		o.Network = cinp.StringAddr(r.GetURI())
		o.Vlan = &detailVlan

		if detailAddressBlock != 0 {
			r, err := contractorClient.UtilitiesAddressBlockGet(ctx, detailAddressBlock)
			if err != nil {
				return err
			}
			o.AddressBlock = cinp.StringAddr(r.GetURI())
		}

		err = o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI | extractID}}
Id:            {{.ID}}
Network:       {{.Network | extractID}}
Address Block: {{.AddressBlock | extractID}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)

		return nil
	},
}

var networkAddressBlockUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Network AddressBlock Link",
	Args:  networkAddressBlockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		linkID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o := contractorClient.UtilitiesNetworkAddressBlockNewWithID(linkID)

		if detailVlan != -1 {
			o.Vlan = &detailVlan
		}

		err = o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI | extractID}}
Id:            {{.ID}}
Network:       {{.Network | extractID}}
Address Block: {{.AddressBlock | extractID}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)

		return nil
	},
}

var networkAddressBlockDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Network AddressBlock Link",
	Args:  networkAddressBlockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		linkID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o, err := contractorClient.UtilitiesNetworkAddressBlockGet(ctx, linkID)
		if err != nil {
			return err
		}
		if err := o.Delete(ctx); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	networkCreateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Name of New Network")
	networkCreateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Site of New Network")
	networkCreateCmd.Flags().IntVarP(&detailMTU, "mtu", "m", 0, "MTU of New Network")

	networkUpdateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Update the Name of the Network with value")
	networkUpdateCmd.Flags().StringVarP(&detailSite, "site", "s", "", "Update the Site of the Network with value")
	networkUpdateCmd.Flags().IntVarP(&detailMTU, "mtu", "m", 0, "Update the MTU of the Network with the value")

	networkAddressBlockCreateCmd.Flags().IntVarP(&detailAddressBlock, "addressblock", "a", 0, "AddressBlock to Link to")
	networkAddressBlockCreateCmd.Flags().IntVarP(&detailVlan, "vlan", "v", 0, "VLan the Addressblock is tagged as")

	networkAddressBlockUpdateCmd.Flags().IntVarP(&detailVlan, "vlan", "v", -1, "VLan the Addressblock is tagged as")

	rootCmd.AddCommand(networkCmd)
	networkCmd.AddCommand(networkListCmd, networkGetCmd, networkCreateCmd, networkUpdateCmd, networkDeleteCmd)

	networkCmd.AddCommand(networkAddressBlockCmd)
	networkAddressBlockCmd.AddCommand(networkAddressBlockCreateCmd, networkAddressBlockUpdateCmd, networkAddressBlockDeleteCmd)
}
