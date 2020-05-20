/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"github.com/spf13/cobra"
	"if0/common"
	"if0/environments"
	"os"
	"path/filepath"
)

const (
	addArg  = "add"
	syncArg = "sync"
	planArg = "plan"
)

// environmentCmd represents the environment command
var environmentCmd = &cobra.Command{
	Use:   "environment",
	Short: "add zero config repository to environments",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide valid arguments.")
			fmt.Println("accepted args: 'add', 'sync', 'load'")
			return
		}
		// cloning
		if args[0] == addArg {
			if len(args) < 2 {
				fmt.Println("Please provide valid arguments.")
				fmt.Println("example command: if0 environment add git@gitlab.com:abc/def.git")
				return
			}
			repoUrl := args[1]
			err := environments.AddEnv(repoUrl)
			if err != nil {
				fmt.Println("Error: Adding repo - ", err)
				return
			}
		}

		// syncing
		if args[0] == syncArg {
			if len(args) < 2 {
				fmt.Println("Please provide valid arguments.")
				fmt.Println("example command: if0 environment sync if0-config")
				return
			}
			repoName := args[1]
			err := environments.SyncEnv(repoName)
			if err != nil {
				fmt.Println("Error: Syncing repo - ", err)
				return
			}
		}

		// loading Environment
		if args[0] == planArg {
			var envDir string
			if len(args) < 2 {
				envDir, _ = os.Getwd()
			} else {
				envDir = filepath.Join(common.EnvDir, args[1])
			}
			err := environments.PlanEnv(envDir)
			if err != nil {
				fmt.Println("Error: Planning env - ", err)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(environmentCmd)
	//environmentCmd.Flags().StringVar(&add, "add", "", "add a new environments config git repository")
}
