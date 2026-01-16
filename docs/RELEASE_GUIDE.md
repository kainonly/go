# v1.0.0 发布指南

## 🎯 当前状态

已在分支 `claude/release-v1.0.0-aFZct` 上完成所有发布准备工作：

✅ **已完成的工作**
- 创建完整的发布说明文档 (RELEASE_NOTES_v1.0.0.md)
- 创建项目变更日志 (CHANGELOG.md)
- 创建发布流程模板 (.github/RELEASE_TEMPLATE.md)
- 创建自动化发布脚本 (scripts/release.sh)
- 创建 GitHub Actions 自动发布工作流 (.github/workflows/release.yml)
- 创建本地 git 标签 v1.0.0
- 所有文件已推送到远程分支

## 📋 发布步骤

### 第一步：创建 Pull Request

访问 GitHub 创建 PR：
```
https://github.com/kainonly/go/pull/new/claude/release-v1.0.0-aFZct
```

**PR 标题**: Release v1.0.0

**PR 描述**: 使用 `PR_DESCRIPTION.md` 的内容

### 第二步：合并 Pull Request

1. 在 GitHub 上审核 PR
2. 确认所有 CI 检查通过
3. 合并 PR 到 main 分支

### 第三步：创建 GitHub Release

合并 PR 后，有三种方式创建 Release：

#### 方式 A：GitHub Actions 自动发布（最简单）

1. 在本地切换到 main 分支并拉取最新代码：
   ```bash
   git checkout main
   git pull origin main
   ```

2. 推送标签（会触发自动发布）：
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

3. GitHub Actions 会自动：
   - 运行测试
   - 创建 GitHub Release
   - 更新 Go Proxy 索引

#### 方式 B：使用发布脚本

```bash
git checkout main
git pull origin main
./scripts/release.sh v1.0.0
```

#### 方式 C：手动在 GitHub Web 界面创建

1. 访问: https://github.com/kainonly/go/releases/new
2. 创建标签: `v1.0.0`
3. 标题: `v1.0.0`
4. 复制 `RELEASE_NOTES_v1.0.0.md` 的内容作为发布说明
5. 点击 "Publish release"

## 📦 发布内容概览

### 核心功能模块（9个）

**安全认证** (4个)
- passport - JWT 认证
- csrf - CSRF 防护
- passlib - 密码哈希
- totp - 一次性密码

**验证处理** (2个)
- vd - 数据验证器
- captcha - 验证码管理

**加密安全** (2个)
- cipher - 对称加密
- help - 工具函数（含国密）

**流量控制** (1个)
- locker - 限流器

## ✅ 发布后验证

1. **验证 Release**
   - 访问 https://github.com/kainonly/go/releases
   - 确认 v1.0.0 显示正确

2. **验证 Go Proxy**
   - 等待 5-10 分钟
   - 访问 https://pkg.go.dev/github.com/kainonly/go@v1.0.0
   - 确认包文档可见

3. **验证安装**
   ```bash
   go get github.com/kainonly/go@v1.0.0
   ```

## 🔗 快速链接

- **创建 PR**: https://github.com/kainonly/go/pull/new/claude/release-v1.0.0-aFZct
- **创建 Release**: https://github.com/kainonly/go/releases/new
- **查看 Releases**: https://github.com/kainonly/go/releases
- **Go 包文档**: https://pkg.go.dev/github.com/kainonly/go

## 📧 发布后通知

可选：发布完成后可以：
- 更新项目主页
- 在社交媒体分享
- 通知早期用户
- 更新相关文档链接

---

**发布版本**: v1.0.0
**准备日期**: 2026-01-16
**状态**: 等待创建 PR 和发布
