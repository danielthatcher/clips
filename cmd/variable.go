package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set [variable] [value]",
	Short: "Set a variable",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		variables[args[0]] = args[1]
		viper.Set("variables", variables)
		err := viper.WriteConfigAs(cfgFile)
		if err != nil {
			log.Fatalf("Failed to write configuration file to %s: %v\n", cfgFile, err)
		}
	},
}

// unsetCmd represents the unset command
var unsetCmd = &cobra.Command{
	Use:   "unset [variable]",
	Short: "Unset a variable",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("unset called")
	},
}

var variablesCmd = &cobra.Command{
	Use:   "variables",
	Short: "List all set variables",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(unsetCmd)
	rootCmd.AddCommand(variablesCmd)
}
