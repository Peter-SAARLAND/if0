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
	"if0/config"
	"os"
	"strings"
)

var (
	// zero flag: used to distinguish between if0 and zero-cluster configuration files.
	// default: false, for if0
	// set to true for zero-cluster configurations when the command is called with -z or --zero flag
	//zero bool

	// add flag: used to add new or update configuration files.
	add string

	//sync flag: used to sync with an external repository.
	//used in conjunction with REMOTE_STORAGE (git repository link) and
	//REPO_SYNC (bool value) variables
	sync bool

	// merge flag: used to merge the new configuration with current running configuration
	// default: false
	// set to true to merge and replace current configuration
	// the file in --src is merged with file in --dst.
	// --dst is optional
	merge bool

	// src flag: used to input the configuration file that needs to be merged
	src string

	// dst file: destination configuration file (in .environments dir, or if0.env file) to be merged with
	dst string

	// set flag: used to set environment variables.
	// accepts comma separated values. Example: if0 addConfig --set "TESTVAR1=testval1, TESTVAR2-testval2"
	set []string

	// configCmd represents the addConfig command
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "adds/modifies running configuration files of if0 or zero clusters",
		Long: `if0 is a CLI tool for zero. 
		It has features to add or update configuration for if0 or for zero-clusters`,
		Run: func(cmd *cobra.Command, args []string) {
			if set != nil {
				fmt.Println("Reading environment variables from flag --set")
				loadConfigFromFlags(set)
			}

			// if --merge is true, the file in --src is merged with file in --dst.
			// --dst is optional
			if merge {
				err := config.MergeConfigFiles(src, dst)
				if err != nil {
					fmt.Println("Error: Merging config files - ", err)
					return
				}
			}

			// checking if a configuration file has been provided in the command
			if add != "" {
				fmt.Println("Updating configuration")
				loadConfigFromFile(add)
			}

			// printing current running configuration to the stdout.
			fmt.Println("Current Running Configuration")
			config.PrintCurrentRunningConfig()

			// automatic garbage collection
			config.GarbageCollection()

			//
			//if sync {
			//	err := config.RepoSync()
			//	if err != nil {
			//		fmt.Errorln("Error while syncing with remote repo: ", err)
			//		return
			//	}
			//}
			if sync {
				fmt.Println("if0 config --sync functionality temporarily disabled.")
				return
			}
		},
	}
)

// loadConfigFromFlags is called when the if0 command is called with --set flag
// it is used to set config variables for that particular run
// example: `if0 config --set var1=val1` sets var1 with value val1
func loadConfigFromFlags(configParams []string) {
	for _, param := range configParams {
		set := strings.Split(param, "=")
		config.SetEnvVariable(set[0], set[1])
	}
}

func loadConfigFromFile(configFile string) {
	// validating the configuration file
	isValid, err := config.IsConfigFileValid(configFile)
	if !isValid {
		fmt.Println("Error: Terminating config update: ", err)
		return
	}

	// checking if the provided configuration file is present.
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("Error: The provided configuration file %s is not found. \n", configFile)
		return
	}

	// adding/updating the config file
	err = config.AddConfigFile(configFile)
	if err != nil {
		fmt.Println("Error: Adding config file - ", err)
	}
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringSliceVar(&set, "set", nil, "sets env variables via CLI")
	//configCmd.Flags().BoolVarP(&zero, "zero", "z",
	//	false, "updates zero cluster configuration")
	configCmd.Flags().BoolVarP(&merge, "merge", "m",
		false, "merges the new configuration with running configuration")
	configCmd.Flags().StringVar(&add, "add", "", "configuration file to be added or updated")
	configCmd.Flags().StringVar(&src, "src", "", "source configuration file for merge")
	configCmd.Flags().StringVar(&dst, "dst", "", "destination configuration file to merge with")
	configCmd.Flags().BoolVar(&sync, "sync",
		false, "syncs configuration files with an external repository")

}
