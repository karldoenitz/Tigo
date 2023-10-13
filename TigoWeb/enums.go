package TigoWeb

const Version = "1.6.5"

// MethodEnum 根据http请求方式获取对应的函数名
//   - httpMethod: http请求方式
func MethodEnum(httpMethod string) string {
	var MethodMapping = map[string]string{
		"get":     "Get",
		"head":    "Head",
		"post":    "Post",
		"put":     "Put",
		"delete":  "Delete",
		"connect": "Connect",
		"options": "Options",
		"trace":   "Trace",
		"GET":     "Get",
		"HEAD":    "Head",
		"POST":    "Post",
		"PUT":     "Put",
		"DELETE":  "Delete",
		"CONNECT": "Connect",
		"OPTIONS": "Options",
		"TRACE":   "Trace",
	}
	return MethodMapping[httpMethod]
}
