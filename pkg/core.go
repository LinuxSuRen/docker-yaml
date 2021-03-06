package pkg

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

func GetApplications(path string)  (apps []Application, err error){
    var data []byte
	//apps = make([]Application, 0)
	if data, err = ioutil.ReadFile(path); err == nil {
		err = yaml.Unmarshal(data, &apps)
	}
	return
}
