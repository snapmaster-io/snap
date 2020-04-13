package config

import (
	"io/ioutil"
	"os"
)

// ReadConfigFile reads config data from filename
func ReadConfigFile(filename string, data []byte) ([]byte, error) {
	// construct the config directory path as $HOME/.config/snap/
	homedir, err := os.UserHomeDir()
	if err != nil {
		return []byte(""), err
	}
	configPath := homedir + "/.config/snap"

	// construct the config file path
	configFile := configPath + "/" + filename
	contents, err := ioutil.ReadFile(configFile)
	if err != nil {
		return []byte(""), err
	}

	return contents, nil
}

// WriteConfigFile writes data into filename in the config directory, and returns filename
func WriteConfigFile(filename string, data []byte) (string, error) {
	// construct the config directory path as $HOME/.config/snap/
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := homedir + "/.config/snap"
	err = os.MkdirAll(configPath, 0755)
	if err != nil {
		return "", err
	}

	// construct the config file path
	configFile := configPath + "/" + filename
	err = ioutil.WriteFile(configFile, data, 0600)
	if err != nil {
		return "", err
	}

	return configFile, nil
}
