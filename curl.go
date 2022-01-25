package rushtool

import (
	"bytes"
	"crypto/tls"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

func CurlG(iUrl string, iHeader gjson.Result) string {
	client := http.Client{}
	reqID, err := http.NewRequest("GET", iUrl, nil)
	CheckErr(err)
	for k, v := range iHeader.Map() {
		reqID.Header.Add(k, v.Str)
	}
	respID, err := client.Do(reqID)
	CheckErr(err)
	bodyID, err := ioutil.ReadAll(respID.Body)
	CheckErr(err)
	return string(bodyID)
}
func CurlP(iUrl string, iData string, iHeader gjson.Result) string {
	jsonStr := bytes.NewBuffer([]byte(iData))
	req, err := http.NewRequest("POST", iUrl, jsonStr)
	CheckErr(err)
	BodyType := "application/json;charset=utf-8"
	req.Header.Set("Content-Type", BodyType)
	for k, v := range iHeader.Map() {
		req.Header.Add(k, v.Str)
	}
	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	//http cookie接口
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{tr, nil, cookieJar, 0}
	resp, err := client.Do(req)
	CheckErr(err)
	body, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	repoB := string(body)
	return repoB
}

func CheckErr(err interface{}) interface{} {
	if err != nil {
		panic(err)
	}
	return nil
}
func Get(iUrl string) string {
	client := http.Client{}
	if req, err := http.NewRequest("GET", iUrl, nil); err == nil {
		if resp, err := client.Do(req); err == nil {
			if body, err := ioutil.ReadAll(resp.Body); err == nil {
				return string(body)
			}
		}
	}
	return ""
}
func Gets(iUrl, iToken string) string {
	client := http.Client{}
	reqID, err := http.NewRequest("GET", iUrl, nil)
	CheckErr(err)
	reqID.Header.Set("Authorization", iToken)
	respID, err := client.Do(reqID)
	CheckErr(err)
	bodyID, err := ioutil.ReadAll(respID.Body)
	CheckErr(err)
	return string(bodyID)
}
func Post(iUrl string, iData string) string {
	jsonStr := bytes.NewBuffer([]byte(iData))
	req, err := http.NewRequest("POST", iUrl, jsonStr)
	CheckErr(err)
	BodyType := "application/json;charset=utf-8"
	req.Header.Set("Content-Type", BodyType)
	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	//http cookie接口
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{tr, nil, cookieJar, 0}
	resp, err := client.Do(req)
	CheckErr(err)
	body, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	repoB := string(body)
	return repoB

}
func IGets(iUrl string, coo map[string]string) string {
	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{tr, nil, cookieJar, 0}
	//client := http.Client{}
	reqID, err := http.NewRequest("GET", iUrl, nil)
	CheckErr(err)
	for k, v := range coo {
		reqID.Header.Add(k, v)
	}
	//reqID.Header.Add("Content-Type", "application/json;charset=utf-8")
	//reqID.Header.Add("X-Auth-Token", iToken)
	respID, err := client.Do(reqID)
	//h,_:=json.Marshal(respID.Header)
	//fmt.Println(string(h))
	CheckErr(err)
	bodyID, err := ioutil.ReadAll(respID.Body)
	CheckErr(err)
	return string(bodyID)
}
func StrToCookies(cookieStr string) map[string]string {
	cooDict := make(map[string]string)
	for n, i := range strings.Split(cookieStr, "\n") {
		if n == 0 {
			continue
		} else {
			//fmt.Println(i)
			a := strings.Split(i, ": ")
			k, v := a[0], a[1]
			cooDict[string(k)] = string(v)
		}
	}
	return cooDict
}
