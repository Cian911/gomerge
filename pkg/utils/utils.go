package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func ReadConfigFile(path string) {
	name, ext := parseConfigFile(path)
	viper.SetConfigName(name)

	// Set the configuration file type
	viper.SetConfigType(ext)

	// Set the path to look for the configurations file
	viper.AddConfigPath(filepath.Dir(path))

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Could not read config file: %s\n", err))
		os.Exit(1)
	}
}

func parseConfigFile(path string) (filename, ext string) {
	filename = filepath.Base(path)
	ext = filepath.Ext(filename)
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	ext = ext[1:len(ext)]

	return
}
