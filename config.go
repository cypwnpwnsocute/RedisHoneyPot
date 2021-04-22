// @Title  config.go
// @Description High Interaction Honeypot Solution for Redis protocol
// @Author  Cy 2021.04.08
package main

import "gopkg.in/ini.v1"

func LoadConfig(filename string) (cfg *ini.File, err error) {
	cfg, err = ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: true,
	}, filename)
	if err != nil {
		return nil, err
	} else {
		return cfg, nil
	}
}
