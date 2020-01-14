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
package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cinp/go"
	"github.com/spf13/cobra"
)

var sitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "Work with sites",
}

var siteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Sites",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("List Sites")
	},
}

var siteGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Site",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a site id argument")
		}
		if _, err := strconv.Atoi(args[0]); err != nil {
			return errors.New("Invalid Id, must be a valid int")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		siteID, _ := strconv.Atoi(args[0])
		fmt.Printf("Getting Site %d\n", siteID)
		golang.NewClient()
	},
}

func init() {
	rootCmd.AddCommand(sitesCmd)
	sitesCmd.AddCommand(siteListCmd)
	sitesCmd.AddCommand(siteGetCmd)
}
