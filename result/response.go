package result

type Response struct {
	Code Codes  `json:"code" xml:"code"`
	Msg  string `json:"msg" xml:"msg"`
	Data any    `json:"data" xml:"data"`
}

type ResponseErrForm struct {
	Code Codes  `form:"code"`
	Msg  string `form:"msg"`
}
