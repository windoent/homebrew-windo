# Windo CLI 工具发布指南

## 一、发布流程

### 1. 创建 Tap 仓库

在 GitLab 创建 `homebrew-windo` 仓库，克隆到本地：

```bash
git clone git@git.code.tencent.com:windo-/homebrew-windo.git
cd homebrew-windo
```

### 2. 创建 Formula 目录和文件

```bash
mkdir -p Formula
```

创建 `Formula/windo.rb` 文件（参考当前仓库中的文件）：

```ruby
class Windo < Formula
  desc "A CLI tool for generating kratos project templates"
  homepage "https://github.com/windoent/homebrew-windo"

  version "0.1.0"

  on_macos do
    on_intel do
      url "https://github.com/windoent/homebrew-windo/-/releases/v0.1.0/downloads/windo-darwin-amd64"
      sha256 "14381d1aaa3fbc46ad624a255a021f0c1610c9de0f8c26c7e4202abe1dd6cad1"
    end

    on_arm do
      url "https://github.com/windoent/homebrew-windo/-/releases/v0.1.0/downloads/windo-darwin-arm64"
      sha256 "09f0a3a3e550783e1881471e89c7147bccd5f6cc76d368c5c2b1f7cd220504e6"
    end
  end

  on_linux do
    url "https://github.com/windoent/homebrew-windo-/releases/v0.1.0/downloads/windo-linux-amd64"
    sha256 "555cfd52a27eba327fa6223653295de4b0c6288248cbabf58f93bc201d6ce5ba"
  end

  def install
    if OS.mac? && Hardware::CPU.intel?
      bin.install "windo-darwin-amd64" => "windo"
    elsif OS.mac? && Hardware::CPU.arm?
      bin.install "windo-darwin-arm64" => "windo"
    elsif OS.linux?
      bin.install "windo-linux-amd64" => "windo"
    end
  end

  test do
    system "#{bin}/windo", "version"
  end
end
```

### 3. 提交 Tap 仓库

```bash
git add .
git commit -m "Add windo formula"
git push origin main
```

### 4. 构建并发布二进制文件

#### 4.1 构建跨平台二进制文件

使用项目中的 Makefile：

```bash
cd /path/to/wd
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

#### 4.2 计算 SHA256

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

#### 4.3 创建 GitLab Release

1. 进入 `wd` 仓库页面：https://git.code.tencent.com/windo-/wd
2. 点击 **Releases** → **New release**
3. 填写：
   - **Tag name**: `v0.1.0`
   - **Release title**: `v0.1.0`
   - **Description**: 版本说明
4. 上传二进制文件（点击 **Attach files**）：
   - `windo-darwin-arm64`
   - `windo-darwin-amd64`
   - `windo-windows.exe`
   - `windo-linux-amd64` (可选)
5. 点击 **Create release**

> **注意**：腾讯工蜂（GitLab）的下载链接格式为 `/releases/v{version}/downloads/{filename}`

#### 4.4 更新 Formula 中的 SHA256

发布后，更新 `homebrew-windo/Formula/windo.rb` 中的 sha256 值。

#### 4.5 提交并推送更新

```bash
cd /path/to/homebrew-windo
git add .
git commit -m "Update sha256 for v0.1.0"
git push origin main
```

---

## 二、用户安装指南

### macOS

直接下载二进制文件：

```bash
# 下载 (根据架构选择)
curl -L https://git.code.tencent.com/windo-/homebrew-windo/-/releases/v0.0.1/downloads/windo-darwin-arm64 -o windo
# 或
curl -L https://git.code.tencent.com/windo-/homebrew-windo/-/releases/v0.0.1/downloads/windo-darwin-amd64 -o windo

# 添加执行权限
chmod +x windo

# 移动到 PATH
sudo mv windo /usr/local/bin/

# 验证
windo version
```

### Linux (Homebrew)

```bash
# 添加 tap
brew tap windo-windo/windo https://git.code.tencent.com/windo-/homebrew-windo

# 安装
brew install windo-windo/windo/windo

# 验证
windo version
```

### Windows

```powershell
irm https://git.code.tencent.com/windo-/homebrew-windo/-/releases/v0.0.1/downloads/windo-windows.exe -OutFile windo.exe
```

---

## 三、版本更新流程

当有新版本时：

1. 更新代码并测试
2. 在 `wd` 仓库执行 `make build-all` 构建二进制文件
3. 在 GitLab 创建新 Release（标签如 `v0.2.0`）
4. 上传新的二进制文件
5. 执行 `make sha256` 获取新的 SHA256 值
6. 更新 `homebrew-windo/Formula/windo.rb` 中的 `version` 和 `sha256`
7. 提交并推送 tap 仓库

```bash
# 用户更新
brew upgrade windo
```

---

## 四、项目结构

```
wd/                          # 主项目仓库
├── main.go                  # 程序入口
├── cmd/
│   ├── root.go              # 根命令
│   ├── new.go               # new 子命令
│   └── templates.go         # 项目模板
├── Makefile                 # 构建脚本 (make build-all, make sha256)
├── windo-darwin-arm64       # 编译后的二进制文件
├── windo-darwin-amd64
├── windo-windows.exe
└── windo-linux-amd64

homebrew-windo/              # 独立的 Homebrew tap 仓库
└── Formula/
    └── windo.rb             # Homebrew Formula
```

---

## 五、故障排除

### Tap 添加失败

```bash
# 检查 tap 是否正确添加
brew tap

# 手动添加
git clone https://git.code.tencent.com/windo-/homebrew-windo.git
cp -r homebrew-windo "$(brew --repository)/Library/Taps/homebrew/homebrew-cask/Casks"
```

### 安装失败

```bash
# 诊断
brew doctor
brew install -v windo

# 重新安装
brew uninstall windo
brew install windo
```

### 二进制文件下载失败

检查 Release 页面是否正确上传了对应平台的二进制文件，以及 URL 是否正确。

---

## 六、快速发布命令

```bash
# 1. 在 wd 仓库构建
cd /path/to/wd
make build-all && make sha256

# 2. 在 GitLab 上传文件并创建 Release

# 3. 更新 homebrew-windo
cd /path/to/homebrew-windo
# 编辑 Formula/windo.rb 更新 version 和 sha256
git add . && git commit -m "Release v0.x.x" && git push
```