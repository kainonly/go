# AGENTS.md

## 项目

面向 Web 服务的 Go 工具库（`github.com/kainonly/go`），提供 Hertz 友好的集成。每个目录是一个独立包，用户按需导入。要求 Go 1.25+。

## 常用命令

```bash
go build ./...    # 编译所有包
go vet ./...      # 静态检查
go test ./...     # 运行测试（未设置 DATABASE_REDIS 时跳过 Redis 相关测试）
```

`captcha` 和 `locker` 包依赖 Redis，未设置 `DATABASE_REDIS` 环境变量时会静默跳过测试。运行完整测试：

```bash
docker run -d --name redis -p 6379:6379 redis:7-alpine
DATABASE_REDIS=redis://127.0.0.1:6379 go test -race ./...
```

## 文档约定

- 所有注释与文档使用中文，包括包 doc 注释（godoc）、README、AGENTS.md、Issue/PR 等；代码中的错误信息等字符串字面量保持英文。
- 详细用法保留在包 doc 注释中，不建立 docs 目录。新增功能时用 `# 标题` godoc 语法补充可运行示例；涉及前端对接的包（`captcha`、`csrf`、`passport`）应同时给出 Hertz 后端与 Angular 前端两套示例。
- `Deprecated:` 弃用标记保持英文原词（gopls/staticcheck 依赖它识别），其后的描述用中文。
- 仓库不维护 CHANGELOG，Release Notes 由 release workflow 自动生成。
- 提交信息使用英文 type 前缀 + 中文描述（如 `chore: 精简项目结构`）。

## 代码约定

- Options 模式：构造函数接受可变参数 `Option`（如 `csrf.SetKey`、`captcha.SetPrefix`）。
- 错误定义为带包名前缀的哨兵值（如 `locker.ErrLocked`、`passport.ErrInvalidIssuer`）。
- 基于 Redis 的包通过可配置前缀对键做命名空间隔离。
- `vd` 中的验证规则是返回 `Rule` 结构的纯函数，正则通过 `sync.OnceValue` 惰性编译。

## CI

GitHub Actions（`.github/workflows/`）：

- `testing.yml` — main push 和 PR 触发；docker compose 启动 Redis，运行 `go test -race` 并上报 Coveralls 覆盖率。
- `release.yml` — `v*` tag 触发；跑测试、创建 GitHub Release（自动生成 notes）、通知 Go module proxy 索引新版本。
