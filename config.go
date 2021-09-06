package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var conf struct {
	Rescan int `yaml:"rescan"`
	Zabbix struct {
		Server  string `yaml:"server"`
		Active  bool   `yaml:"active"`
		Trapper bool   `yaml:"trapper"`
	} `yaml:"zabbix"`
	RBL struct {
		IPv4 []string `yaml:"ipv4"`
	} `yaml:"rbl"`
	IP []string `yaml:"ip"`
}

func confLoad() {
	c, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	yaml.Unmarshal(c, &conf)
}
