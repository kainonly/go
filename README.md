# Go Utils

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/kainonly/go/testing.yml?style=flat-square)](https://github.com/kainonly/go/actions/workflows/testing.yml)
[![Coveralls github](https://img.shields.io/coveralls/github/kainonly/go.svg?style=flat-square)](https://coveralls.io/github/kainonly/go)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kainonly/go?style=flat-square)](https://github.com/kainonly/go)
[![Go Report Card](https://goreportcard.com/badge/github.com/kainonly/go?style=flat-square)](https://goreportcard.com/report/github.com/kainonly/go)
[![Release](https://img.shields.io/github/v/release/kainonly/go.svg?style=flat-square)](https://github.com/kainonly/go)
[![GitHub license](https://img.shields.io/github/license/kainonly/go?style=flat-square)](https://raw.githubusercontent.com/kainonly/go/main/LICENSE)

Common Go utilities for web services, with Hertz-friendly integrations.

## Install

```bash
go get github.com/kainonly/go
```

## Packages

| Package | Purpose |
| --- | --- |
| `vd` | Validator wrapper and Hertz integration |
| `passport` | JWT auth helpers |
| `csrf` | CSRF protection middleware |
| `captcha` | Redis-backed captcha verification |
| `locker` | Redis-backed counters and lockout helpers |
| `passlib` | Password hashing and verification |
| `totp` | TOTP secret generation and validation |
| `cipher` | Symmetric encryption helpers |
| `help` | Small utility helpers |

## Usage

Import only the package you need:

```go
import "github.com/kainonly/go/vd"

validate := vd.Default()
_ = validate.Validate(req)
```

Detailed usage is kept in package comments:

```bash
go doc github.com/kainonly/go/vd
go doc github.com/kainonly/go/passport
go doc github.com/kainonly/go/csrf
```

## License

[BSD-3-Clause License](https://github.com/kainonly/go/blob/main/LICENSE)
