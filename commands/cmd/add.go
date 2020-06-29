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
	"if0/environments"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "used to add an environment",
	Long: `Example: if0 add test-env [git@gitlab.com:test-env.git]
This command adds the environment locally at ~/.if0/.environments/gitlab.com/test-env.

If GL_TOKEN is set in ~/.if0/if0.env, a private GitLab project is created, 
and the local environment is synced with it. In this case a repo url is not needed.

If a repo url is provided, the repository present at the remote repository url is cloned locally.
If there is no GL_TOKEN and no repo url, only a local copy is created at ~/.if0/.environments`,

	Run: func(cmd *cobra.Command, args []string) {
		//if len(args) < 1 {
		//	fmt.Println("Please provide valid arguments.")
		//	fmt.Println("example command: if0 add repo-name [git@gitlab.com:repo-name.git]")
		//	return
		//}
		err := environments.AddEnv(args)
		if err != nil {
			fmt.Println("Error: Adding repo - ", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
