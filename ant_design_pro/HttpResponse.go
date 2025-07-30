package ant_design_pro

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//https://pro.ant.design/zh-CN/docs/request

type HttpRespData struct {
	List     interface{} `json:"list,omitempty"`
	Current  int         `json:"current,omitempty"`
	PageSize int         `json:"pageSize,omitempty"`
	Total    int         `json:"total,omitempty"`
}

type HttpResp struct {
	// if request is success
	Success bool `json:"success,omitempty"`

	// response data
	Data HttpRespData `json:"data,omitempty"`

	// code for errorType
	ErrorCode string `json:"errorCode,omitempty"`

	// message display to user
	ErrMsg string `json:"errMsg,omitempty"`

	// error display type：
	//0 silent; 1 message.warn; 2 message.error; 4 notification; 9 page
	ShowType int `json:"showType,omitempty"`

	// Convenient for back-end Troubleshooting: unique request ID
	TraceId string `json:"traceId,omitempty"`

	// onvenient for backend Troubleshooting: host of current access server
	Host string `json:"host,omitempty"`
}

func (a *HttpResp) DoResponse(w http.ResponseWriter) {
	jsonBytes, _ := json.Marshal(*a)
	w.Write(jsonBytes)
}

func RespondError(w http.ResponseWriter, err interface{}) {
	resp := HttpResp{
		Success: false,
		ErrMsg:  fmt.Sprint(err),
	}
	resp.DoResponse(w)
}

func RespondSuccess(w http.ResponseWriter) {
	resp := HttpResp{
		Success: true,
	}
	resp.DoResponse(w)
}

// return HttpRespData.List
func ParseDataListFromHttpResp(respBytes []byte) (interface{}, error) {
	var resp HttpResp
	err := json.Unmarshal(respBytes, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.ErrMsg)
	}
	return resp.Data.List, nil
}
