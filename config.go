package main

import (
	"flag"
	"gopkg.in/ini.v1"
	"path/filepath"
)

type Configuration struct {

	// Specify bind address for web socket
	BindAddress string `ini:"bind"`

	// Allowed origins
	Origins []string `ini:"origins"`

	// Allowed ips
	Allow []string `ini:"allow"`
}

var config = new(Configuration)

func configSetup() {

	var configFileFlag = flag.String("config", "/etc/omeetings/ows.ini", "Config file")

	var bindFlag = flag.String("bind", "", "Bind address")

	flag.Parse()

	configFile, err := filepath.Abs(*configFileFlag)
	onError(err, "Fail to get absolute path from ", *configFileFlag)
	ini, err := ini.Load(configFile)
	onError(err, "Config file not exist", configFile)

	ini.Section("").MapTo(config)

	if *bindFlag != "" {
		config.BindAddress = *bindFlag
	}
}
