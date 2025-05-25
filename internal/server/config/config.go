// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/grd/FreePDM/internal/util"
)

var (
	AppDir                string
	configName, configDir string
	Conf                  = Config{}
)

type Config struct {
	VaultsDirectory string
	LogFile         string
	LogLevel        string
	Users           map[string]int
}

// GetUid returns the uid for a given user name or -1 if not found.
func GetUid(name string) int {
	if uid, ok := Conf.Users[name]; ok {
		return uid
	}
	return -1
}

// ReadConfig reads the configuration file into Conf, handling errors appropriately.
func ReadConfig() error {
	if !util.FileExists(configName) {
		return fmt.Errorf("config file %s does not exist", configName)
	}

	_, err := toml.DecodeFile(configName, &Conf)
	if err != nil {
		return fmt.Errorf("error decoding config file %s: %v", configName, err)
	}

	if len(Conf.Users) == 0 {
		log.Printf("Warning: No users found in configuration file %s", configName)
	}

	if !util.DirExists(Conf.VaultsDirectory) {
		log.Fatal("vaults directory doesn't exist. See the installation manual")
	}

	return nil
}

// WriteConfig writes the current Conf structure to the config file.
func WriteConfig() error {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(&Conf)
	if err != nil {
		return fmt.Errorf("error encoding config to TOML: %v", err)
	}

	if err := os.WriteFile(configName, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing config file %s: %v", configName, err)
	}

	return nil
}

// String returns the configuration as a formatted string.
func (cfg *Config) String() string {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(cfg); err != nil {
		log.Fatal(err) // Here it’s fine to be fatal as it's just the String method.
	}
	return buf.String()
}

func init() {
	AppDir, ok := os.LookupEnv("FREEPDM_DIR")
	if !ok {
		log.Fatal("The environment FREEPDM_DIR is not set. Please read the installation page.")
	}
	configDir = path.Join(AppDir, "data")
	configName = path.Join(configDir, "FreePDM.toml")

	// Ensure the config directory exists
	if !util.DirExists(configDir) {
		if err := os.Mkdir(configDir, 0700); err != nil {
			log.Fatalf("error creating config directory %s: %v", configDir, err)
		}
	}

	// Create a new config file if it doesn't exist
	if !util.FileExists(configName) {
		if err := WriteConfig(); err != nil {
			log.Fatalf("error creating initial config file %s: %v", configName, err)
		}
	}

	// Read the configuration file
	if err := ReadConfig(); err != nil {
		log.Fatalf("error reading config file: %v", err)
	}
}
