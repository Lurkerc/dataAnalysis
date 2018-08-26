package controllers

import (
	"dataAnalysis/models"
	"dataAnalysis/utils"
	"encoding/json"
	"github.com/astaxie/beego"
)

// Operations about object
type ExcelController struct {
	beego.Controller
}

type dateJsonData struct {
	Data  []string    `json:"data"`
	State utils.State `json:"state"`
}

type tableData struct {
	Data  models.TableDataObj `json:"data"`
	State utils.State         `json:"state"`
}

// @Title get data by folder
// @Description get data by folder in excel conf
// @Success 0 {jsonData}
// @Failure 1 {jsonData}
// @router /getDate [get]
func (e *ExcelController) GetDate() {
	dataArr, err := models.GetDateByFolder()
	state := utils.GetState(err, utils.State{"", 0, false})
	jsonData := dateJsonData{dataArr, state}
	e.Data["json"] = jsonData
	e.ServeJSON()
}

// @Title 获取表数据
// @Description 根据表名和日期获取数据
// @Success 0 {jsonData}
// @Failure 1 {jsonData}
// @router /getDataByTableName [post]
func (e *ExcelController) GetDataByTableName() {
	var paramsData models.GetDataParams
	json.Unmarshal(e.Ctx.Input.RequestBody, &paramsData)
	tableDataObj, err := models.GetTableData(paramsData)
	state := utils.GetState(err, utils.State{"", 0, false})
	jsonData := tableData{tableDataObj, state}
	e.Data["json"] = jsonData
	e.ServeJSON()
}

// @Title 获取表数据
// @Description 根据表名和日期获取数据
// @Success 0 {jsonData}
// @Failure 1 {jsonData}
// @router /createExcel [post]
func (e *ExcelController) GetExcelByTable() {
	var paramsData models.GetExcelParams
	json.Unmarshal(e.Ctx.Input.RequestBody, &paramsData)
	tablePath, err := models.CareatTableExcel(paramsData)
	state := utils.GetState(err, utils.State{"", 0, false})
	jsonData := utils.JsonData{tablePath, state}
	e.Data["json"] = jsonData
	e.ServeJSON()
}
