package rushtool

import (
	"fmt"
	"strings"
)

func expendSlice(someStringSlice []string, largeThan int) []string {
	e := make([]string, largeThan)
	d := append(someStringSlice, e...)
	return d
}

func transposeTwoDimensionArray(inputSlice [][]string) [][]string {
	var b int
	var c, d [][]string
	for _, i := range inputSlice {
		//fmt.Println(len(i))
		if len(i) > b {
			*&b = len(i)
		}
	}
	for _, i := range inputSlice {
		c = append(c, expendSlice(i, b-len(i)))
	}
	for i := 0; i < b; i++ {
		var e []string
		for j := 0; j < len(c); j++ {
			e = append(e, c[j][i])
		}
		d = append(d, e)
	}
	return d

}
func maxLength(inputSlice []string) int {
	var maxLen int
	for _, i := range inputSlice {
		iRune := []rune(i)
		currentLen := len(iRune) + extractChineseNum(i)
		if currentLen > maxLen {
			*&maxLen = currentLen
		}
	}
	return maxLen
}
func addNumStr(someString string, someNumber int) string {
	exNum := strings.Count(someString, "_")
	for i := 0; i < someNumber; i++ {
		someString = someString + "  "
	}
	for i := 0; i < exNum; i++ {
		someString = someString + " "
	}
	return someString
}

func TurnSliceToStringAddFrame(inputSlice [][]string) string {
	tranSlice := transposeTwoDimensionArray(inputSlice)
	//fmt.Println(a)
	for _, i := range tranSlice {
		max := maxLength(i)
		for j := 0; j < len(i); j++ {
			if len(i[j]) < max {
				i[j] = addNumStr(i[j], max-len(i[j]))
			}
		}
	}
	//fmt.Println(tranSlice)
	d := transposeTwoDimensionArray(tranSlice)
	var e string
	var c int
	for n, i := range d {
		switch n {
		case 0:
			e = e + strings.Join(i, " | ") + "\n"
			*&c = len(e)
			e = "\n" + e
			for j := 0; j < c; j++ {
				e = "-" + e
			}
			for j := 0; j < c; j++ {
				e += "-"
			}
			e += "\n"
		case len(d) - 1:
			e = e + strings.Join(i, " | ") + "\n"
			for j := 0; j < c; j++ {
				e += "-"
			}
			e += "\n"
		default:
			e = e + strings.Join(i, " | ") + "\n"
		}
	}
	return e
}

func extractChineseNum(someStr string) int {
	runeStr := []rune(someStr)
	num := 0
	for i := 0; i < len(runeStr); i++ {
		charNum := runeStr[i]
		if charNum <= 40869 && charNum >= 19968 {
			num++
		}
	}
	return num
}
func TurnSliceToHTML(inputSlice [][]string) string {
	result := `<!-- CSS goes in the document HEAD or added to your external stylesheet -->
<style type="text/css">
table.gridtable {
	font-family: verdana,arial,sans-serif;
	font-size:11px;
	color:#333333;
	border-width: 1px;
	border-color: #666666;
	border-collapse: collapse;
}
table.gridtable th {
	border-width: 1px;
	padding: 8px;
	border-style: solid;
	border-color: #666666;
	background-color: #dedede;
}
table.gridtable td {
	border-width: 1px;
	padding: 8px;
	border-style: solid;
	border-color: #666666;
	background-color: #ffffff;
}
</style>

<!-- Table goes in the document BODY -->`
	for n, i := range inputSlice {
		if n == 0 {
			for _, j := range i {

				result = result + fmt.Sprintf("<th>%s</th>", j)
			}
			result = fmt.Sprintf("<tr>%s</tr>", result)
		} else {

			for _, j := range i {
				result = result + fmt.Sprintf("<td>%s</td>", j)
			}
			result = fmt.Sprintf("<tr>%s</tr>", result)
		}
	}
	result = fmt.Sprintf(`<table class="gridtable">%s</table>`, result)
	return result
}
