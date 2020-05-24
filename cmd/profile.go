package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addprofileCmd represents the addprofile command
var addprofileCmd = &cobra.Command{
	Use:   "addprofile",
	Short: "Add a new profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		d := path.Join(configDir, args[0])
		err := os.Mkdir(d, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create profile directory: %v\n", err)
		}
	},
}

// removeprofileCmd represents the removeprofile command
var removeprofileCmd = &cobra.Command{
	Use:   "removeprofile",
	Short: "remove a profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		d := path.Join(configDir, args[0])
		err := os.RemoveAll(d)
		if err != nil {
			log.Fatalf("Failed to remove profile directory: %v\n", err)
		}
	},
}

// listProfilesCmd represents the listprofiles command
var listProfilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Show all available profiles",
	Run: func(cmd *cobra.Command, args []string) {
		profiles, err := GetProfiles()
		if err != nil {
			log.Fatalf("Failed to get profiles: %v\n", err)
		}
		for _, p := range profiles {
			fmt.Println(p)
		}
	},
}

// setprofileCmd sets the profile in use
var setprofileCmd = &cobra.Command{
	Use:   "setprofile",
	Short: "Set the profile to use",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		p := args[0]

		// Check the profile exists
		profiles, err := GetProfiles()
		if err != nil {
			log.Fatalf("Failed to get profiles: %v\n", err)
		}
		exists := false
		for _, pr := range profiles {
			if p == pr {
				exists = true
				break
			}
		}
		if !exists {
			log.Fatalf("Profile %s doesn't exist\n", p)
		}

		// Write the profile change to the config file
		viper.Set("profile", p)
		err = viper.WriteConfigAs(cfgFile)
		if err != nil {
			log.Fatalf("Failed write configuration to %s: %v\n", cfgFile, err)
		}
	},
}

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Print the current profile",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(profile)
	},
}

func init() {
	rootCmd.AddCommand(addprofileCmd)
	rootCmd.AddCommand(removeprofileCmd)
	rootCmd.AddCommand(listProfilesCmd)
	rootCmd.AddCommand(setprofileCmd)
	rootCmd.AddCommand(profileCmd)
}
