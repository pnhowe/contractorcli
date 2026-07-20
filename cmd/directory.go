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

	cinp "github.com/cinp/go/v2"
	"github.com/spf13/cobra"
)

func directoryArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a Directory Zone Id argument")
	}
	return nil
}

func directoryEntryArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a Directory Entry Id argument")
	}
	return nil
}

func directoryEntryZoneArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return errors.New("takes at most one Directory Zone Id argument, omit it to work with global Entries")
	}
	return nil
}

var directoryCmd = &cobra.Command{
	Use:   "directory",
	Short: "Work with Directory Zones",
}

var directoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Directory Zones",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		rl := []cinp.Object{}
		vchan, err := contractorClient.DirectoryZoneList(ctx, "", map[string]interface{}{})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Id", "Name", "Fqdn", "Parent", "Created", "Updated"}, "{{.GetURI | extractID}}	{{.Name}}	{{.Fqdn}}	{{or .Parent \":<None>\" | extractID}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var directoryGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Directory Zone",
	Args:  directoryArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		zoneID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o, err := contractorClient.DirectoryZoneGet(ctx, zoneID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:            {{.GetURI | extractID}}
Name:          {{.Name}}
Fqdn:          {{.Fqdn}}
Parent:        {{or .Parent ":<None>" | extractID}}
TTL:           {{.TTL}}
Refresh:       {{.Refresh}}
Retry:         {{.Retry}}
Expire:        {{.Expire}}
Minimum:       {{.Minimum}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)

		rl := []cinp.Object{}
		vchan, err := contractorClient.DirectoryEntryList(ctx, "zone", map[string]interface{}{"zone": o.GetURI()})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Id", "Type", "Name", "Priority", "Weight", "Port", "Target"}, "{{.GetURI | extractID}}	{{.Type}}	{{.Name}}	{{.Priority}}	{{.Weight}}	{{.Port}}	{{.Target}}\n")

		return nil
	},
}

var directoryCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Directory Zone",
	RunE: func(cmd *cobra.Command, args []string) error {
		o := contractorClient.DirectoryZoneNew()

		if detailName != "" {
			o.Name = &detailName
		}

		if detailTTL != 0 {
			o.TTL = &detailTTL
		}

		if detailRefresh != 0 {
			o.Refresh = &detailRefresh
		}

		if detailRetry != 0 {
			o.Retry = &detailRetry
		}

		if detailExpire != 0 {
			o.Expire = &detailExpire
		}

		if detailMinimum != 0 {
			o.Minimum = &detailMinimum
		}

		ctx := cmd.Context()

		if detailZone != 0 {
			r, err := contractorClient.DirectoryZoneGet(ctx, detailZone)
			if err != nil {
				return err
			}
			o.Parent = cinp.StringAddr(r.GetURI())
		}

		err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI | extractID}}
Name:          {{.Name}}
Fqdn:          {{.Fqdn}}
Parent:        {{or .Parent ":<None>" | extractID}}
TTL:           {{.TTL}}
Refresh:       {{.Refresh}}
Retry:         {{.Retry}}
Expire:        {{.Expire}}
Minimum:       {{.Minimum}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)

		return nil
	},
}

var directoryUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Directory Zone",
	Args:  directoryArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		zoneID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o := contractorClient.DirectoryZoneNewWithID(zoneID)

		if detailName != "" {
			o.Name = &detailName
		}

		if detailTTL != 0 {
			o.TTL = &detailTTL
		}

		if detailRefresh != 0 {
			o.Refresh = &detailRefresh
		}

		if detailRetry != 0 {
			o.Retry = &detailRetry
		}

		if detailExpire != 0 {
			o.Expire = &detailExpire
		}

		if detailMinimum != 0 {
			o.Minimum = &detailMinimum
		}

		if detailZone != 0 {
			r, err := contractorClient.DirectoryZoneGet(ctx, detailZone)
			if err != nil {
				return err
			}
			o.Parent = cinp.StringAddr(r.GetURI())
		}

		err = o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI | extractID}}
Name:          {{.Name}}
Fqdn:          {{.Fqdn}}
Parent:        {{or .Parent ":<None>" | extractID}}
TTL:           {{.TTL}}
Refresh:       {{.Refresh}}
Retry:         {{.Retry}}
Expire:        {{.Expire}}
Minimum:       {{.Minimum}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)
		return nil
	},
}

var directoryDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Directory Zone",
	Args:  directoryArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		zoneID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		ctx := cmd.Context()

		o, err := contractorClient.DirectoryZoneGet(ctx, zoneID)
		if err != nil {
			return err
		}
		if err := o.Delete(ctx); err != nil {
			return err
		}

		return nil
	},
}

var directoryEntryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Work with Directory Entries",
}

var directoryEntryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Entries attached to a Directory Zone, or global Entries if no Zone Id is given",
	Args:  directoryEntryZoneArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		rl := []cinp.Object{}

		if len(args) == 1 {
			zoneID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			o, err := contractorClient.DirectoryZoneGet(ctx, zoneID)
			if err != nil {
				return err
			}

			vchan, err := contractorClient.DirectoryEntryList(ctx, "zone", map[string]interface{}{"zone": o.GetURI()})
			if err != nil {
				return err
			}
			for v := range vchan {
				rl = append(rl, v)
			}
		} else {
			vchan, err := contractorClient.DirectoryEntryList(ctx, "", map[string]interface{}{})
			if err != nil {
				return err
			}
			for v := range vchan {
				if v.Zone == nil {
					rl = append(rl, v)
				}
			}
		}

		outputList(rl, []string{"Id", "Zone", "Type", "Name", "Target", "Created", "Updated"}, "{{.GetURI | extractID}}	{{or .Zone \":<None>\" | extractID}}	{{.Type}}	{{.Name}}	{{.Target}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var directoryEntryGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Directory Entry",
	Args:  directoryEntryArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		entryID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o, err := contractorClient.DirectoryEntryGet(ctx, entryID)
		if err != nil {
			return err
		}
		outputDetail(o, `Id:            {{.GetURI | extractID}}
Zone:          {{or .Zone ":<None>" | extractID}}
Type:          {{.Type}}
Name:          {{.Name}}
Priority:      {{.Priority}}
Weight:        {{.Weight}}
Port:          {{.Port}}
Target:        {{.Target}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)
		return nil
	},
}

var directoryEntryCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create New Directory Entry in a Zone, or as a global Entry if no Zone Id is given",
	Args:  directoryEntryZoneArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		o := contractorClient.DirectoryEntryNew()

		if len(args) == 1 {
			zoneID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			r, err := contractorClient.DirectoryZoneGet(ctx, zoneID)
			if err != nil {
				return err
			}
			o.Zone = cinp.StringAddr(r.GetURI())
		}

		if detailType != "" {
			o.Type = &detailType
		}

		if detailName != "" {
			o.Name = &detailName
		}

		if detailPriority != 0 {
			o.Priority = &detailPriority
		}

		if detailWeight != 0 {
			o.Weight = &detailWeight
		}

		if detailPort != 0 {
			o.Port = &detailPort
		}

		if detailTarget != "" {
			o.Target = &detailTarget
		}

		err := o.Create(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI | extractID}}
Zone:          {{or .Zone ":<None>" | extractID}}
Type:          {{.Type}}
Name:          {{.Name}}
Priority:      {{.Priority}}
Weight:        {{.Weight}}
Port:          {{.Port}}
Target:        {{.Target}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)

		return nil
	},
}

var directoryEntryUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Directory Entry",
	Args:  directoryEntryArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		entryID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		o := contractorClient.DirectoryEntryNewWithID(entryID)

		if detailType != "" {
			o.Type = &detailType
		}

		if detailName != "" {
			o.Name = &detailName
		}

		if detailPriority != 0 {
			o.Priority = &detailPriority
		}

		if detailWeight != 0 {
			o.Weight = &detailWeight
		}

		if detailPort != 0 {
			o.Port = &detailPort
		}

		if detailTarget != "" {
			o.Target = &detailTarget
		}

		if detailZone != 0 {
			r, err := contractorClient.DirectoryZoneGet(ctx, detailZone)
			if err != nil {
				return err
			}
			o.Zone = cinp.StringAddr(r.GetURI())
		}

		err = o.Update(ctx)
		if err != nil {
			return err
		}

		outputDetail(o, `Id:            {{.GetURI | extractID}}
Zone:          {{or .Zone ":<None>" | extractID}}
Type:          {{.Type}}
Name:          {{.Name}}
Priority:      {{.Priority}}
Weight:        {{.Weight}}
Port:          {{.Port}}
Target:        {{.Target}}
Created:       {{.Created}}
Updated:       {{.Updated}}
`)
		return nil
	},
}

var directoryEntryDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Directory Entry",
	Args:  directoryEntryArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		entryID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		ctx := cmd.Context()

		o, err := contractorClient.DirectoryEntryGet(ctx, entryID)
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
	directoryCreateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Name of the New Directory Zone")
	directoryCreateCmd.Flags().IntVarP(&detailZone, "parent", "p", 0, "Parent Zone of the New Directory Zone")
	directoryCreateCmd.Flags().IntVarP(&detailTTL, "ttl", "t", 0, "TTL of the New Directory Zone")
	directoryCreateCmd.Flags().IntVarP(&detailRefresh, "refresh", "r", 0, "Refresh of the New Directory Zone")
	directoryCreateCmd.Flags().IntVarP(&detailRetry, "retry", "y", 0, "Retry of the New Directory Zone")
	directoryCreateCmd.Flags().IntVarP(&detailExpire, "expire", "e", 0, "Expire of the New Directory Zone")
	directoryCreateCmd.Flags().IntVarP(&detailMinimum, "minimum", "m", 0, "Minimum of the New Directory Zone")

	directoryUpdateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Update the Name of the Directory Zone")
	directoryUpdateCmd.Flags().IntVarP(&detailZone, "parent", "p", 0, "Update the Parent Zone of the Directory Zone")
	directoryUpdateCmd.Flags().IntVarP(&detailTTL, "ttl", "t", 0, "Update the TTL of the Directory Zone")
	directoryUpdateCmd.Flags().IntVarP(&detailRefresh, "refresh", "r", 0, "Update the Refresh of the Directory Zone")
	directoryUpdateCmd.Flags().IntVarP(&detailRetry, "retry", "y", 0, "Update the Retry of the Directory Zone")
	directoryUpdateCmd.Flags().IntVarP(&detailExpire, "expire", "e", 0, "Update the Expire of the Directory Zone")
	directoryUpdateCmd.Flags().IntVarP(&detailMinimum, "minimum", "m", 0, "Update the Minimum of the Directory Zone")

	directoryEntryCreateCmd.Flags().StringVarP(&detailType, "type", "t", "", "Type of the New Directory Entry (MX, SRV, CNAME, TXT)")
	directoryEntryCreateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Name of the New Directory Entry")
	directoryEntryCreateCmd.Flags().IntVarP(&detailPriority, "priority", "i", 0, "Priority of the New Directory Entry (MX, SRV)")
	directoryEntryCreateCmd.Flags().IntVarP(&detailWeight, "weight", "w", 0, "Weight of the New Directory Entry (SRV)")
	directoryEntryCreateCmd.Flags().IntVarP(&detailPort, "port", "o", 0, "Port of the New Directory Entry (SRV)")
	directoryEntryCreateCmd.Flags().StringVarP(&detailTarget, "target", "g", "", "Target of the New Directory Entry (MX, SRV, CNAME, TXT)")

	directoryEntryUpdateCmd.Flags().StringVarP(&detailType, "type", "t", "", "Update the Type of the Directory Entry (MX, SRV, CNAME, TXT)")
	directoryEntryUpdateCmd.Flags().StringVarP(&detailName, "name", "n", "", "Update the Name of the Directory Entry")
	directoryEntryUpdateCmd.Flags().IntVarP(&detailPriority, "priority", "i", 0, "Update the Priority of the Directory Entry (MX, SRV)")
	directoryEntryUpdateCmd.Flags().IntVarP(&detailWeight, "weight", "w", 0, "Update the Weight of the Directory Entry (SRV)")
	directoryEntryUpdateCmd.Flags().IntVarP(&detailPort, "port", "o", 0, "Update the Port of the Directory Entry (SRV)")
	directoryEntryUpdateCmd.Flags().StringVarP(&detailTarget, "target", "g", "", "Update the Target of the Directory Entry (MX, SRV, CNAME, TXT)")
	directoryEntryUpdateCmd.Flags().IntVarP(&detailZone, "zone", "z", 0, "Update the Zone of the Directory Entry")

	rootCmd.AddCommand(directoryCmd)
	directoryCmd.AddCommand(directoryListCmd, directoryGetCmd, directoryCreateCmd, directoryUpdateCmd, directoryDeleteCmd)

	directoryCmd.AddCommand(directoryEntryCmd)
	directoryEntryCmd.AddCommand(directoryEntryListCmd, directoryEntryGetCmd, directoryEntryCreateCmd, directoryEntryUpdateCmd, directoryEntryDeleteCmd)
}
