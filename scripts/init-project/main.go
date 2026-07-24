// Command init-project 将 go-api-template 转换为可直接开发的项目。
//
// 初始化过程按以下顺序执行：
//  1. 校验项目名称、Go Module、模板根目录和 Git 工作区状态；
//  2. 删除所有 Profile 都不保留的模板资源，thin Profile 再额外删除 MTL 实现；
//  3. 处理 Profile 条件代码，替换项目名称和 Go Module；
//  4. 重写 README，删除初始化脚本、Profile 清单和失效的 Make 目标；
//  5. 对最终项目执行 gofmt 和 go mod tidy。
//
// wip 目录是模板维护者的私有工作区，初始化过程不读写、删除或格式化其中内容。
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	// templateModule 同时是未初始化模板的 Go Module 和默认项目名。
	templateModule = "go-api-template"
	// MTL 标记包围仅完整 Profile 保留的 Metrics、Trace 和 OTEL 代码。
	mtlStartMarker = "profile:mtl:start"
	mtlEndMarker   = "profile:mtl:end"
	// init 标记包围只在模板初始化阶段有效的 Make 变量和目标。
	initStartMarker = "profile:init:start"
	initEndMarker   = "profile:init:end"
	// common 清单对 full/thin 都生效，thin 清单只记录 thin 额外删除的文件。
	commonRemoveManifest = "profiles/common/remove.txt"
	thinRemoveManifest   = "profiles/thin/remove.txt"
)

var projectNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

type options struct {
	profile string // full 保留 MTL，thin 移除 MTL。
	name    string // 项目展示名称，同时用于替换模板中的服务名。
	module  string // go.mod 中的 Module Path；未显式传入时使用 name。
}

func main() {
	var opts options
	flag.StringVar(&opts.profile, "profile", "full", "初始化 Profile：full 或 thin")
	flag.StringVar(&opts.name, "name", "", "项目名称")
	flag.StringVar(&opts.module, "module", "", "Go Module；为空时使用项目名称")
	flag.Parse()

	if err := run(opts); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "初始化项目失败：%v\n", err)
		os.Exit(1)
	}
}

// run 编排完整初始化流程。
//
// 所有破坏性删除都在参数、模板根目录和 Git 工作区校验通过后才会发生。
// 流程不尝试回滚已完成的文件操作，因此只允许在干净的模板工作区执行；
// 任一步失败时，用户应重新克隆模板后再次初始化。
func run(opts options) error {
	// module 为可选参数，便于内部项目直接使用项目名作为 Module Path。
	if opts.module == "" {
		opts.module = opts.name
	}
	if err := validateOptions(opts); err != nil {
		return err
	}

	root, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录: %w", err)
	}
	if err := ensureTemplateRoot(root); err != nil {
		return err
	}
	if err := ensureCleanWorktree(root); err != nil {
		return err
	}

	// 先执行公共清理：README 中的演示图、MTL 专题文档等不进入任何派生项目。
	if err := removeManifestFiles(root, filepath.Join(root, commonRemoveManifest)); err != nil {
		return err
	}
	// thin 在公共清理的基础上，再删除独立的 Metrics、Trace、Telemetry 代码和联调资源。
	if opts.profile == "thin" {
		if err := removeManifestFiles(root, filepath.Join(root, thinRemoveManifest)); err != nil {
			return err
		}
	}
	// 删除完独立文件后，再处理与基础代码共存的条件片段和项目标识。
	if err := rewriteTemplateFiles(root, opts); err != nil {
		return err
	}
	// thin 已无 Trace/OTEL Handler，将 Context 日志调用收敛为普通日志调用，避免保留无效 ctx 参数。
	if opts.profile == "thin" {
		if err := simplifyContextLogging(root); err != nil {
			return err
		}
	}
	// 模板 README 只用于指导初始化；派生项目从一行标题开始自行维护项目文档。
	if err := writeProjectReadme(root, opts.name); err != nil {
		return err
	}
	// 初始化工具已完成使命，先删除模板专用文件，再根据最终项目内容格式化和整理依赖。
	if err := removeTemplateFiles(root); err != nil {
		return err
	}
	if err := runGoFmt(root); err != nil {
		return err
	}
	if err := runCommand(root, "go", "mod", "tidy"); err != nil {
		return err
	}

	fmt.Printf("项目初始化完成：name=%s module=%s profile=%s\n", opts.name, opts.module, opts.profile)
	return nil
}

// runGoFmt 只格式化派生项目的 Go 源码。
// .git、构建产物、本地数据、vendor 和 wip 不属于初始化输出的可改写范围。
func runGoFmt(root string) error {
	files := make([]string, 0)
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", "bin", "data", "vendor", "wip":
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) == ".go" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("收集 Go 文件: %w", err)
	}
	if len(files) == 0 {
		return nil
	}
	return runCommand(root, "gofmt", append([]string{"-w"}, files...)...)
}

// validateOptions 在任何文件操作前校验输入。
// name 会出现在文件内容和 README 中，因此限制为稳定、可预期的字符集；
// module 允许内部仓库路径，但不允许空白字符和首尾斜杠。
func validateOptions(opts options) error {
	if opts.profile != "full" && opts.profile != "thin" {
		return fmt.Errorf("不支持的 profile %q，仅支持 full 或 thin", opts.profile)
	}
	if !projectNamePattern.MatchString(opts.name) {
		return errors.New("name 只能包含字母、数字、点、下划线和连字符，且必须以字母或数字开头")
	}
	if opts.name == templateModule {
		return errors.New("name 不能继续使用模板名称 go-api-template")
	}
	if opts.module == "" || strings.ContainsAny(opts.module, " \\\t\r\n") ||
		strings.HasPrefix(opts.module, "/") || strings.HasSuffix(opts.module, "/") {
		return fmt.Errorf("无效的 Go Module %q", opts.module)
	}
	if opts.module == templateModule {
		return errors.New("module 不能继续使用模板模块名 go-api-template")
	}
	return nil
}

// ensureTemplateRoot 通过 go.mod 确认当前目录仍是未初始化的模板。
// 初始化后 Module Path 已被替换，即使用户手工保留了脚本，也不会重复执行。
func ensureTemplateRoot(root string) error {
	content, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return fmt.Errorf("读取 go.mod: %w", err)
	}
	if !strings.Contains(string(content), "module "+templateModule) {
		return errors.New("当前目录不是未初始化的 go-api-template 根目录")
	}
	return nil
}

// ensureCleanWorktree 防止删除清单或全局替换覆盖用户未保存的修改。
// 通过下载压缩包等方式获取、不含 .git 的模板仍允许初始化；
// 存在 .git 时则要求 tracked/untracked 状态全部干净。
func ensureCleanWorktree(root string) error {
	if _, err := os.Stat(filepath.Join(root, ".git")); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return fmt.Errorf("检查 .git: %w", err)
	}

	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = root
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("检查 Git 工作区: %w", err)
	}
	if len(bytes.TrimSpace(output)) != 0 {
		return errors.New("Git 工作区不干净；请先提交或清理变更后再初始化")
	}
	return nil
}

// removeManifestFiles 逐行执行删除清单。
// 空行和 # 开头的注释会被忽略；文件和目录统一使用 RemoveAll，
// 因此清单维护者必须使用项目根目录下的相对路径。
func removeManifestFiles(root, manifestPath string) error {
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("读取删除清单 %s: %w", manifestPath, err)
	}
	for _, rawLine := range strings.Split(string(content), "\n") {
		path := strings.TrimSpace(rawLine)
		if path == "" || strings.HasPrefix(path, "#") {
			continue
		}
		target, err := safeTarget(root, path)
		if err != nil {
			return err
		}
		if err := os.RemoveAll(target); err != nil {
			return fmt.Errorf("删除文件 %s: %w", path, err)
		}
	}
	return nil
}

// safeTarget 将清单中的相对路径转换为可删除的绝对路径。
// 该校验拒绝绝对路径和 ../ 越界，避免错误的清单删除项目目录之外的文件。
func safeTarget(root, relativePath string) (string, error) {
	if filepath.IsAbs(relativePath) {
		return "", fmt.Errorf("删除清单不允许绝对路径 %q", relativePath)
	}
	target := filepath.Clean(filepath.Join(root, relativePath))
	rel, err := filepath.Rel(root, target)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("删除路径越出项目目录 %q", relativePath)
	}
	return target, nil
}

// rewriteTemplateFiles 遍历项目中的普通文本文件，依次完成：
//  1. 保留或删除 MTL Profile 片段；
//  2. 删除只为初始化脚本服务的片段；
//  3. 替换模板项目名和 Go Module。
//
// README 由 writeProjectReadme 直接重写，脚本源码在流程末尾自删除，二者都无需参与全局替换。
// wip 作为模板维护文档始终跳过；二进制文件和非 UTF-8 文本也不作内容替换。
func rewriteTemplateFiles(root string, opts options) error {
	return filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", "bin", "data", "vendor", "wip":
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.Type().IsRegular() {
			return nil
		}

		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if relativePath == "README.md" || strings.HasPrefix(relativePath, filepath.Join("scripts", "init-project")) {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("读取 %s: %w", relativePath, err)
		}
		// NUL 字节是常见的二进制文件信号；与 UTF-8 校验结合，避免损坏图片或其他非文本资源。
		if bytes.IndexByte(content, 0) >= 0 || !utf8.Valid(content) {
			return nil
		}

		processed, err := applyProfileMarkers(string(content), opts.profile)
		if err != nil {
			return fmt.Errorf("处理 %s: %w", relativePath, err)
		}
		processed, err = removeMarkedBlock(processed, initStartMarker, initEndMarker, "项目初始化")
		if err != nil {
			return fmt.Errorf("处理 %s: %w", relativePath, err)
		}
		processed = replaceTemplateIdentity(processed, opts)
		if processed == string(content) {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("读取 %s 权限: %w", relativePath, err)
		}
		if err := os.WriteFile(path, []byte(processed), info.Mode().Perm()); err != nil {
			return fmt.Errorf("写入 %s: %w", relativePath, err)
		}
		return nil
	})
}

// simplifyContextLogging 将 thin 项目中的 Context 日志 API 改为普通日志 API。
//
// 转换通过 Go AST 完成，而不是使用字符串替换：
//
//	logger.ErrorContext(ctx, "query failed", "error", err)
//
// 会生成：
//
//	logger.Error("query failed", "error", err)
//
// AST 能正确删除 ctx、c.Request.Context() 等任意形式的第一个参数，
// 并且不会改动注释和字符串。只有导入项目 common/logger 包的文件才会执行转换，
// 避免误改局部变量或其他依赖中同样命名为 logger 的对象。
func simplifyContextLogging(root string) error {
	methodNames := map[string]string{
		"DebugContext": "Debug",
		"InfoContext":  "Info",
		"WarnContext":  "Warn",
		"ErrorContext": "Error",
	}

	return filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", "bin", "data", "vendor", "wip", "init-project":
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		fileSet := token.NewFileSet()
		file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("解析 Go 文件 %s: %w", path, err)
		}
		if !importsProjectLogger(file) {
			return nil
		}

		changed := false
		var transformErr error
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			packageName, ok := selector.X.(*ast.Ident)
			if !ok || packageName.Name != "logger" {
				return true
			}
			methodName, ok := methodNames[selector.Sel.Name]
			if !ok {
				return true
			}
			if len(call.Args) == 0 {
				transformErr = fmt.Errorf("日志调用 %s.%s 缺少 context 参数", packageName.Name, selector.Sel.Name)
				return false
			}

			selector.Sel.Name = methodName
			call.Args = call.Args[1:]
			changed = true
			return true
		})
		if transformErr != nil {
			return fmt.Errorf("处理 Go 文件 %s: %w", path, transformErr)
		}
		if !changed {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("读取 Go 文件权限 %s: %w", path, err)
		}
		var output bytes.Buffer
		if err := format.Node(&output, fileSet, file); err != nil {
			return fmt.Errorf("格式化 Go 文件 %s: %w", path, err)
		}
		if err := os.WriteFile(path, output.Bytes(), info.Mode().Perm()); err != nil {
			return fmt.Errorf("写入 Go 文件 %s: %w", path, err)
		}
		return nil
	})
}

// importsProjectLogger 判断当前文件的 logger 标识符是否指向项目 common/logger 包。
// 默认导入名和显式 `logger` 别名都受支持，dot import 和其他别名不会误命中。
func importsProjectLogger(file *ast.File) bool {
	for _, importSpec := range file.Imports {
		importPath, err := strconv.Unquote(importSpec.Path.Value)
		if err != nil || !strings.HasSuffix(importPath, "/common/logger") {
			continue
		}
		if importSpec.Name == nil || importSpec.Name.Name == "logger" {
			return true
		}
	}
	return false
}

// writeProjectReadme 丢弃模板使用说明，仅为派生项目保留一级标题。
func writeProjectReadme(root, projectName string) error {
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("# "+projectName+"\n"), 0o644); err != nil {
		return fmt.Errorf("重写 README.md: %w", err)
	}
	return nil
}

// removeTemplateFiles 删除初始化完成后不再属于业务项目的工具资源。
// profiles 和 scripts/init-project 必定删除；如果 scripts 下还有用户脚本，则仅删除初始化子目录并保留父目录。
// Makefile 中的 init/init-thin 目标已在 rewriteTemplateFiles 阶段通过 init 标记移除。
func removeTemplateFiles(root string) error {
	for _, path := range []string{"profiles", filepath.Join("scripts", "init-project")} {
		target, err := safeTarget(root, path)
		if err != nil {
			return err
		}
		if err := os.RemoveAll(target); err != nil {
			return fmt.Errorf("清理模板文件 %s: %w", path, err)
		}
	}
	// scripts 中没有其他项目脚本时一并移除空目录；存在其他文件则保留。
	scriptsDir := filepath.Join(root, "scripts")
	entries, err := os.ReadDir(scriptsDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("检查 scripts 目录: %w", err)
	}
	if len(entries) == 0 {
		if err := os.Remove(scriptsDir); err != nil {
			return fmt.Errorf("清理空 scripts 目录: %w", err)
		}
	}
	return nil
}

// applyProfileMarkers 根据 Profile 处理 MTL 条件片段。
// full 仅删除标记行并保留其中代码；thin 连同标记之间的内容一并删除。
func applyProfileMarkers(content, profile string) (string, error) {
	if profile == "full" {
		return removeMarkerLines(content, mtlStartMarker, mtlEndMarker, "MTL Profile")
	}
	return removeMarkedBlock(content, mtlStartMarker, mtlEndMarker, "MTL Profile")
}

// removeMarkerLines 保留标记区间的内容，只移除开始和结束标记。
// 即使不删除区间内容，仍校验标记必须成对且不允许嵌套，防止生成包含半个条件块的项目。
func removeMarkerLines(content, startMarker, endMarker, markerName string) (string, error) {
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))
	inBlock := false
	for _, line := range lines {
		switch {
		case strings.Contains(line, startMarker):
			if inBlock {
				return "", fmt.Errorf("%s 标记不能嵌套", markerName)
			}
			inBlock = true
			continue
		case strings.Contains(line, endMarker):
			if !inBlock {
				return "", fmt.Errorf("存在没有开始标记的 %s 结束标记", markerName)
			}
			inBlock = false
			continue
		}
		result = append(result, line)
	}
	if inBlock {
		return "", fmt.Errorf("%s 开始标记缺少结束标记", markerName)
	}
	return strings.Join(result, "\n"), nil
}

// removeMarkedBlock 删除开始标记、结束标记及二者之间的全部内容。
// 该方法同时用于 thin 的 MTL 片段和两种 Profile 都需移除的初始化专用片段。
func removeMarkedBlock(content, startMarker, endMarker, markerName string) (string, error) {
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))
	inBlock := false
	for _, line := range lines {
		switch {
		case strings.Contains(line, startMarker):
			if inBlock {
				return "", fmt.Errorf("%s 标记不能嵌套", markerName)
			}
			inBlock = true
			continue
		case strings.Contains(line, endMarker):
			if !inBlock {
				return "", fmt.Errorf("存在没有开始标记的 %s 结束标记", markerName)
			}
			inBlock = false
			continue
		}
		if !inBlock {
			result = append(result, line)
		}
	}
	if inBlock {
		return "", fmt.Errorf("%s 开始标记缺少结束标记", markerName)
	}
	return strings.Join(result, "\n"), nil
}

// replaceTemplateIdentity 替换模板标识。
// 替换顺序从更具体的仓库 URL、Module 声明和 import path 开始，最后再替换剩余的模板名，
// 避免先替换 go-api-template 后无法识别完整路径。Module 不包含斜杠时，不推断仓库 URL。
func replaceTemplateIdentity(content string, opts options) string {
	if strings.Contains(opts.module, "/") {
		content = strings.ReplaceAll(content, "https://github.com/Hargeek/go-api-template.git", "https://"+opts.module+".git")
		content = strings.ReplaceAll(content, "https://github.com/hargeek/go-api-template", "https://"+opts.module)
		content = strings.ReplaceAll(content, "github.com/hargeek/go-api-template", opts.module)
	}
	content = strings.ReplaceAll(content, "module "+templateModule, "module "+opts.module)
	content = strings.ReplaceAll(content, "SERVICE_NAME ?= "+templateModule, "SERVICE_NAME ?= "+opts.module)
	content = strings.ReplaceAll(content, templateModule+"/", opts.module+"/")
	content = strings.ReplaceAll(content, "Go API Template", opts.name)
	return strings.ReplaceAll(content, templateModule, opts.name)
}

// runCommand 在项目根目录执行外部工具，并将 stdout/stderr 原样透传给用户。
// 上层通过包装后的错误获知具体是 gofmt 还是 go mod tidy 失败。
func runCommand(root, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("执行 %s %s: %w", name, strings.Join(args, " "), err)
	}
	return nil
}
