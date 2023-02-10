package initialize

import (
	"io/ioutil"

	"study/internal/admin_api/global"

	"gopkg.in/yaml.v2"
)

func Config(filename string) {
	c, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(c, &global.ServerConfig)
	if err != nil {
		panic(err)
	}
}
