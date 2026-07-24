package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyProfileMarkers(t *testing.T) {
	content := "base\n// profile:mtl:start\nmtl\n// profile:mtl:end\nend"

	full, err := applyProfileMarkers(content, "full")
	require.NoError(t, err)
	assert.Equal(t, "base\nmtl\nend", full)

	thin, err := applyProfileMarkers(content, "thin")
	require.NoError(t, err)
	assert.Equal(t, "base\nend", thin)
}

func TestApplyProfileMarkersRejectsInvalidMarkers(t *testing.T) {
	_, err := applyProfileMarkers("base\n// profile:mtl:start\nmtl", "thin")
	require.Error(t, err)
}

func TestRemoveMarkedBlock(t *testing.T) {
	content := "base\n# profile:init:start\ninit\n# profile:init:end\nend"

	got, err := removeMarkedBlock(content, initStartMarker, initEndMarker, "项目初始化")
	require.NoError(t, err)
	assert.Equal(t, "base\nend", got)
}

func TestWriteProjectReadme(t *testing.T) {
	root := t.TempDir()

	require.NoError(t, writeProjectReadme(root, "demo"))
	content, err := os.ReadFile(filepath.Join(root, "README.md"))
	require.NoError(t, err)
	assert.Equal(t, "# demo\n", string(content))
}

func TestRemoveTemplateFiles(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "profiles", "thin"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "scripts", "init-project"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "profiles", "thin", "remove.txt"), nil, 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "scripts", "init-project", "main.go"), nil, 0o644))

	require.NoError(t, removeTemplateFiles(root))
	_, err := os.Stat(filepath.Join(root, "profiles"))
	require.ErrorIs(t, err, os.ErrNotExist)
	_, err = os.Stat(filepath.Join(root, "scripts"))
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestSimplifyContextLogging(t *testing.T) {
	root := t.TempDir()
	source := `package service

	import (
		"context"
		"net/http"

		"example.com/demo/common/logger"
	)

func logFailure(ctx context.Context, request *http.Request, err error) {
	logger.ErrorContext(ctx, "query failed", "error", err)
	logger.InfoContext(request.Context(), "query finished")
}
`
	path := filepath.Join(root, "service.go")
	require.NoError(t, os.WriteFile(path, []byte(source), 0o644))

	require.NoError(t, simplifyContextLogging(root))
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(content), `logger.Error("query failed", "error", err)`)
	assert.Contains(t, string(content), `logger.Info("query finished")`)
	assert.NotContains(t, string(content), "ErrorContext")
	assert.NotContains(t, string(content), "InfoContext")
}

func TestSimplifyContextLoggingIgnoresOtherLogger(t *testing.T) {
	root := t.TempDir()
	source := `package service

import "example.com/external/logger"

func logFailure(ctx context.Context) {
	logger.ErrorContext(ctx, "query failed")
}
`
	path := filepath.Join(root, "service.go")
	require.NoError(t, os.WriteFile(path, []byte(source), 0o644))

	require.NoError(t, simplifyContextLogging(root))
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(content), `logger.ErrorContext(ctx, "query failed")`)
}

func TestReplaceTemplateIdentity(t *testing.T) {
	content := `module go-api-template
import "go-api-template/common/logger"
const ServiceName = "go-api-template"
SERVICE_NAME ?= go-api-template
`

	got := replaceTemplateIdentity(content, options{
		name:   "demo-service",
		module: "github.com/example/demo-service",
	})

	assert.Contains(t, got, "module github.com/example/demo-service")
	assert.Contains(t, got, `"github.com/example/demo-service/common/logger"`)
	assert.Contains(t, got, `const ServiceName = "demo-service"`)
	assert.Contains(t, got, "SERVICE_NAME ?= github.com/example/demo-service")
}

func TestSafeTarget(t *testing.T) {
	root := t.TempDir()

	_, err := safeTarget(root, "../outside")
	require.Error(t, err)

	target, err := safeTarget(root, "pkg/telemetry")
	require.NoError(t, err)
	assert.Equal(t, root+"/pkg/telemetry", target)
}
