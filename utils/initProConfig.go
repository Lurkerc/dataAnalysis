package utils

import (
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"os"
	"strings"
	"regexp"
)

type tableNameObj struct {
	TableName string `json:"tableName"` // 表名
	Text      string `json:"text"`      // 表中文名称
	Form      string `json:"form"`      // 全量/增量形式
}

func init() {
	logs.SetLogger(logs.AdapterFile, `{"filename":"log/project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
	logs.Info("系统启动")
	logs.Info("初始化")
	beego.LoadAppConfig("ini", "conf/static.conf")
	createClassNameByExcel()
}

// 根据excel创建分类名称
func createClassNameByExcel() {
	logs.SetLogger(logs.AdapterFile, `{"filename":"log/project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
	logs.Info("根据Excel解析生成表配置")
	xis, err := excelize.OpenFile("static/data/table.xlsx")
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Info("Excle文件读取完成")
	tableName := make(map[string]tableNameObj)
	// 获取工作表中sheet1中的值
	rows := xis.GetRows("Sheet1")
	re, _ := regexp.Compile("表$")
	for _, row := range rows[2:] {
		// 表明（英文）转换为小写
		tableNameKey := strings.ToLower(row[0])
		// 移除最后一个“表”字
		tableNameText := re.ReplaceAllString(row[1], "")
		tableName[tableNameKey] = tableNameObj{tableNameKey, tableNameText, strings.ToLower(row[2])}
	}
	// map 转 json
	tableNameStr, err := json.Marshal(tableName)
	if err != nil {
		logs.Error(err)
	}
	logs.Info("表配置解析完成，已经存放在static/config/tableName.json")
	// 存储 json
	jsonFile, err := os.Create("static/config/tableName.json")
	if err != nil {
		logs.Error(err)
	}
	defer jsonFile.Close()
	jsonFile.Write(tableNameStr)
	logs.Info("初始化完成")
}
