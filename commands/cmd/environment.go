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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"if0/environments"
)

const (
	addArg  = "add"
	syncArg = "sync"
)

// environmentCmd represents the environment command
var environmentCmd = &cobra.Command{
	Use:   "environment",
	Short: "add zero config repository to environments",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// checking if the required arguments are present
		if len(args) < 2 {
			fmt.Println("Please provide valid arguments.")
			fmt.Println("example command: if0 environment add git@gitlab.com:abc/def.git")
			return
		}
		// cloning
		clone := args[0] == addArg
		repoUrl := args[1]
		if clone {
			err := environments.AddEnv(repoUrl)
			if err != nil {
				log.Errorln("Error adding repo - ", err)
				return
			}
		}

		// syncing
		sync := args[0] == syncArg
		repoName := args[1]
		if sync {
			err := environments.SyncEnv(repoName)
			if err != nil {
				log.Errorln("Error syncing repo - ", err)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(environmentCmd)
	//environmentCmd.Flags().StringVar(&add, "add", "", "add a new environments config git repository")
}
