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

func jobArgCheck(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a Job Id Argument")
	}
	return nil
}

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Work with Jobs",
}

var jobFoundationCmd = &cobra.Command{
	Use:   "foundation",
	Short: "Work with Foundation Jobs",
}

var jobFoundationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Foundation Jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		rl := []cinp.Object{}
		vchan, err := contractorClient.ForemanFoundationJobList(ctx, "", map[string]interface{}{})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Id", "Foundation", "State", "Status", "Message", "Script", "Updated", "Created"}, "{{.GetID | extractID}}	{{.Foundation | extractID}}	{{.State}}	{{.Status}}	{{.Message}}	{{.ScriptName}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var jobFoundationGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Foundation Job",
	Args:  jobArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		jobID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		r, err := contractorClient.ForemanFoundationJobGet(ctx, jobID)
		if err != nil {
			return err
		}
		outputDetail(r, `Id:          {{.GetID | extractID}}
Site:        {{.Site}}
Foundation:  {{.Foundation | extractID}}
Script:      {{.ScriptName}}
State:       {{.State}}
Status:      {{.Status}}
Message:     {{.Message}}
Status:      {{.Status}}
CanStart:    {{.CanStart}}
Created:     {{.Created}}
Updated:     {{.Updated}}
`)
		return nil
	},
}

var jobFoundationPauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause Foundation Job",
	Args:  jobArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		jobID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		r, err := contractorClient.ForemanFoundationJobGet(ctx, jobID)
		if err != nil {
			return err
		}

		if err = r.CallPause(ctx); err != nil {
			return err
		}

		return nil
	},
}

var jobFoundationResumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume Foundation Job",
	Args:  jobArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		jobID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		r, err := contractorClient.ForemanFoundationJobGet(ctx, jobID)
		if err != nil {
			return err
		}

		if err = r.CallResume(ctx); err != nil {
			return err
		}

		return nil
	},
}

var jobFoundationRestCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset Foundation Job",
	Args:  jobArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		jobID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		r, err := contractorClient.ForemanFoundationJobGet(ctx, jobID)
		if err != nil {
			return err
		}

		if err = r.CallReset(ctx); err != nil {
			return err
		}

		return nil
	},
}

var jobFoundationRollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback Foundation Job",
	Args:  jobArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		jobID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		r, err := contractorClient.ForemanFoundationJobGet(ctx, jobID)
		if err != nil {
			return err
		}

		if err = r.CallRollback(ctx); err != nil {
			return err
		}

		return nil
	},
}

var jobStructureCmd = &cobra.Command{
	Use:   "structure",
	Short: "Work with Structure Jobs",
}

var jobStructureListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Structure Jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		rl := []cinp.Object{}
		vchan, err := contractorClient.ForemanStructureJobList(ctx, "", map[string]interface{}{})
		if err != nil {
			return err
		}
		for v := range vchan {
			rl = append(rl, v)
		}
		outputList(rl, []string{"Id", "Structure", "State", "Status", "Message", "Script", "Updated", "Created"}, "{{.GetID | extractID}}	{{.Structure | extractID}}	{{.State}}	{{.Status}}	{{.Message}}	{{.ScriptName}}	{{.Created}}	{{.Updated}}\n")

		return nil
	},
}

var jobStructureGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Structure Job",
	Args:  jobArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		jobID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		r, err := contractorClient.ForemanStructureJobGet(ctx, jobID)
		if err != nil {
			return err
		}
		outputDetail(r, `Id:          {{.GetID | extractID}}
Site:        {{.Site}}
Structure:   {{.Structure | extractID}}
Script:      {{.ScriptName}}
State:       {{.State}}
Status:      {{.Status}}
Message:     {{.Message}}
Status:      {{.Status}}
CanStart:    {{.CanStart}}
Created:     {{.Created}}
Updated:     {{.Updated}}
`)
		return nil
	},
}

var jobStructurePauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause Structure Job",
	Args:  jobArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		jobID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		r, err := contractorClient.ForemanStructureJobGet(ctx, jobID)
		if err != nil {
			return err
		}

		if err = r.CallPause(ctx); err != nil {
			return err
		}

		return nil
	},
}

var jobStructureResumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume Structure Job",
	Args:  jobArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		jobID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		r, err := contractorClient.ForemanStructureJobGet(ctx, jobID)
		if err != nil {
			return err
		}

		if err = r.CallResume(ctx); err != nil {
			return err
		}

		return nil
	},
}

var jobStructureRestCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset Structure Job",
	Args:  jobArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		jobID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		r, err := contractorClient.ForemanStructureJobGet(ctx, jobID)
		if err != nil {
			return err
		}

		if err = r.CallReset(ctx); err != nil {
			return err
		}

		return nil
	},
}

var jobStructureRollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback Structure Job",
	Args:  jobArgCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		jobID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		r, err := contractorClient.ForemanStructureJobGet(ctx, jobID)
		if err != nil {
			return err
		}

		if err = r.CallRollback(ctx); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(jobCmd)
	jobCmd.AddCommand(jobFoundationCmd)
	jobFoundationCmd.AddCommand(jobFoundationListCmd, jobFoundationGetCmd, jobFoundationPauseCmd, jobFoundationResumeCmd, jobFoundationRestCmd, jobFoundationRollbackCmd)

	jobCmd.AddCommand(jobStructureCmd)
	jobStructureCmd.AddCommand(jobStructureListCmd, jobStructureGetCmd, jobStructurePauseCmd, jobStructureResumeCmd, jobStructureRestCmd, jobStructureRollbackCmd)
}
