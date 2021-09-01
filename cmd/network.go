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

	cinp "github.com/cinp/go"
	"github.com/spf13/cobra"
)

var detailAddressBlock string
var detailVlan, detailMTU int

func networkArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Requires a Network Id Argument")
	}
	return nil
}

func networkAddressBlockArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Requires a Network AddressBlock Link Id Argument")
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
		c := getContractor()
		defer c.Logout()

		rl := []cinp.Object{}
		for v := range c.UtilitiesNetworkList("", map[string]interface{}{}) {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Id", "Name", "Site", "Created", "Updated"}, "{{.GetID | extractID}}	{{.Name}}	{{.Site | extractID}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var networkGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Network",
	Args:  networkArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		networkID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.UtilitiesNetworkGet(networkID)
		if err != nil {
			return err
		}
		outputDetail(r, `Name:          {{.Name}}
Site:          {{.Site | extractID}}
MTU:           {{.Mtu}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)
		rl := []cinp.Object{}
		for v := range c.UtilitiesNetworkAddressBlockList("network", map[string]interface{}{"network": r.GetID()}) {
			rl = append(rl, v)
		}
		outputList(rl, []string{"link id", "Address Block", "vlan id", "Created", "Update"}, "{{.GetID | extractID}}	{{.AddressBlock | extractID}}	{{.Vlan}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var networkCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Network",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getContractor()
		defer c.Logout()

		o := c.UtilitiesNetworkNew()
		o.Name = detailName
		o.Mtu = detailMTU

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

var networkUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Network",
	Args:  networkArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		networkID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.UtilitiesNetworkGet(networkID)
		if err != nil {
			return err
		}

		if detailName != "" {
			o.Name = detailName
			fieldList = append(fieldList, "name")
		}

		if detailSite != "" {
			r, err := c.SiteSiteGet(detailSite)
			if err != nil {
				return err
			}
			o.Site = r.GetID()
			fieldList = append(fieldList, "site")
		}

		if detailMTU != 0 {
			o.Mtu = detailMTU
			fieldList = append(fieldList, "mtu")
		}

		if err := o.Update(fieldList); err != nil {
			return err
		}

		return nil
	},
}

var networkDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Network",
	Args:  networkArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		networkID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.UtilitiesNetworkGet(networkID)
		if err != nil {
			return err
		}
		if err := r.Delete(); err != nil {
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
		networkID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.UtilitiesNetworkGet(networkID)
		if err != nil {
			return err
		}

		o := c.UtilitiesNetworkAddressBlockNew()
		o.Network = r.GetID()
		o.Vlan = detailVlan

		if detailAddressBlock != "" {
			r, err := c.UtilitiesAddressBlockGet(detailAddressBlock)
			if err != nil {
				return err
			}
			o.AddressBlock = r.GetID()
		}

		if err := o.Create(); err != nil {
			return err
		}

		outputKV(map[string]interface{}{"id": o.GetID()})

		return nil
	},
}

var networkAddressBlockUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Network AddressBlock Link",
	Args:  networkAddressBlockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldList := []string{}
		linkID := args[0]
		c := getContractor()
		defer c.Logout()

		o, err := c.UtilitiesNetworkAddressBlockGet(linkID)
		if err != nil {
			return err
		}

		if detailVlan != -1 {
			o.Vlan = detailVlan
			fieldList = append(fieldList, "vlan")
		}

		if err := o.Update(fieldList); err != nil {
			return err
		}

		return nil
	},
}

var networkAddressBlockDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Network AddressBlock Link",
	Args:  networkAddressBlockArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		linkID := args[0]
		c := getContractor()
		defer c.Logout()

		r, err := c.UtilitiesNetworkAddressBlockGet(linkID)
		if err != nil {
			return err
		}
		if err := r.Delete(); err != nil {
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

	networkAddressBlockCreateCmd.Flags().StringVarP(&detailAddressBlock, "addressblock", "a", "", "AddressBlock to Link to")
	networkAddressBlockCreateCmd.Flags().IntVarP(&detailVlan, "vlan", "v", 0, "VLan the Addressblock is tagged as")

	networkAddressBlockUpdateCmd.Flags().IntVarP(&detailVlan, "vlan", "v", -1, "VLan the Addressblock is tagged as")

	rootCmd.AddCommand(networkCmd)
	networkCmd.AddCommand(networkListCmd)
	networkCmd.AddCommand(networkGetCmd)
	networkCmd.AddCommand(networkCreateCmd)
	networkCmd.AddCommand(networkUpdateCmd)
	networkCmd.AddCommand(networkDeleteCmd)

	networkCmd.AddCommand(networkAddressBlockCmd)
	networkAddressBlockCmd.AddCommand(networkAddressBlockCreateCmd)
	networkAddressBlockCmd.AddCommand(networkAddressBlockUpdateCmd)
	networkAddressBlockCmd.AddCommand(networkAddressBlockDeleteCmd)
}
