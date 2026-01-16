# Changelog

本文档记录了项目的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [1.0.1] - 2026-01-16

### 新增功能 (Added)
- **help** - 添加 UUIDv7 生成与时间提取功能
  - 支持生成基于时间戳的 UUIDv7
  - 支持从 UUIDv7 提取时间戳信息

### 改进 (Changed)
- **依赖** - 更新 Go 依赖版本以修复安全漏洞和兼容性问题
  - 更新核心依赖包至最新稳定版本
  - 修复已知安全漏洞

### CI/CD (Infrastructure)
- 统一工作流中的 Go 版本设置格式
- 简化测试工作流中的 go 模块命令
- 简化 go 版本设置并移除缓存配置
- 在测试前清理 go 模块缓存以提高可靠性

## [1.0.0] - 2026-01-16

### 首次发布 🎉

这是 Go Utils 工具库的首次正式发布，提供了一套完整的 Web 开发工具集。

### 新增功能 (Added)

#### 安全认证模块
- **passport** - JWT 认证库，基于 HS256 算法
  - 支持自定义声明
  - 内置安全验证
  - 灵活的配置选项

- **csrf** - CSRF 防护中间件
  - 双重提交 Cookie 验证
  - 支持 Hertz 框架集成
  - 可配置的 Cookie 选项

- **passlib** - 密码哈希工具
  - 基于 Argon2id 算法
  - 自动参数调优
  - 支持哈希升级检查

- **totp** - TOTP 一次性密码
  - 符合 RFC 6238 标准
  - 支持多种时间窗口
  - 密钥生成和验证

#### 验证与处理模块
- **vd** - 数据验证器
  - 集成 Hertz 框架
  - 中国本地化验证规则（手机号、身份证、银行卡等）
  - 密码强度验证（弱/中/强）
  - 命名规范验证（snake_case、camelCase 等）
  - 支持自定义验证规则

- **captcha** - 验证码管理
  - 基于 Redis 存储
  - 自动过期管理
  - 防暴力破解

#### 加密与安全模块
- **cipher** - 对称加密工具
  - AES-256-GCM 加密
  - 简洁的 API 接口
  - 自动 IV 生成

- **help** - 工具函数库
  - 国密 SM2 签名验签
  - 国密 SM4 加解密
  - SHA256 和 HMAC-SHA256
  - UUID 和雪花 ID 生成
  - 随机字符串生成器

#### 流量控制模块
- **locker** - 限流器
  - 基于 Redis 的分布式计数
  - 尝试次数限制
  - 自动锁定机制
  - 灵活的时间窗口

### 测试 (Tests)
- 为所有模块添加了完整的单元测试
- 代码覆盖率监控
- CI/CD 集成

### 文档 (Documentation)
- 完整的 README 文档
- 每个模块的使用示例
- Hertz 框架集成指南
- 安全策略文档 (SECURITY.md)
- BSD-3-Clause 开源协议

### 基础设施 (Infrastructure)
- GitHub Actions 工作流配置
- Coveralls 代码覆盖率集成
- Go Report Card 集成
- 自动化测试流程

---

## 版本说明

### 语义化版本规则

- **主版本号 (MAJOR)**: 不兼容的 API 变更
- **次版本号 (MINOR)**: 向后兼容的功能新增
- **修订号 (PATCH)**: 向后兼容的问题修正

### 变更类型说明

- **Added**: 新增功能
- **Changed**: 现有功能的变更
- **Deprecated**: 即将废弃的功能
- **Removed**: 已删除的功能
- **Fixed**: 问题修复
- **Security**: 安全性相关的修复

[1.0.1]: https://github.com/kainonly/go/releases/tag/v1.0.1
[1.0.0]: https://github.com/kainonly/go/releases/tag/v1.0.0
