package pkg

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

func GetApplications(path string)  (apps []Application, err error){
    var data []byte
	apps = make([]Application, 0)
	if data, err = ioutil.ReadFile(path); err == nil {
		fmt.Println(string(data))
		err = yaml.Unmarshal(data, apps)
	}
	return
}
