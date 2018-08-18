package controllers

import (
	"github.com/astaxie/beego"
	"dataAnalysis/utils"
	"dataAnalysis/models"
	"encoding/json"
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
	state := utils.State{}
	dataArr, err := models.GetDateByFolder()
	if err != nil {
		state.Msg = "失败"
		state.Code = 1
		state.Success = false
	} else {
		state.Msg = "成功"
		state.Code = 0
		state.Success = true
	}
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
	state := utils.State{}
	tableDataObj, err := models.GetTableData(paramsData)
	if err != nil {
		state.Msg = "失败"
		state.Code = 1
		state.Success = false
	} else {
		state.Msg = "成功"
		state.Code = 0
		state.Success = true
	}
	jsonData := tableData{tableDataObj, state}
	e.Data["json"] = jsonData
	e.ServeJSON()
}
