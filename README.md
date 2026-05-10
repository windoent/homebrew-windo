# Windo

A CLI tool for generating kratos project templates.

## 安装

### macOS / Linux (Homebrew)

```bash
# 添加 tap
brew tap windoent/windo https://github.com/windoent/homebrew-windo

# 安装
brew install windo

# 卸载
brew untap windoent/windo
```



### macOS / Linux / Windows (直接下载)

直接下载二进制文件：

| 平台 | 下载地址 |
|------|----------|
| macOS ARM64 | [windo-darwin-arm64](https://github.com/windoent/homebrew-windo/releases/download/v0.0.1/windo-darwin-arm64) |
| macOS AMD64 | [windo-darwin-amd64](https://github.com/windoent/homebrew-windo/releases/download/v0.0.1/windo-darwin-amd64) |
| Linux AMD64 | [windo-linux-amd64](https://github.com/windoent/homebrew-windo/releases/download/v0.0.1/windo-linux-amd64) |
| Windows | [windo-windows.exe](https://github.com/windoent/homebrew-windo/releases/download/v0.0.1/windo-windows.exe) |

#### macOS / Linux 安装示例

```bash
# 下载 (根据架构选择)
curl -L https://github.com/windoent/homebrew-windo/releases/download/v0.0.1/windo-darwin-arm64 -o windo

# 添加执行权限
chmod +x windo

# 移动到 PATH
sudo mv windo /usr/local/bin/

# 验证
windo version
```

#### Windows 安装示例

```powershell
irm https://github.com/windoent/homebrew-windo/releases/download/v0.0.1/windo-windows.exe -OutFile windo.exe
```

## 验证安装

```bash
windo version
```

## 卸载

```bash
brew uninstall windo
```