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
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/olekukonko/tablewriter"
	toml "github.com/pelletier/go-toml/v2"
	contractor "github.com/t3kton/contractor_goclient"

	cinp "github.com/cinp/go/v2"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var asJSON, asTOML, debug, showSecrets bool
var version = "development"
var gitVersion = "none"
var contractorClient *contractor.Contractor = nil

var rootCmd = &cobra.Command{
	Use:   "contractorcli",
	Short: "A CLI utility to work with Contractor",
	Long: `contractorcli allows you to do some basic manipulation
of contractor without having to write your own small app, or use the API`,
	SilenceUsage:  true,
	SilenceErrors: false,
}

var versionCmd = &cobra.Command{
	Use:                "version",
	Short:              "Show Version",
	PersistentPreRunE:  func(cmd *cobra.Command, args []string) error { return nil },
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error { return nil },
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
	cobra.OnInitialize(doInit)

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if asJSON && asTOML {
			return fmt.Errorf("only one of --json or --toml may be specified")
		}
		return connectContractor(cmd)
	}
	rootCmd.PersistentPostRunE = func(cmd *cobra.Command, args []string) error {
		doFinalize(cmd)
		return nil
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.contractorcli.ini)")
	rootCmd.PersistentFlags().BoolVar(&asJSON, "json", false, "Output as JSON")
	rootCmd.PersistentFlags().BoolVar(&asTOML, "toml", false, "Output as TOML")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "Debug Output(will interfere with JSON output)")
	rootCmd.PersistentFlags().BoolVar(&showSecrets, "show-secrets", false, "Show passwords/secrets in plain text in table output (JSON/TOML output always includes them)")

	rootCmd.AddCommand(versionCmd)
}

func doInit() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".contractorcli")
	}

	viper.SetDefault("contractor.host", "http://contractor")
	viper.SetDefault("contractor.proxy", "")
	viper.SetEnvPrefix("CONTRACTOR")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Error reading config file: '%s'\n", err)
			os.Exit(1)
		}
	}
}

func connectContractor(cmd *cobra.Command) error {
	handlerOptions := &slog.HandlerOptions{}
	if debug {
		handlerOptions.Level = slog.LevelDebug
	} else {
		handlerOptions.Level = slog.LevelWarn
	}
	log := slog.New(slog.NewTextHandler(os.Stderr, handlerOptions))

	var err error
	contractorClient, err = contractor.NewContractor(cmd.Context(), log, viper.GetString("contractor.host"), viper.GetString("contractor.proxy"), viper.GetString("contractor.username"), viper.GetString("contractor.password"))
	if err != nil {
		return err
	}
	return nil
}

func doFinalize(cmd *cobra.Command) {
	if contractorClient == nil {
		return
	}
	contractorClient.Logout(cmd.Context())
	contractorClient = nil
}

func extractID(value string) string {
	if value == "" {
		return ""
	}
	return strings.Split(value, ":")[1]
}

func maskSecret(value string) string {
	if value == "" || showSecrets {
		return value
	}
	return "<hidden, use --show-secrets to display>"
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
			return "", fmt.Errorf("unable to detect or find an editor, set CONTRACTORCLI_EDITOR or EDITOR")
		}
	}

	tmpfile, err := os.CreateTemp("", "contractorcli")
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

	editorArgs := strings.Fields(editor)
	cmd := exec.Command(editorArgs[0], append(editorArgs[1:], tmpfile.Name())...)
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

	buf, err := io.ReadAll(fd)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(buf)), nil
}

func outputList(valueList []cinp.Object, header []string, itemTemplate string) {
	if asJSON {
		buff, err := json.MarshalIndent(valueList, "", " ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Stdout.Write(buff)
		os.Stdout.Write([]byte("\n"))
	} else if asTOML {
		buff, err := toml.Marshal(map[string]interface{}{"items": valueList})
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
		t.Funcs(template.FuncMap{"extractID": extractID, "extractIDList": extractIDList, "maskSecret": maskSecret})
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
		os.Stdout.Write([]byte("\n"))
	} else if asTOML {
		buff, err := toml.Marshal(value)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Stdout.Write(buff)
	} else {
		t := template.New("output")
		t.Funcs(template.FuncMap{"extractID": extractID, "extractIDList": extractIDList, "maskSecret": maskSecret})
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
	} else if asTOML {
		buff, err := toml.Marshal(valueMap)
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

// readConfigFile reads a config values file (or stdin if path is "-") and decodes
// it as either JSON or TOML, detecting the format from the file extension, falling
// back to sniffing the content if the extension is absent or unrecognized.
func readConfigFile(path string) (map[string]interface{}, error) {
	var reader io.Reader
	if path == "-" {
		reader = os.Stdin
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		reader = f
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	useJSON := false
	switch {
	case strings.HasSuffix(strings.ToLower(path), ".json"):
		useJSON = true
	case strings.HasSuffix(strings.ToLower(path), ".toml"):
		useJSON = false
	default:
		trimmed := bytes.TrimLeft(data, " \t\r\n")
		useJSON = len(trimmed) > 0 && trimmed[0] == '{'
	}

	newValues := map[string]interface{}{}
	if useJSON {
		if err := json.Unmarshal(data, &newValues); err != nil {
			return nil, fmt.Errorf("unable to parse '%s' as JSON: %w", path, err)
		}
	} else {
		if err := toml.Unmarshal(data, &newValues); err != nil {
			return nil, fmt.Errorf("unable to parse '%s' as TOML: %w", path, err)
		}
	}

	return newValues, nil
}
