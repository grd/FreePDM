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
	extras "github.com/grd/FreePDM/src/utils"
)

var (
	appName    = "FreePDM"
	configDir  = path.Join("/vault", ".config")
	configName = path.Join(configDir, appName+".toml")
	Conf       = Config{}
)

type Config struct {
	// show_fc_files_only      = 1
	// hide_versioned_fc_files = 2
	StartupDirectory string
	LogFile          string
	LogLevel         string
	// Users            []UserUid
	Users map[string]int
}

// type UserUid struct {
// 	UserName string
// 	Uid      int
// }

// Returns the uid from a name or when not found -1
func GetUid(name string) int {
	for k, v := range Conf.Users {
		if name == k {
			return v
		}
	}
	return -1
}

func ReadConfig() {
	_, err := toml.DecodeFile(configName, &Conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func WriteConfig() {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(&Conf)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(configName, buf.Bytes(), 0644)
	extras.CheckErr(err)

}

func (self *Config) String() string {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(self)
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

func init() {
	// create the new directory if it doesn't exist
	if extras.DirExists(configDir) == false {
		os.Mkdir(configDir, 0700)
	}

	// create a new config file when it doesn't exist
	if extras.FileExists(configName) == false {
		WriteConfig()
	}

	// Reading the configuration file
	ReadConfig()
}
