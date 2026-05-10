# Windo 发布指南

## 一、发布流程

### 1. 构建跨平台二进制文件

使用项目中的 Makefile：

```bash
make build-all
```

或手动构建：

```bash
# macOS ARM64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o windo-darwin-arm64 .

# macOS AMD64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o windo-darwin-amd64 .

# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o windo-windows.exe .

# Linux AMD64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o windo-linux-amd64 .
```

### 2. 计算 SHA256

```bash
make sha256
```

输出示例：

```
=== SHA256 Checksums ===
darwin-arm64: 09f0a3a3e550783e1881471e89c7147bccd5f6cc76d368c5c2b1f7cd220504e6
darwin-amd64: 14381d1aaa3fbc46ad624a255a021f0c1610c9de0f8c26c7e4202abe1dd6cad1
windows: a4f5f95e87f806f3648242b1b41b95a5ccc9737840670662269a6e71f2626454
linux-amd64: 555cfd52a27eba327fa6223653295de4b0c6288248cbabf58f93bc201d6ce5ba
```

### 3. 创建 GitHub Release

1. 进入 GitHub 仓库，点击 **Releases** → **Draft a new release**
2. 填写：
   - **Tag**: `v0.1.0`
   - **Release title**: `v0.1.0`
   - **Description**: 版本说明
3. 上传二进制文件：
   - `windo-darwin-arm64`
   - `windo-darwin-amd64`
   - `windo-windows.exe`
   - `windo-linux-amd64`
4. 点击 **Publish release**

> **注意**：GitHub Release 下载链接格式为 `/releases/download/v{version}/{filename}`

### 4. 更新 Formula

更新 `homebrew-windo/Formula/windo.rb` 中的 `version` 和 `sha256`。

### 5. 提交更新

```bash
git add .
git commit -m "Release v0.x.x"
git push origin main
```

---

## 二、快速发布命令

```bash
# 1. 构建并计算 SHA256
make build-all && make sha256

# 2. 在 GitHub 上传文件并创建 Release

# 3. 更新 homebrew-windo
# 编辑 Formula/windo.rb 更新 version 和 sha256
git add . && git commit -m "Release v0.x.x" && git push
```