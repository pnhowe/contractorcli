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
	"encoding/json"
	"fmt"
	cinp "github.com/cinp/go"
	"github.com/spf13/cobra"
	"os"
	"text/template"

	"github.com/t3kton/contractorcli/contractor"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var asJSON bool

var rootCmd = &cobra.Command{
	Use:   "contractorcli",
	Short: "A CLI utility to work with Contractor",
	Long: `contractorcli allows you to do some basic maniplutation
of contractor without having to write your own small app, or use the API`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.contractorcli.ini)")
	rootCmd.PersistentFlags().BoolVarP(&asJSON, "json", "j", false, "Output as JSON")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("json", "j", false, "Output as JSON")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".contractorcli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".contractorcli")
	}

	viper.SetDefault("contractor.host", "http://contractor")
	viper.SetEnvPrefix("CONTRACTOR")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Error reading config file: '%s'\n", err)
			os.Exit(1)
		}
	} else {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getContractor() *contractor.Contractor {
	c, err := contractor.NewContractor(viper.GetString("contractor.host"), "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return c
}

func outputList(valueList []cinp.Object, header string, itemTemplate string) {
	if asJSON {
		buff, err := json.MarshalIndent(valueList, "", " ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Stdout.Write(buff)
	} else {
		t, err := template.New("output").Parse(itemTemplate)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Stdout.WriteString(header)
		for _, value := range valueList {
			err = t.Execute(os.Stdout, value)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}

func outputDetail(value interface{}, detailTemplate string) {
	if asJSON {
		buff, err := json.MarshalIndent(value, "", " ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Stdout.Write(buff)
	} else {
		t, err := template.New("output").Parse(detailTemplate)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = t.Execute(os.Stdout, value)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func outputKV(valueMap map[string]interface{}) {
	if asJSON {
		buff, err := json.MarshalIndent(valueMap, "", " ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Stdout.Write(buff)
	} else {
		for k, v := range valueMap {
			fmt.Printf("%s:\t%+v\n", k, v)
		}
	}
}
