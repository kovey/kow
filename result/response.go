package result

type Response struct {
	Code Codes  `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
