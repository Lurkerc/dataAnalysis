package utils

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type JsonData struct {
	Data  string `json:"data"`
	State State  `json:"state"`
}

type State struct {
	Msg     string `json:"msg"`
	Code    int64  `json:"code"`
	Success bool   `json:"success"`
}

// 获取状态数据
func GetState(err error, args State) State {
	var state State
	if err != nil {
		state.Msg = "失败"
		state.Code = 0
		state.Success = false
	} else {
		state.Msg = "成功"
		state.Code = 0
		state.Success = true
	}
	if args.Msg != "" {
		state.Msg = args.Msg
	}
	if args.Code != 0 {
		state.Code = args.Code
	}
	if beego.AppConfig.String("debug") == "true" && err != nil {
		logs.Error(err)
		state.Msg = err.Error()
	}
	return state
}
