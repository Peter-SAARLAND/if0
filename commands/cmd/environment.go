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
	provisionArg = "provision"
	zeroArg = "zero"
	destroyArg = "destroy"
)

// environmentCmd represents the environment command
var environmentCmd = &cobra.Command{
	Use:   "env",
	Short: "add zero config repository to environments",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide valid arguments.")
			fmt.Println("accepted args: 'add', 'sync', 'plan', 'zero', 'provision'")
			return
		}

		switch args[0] {
		case addArg:
			if len(args) < 2 {
				fmt.Println("Please provide valid arguments.")
				fmt.Println("example command: if0 environment add repo-name [git@gitlab.com:repo-name.git]")
				return
			}
			err := environments.AddEnv(args)
			if err != nil {
				fmt.Println("Error: Adding repo - ", err)
				return
			}
		case syncArg:
			envDir := getEnvDir(args)
			err := environments.SyncEnv(envDir)
			if err != nil {
				fmt.Println("Error: Syncing repo - ", err)
				return
			}
		case planArg:
			envDir := getEnvDir(args)
			err := environments.Dash1Plan(envDir)
			if err != nil {
				fmt.Println("Error: dash1 plan - ", err)
				return
			}
		case provisionArg:
			envDir := getEnvDir(args)
			err := environments.ZeroPlatform(envDir)
			if err != nil {
				fmt.Println("Error: zero provision - ", err)
				return
			}
		case zeroArg:
			envDir := getEnvDir(args)
			err := environments.Dash1Infrastructure(envDir)
			if err != nil {
				fmt.Println("Error: dash1 zero - ", err)
				return
			}
		case destroyArg:
			envDir := getEnvDir(args)
			err := environments.Dash1Destroy(envDir)
			if err != nil {
				fmt.Println("Error: dash1 destroy - ", err)
				return
			}
		}
	},
}

func getEnvDir(args []string) string {
	var envDir string
	if len(args) < 2 {
		envDir, _ = os.Getwd()
	} else {
		envDir = filepath.Join(common.EnvDir, args[1])
	}
	return envDir
}

func init() {
	//rootCmd.AddCommand(environmentCmd)
}
