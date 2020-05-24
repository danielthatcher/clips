package cmd

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"path/filepath"
)

type Template struct {
	Line      string            `json:"line"`
	Variables map[string]string `json:"variables"`
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