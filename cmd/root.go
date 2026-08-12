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
	"sort"
	"strconv"
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

// The Map fields (config_values, id_map, validation_template, script_map, ...) are arbitrary
// nested structures, and Go's default map formatting turns anything deeper than one level into an
// unreadable run-on.  prettyMap flattens them to one sorted "dotted.path = value" per line -- the
// same path notation contractor itself uses in id_map validation errors, so a mismatch it reports
// can be looked up here directly.
func prettyValue(value interface{}) (string, bool) { // ( rendering, is a leaf ) -- a leaf is anything that fits on one line
	switch item := value.(type) {
	case map[string]interface{}:
		return "", false

	case []interface{}:
		partList := []string{}
		for _, entry := range item {
			text, isLeaf := prettyValue(entry)
			if !isLeaf { // a list with a map in it gets a line per element instead, see appendPretty
				return "", false
			}
			partList = append(partList, text)
		}
		return "[" + strings.Join(partList, ", ") + "]", true

	case string:
		return item, true

	case float64: // everything numeric arrives as a float64 from the json decode, so format it back without the exponent %v would use for anything large
		return strconv.FormatFloat(item, 'f', -1, 64), true

	case nil:
		return "<null>", true

	default:
		return fmt.Sprintf("%v", item), true
	}
}

func appendPretty(path string, value interface{}, lineList *[]string) {
	if text, isLeaf := prettyValue(value); isLeaf {
		*lineList = append(*lineList, path+" = "+text)
		return
	}

	switch item := value.(type) {
	case map[string]interface{}:
		if len(item) == 0 {
			*lineList = append(*lineList, path+" = {}")
			return
		}
		appendPrettyMap(path, item, lineList)

	case []interface{}:
		if len(item) == 0 {
			*lineList = append(*lineList, path+" = []")
			return
		}
		for index, entry := range item {
			appendPretty(fmt.Sprintf("%s.%d", path, index), entry, lineList)
		}
	}
}

func appendPrettyMap(prefix string, value map[string]interface{}, lineList *[]string) {
	keyList := make([]string, 0, len(value))
	for key := range value {
		keyList = append(keyList, key)
	}
	sort.Strings(keyList)

	for _, key := range keyList {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}
		appendPretty(path, value[key], lineList)
	}
}

func prettyMap(indent int, value *map[string]interface{}) string {
	if value == nil || len(*value) == 0 {
		return ""
	}

	lineList := []string{}
	appendPrettyMap("", *value, &lineList)

	return strings.Join(lineList, "\n"+strings.Repeat(" ", indent))
}

// A job's Status is a list of { "percent", "operation", "parameters" } maps, outermost scope
// first, getting more specific as it goes.  Printed raw that is a wall of Go map formatting, so
// jobStatus/jobStatusShort render it the way the ui does.  Every lookup below is a comma-ok type
// assertion and "parameters" is often null, which is fine -- indexing a nil map yields the zero
// value rather than panicking, so a status entry that is not shaped as expected degrades to a
// shorter line instead of taking the whole command down.
func jobStatusLabel(item map[string]interface{}) string {
	operation, _ := item["operation"].(string)
	parameters, _ := item["parameters"].(map[string]interface{})

	switch operation {
	case "Scope":
		if description, ok := parameters["description"].(string); ok && description != "" {
			return description
		}

	case "Function":
		name, _ := parameters["name"].(string)
		module, _ := parameters["module"].(string)
		label := operation
		if module != "" && name != "" {
			label = module + "." + name
		}
		if dispatched, _ := parameters["dispatched"].(bool); dispatched {
			label += " [dispatched]"
		}
		return label

	default:
		if doing, ok := parameters["doing"].(string); ok && doing != "" {
			return operation + "(" + doing + ")"
		}
	}

	return operation
}

func jobStatusLine(item map[string]interface{}) string {
	percent, _ := item["percent"].(float64)
	parameters, _ := item["parameters"].(map[string]interface{})

	partList := []string{fmt.Sprintf("%6.2f%%", percent), jobStatusLabel(item)}

	if elapsed, ok := parameters["time_elapsed"].(string); ok && elapsed != "" {
		partList = append(partList, "Elapsed: "+elapsed)
	}
	if remaining, ok := parameters["time_remaining"].(string); ok && remaining != "" {
		partList = append(partList, "Remaining: "+remaining)
	}

	return strings.Join(partList, " ")
}

func jobStatus(indent int, status *[]map[string]interface{}) string { // one line per entry, continuations indented to line up under the template's label column
	if status == nil || len(*status) == 0 {
		return ""
	}

	lineList := []string{}
	for _, item := range *status {
		lineList = append(lineList, jobStatusLine(item))
	}

	return strings.Join(lineList, "\n"+strings.Repeat(" ", indent))
}

func jobStatusShort(status *[]map[string]interface{}) string { // a single cell for the table output: overall progress, plus whatever it is currently doing
	if status == nil || len(*status) == 0 {
		return ""
	}

	itemList := *status
	percent, _ := itemList[0]["percent"].(float64) // entry 0 is the outermost scope, ie. the job as a whole
	result := fmt.Sprintf("%.2f%%", percent)

	if label := jobStatusLabel(itemList[len(itemList)-1]); label != "" { // the last entry is the most specific
		result += " " + label
	}

	return result
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
		t.Funcs(template.FuncMap{"extractID": extractID, "extractIDList": extractIDList, "maskSecret": maskSecret, "jobStatus": jobStatus, "jobStatusShort": jobStatusShort, "prettyMap": prettyMap})
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
		t.Funcs(template.FuncMap{"extractID": extractID, "extractIDList": extractIDList, "maskSecret": maskSecret, "jobStatus": jobStatus, "jobStatusShort": jobStatusShort, "prettyMap": prettyMap})
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
	} else { // same rendering the Map fields get in the detail output, see prettyMap -- ranging the map directly also meant the order changed from run to run, which made two runs impossible to diff
		if text := prettyMap(0, &valueMap); text != "" {
			fmt.Println(text)
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
