# Go Utils

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/kainonly/go/testing.yml?style=flat-square)](https://github.com/kainonly/go/actions/workflows/testing.yml)
[![Coveralls github](https://img.shields.io/coveralls/github/kainonly/go.svg?style=flat-square)](https://coveralls.io/github/kainonly/go)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kainonly/go?style=flat-square)](https://github.com/kainonly/go)
[![Go Report Card](https://goreportcard.com/badge/github.com/kainonly/go?style=flat-square)](https://goreportcard.com/report/github.com/kainonly/go)
[![Release](https://img.shields.io/github/v/release/kainonly/go.svg?style=flat-square)](https://github.com/kainonly/go)
[![GitHub license](https://img.shields.io/github/license/kainonly/go?style=flat-square)](https://raw.githubusercontent.com/kainonly/go/main/LICENSE)

面向 Web 服务的 Go 工具库，提供 Hertz 友好的集成。

## 安装

```bash
go get github.com/kainonly/go
```

## 包列表

| 包 | 用途 |
| --- | --- |
| `vd` | 校验器封装与 Hertz 集成 |
| `passport` | JWT 认证辅助 |
| `csrf` | CSRF 防护中间件 |
| `captcha` | 基于 Redis 的验证码管理 |
| `locker` | 基于 Redis 的计数与锁定辅助 |
| `passlib` | 密码哈希与校验 |
| `totp` | TOTP 密钥生成与校验 |
| `cipher` | 对称加密辅助 |
| `help` | 常用小工具函数 |

## 使用

按需导入所需的包：

```go
import "github.com/kainonly/go/vd"

validate := vd.Default()
_ = validate.Validate(req)
```

详细用法保留在各包的文档注释中：

```bash
go doc github.com/kainonly/go/vd
go doc github.com/kainonly/go/passport
go doc github.com/kainonly/go/csrf
```

## 许可证

[BSD-3-Clause License](https://github.com/kainonly/go/blob/main/LICENSE)
