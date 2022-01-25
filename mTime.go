package rushtool

import "time"

func GetNowStr() string {
	return time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")

}

type StrNow struct {
	name    string
	Str     string
	TFormat string
	Unix    int64
	OneDay  int64
	Proper  []StrTime
}
type StrTime struct {
	name  string
	value int64
}

func (sn *StrNow) Init() {
	sn.OneDay = 86400
	if sn.TFormat == "" {
		sn.TFormat = "2006-01-02 15:04:05"
	}
	sn.Unix = time.Now().Unix()
	sn.Str = time.Unix(sn.Unix, 0).Format(sn.TFormat)
	sn.Proper = []StrTime{
		{"S", 1},
		{"M", 60},
		{"H", 3600},
		{"d", 86400},
		{"w", 604800},
		{"m", 18144000},
		{"Y", 31536000},
	}
}
func (sn *StrNow) Past(num int64, strTime string) string {
	result := ""
	for _, k := range sn.Proper {
		if k.name == strTime {
			result = time.Unix(sn.Unix+(num*k.value), 0).Format(sn.TFormat)
		}
	}
	return result
}
