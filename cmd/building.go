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

	"github.com/spf13/cobra"
)

var structureCmd = &cobra.Command{
	Use:   "structure",
	Short: "Work with structures",
}

var structureListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Structures",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("List Structures")
	},
}

var structureGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Structure",
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
		id, _ := strconv.Atoi(args[0])
		fmt.Printf("Getting Site %d\n", id)
		c := getContractor()
		s, err := c.Building.Structure.Get(id)
		fmt.Println(err)
		fmt.Printf("%+v\n", s)

	},
}

func init() {
	rootCmd.AddCommand(structureCmd)
	structureCmd.AddCommand(structureListCmd)
	structureCmd.AddCommand(structureGetCmd)
}
