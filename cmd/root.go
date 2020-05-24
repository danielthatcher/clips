package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var profile string
var profileDir string
var varArgs []string
var variables map[string]string
var copy bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bline [template]",
	Args:  cobra.ExactArgs(1),
	Short: "Generate one-liners from templates",
	Run: func(cmd *cobra.Command, args []string) {
		tpl := args[0]

		// Check the requested template exists
		templates, err := GetTemplates()
		if err != nil {
			log.Fatalf("Failed to get templates: %v\n", err)
		}
		exists := false
		for _, t := range templates {
			if t == tpl {
				exists = true
				break
			}
		}
		if !exists {
			log.Fatalf("Template %s doesn't exist\n", tpl)
		}

		// Fill in the template
		template, err := GetTemplate(tpl)
		if err != nil {
			log.Fatalf("Failed to load template %s: %v\n", tpl, err)
		}
		line := template.Line
		for repl, v := range template.Variables {
			line = strings.ReplaceAll(line, repl, variables[v])
		}

		// Output to user
		fmt.Println(line)

		// Write to clipboard if requested
		if copy {
			clipboard.WriteAll(line)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "", "config file (default is $HOME/.config/bline/config.yaml)")
	rootCmd.PersistentFlags().StringP("profile", "p", "", "profile name")

	// root command flags
	rootCmd.Flags().StringSliceVarP(&varArgs, "set", "s", make([]string, 0), "set a variable using varname=value (can be specified multiple times)")
	rootCmd.Flags().BoolVarP(&copy, "copy", "c", false, "copy command to the clipboard")

	// Bind to viper
	viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))

	// Defaults
	viper.SetDefault("profile", "default")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("json")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		configDir = path.Join(home, ".config", "bline")

		// Create the config dir if it doesn't exist
		err = os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to make config directory: %v\n", err)
		}

		// Create the default profile dir if no profiles exist
		profiles, err := GetProfiles()
		if err != nil {
			log.Fatalf("Failed to get profiles: %v\n", err)
		}
		if len(profiles) == 0 {
			d := path.Join(configDir, "default")
			err = os.MkdirAll(d, os.ModePerm)
			if err != nil {
				log.Fatalf("Failed to created default profile directory: %v\n", err)
			}
		}

		viper.AddConfigPath(configDir)
		viper.SetConfigName("config.json")
		cfgFile = path.Join(configDir, "config.json")
	}

	// If a config file is found, read it in.
	viper.ReadInConfig()

	// Check the current profile is valid
	profile = viper.Get("profile").(string)
	profiles, err := GetProfiles()
	if err != nil {
		log.Fatalf("Failed to get profiles: %v\n", err)
	}
	valid := false
	for _, p := range profiles {
		if profile == p {
			valid = true
			break
		}
	}
	if !valid {
		log.Fatalf("Unknown profile: %s\n", profile)
	}

	profileDir = path.Join(configDir, profile)

	// Process the passed variables, overwriting values from the config file with those specified on the command line
	variables = GetConfigVariables()
	for _, s := range varArgs {
		split := strings.SplitN(s, "=", 2)
		if len(split) != 2 {
			log.Fatalf("Error processing variable setting '%s'. Use format 'var=value'\n", s)
		}
		variables[split[0]] = split[1]
	}
}
