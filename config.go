package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/vaughan0/go-ini"
)

func readConf() (ini.File, error) {
	confPath, err := confPath()
	if err != nil {
		return nil, err
	}

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

func writeConf(conf string) error {
	confPath, err := confPath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(confPath, []byte(conf), 0644)
	return err
}

func confPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	confPath := path.Join(usr.HomeDir, ".necdrc")
	return confPath, nil
}
