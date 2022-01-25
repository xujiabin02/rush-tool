package rushtool

import "strconv"

func TriggerCondition(fileName string, threshold int) bool {
	//const threshold = 3 //阈值
	result := false
	if readNum, err := strconv.Atoi(ReadFile(fileName)); err != nil {
		panic(err)
	} else {
		if readNum == 0 || readNum < threshold {
			readNum++
			WriteFile(fileName, strconv.Itoa(readNum))
			return result
		}
		if readNum >= threshold {
			WriteFile(fileName, strconv.Itoa(0))
			result = true
			return result
		}
	}
	return result
}
