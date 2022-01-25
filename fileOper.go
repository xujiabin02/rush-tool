package rushtool

import (
	"io/ioutil"
	"os"
)

func ReadFile(fileName string) string {
	var resultReturn string
	//fileName := "./logconfig.json"
	if !checkFileIsExist(fileName) {
		WriteFile(fileName, "0")
	}
	if fileObj, err := os.Open(fileName); err == nil {
		defer fileObj.Close()
		if contents, err := ioutil.ReadAll(fileObj); err == nil {
			resultReturn = string(contents)
			//fmt.Println("Use os.Open family functions and ioutil.ReadAll to read a file contents:",result)
		}

	}
	return resultReturn
}
func WriteFile(fileName string, content string) string {
	statusStr := "OK"
	var d1 = []byte(content)
	err2 := ioutil.WriteFile(fileName, d1, 0666) //写入文件(字节数组)
	if err2 != nil {
		panic(err2)
	}
	return statusStr
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
