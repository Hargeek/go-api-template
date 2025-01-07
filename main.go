package main

import "go-api-template/cmd"

// @title                      go-api-template
// @version                    1.0
// // @contact.name               go-api-template
// // @contact.url                http://www.swagger.io/support
// // @license.url                http://www.apache.org/licenses/LICENSE-2.0.html
// // @contact.email              ssgeek@hotmail.com
// // @license.name               Apache 2.0
// @BasePath                   /
// @Schemes                    http https
// @securityDefinitions.apiKey Authorization
// @in                         header
// @name                       Authorization
func main() {
	// 启动服务
	cmd.RunServer()
}
