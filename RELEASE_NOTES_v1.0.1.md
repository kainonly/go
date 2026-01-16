# v1.0.1 Release Notes

**发布日期**: 2026-01-16
**版本类型**: 补丁版本 (Patch Release)

## 🐛 问题修复

### 构建错误修复

修复了 Go 1.24.7 环境下由 `bytedance/sonic` 库引起的严重链接错误，该错误导致 `csrf` 和 `help` 模块无法编译：

```
link: github.com/bytedance/sonic/loader: invalid reference to runtime.lastmoduledatap
FAIL	github.com/kainonly/go/csrf [build failed]
FAIL	github.com/kainonly/go/help [build failed]
```

**解决方案**:
- 升级 `github.com/bytedance/sonic` 从 v1.14.0 到 v1.14.2
- 升级相关依赖库以确保兼容性

### 测试失败修复

修复了 CI 环境中 `help` 包的 Sonyflake ID 生成测试失败问题：

```
TestSID: sonyflake: not initialized, SF is nil
TestSIDWithError: sonyflake: not initialized, SF is nil
```

**原因**: 在没有私网 IP 地址的测试环境（如容器化 CI）中，Sonyflake 无法自动获取机器标识。

**解决方案**: 为测试环境添加固定的 MachineID 配置，确保在任何环境下都能正确初始化。

## 📦 依赖更新

| 依赖包 | 旧版本 | 新版本 |
|--------|--------|--------|
| github.com/bytedance/sonic | v1.14.0 | v1.14.2 |
| github.com/bytedance/sonic/loader | v0.3.0 | v0.4.0 |
| github.com/cloudwego/base64x | v0.1.5 | v0.1.6 |
| github.com/klauspost/cpuid/v2 | v2.0.9 | v2.2.9 |

## ✅ 测试状态

所有模块测试已通过：

```
✓ captcha   - 94.4% coverage
✓ cipher    - 93.8% coverage
✓ csrf      - 100.0% coverage
✓ help      - 84.6% coverage
✓ locker    - 88.5% coverage
✓ passlib   - 97.5% coverage
✓ passport  - 100.0% coverage
✓ totp      - 68.4% coverage
✓ vd        - 85.9% coverage
```

## 📥 安装与升级

### 升级到 v1.0.1

```bash
go get -u github.com/kainonly/go@v1.0.1
```

### 安装特定模块

```bash
go get github.com/kainonly/go/csrf@v1.0.1
go get github.com/kainonly/go/help@v1.0.1
```

## 🔧 兼容性

- **Go 版本**: 1.24.0+
- **向后兼容**: 完全兼容 v1.0.0，无破坏性变更
- **API 变更**: 无

## 📝 完整变更日志

详见 [CHANGELOG.md](CHANGELOG.md)

## 🔗 相关链接

- **GitHub Release**: https://github.com/kainonly/go/releases/tag/v1.0.1
- **Go Package**: https://pkg.go.dev/github.com/kainonly/go@v1.0.1
- **问题反馈**: https://github.com/kainonly/go/issues

---

**感谢使用 Go Utils 工具库！**
