package resource

import (
	"embed"
)

//go:embed static_resource
var StaticResourceFS embed.FS

//go:embed error_code
var ErrorCodeFS embed.FS

func GetStaticResourceEmbed() embed.FS {
	return StaticResourceFS
}

func GetErrorCodeEmbed() embed.FS {
	return ErrorCodeFS
}
