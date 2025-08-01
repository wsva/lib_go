package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//https://pro.ant.design/zh-CN/docs/request

type ResponseData struct {
	List     any `json:"list,omitempty"`
	Current  int `json:"current,omitempty"`
	PageSize int `json:"pageSize,omitempty"`
	Total    int `json:"total,omitempty"`
}

type Response struct {
	// if request is success
	Success bool `json:"success,omitempty"`

	// response data
	Data ResponseData `json:"data,omitempty"`

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

func (r *Response) DoResponse(w http.ResponseWriter) {
	jsonBytes, _ := json.Marshal(*r)
	w.Write(jsonBytes)
}

func RespondError(w http.ResponseWriter, err any) {
	resp := Response{
		Success: false,
		ErrMsg:  fmt.Sprint(err),
	}
	resp.DoResponse(w)
}

func RespondSuccess(w http.ResponseWriter) {
	resp := Response{
		Success: true,
	}
	resp.DoResponse(w)
}

// return HttpRespData.List
func ParseDataListFromResponse(respBytes []byte) (any, error) {
	var resp Response
	err := json.Unmarshal(respBytes, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.ErrMsg)
	}
	return resp.Data.List, nil
}
