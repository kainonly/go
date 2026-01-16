<!--
发布检查清单 - Release Checklist

使用此模板创建新版本的 GitHub Release
-->

## 版本信息
- 版本号: v1.0.0
- 发布日期: 2026-01-16

## 发布前检查清单

- [x] 所有测试通过
- [x] 文档已更新
- [x] CHANGELOG 已更新
- [x] 版本号已确认
- [x] 安全策略文档已创建

## 发布步骤

1. 合并此 PR 到 main 分支
2. 在 main 分支上创建标签:
   ```bash
   git checkout main
   git pull origin main
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
3. 在 GitHub 上创建 Release:
   - 访问: https://github.com/kainonly/go/releases/new
   - 选择标签: v1.0.0
   - 填写发布说明 (使用 RELEASE_NOTES_v1.0.0.md 的内容)
   - 发布

## 发布后检查

- [ ] Release 已在 GitHub 上可见
- [ ] 标签已正确创建
- [ ] Go Pkg 文档已更新
- [ ] 通知相关用户
