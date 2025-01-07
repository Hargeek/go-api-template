[This is an API documentation for the Go API Template](https://github.com/hargeek/go-api-template)

### 请求认证方式

需要认证的接口，要加token作为特定请求头

	请求头名称: Authorization
	类型: Bearer Token
	值: Bearer ${Token}

### 响应体错误码

接口如果请求出错且返回内容包含了自定义的错误码，可以根据错误码来判断错误类型

	错误码名称: code
	类型: int
    描述: 自定义错误码共六位，前三位表示功能模块，后三位表示具体错误类型

### 错误码对照表

| 错误码Code | 描述Msg                     | 中文释义      |
|---------|---------------------------|-----------|
| 0       | no error                  | 请求成功无错误   |
| 101001  | invalid request parameter | 请求参数无效    |
| 101002  | missing request parameter | 请求参数不足    |
