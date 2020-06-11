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

// zeroCmd represents the zero command
var zeroCmd = &cobra.Command{
	Use:   "zero",
	Short: "",
	Long: `Example: if0 zero [env-name]`,
	Run: func(cmd *cobra.Command, args []string) {
		envDir := getEnvDir(args)
		err := environments.Dash1Infrastructure(envDir)
		if err != nil {
			fmt.Println("Error: dash1 zero - ", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(zeroCmd)
}
