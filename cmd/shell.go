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

/*
TODO: watch https://github.com/c-bata/go-prompt/issues/94 for updates on simplifying history load

*/

import (
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	cobraprompt "github.com/stromland/cobra-prompt"
	"golang.org/x/term"
)

var shellExiting bool

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Start Interactive Shell",
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.RemoveCommand(cmd)

		fd := int(os.Stdin.Fd())
		savedState, err := term.GetState(fd)
		if err != nil {
			savedState = nil
		}

		shell := &cobraprompt.CobraPrompt{
			RootCmd: rootCmd,
			GoPromptOptions: []prompt.Option{
				prompt.OptionTitle("contractor"),
				prompt.OptionPrefix("contractor> "),
				prompt.OptionMaxSuggestion(10),
				prompt.OptionSetExitCheckerOnInput(func(in string, breakline bool) bool {
					return shellExiting
				}),
			},
		}
		shell.Run()

		// Restore terminal state that go-prompt left in raw mode when
		// exiting via ExitChecker (its teardown skips input restoration).
		if savedState != nil {
			term.Restore(fd, savedState)
		}
		fmt.Println()
	},
}

var exitCmd = &cobra.Command{
	Use:               "exit",
	Short:             "Exit the shell",
	PersistentPreRunE:  func(cmd *cobra.Command, args []string) error { return nil },
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		shellExiting = true
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
	rootCmd.AddCommand(exitCmd)
}
