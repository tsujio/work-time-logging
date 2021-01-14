package configuration

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

type ConfigFile struct {
	Spreadsheets []*struct {
		Id   string
		Name string
	}
}

type Config struct {
	ConfigFile
	Path string
	Dir  string
}

func Load(settingsDir string) *Config {
	configPath := filepath.Join(settingsDir, "config.json")
	f, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var config Config
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	config.Path = configPath
	config.Dir = settingsDir
	return &config
}

func (this *Config) FindSpreadsheetId(name string) (string, error) {
	for _, sheet := range this.Spreadsheets {
		if sheet.Name == name {
			return sheet.Id, nil
		}
	}
	return "", xerrors.Errorf("Spreadsheet not found: %s", name)
}
