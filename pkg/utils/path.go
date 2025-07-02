package utils

import (
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

// GetBasePathFromURL 从URL获取最后一个路径
func GetBasePathFromURL(URL string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	parts := strings.Split(u.Path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1], nil
	}
	return "", errors.New("invalid URL")
}
