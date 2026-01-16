# v1.0.0 - Go 常用工具库集合首次发布

## 🎉 首次发布

Go Utils 是一个专为 Web 开发设计的 Go 常用工具库集合，提供开箱即用的安全组件和实用工具。

## ✨ 核心功能

### 🛡️ 安全认证模块
- **passport** - JWT 认证库 (HS256)，支持自定义声明和安全验证
- **csrf** - 基于双重提交 Cookie 的 CSRF 防护中间件
- **passlib** - 基于 Argon2id 的密码哈希工具
- **totp** - TOTP 一次性密码生成与验证

### 📝 验证与处理
- **vd** - 强大的验证器，集成 Hertz 框架
  - 支持中国本地化验证（手机号、身份证、银行卡、车牌号等）
  - 多种密码强度验证
  - 命名规范验证（snake_case、camelCase、PascalCase 等）
  - 支持自定义验证规则
- **captcha** - 基于 Redis 的验证码管理系统

### 🔐 加密与安全
- **cipher** - 对称加密工具，支持数据加解密
- **help** - 包含国密 SM2/SM4 支持的工具函数库
  - SM2 签名验签
  - SM4 加解密
  - SHA256、HMAC-SHA256
  - 随机字符串生成

### 🚦 流量控制
- **locker** - 基于 Redis 的限流器
  - 尝试次数计数
  - 自动锁定机制
  - 灵活的时间窗口配置

### 🛠️ 实用工具
- **help** - 丰富的工具函数集
  - UUID 和雪花 ID 生成
  - 随机数生成（数字、字母、混合）
  - 切片操作（反转、打乱）
  - 空值检查和指针工具

## 📦 安装

```bash
go get github.com/kainonly/go
```

## 🧪 质量保证

- ✅ 完整的单元测试覆盖
- ✅ 持续集成 (CI/CD)
- ✅ 代码质量检查
- ✅ 安全策略文档 (SECURITY.md)

## 📚 文档

详细的使用文档和示例请参阅 [README.md](https://github.com/kainonly/go/blob/main/README.md)

## 🔄 主要更新

- feat(vd): 添加自定义验证器模块并替换原有验证器实现
- feat(help): 重构帮助库，增强加密与工具函数
- feat(locker): 实现基于 Redis 的限流与尝试次数计数功能
- feat(captcha): 优化验证码管理功能并完善单元测试
- feat(passport): 重构 JWT 认证库并增强安全验证
- feat(csrf): 实现基于双重提交 Cookie 的 CSRF 防护中间件
- feat(totp): 重构 TOTP 实现并优化测试
- feat(cipher): 优化加密库接口并增强错误处理
- docs: 完善文档，添加 Hertz 框架集成说明和使用示例
- test: 为所有模块添加全面的单元测试
- ci: 配置 GitHub Actions 工作流和代码覆盖率检查

## 📄 开源协议

本项目采用 [BSD-3-Clause License](https://github.com/kainonly/go/blob/main/LICENSE)
