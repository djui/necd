package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/vaughan0/go-ini"
)

func readConf() (ini.File, error) {
	confPath := confPath()
	conf, err := ini.LoadFile(confPath)
	if os.IsNotExist(err) {
		if err := ioutil.WriteFile(confPath, []byte{}, 0644); err != nil {
			return nil, err
		}
		return readConf()
	}

	if err != nil {
		return nil, err
	}

	return conf, nil
}

func confPath() string {
	// .
	confPath := ".necdrc"
	if _, err := os.Stat(confPath); err == nil {
		return confPath
	}

	// CWD
	wDir, err := os.Getwd()
	if err == nil {
		confPath := path.Join(wDir, ".necdrc")
		if _, err := os.Stat(confPath); err == nil {
			return confPath
		}
	}

	// $HOME
	usr, err := user.Current()
	if err == nil {
		confPath := path.Join(usr.HomeDir, ".necdrc")
		return confPath
	}

	// Fallback
	return ""
}
