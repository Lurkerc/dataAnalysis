package utils

type JsonData struct {
	Data  string `json:"data"`
	State State  `json:"state"`
}

type State struct {
	Msg     string `json:"msg"`
	Code    int64  `json:"code"`
	Success bool   `json:"success"`
}
