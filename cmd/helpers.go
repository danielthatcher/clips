package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/spf13/viper"
)

type Template struct {
	Line      string            `json:"line"`
	Variables map[string]string `json:"variables"`
	Defaults  map[string]string `json:"defaults"`
}

type VariablesConfig struct {
	Variables map[string]string `json:"variables"`
}

// GetProfiles returns a list of all the profiles
func GetProfiles() ([]string, error) {
	files, err := ioutil.ReadDir(configDir)
	if err != nil {
		return nil, err
	}

	profiles := make([]string, 0)
	for _, f := range files {
		if f.IsDir() {
			profiles = append(profiles, f.Name())
		}
	}
	return profiles, nil
}

// GetTemplates returns a list of all the templates in the current profile
func GetTemplates() ([]string, error) {
	g := path.Join(profileDir, "*.json")
	matches, err := filepath.Glob(g)
	if err != nil {
		return nil, err
	}

	templates := make([]string, 0)
	for _, f := range matches {
		ext := filepath.Ext(f)
		t := f[0 : len(f)-len(ext)]
		templates = append(templates, filepath.Base(t))
	}
	return templates, nil
}

// GetTemplate returns a Template object loaded using the specified template
func GetTemplate(t string) (*Template, error) {
	f := path.Join(profileDir, t+".json")
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	tpl := &Template{}
	json.Unmarshal(b, tpl)
	return tpl, nil

}

// GetConfigVariables gets the variables set in the configuration file
func GetConfigVariables() map[string]string {
	// Unfortunately, viper is case insensitive, so we have to load the config file manually
	b, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return make(map[string]string, 0)
	}
	vc := &VariablesConfig{}
	json.Unmarshal(b, vc)
	return vc.Variables
}

// SaveVariables saves the variables the config file, respecting case
func SaveVariables() error {
	viper.Set("variables", variables)
	err := viper.WriteConfigAs(cfgFile)
	return err
}

// EditFile opens the specified in an editor, respecting $EDITOR if it is set, otherwise using vim
func EditFile(f string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	c := exec.Command(editor, f)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
