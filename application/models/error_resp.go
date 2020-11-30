package models

type RespError struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

type RespOk struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body,omitempty"`
}