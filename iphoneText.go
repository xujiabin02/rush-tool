package rushtool

import (
	"github.com/tidwall/sjson"
)

func IphoneText(apiUrl string, personalMobile []string, iText string) string {
	iData := `{
    "mobiles":[
    ],
    "context":"test for alarm"
}`
	iData, _ = sjson.Set(iData, "mobiles", personalMobile)
	iData, _ = sjson.Set(iData, "context", iText)
	result := Post(apiUrl, iData)
	//fmt.Println(result)
	return result
}
