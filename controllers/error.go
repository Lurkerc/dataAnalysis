package controllers

import (
	"dataAnalysis/utils"
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

// 404
func (e *ErrorController) Error404() {
	jsonData := utils.JsonData{"", utils.State{"not found", 404, false,},}
	e.Data["json"] = jsonData
	e.ServeJSON()
}

// 501
func (e *ErrorController) Error501() {
	jsonData := utils.JsonData{"", utils.State{"server error", 501, false}}
	e.Data["json"] = jsonData
	e.ServeJSON()
}

// db connect error
func (e *ErrorController) ErrorDb() {
	jsonData := utils.JsonData{"", utils.State{"database is now down", 500, false}}
	e.Data["json"] = jsonData
	e.ServeJSON()
}
