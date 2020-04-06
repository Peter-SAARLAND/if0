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
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"if0/config"
	"os"
)

const (
	IFO_VERSION  = "IF0_VERSION"
	ZERO_VERSION = "ZERO_VERSION"
)

var (
	// zero flag: used to distinguish between if0 and zero-cluster configuration files.
	// default: false, for if0
	// set to true for zero-cluster configurations when the command is called with -z or --zero flag
	zero bool

	// addConfigCmd represents the addConfig command
	addConfigCmd = &cobra.Command{
		Use:   "addConfig",
		Short: "adds/modifies running configuration files of if0 or zero clusters",
		Long: `if0 is a CLI tool for zero. 
		It has features to add or update configuration for if0 or for zero-clusters`,
		Run: func(cmd *cobra.Command, args []string) {

			// printing current running configuration to the stdout.
			log.Println("Current Running Configuration")
			config.PrintCurrentRunningConfig()

			// checking if a configuration file has been provided in the command
			if len(args) == 0 {
				log.Fatalln("Configuration file missing. Please provide a valid configuration file.")
			}

			configFile := args[0]
			// validating the configuration file
			isValid, err := isConfigFileValid(configFile)
			if !isValid {
				log.Fatalln("Terminating config update: ", err)
			}

			// checking if the provided configuration file is present.
			rootPath, _ := os.Getwd()
			filePath := rootPath + string(os.PathSeparator) + configFile
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				log.Fatalf("The provided configuration file %s is not found.", filePath)
			}

			// adding/updating the config file
			config.AddConfigFile(configFile, zero)
		},
	}
)

func isConfigFileValid(configFile string) (bool, error) {
	// read IF0_VERSION, ZERO_VERSION
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Error while reading config file: ", err)
	}
	if0Version := viper.IsSet(IFO_VERSION)
	zeroVersion := viper.IsSet(ZERO_VERSION)

	if !if0Version && !zeroVersion {
		return false, errors.New("no valid versions (IF0_VERSION or ZERO_VERSION) found in the config file")
	} else if if0Version && zero {
		return false, errors.New("zero-cluster config update invoked with if0.env config file")
	} else if zeroVersion && !zero {
		return false, errors.New("if0.env update invoked with zero-cluster config file")
	}
	return true, nil
}

func init() {
	rootCmd.AddCommand(addConfigCmd)

	// adding a 'zero' flag to the addConfig command.
	// default value: false. By default, the configuration is updated to if0.env
	// upon setting the zero flag, zero cluster configuration is updated
	addConfigCmd.Flags().BoolVarP(&zero, "zero", "z",
		false, "updates zero cluster configuration")
}
