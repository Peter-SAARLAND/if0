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
	"if0/environments"

	"github.com/spf13/cobra"
)

// provisionCmd represents the provision command
var provisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "A brief description of your command",
	Long: `Example: if0 provision [env-name]`,
	Run: func(cmd *cobra.Command, args []string) {
		envDir := getEnvDir(args)
		err := environments.ZeroProvision(envDir)
		if err != nil {
			fmt.Println("Error: zero provision - ", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(provisionCmd)
}
