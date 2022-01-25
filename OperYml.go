package rushtool

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)


func VarsFromYml(filename string) map[string]interface{} {
	yamlFile,err:=ioutil.ReadFile(filename)
	resultMap:=make(map[string]interface{})
	if err!=nil{
		fmt.Println(err)
	}
	err = yaml.Unmarshal(yamlFile, &resultMap)
	return resultMap
}
