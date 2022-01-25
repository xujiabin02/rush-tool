package rushtool

func StrCopy(inputStr string, num int) []interface{} {
	var b []interface{}
	for i := 0; i < num; i++ {
		b = append(b, inputStr)
	}
	return b

}
