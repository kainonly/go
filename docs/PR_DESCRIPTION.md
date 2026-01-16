# Release v1.0.0

## 📋 变更概述

为 Go Utils 工具库准备 v1.0.0 正式版本发布，添加了完整的发布文档和自动化流程。

## 📝 新增文件

1. **RELEASE_NOTES_v1.0.0.md** - 详细的 v1.0.0 版本发布说明
   - 完整的功能列表
   - 使用示例和安装指南
   - 质量保证说明

2. **CHANGELOG.md** - 项目变更日志
   - 遵循 Keep a Changelog 格式
   - 语义化版本说明
   - 完整的 v1.0.0 变更记录

3. **.github/RELEASE_TEMPLATE.md** - 发布流程模板
   - 发布前检查清单
   - 详细的发布步骤
   - 发布后验证指南

4. **scripts/release.sh** - 自动化发布脚本
   - 自动运行测试
   - 创建并推送标签
   - 支持 gh CLI 创建 Release

5. **.github/workflows/release.yml** - GitHub Actions 发布工作流
   - 标签推送自动触发
   - 自动运行测试验证
   - 自动创建 GitHub Release

## 🚀 发布流程

### 方式一：使用自动化脚本（推荐）

合并此 PR 后，在 main 分支执行：

```bash
git checkout main
git pull origin main
./scripts/release.sh v1.0.0
```

脚本会自动：
- ✅ 运行所有测试
- ✅ 创建 git 标签
- ✅ 推送标签到远程
- ✅ 创建 GitHub Release（需要 gh CLI）

### 方式二：手动发布

如果自动化脚本无法使用，可以手动执行：

```bash
# 1. 确保在 main 分支
git checkout main
git pull origin main

# 2. 运行测试
go test ./... -v

# 3. 创建标签
git tag -a v1.0.0 -m "Release v1.0.0"

# 4. 推送标签
git push origin v1.0.0
```

推送标签后，GitHub Actions 会自动创建 Release。

### 方式三：GitHub Web 界面

1. 合并此 PR
2. 访问 https://github.com/kainonly/go/releases/new
3. 选择或创建标签 `v1.0.0`
4. 标题: `v1.0.0`
5. 复制 `RELEASE_NOTES_v1.0.0.md` 内容到发布说明
6. 点击 "Publish release"

## ✅ 发布前检查

- [x] 所有单元测试通过
- [x] CI/CD 工作流正常
- [x] 文档完整且最新
- [x] README.md 准确描述所有功能
- [x] 安全策略文档已创建
- [x] 开源协议清晰
- [x] 发布说明详尽

## 📊 版本亮点

这是 Go Utils 的**首次正式发布**，包含：

- 🛡️ **4个安全认证模块**: passport, csrf, passlib, totp
- 📝 **2个验证处理模块**: vd, captcha
- 🔐 **2个加密安全模块**: cipher, help (含国密支持)
- 🚦 **1个流量控制模块**: locker
- 🧪 **完整的测试覆盖**
- 📚 **详尽的文档和示例**

## 🔗 相关链接

- [RELEASE_NOTES_v1.0.0.md](./RELEASE_NOTES_v1.0.0.md) - 发布说明
- [CHANGELOG.md](./CHANGELOG.md) - 变更日志
- [.github/RELEASE_TEMPLATE.md](./.github/RELEASE_TEMPLATE.md) - 发布模板
- [scripts/release.sh](./scripts/release.sh) - 发布脚本
- [.github/workflows/release.yml](./.github/workflows/release.yml) - 发布工作流

## 📦 发布后续步骤

1. 验证 Release 在 GitHub 上正确显示
2. 确认 Go Proxy 已索引新版本: https://pkg.go.dev/github.com/kainonly/go@v1.0.0
3. 更新项目徽章（如果需要）
4. 通知用户和社区

---

**准备发布人**: Claude Code
**发布日期**: 2026-01-16
**版本**: v1.0.0
