// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package workcache

import (
	"fmt"
	"log"
	"github.com/spf13/viper"
)

var configRead bool = false
var configFileName string = "config"

func SetConfigFile(configFile string) {
	configFileName = configFile
}

func readConfigIfNeeded() {
	if (configRead) { return }

	// set defaults
	viper.SetDefault("Main.ListenIpPort", ":7176")

	// read config file
	viper.SetConfigName(configFileName) // name of config file (without extension)
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.AddConfigPath("/")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	configRead = true
	log.Println("Config file has been read")
}

func ConfigGetString(keyName string) string {
	readConfigIfNeeded()
	return viper.GetString(keyName)
}
