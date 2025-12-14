# Claude Code Switcher

[English](#english) | [中文](#中文)

---

## English

A CLI tool for managing and quickly switching between multiple Claude Code API providers.

### Features

- Manage multiple API providers (Anthropic, compatible third-party services)
- Quick switching between providers
- Automatic Claude Code settings.json synchronization
- Support for JSONC format (comments in settings.json)
- Interactive CLI with intuitive prompts

### Installation

#### From Release (Recommended)

Download the latest release for your platform:

```bash
chmod +x ccs
sudo mv ccs /usr/local/bin/
```

#### From Source

```bash
git clone https://github.com/freekatz/claude-code-switch.git
cd claude-code-switch
./scripts/build.sh
sudo mv ccs /usr/local/bin/
```

### Usage

```bash
ccs [command]

Available Commands:
  add (a)       Add a new provider
  list (ls)     List providers or show provider details
  use (u)       Switch to a provider
  edit (e)      Edit a provider
  remove (rm)   Remove a provider
  help (h)      Help about any command

Flags:
  -h, --help      help for ccs
  -v, --version   version for ccs
```

### Quick Start

#### 1. Add a Provider

```bash
ccs add
# or use alias
ccs a
```

Follow the interactive prompts to configure:

- Provider name
- Provider alias (short name)
- API Base URL
- API Key
- Model names (main model, sonnet, opus, haiku, small/fast)
- Timeout (optional)

**Note**: If you don't specify individual model names, they will default to the main model.

#### 2. List Providers

```bash
ccs ls
```

Shows all configured provider aliases, with the current one marked with `*`.

```bash
ccs ls <alias>
```

Shows detailed configuration for a specific provider.

#### 3. Switch Provider

```bash
ccs use <alias>
# or
ccs u <alias>
```

Automatically updates `~/.claude/settings.json` with the selected provider's configuration.

#### 4. Edit Provider

```bash
ccs edit <alias>
# or
ccs e <alias>
```

If the provider being edited is currently active, the settings.json will be updated automatically.

#### 5. Remove Provider

```bash
ccs rm <alias>
```

### Configuration Files

- **CCS config**: `~/.config/ccs/config.json`
- **Claude Code settings**: `~/.claude/settings.json`

### Example Provider Configuration

```json
{
  "name": "Doubao",
  "alias": "db",
  "base_url": "https://ark.cn-beijing.volces.com/api/compatible",
  "api_key": "your-api-key-here",
  "model": "doubao-seed-code-preview-latest",
  "small_model": "doubao-seed-code-preview-latest",
  "sonnet_model": "doubao-seed-code-preview-latest",
  "opus_model": "doubao-seed-code-preview-latest",
  "haiku_model": "doubao-seed-code-preview-latest",
  "timeout_ms": 300000
}
```

### Build

```bash
# Development build
./scripts/build.sh

# Build with version
./scripts/build.sh v1.0.0

# Cross-compile
GOOS=linux GOARCH=amd64 ./scripts/build.sh v1.0.0 ccs-linux-amd64
```

---

## 中文

一个用于管理和快速切换多个 Claude Code API 提供商的命令行工具。

### 功能特性

- 管理多个 API 提供商（Anthropic、兼容的第三方服务）
- 快速切换提供商
- 自动同步 Claude Code settings.json
- 支持 JSONC 格式（settings.json 中的注释）
- 交互式命令行界面

### 安装

#### 从 Release 安装（推荐）

下载适合您平台的最新版本：

```bash
chmod +x ccs
sudo mv ccs /usr/local/bin/
```

#### 从源码安装

```bash
git clone https://github.com/freekatz/claude-code-switch.git
cd claude-code-switch
./scripts/build.sh
sudo mv ccs /usr/local/bin/
```

### 使用方法

```bash
ccs [命令]

可用命令:
  add (a)       添加新提供商
  list (ls)     列出提供商或显示详情
  use (u)       切换到指定提供商
  edit (e)      编辑提供商配置
  remove (rm)   删除提供商
  help (h)      显示帮助

选项:
  -h, --help      显示帮助
  -v, --version   显示版本
```

### 快速开始

#### 1. 添加提供商

```bash
ccs add
# 或使用缩写
ccs a
```

按照交互式提示配置：

- 提供商名称
- 提供商缩写
- API Base URL
- API Key
- 模型名称（主模型、sonnet、opus、haiku、小模型）
- 超时时间（可选）

**注意**：如果不指定各个模型名称，它们将默认使用主模型。

#### 2. 列出提供商

```bash
ccs ls
```

显示所有已配置的提供商别名，当前使用的标记为 `*`。

```bash
ccs ls <alias>
```

显示指定提供商的详细配置。

#### 3. 切换提供商

```bash
ccs use <alias>
# 或
ccs u <alias>
```

自动更新 `~/.claude/settings.json` 为所选提供商的配置。

#### 4. 编辑提供商

```bash
ccs edit <alias>
# 或
ccs e <alias>
```

如果编辑的是当前使用的提供商，settings.json 会自动更新。

#### 5. 删除提供商

```bash
ccs rm <alias>
```

### 配置文件

- **CCS 配置**: `~/.config/ccs/config.json`
- **Claude Code 设置**: `~/.claude/settings.json`

### 提供商配置示例

```json
{
  "name": "豆包",
  "alias": "db",
  "base_url": "https://ark.cn-beijing.volces.com/api/compatible",
  "api_key": "你的-api-key",
  "model": "doubao-seed-code-preview-latest",
  "small_model": "doubao-seed-code-preview-latest",
  "sonnet_model": "doubao-seed-code-preview-latest",
  "opus_model": "doubao-seed-code-preview-latest",
  "haiku_model": "doubao-seed-code-preview-latest",
  "timeout_ms": 300000
}
```

### 构建

```bash
# 开发构建
./scripts/build.sh

# 指定版本构建
./scripts/build.sh v1.0.0

# 交叉编译
GOOS=linux GOARCH=amd64 ./scripts/build.sh v1.0.0 ccs-linux-amd64
```
