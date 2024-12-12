package context

const (
	Content_Type_Key       = "Content-Type"
	Content_Type_Json      = "application/json"
	Content_Type_Xml       = "text/xml"
	Content_Type_Html      = "text/html"
	Content_Type_Form      = "application/x-www-form-urlencoded"
	Content_Type_Binary    = "application/octet-stream"
	Header_X_Real_Ip       = "X-Real-IP"
	Header_X_Forwarded_For = "X-Forwarded-For"
	Header_X_Request_Id    = "X-Request-Id"
)

func IsJson(t string) bool {
	return t == Content_Type_Json
}

func IsForm(t string) bool {
	return t == Content_Type_Form
}

func IsBinary(t string) bool {
	return t == Content_Type_Binary
}

func IsXml(t string) bool {
	return t == Content_Type_Xml
}
