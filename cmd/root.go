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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"

	contractor "github.com/t3kton/contractor_client/go"

	"github.com/olekukonko/tablewriter"

	cinp "github.com/cinp/go"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var asJSON bool
var version = "development"
var gitVersion = "none"

var rootCmd = &cobra.Command{
	Use:   "contractorcli",
	Short: "A CLI utility to work with Contractor",
	Long: `contractorcli allows you to do some basic maniplutation
of contractor without having to write your own small app, or use the API`,
	SilenceUsage:  true,
	SilenceErrors: false,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show Version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("contractorcli\n  Version:\t%s\n  Commit:\t%s\n", version, gitVersion)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// fmt.Println("Error:", err) cobra prints the error message unless SilenceErrors is set
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.contractorcli.ini)")
	rootCmd.PersistentFlags().BoolVarP(&asJSON, "json", "j", false, "Output as JSON")

	rootCmd.AddCommand(versionCmd)
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
	viper.SetDefault("contractor.proxy", "")
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
	c, err := contractor.NewContractor(viper.GetString("contractor.host"), viper.GetString("contractor.proxy"), viper.GetString("contractor.username"), viper.GetString("contractor.password"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return c
}

func extractID(value string) string {
	if value == "" {
		return ""
	}
	return strings.Split(value, ":")[1]
}

func extractIDList(values []string) string {
	if len(values) == 0 {
		return ""
	}
	workList := []string{}

	for _, value := range values {
		workList = append(workList, strings.Split(value, ":")[1])
	}
	return strings.Join(workList, ",")
}

func extractIDInt(value string) (int, error) {
	value = extractID(value)
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func editBuffer(value string) (string, error) {
	editor := os.Getenv("CONTRACTORCLI_EDITOR")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		tmp, err := os.Readlink("/etc/alternatives/editor")
		if err == nil {
			editor = tmp
		}
	}
	if editor == "" {
		editor = "/usr/bin/vi"
		_, err := os.Stat(editor)
		if err != nil {
			return "", fmt.Errorf("Unable to detect or find an editor, set CONTRACTORCLI_EDITOR or EDITOR")
		}
	}

	tmpfile, err := ioutil.TempFile("", "contractorcli")
	if err != nil {
		return "", err
	}
	defer func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}()

	if _, err := tmpfile.Write([]byte(value)); err != nil {
		return "", err
	}

	if err := tmpfile.Close(); err != nil {
		return "", err
	}

	cmd := exec.Command(editor, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	fd, err := os.Open(tmpfile.Name())
	if err != nil {
		return "", err
	}
	defer fd.Close()

	buf := make([]byte, 4096*1024)
	len, err := fd.Read(buf)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(buf[:len])), nil
}

func outputList(valueList []cinp.Object, header []string, itemTemplate string) {
	if asJSON {
		buff, err := json.MarshalIndent(valueList, "", " ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Stdout.Write(buff)
	} else {
		var rederbuff bytes.Buffer
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(header)
		t := template.New("output")
		t.Funcs(template.FuncMap{"extractID": extractID, "extractIDList": extractIDList})
		t, err := t.Parse(itemTemplate)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, value := range valueList {
			rederbuff.Reset()
			err = t.Execute(&rederbuff, value)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			table.Append(strings.Split(rederbuff.String(), "\t"))
		}
		table.Render()
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
		t := template.New("output")
		t.Funcs(template.FuncMap{"extractID": extractID, "extractIDList": extractIDList})
		t, err := t.Parse(detailTemplate)
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
