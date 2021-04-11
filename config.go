// @Title  config.go
// @Description A highly interactive honeypot supporting redis protocol
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
