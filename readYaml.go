package rushtool

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Book struct {
	C string `OpYml:"c"`
	D string `OpYml:"d"`
}
type My struct {
	Book `OpYml:",inline"`
	A string `OpYml:"a"`
	B string `OpYml:"b"`
}
func LoadYaml(src string) My {
	content, err:=ioutil.ReadFile(src)
	if err!=nil{
		fmt.Println(err)
	}
	c:=My{}
	err=yaml.Unmarshal(content, &c)
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(c)
	return c

}

func PrintYaml(f My) {
	b,err:=yaml.Marshal(&f)
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(string(b))
}
