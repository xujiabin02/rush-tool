package rushtool

import "github.com/astaxie/beego/logs"

func LogConsole() *logs.BeeLogger {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	log.EnableFuncCallDepth(true)
	return log
}
