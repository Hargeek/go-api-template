package types

import (
	"runtime"
	"runtime/debug"
)

const ServiceName = "go-api-template"

var (
	ServerInfo, _ = debug.ReadBuildInfo()
	GoVersion     = runtime.Version()
	Branch        string
	Revision      string
	BuildDate     string
	BuildUser     string
)
