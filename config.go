package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	yaml "gopkg.in/yaml.v2"
)

func parseFlags() *Configuration {
	var config Configuration
	var filename string
	flag.StringVar(&filename, "conf", "", "Filename to read configuration from")
	flag.StringVar(&config.ServerRoot, "serverRoot", "", "The base URL of the IIS Anycast server. e.g. https://api.anycast.iis.se")
	flag.StringVar(&config.ApiUser, "apiUser", "", "The API user name")
	flag.StringVar(&config.ApiSecret, "apiSecret", "", "The API user secret")
	flag.StringVar(&config.ApiPasswd, "apiPasswd", "", "The API user password")
	flag.StringVar(&config.InfluxServer, "influxServer", "", "Server with InfluxDB running")
	flag.StringVar(&config.InfluxPort, "influxPort", "", "Port used by InfluxDB")
	flag.StringVar(&config.InfluxDB, "influxDB", "", "Name of InfluxDB database")
	flag.Parse()

	var confFromFile *Configuration
	if filename != "" {
		var err error
		confFromFile, err = readConfigFile(filename)
		if err != nil {
			panic(err)
		}
	}
	return joinConfig(confFromFile, &config)
}

func readConfigFile(filename string) (config *Configuration, error error) {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	source, err := ioutil.ReadFile(path.Join(usr.HomeDir, ".dzone"))
	if err != nil {
		return nil, err
	}
	config = &Configuration{}
	err = yaml.Unmarshal(source, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func readDefaultConfigFiles() (config *Configuration) {

	// .dzone in current directory
	fileconfig, err := readConfigFile(".dzone")
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	config = joinConfig(config, fileconfig)

	// .dzone in user home directory
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	fileconfig, err = readConfigFile(path.Join(usr.HomeDir, ".dzone"))
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	config = joinConfig(config, fileconfig)

	// done
	return
}

func joinConfig(oldConf *Configuration, newConf *Configuration) (config *Configuration) {
	if oldConf == nil && newConf == nil {
		return nil
	}
	if oldConf != nil && newConf == nil {
		return oldConf
	}
	if oldConf == nil && newConf != nil {
		return newConf
	}

	// we have two configs, join them
	config = &Configuration{}
	if newConf.ServerRoot != "" {
		config.ServerRoot = newConf.ServerRoot
	} else {
		config.ServerRoot = oldConf.ServerRoot
	}
	if newConf.ApiUser != "" {
		config.ApiUser = newConf.ApiUser
	} else {
		config.ApiUser = oldConf.ApiUser
	}
	if newConf.ApiSecret != "" {
		config.ApiSecret = newConf.ApiSecret
	} else {
		config.ApiSecret = oldConf.ApiSecret
	}
	if newConf.ApiPasswd != "" {
		config.ApiPasswd = newConf.ApiPasswd
	} else {
		config.ApiPasswd = oldConf.ApiPasswd
	}
	if newConf.InfluxServer != "" {
		config.InfluxServer = newConf.InfluxServer
	} else {
		config.InfluxServer = oldConf.InfluxServer
	}
	if newConf.InfluxPort != "" {
		config.InfluxPort = newConf.InfluxPort
	} else {
		config.InfluxPort = oldConf.InfluxPort
	}
	if newConf.InfluxDB != "" {
		config.InfluxDB = newConf.InfluxDB
	} else {
		config.InfluxDB = oldConf.InfluxDB
	}

	// Done
	return config
}
