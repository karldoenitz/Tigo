package TigoWeb

const Version = "1.6.5"

const (
	httpLowerGet     = "get"
	httpLowerHead    = "head"
	httpLowerPost    = "post"
	httpLowerPut     = "put"
	httpLowerDelete  = "delete"
	httpLowerConnect = "connect"
	httpLowerOptions = "options"
	httpLowerTrace   = "trace"
	httpUpperGet     = "GET"
	httpUpperHead    = "HEAD"
	httpUpperPost    = "POST"
	httpUpperPut     = "PUT"
	httpUpperDelete  = "DELETE"
	httpUpperConnect = "CONNECT"
	httpUpperOptions = "OPTIONS"
	httpUpperTrace   = "TRACE"
)

// MethodEnum 根据http请求方式获取对应的函数名
//   - httpMethod: http请求方式
func MethodEnum(httpMethod string) string {
	var MethodMapping = map[string]string{
		httpLowerGet:     "Get",
		httpLowerHead:    "Head",
		httpLowerPost:    "Post",
		httpLowerPut:     "Put",
		httpLowerDelete:  "Delete",
		httpLowerConnect: "Connect",
		httpLowerOptions: "Options",
		httpLowerTrace:   "Trace",
		httpUpperGet:     "Get",
		httpUpperHead:    "Head",
		httpUpperPost:    "Post",
		httpUpperPut:     "Put",
		httpUpperDelete:  "Delete",
		httpUpperConnect: "Connect",
		httpUpperOptions: "Options",
		httpUpperTrace:   "Trace",
	}
	return MethodMapping[httpMethod]
}
