package rushtool

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/fs"
	"io/ioutil"
)

type DBModel struct {
	Name        string `yaml:"name"`
	Cron        string `yaml:"cron"`
	Description string `yaml:"description"`
	DB          string `yaml:"db"`
	Sql         string `yaml:"sql"`
}
type Prometheus struct {
	Name        string `yaml:"name"`
	Cron        string `yaml:"cron"`
	Description string `yaml:"description"`
	PromQL      string `yaml:"PromQL"`
}
type Loki struct {
	Name        string `yaml:"name"`
	Cron        string `yaml:"cron"`
	Description string `yaml:"description"`
	LogQL       string `yaml:"LogQL"`
}
type SliceDBModel struct {
	list []DBModel
}

type ReportModel struct {
	CH   []DBModel    `yaml:"clickhouse"`
	PG   []DBModel    `yaml:"postgresql"`
	PROM []Prometheus `yaml:"prometheus"`
	LO   []Loki       `yaml:"loki"`
}
type SrcReportModel struct {
	FILE        string `yaml:"src"`
	ReportModel `yaml:"data"`
}

func (d *DBModel) GetName() string {
	return d.Name
}
func (d *Prometheus) GetName() string {
	return d.Name
}
func (d *Loki) GetName() string {
	return d.Name
}

func PrintErr(err error) {
	if err != nil {
		Print(err)
	}
}
func CE(err error) bool {
	if err != nil {
		return true
	} else {
		return false
	}
}

func ReportCollector(src string) (SrcReportModel, error) {
	var resultStruct ReportModel
	var srcResultStruct SrcReportModel
	defer func() {
		if err := recover(); err != nil {
			Print(err)
		}
	}()
	yamlFile, err := ioutil.ReadFile(src)
	if CE(err) {
		return srcResultStruct, err
	}
	err = yaml.Unmarshal(yamlFile, &resultStruct)
	srcResultStruct.FILE = src
	srcResultStruct.ReportModel = resultStruct
	if CE(err) {
		return srcResultStruct, err
	}
	return srcResultStruct, nil
}

func XRangeFiles(x []fs.FileInfo) []string {
	xs := make([]string, 0)
	for _, file := range x {
		xs = append(xs, file.Name())
	}
	return xs
}

func Print(x ...interface{}) {
	fmt.Println(x...)
}
func ScanFiles(srcDir string, exclude []string) []string {
	xs := make([]string, 0)
	files, err := ioutil.ReadDir(srcDir)
	if CE(err) {
		Print(err)
	}
	for _, file := range XRangeFiles(files) {
		if InSliceString(file, exclude) {
			continue
		} else {
			xs = append(xs, fmt.Sprintf("%s/%s", srcDir, file))
		}
	}
	return xs
}

func InSliceString(src string, destSlice []string) bool {
	for _, obj := range destSlice {
		if src == obj {
			return true
		}
	}
	return false

}

type CallName interface {
	GetName() string
}

func (s *SrcReportModel) GetFileName() string {
	return s.FILE
}
func (s *SrcReportModel) GetSlice() ReportModel {
	return s.ReportModel
}

type SrcReportInterface interface {
	GetFileName() string
	GetSlice() ReportModel
}

func AppendDup(x CallName, objYaml SrcReportInterface, tmpDict map[string]string, dupDict map[string][]string) (map[string]string, map[string][]string) {
	if _, ok := tmpDict[x.GetName()]; ok {
		tmpDict[x.GetName()] = objYaml.GetFileName()
		dupDict[x.GetName()] = append(dupDict[x.GetName()], objYaml.GetFileName())
	} else {
		dupDict[x.GetName()] = append(dupDict[x.GetName()], objYaml.GetFileName())
	}
	return tmpDict, dupDict
}

//扫描重复key
func CheckDuplicate(keySrc ...SrcReportModel) map[string][]string {
	tmpDict := map[string]string{}
	dupDict := map[string][]string{}
	for _, objYaml := range keySrc {
		for _, i := range objYaml.ReportModel.CH {
			AppendDup(&i, &objYaml, tmpDict, dupDict)
		}
		for _, i := range objYaml.ReportModel.PG {
			AppendDup(&i, &objYaml, tmpDict, dupDict)
		}
		for _, i := range objYaml.ReportModel.PROM {
			AppendDup(&i, &objYaml, tmpDict, dupDict)
		}
		for _, i := range objYaml.ReportModel.LO {
			AppendDup(&i, &objYaml, tmpDict, dupDict)
		}
	}
	resultDict := map[string][]string{}
	for key, value := range dupDict {
		if len(value) > 1 {
			resultDict[key] = value
		}
	}
	return resultDict
}

func AppendReportModel(filesSrcReport ...SrcReportModel) ReportModel {
	result := ReportModel{}
	for _, i := range filesSrcReport {
		result.CH = append(result.CH, i.CH...)
		result.PG = append(result.PG, i.PG...)
		result.PROM = append(result.PROM, i.PROM...)
		result.LO = append(result.LO, i.LO...)
	}
	return result

}

func ParseToYaml(model ReportModel) string {
	result, err:=yaml.Marshal(model)
	if CE(err){
		Print(err)
	}
	return string(result)
}
func LoadYamlToMap(src string) map[string]string {
	resultMap:=make(map[string]string)
	yamlFile, err := ioutil.ReadFile(src)
	CheckErr(err)
	err = yaml.Unmarshal(yamlFile, &resultMap)
	return resultMap

}
