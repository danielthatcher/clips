package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/spf13/cobra"
)

var configDir string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all templates in the current profile",
	Run: func(cmd *cobra.Command, args []string) {
		templates, err := GetTemplates()
		if err != nil {
			log.Fatalf("Failed to get templates: %v\n", err)
		}
		for _, t := range templates {
			fmt.Println(t)
		}
	},
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Add a new template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tpl := args[0]

		// Check the template doesn't exist
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
		if exists {
			log.Fatalf("Template %s alread exists\n", tpl)
		}

		// Add the new template in the correct file
		f := path.Join(profileDir, tpl+".json")
		base := `{
	"line": "",
	"variables": {
	}
}
`
		err = ioutil.WriteFile(f, []byte(base), 0644)
		if err != nil {
			log.Fatalf("Failed to write file %s: %v\n", f, err)
		}

		// Print template file
		fmt.Printf("Create new template %s in %s\n", tpl, f)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(newCmd)
}
