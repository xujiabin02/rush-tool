package rushtool

//build by XuJiaBin

import (
	"github.com/tidwall/gjson"
	"strings"
)

func ClickHouseQuery(conHttp, SQL string) [][]string {
	var resultSlice [][]string
	result := Post(conHttp, SQL)
	for _, i := range strings.Split(result, "\n") {
		if i != "" {
			resultSlice = append(resultSlice, strings.Split(i, "\t"))
		}
	}
	return resultSlice
}

func JsonClickHouseQuery(conHttp, Sql string) map[string]interface{} {
	return gjson.Parse(Post(conHttp, Sql)).Value().(map[string]interface{})
}
