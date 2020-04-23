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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"if0/config"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync with a remote git repository",
	Long: `This command helps sync .if0 config files with a remote git repository.
The remote git repository to be synced with should be set in the ~/.if0/if0.env file 
as a value for the config variable REMOTE_STORAGE `,
	Run: func(cmd *cobra.Command, args []string) {
		err := config.RepoSync()
		if err != nil {
			log.Errorln("Error while syncing with remote repo: ", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
