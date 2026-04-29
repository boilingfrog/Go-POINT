## Claude Code 使用指南

### 前言

AI 编程最近很火，随之而来的是各家公司裁员的消息。

作为 AI 编程的头部，`Claude Code` 学习如何使用是一件很重要的事情。

本文从安装、命令行参数、交互式命令、配置文件、Hooks、MCP、Subagents、Plugins 等多个维度，系统梳理 Claude Code 的使用方式。

---

### 一、安装与认证

#### 1.1 安装

```bash
# npm 全局安装
npm install -g @anthropic-ai/claude-code

# 或使用 native 构建
claude install stable
```

#### 1.2 认证方式

```bash
# 方式一：交互式登录（推荐，使用 Claude 订阅）
claude auth login

# 方式二：使用长期 token（适合 CI/CD）
claude setup-token

# 方式三：使用 API Key（环境变量）
export ANTHROPIC_API_KEY=sk-ant-xxx
```

#### 1.3 健康检查

```bash
# 检查自动更新器、MCP 服务器健康状态
claude doctor

# 检查更新并安装
claude update
```

---

### 二、基础命令与命令行参数

在 `Claude Code` 中输入 `/` 可以查看所有可用命令，或输入 `/` 后跟任何字母来筛选。

完整 `claude -h` 输出：

```
% claude -h
Usage: claude [options] [command] [prompt]

Claude Code - starts an interactive session by default, use -p/--print for non-interactive output

Arguments:
  prompt                                            Your prompt

Options:
  --add-dir <directories...>                        Additional directories to allow tool access to
  --agent <agent>                                   Agent for the current session. Overrides the 'agent' setting.
  --agents <json>                                   JSON object defining custom agents (e.g. '{"reviewer": {"description": "Reviews code", "prompt": "You are a code reviewer"}}')
  --allow-dangerously-skip-permissions              Enable bypassing all permission checks as an option, without it being enabled by default. Recommended only for sandboxes with no
                                                    internet access.
  --allowedTools, --allowed-tools <tools...>        Comma or space-separated list of tool names to allow (e.g. "Bash(git *) Edit")
  --append-system-prompt <prompt>                   Append a system prompt to the default system prompt
  --bare                                            Minimal mode: skip hooks, LSP, plugin sync, attribution, auto-memory, background prefetches, keychain reads, and CLAUDE.md
                                                    auto-discovery. Sets CLAUDE_CODE_SIMPLE=1. Anthropic auth is strictly ANTHROPIC_API_KEY or apiKeyHelper via --settings (OAuth and
                                                    keychain are never read). 3P providers (Bedrock/Vertex/Foundry) use their own credentials. Skills still resolve via /skill-name.
                                                    Explicitly provide context via: --system-prompt[-file], --append-system-prompt[-file], --add-dir (CLAUDE.md dirs), --mcp-config,
                                                    --settings, --agents, --plugin-dir.
  --betas <betas...>                                Beta headers to include in API requests (API key users only)
  --brief                                           Enable SendUserMessage tool for agent-to-user communication
  --chrome                                          Enable Claude in Chrome integration
  -c, --continue                                    Continue the most recent conversation in the current directory
  --dangerously-skip-permissions                    Bypass all permission checks. Recommended only for sandboxes with no internet access.
  -d, --debug [filter]                              Enable debug mode with optional category filtering (e.g., "api,hooks" or "!1p,!file")
  --debug-file <path>                               Write debug logs to a specific file path (implicitly enables debug mode)
  --disable-slash-commands                          Disable all skills
  --disallowedTools, --disallowed-tools <tools...>  Comma or space-separated list of tool names to deny (e.g. "Bash(git *) Edit")
  --effort <level>                                  Effort level for the current session (low, medium, high, xhigh, max)
  --exclude-dynamic-system-prompt-sections          Move per-machine sections (cwd, env info, memory paths, git status) from the system prompt into the first user message. Improves
                                                    cross-user prompt-cache reuse. Only applies with the default system prompt (ignored with --system-prompt). (default: false)
  --fallback-model <model>                          Enable automatic fallback to specified model when default model is overloaded (only works with --print)
  --file <specs...>                                 File resources to download at startup. Format: file_id:relative_path (e.g., --file file_abc:doc.txt file_def:img.png)
  --fork-session                                    When resuming, create a new session ID instead of reusing the original (use with --resume or --continue)
  --from-pr [value]                                 Resume a session linked to a PR by PR number/URL, or open interactive picker with optional search term
  -h, --help                                        Display help for command
  --ide                                             Automatically connect to IDE on startup if exactly one valid IDE is available
  --include-hook-events                             Include all hook lifecycle events in the output stream (only works with --output-format=stream-json)
  --include-partial-messages                        Include partial message chunks as they arrive (only works with --print and --output-format=stream-json)
  --input-format <format>                           Input format (only works with --print): "text" (default), or "stream-json" (realtime streaming input) (choices: "text",
                                                    "stream-json")
  --json-schema <schema>                            JSON Schema for structured output validation. Example: {"type":"object","properties":{"name":{"type":"string"}},"required":["name"]}
  --max-budget-usd <amount>                         Maximum dollar amount to spend on API calls (only works with --print)
  --mcp-config <configs...>                         Load MCP servers from JSON files or strings (space-separated)
  --mcp-debug                                       [DEPRECATED. Use --debug instead] Enable MCP debug mode (shows MCP server errors)
  --model <model>                                   Model for the current session. Provide an alias for the latest model (e.g. 'sonnet' or 'opus') or a model's full name (e.g.
                                                    'claude-sonnet-4-6').
  -n, --name <name>                                 Set a display name for this session (shown in the prompt box, /resume picker, and terminal title)
  --no-chrome                                       Disable Claude in Chrome integration
  --no-session-persistence                          Disable session persistence - sessions will not be saved to disk and cannot be resumed (only works with --print)
  --output-format <format>                          Output format (only works with --print): "text" (default), "json" (single result), or "stream-json" (realtime streaming) (choices:
                                                    "text", "json", "stream-json")
  --permission-mode <mode>                          Permission mode to use for the session (choices: "acceptEdits", "auto", "bypassPermissions", "default", "dontAsk", "plan")
  --plugin-dir <path>                               Load plugins from a directory for this session only (repeatable: --plugin-dir A --plugin-dir B) (default: [])
  -p, --print                                       Print response and exit (useful for pipes). Note: The workspace trust dialog is skipped when Claude is run with the -p mode. Only use
                                                    this flag in directories you trust.
  --remote-control-session-name-prefix <prefix>     Prefix for auto-generated Remote Control session names (default: hostname)
  --replay-user-messages                            Re-emit user messages from stdin back on stdout for acknowledgment (only works with --input-format=stream-json and
                                                    --output-format=stream-json)
  -r, --resume [value]                              Resume a conversation by session ID, or open interactive picker with optional search term
  --session-id <uuid>                               Use a specific session ID for the conversation (must be a valid UUID)
  --setting-sources <sources>                       Comma-separated list of setting sources to load (user, project, local).
  --settings <file-or-json>                         Path to a settings JSON file or a JSON string to load additional settings from
  --strict-mcp-config                               Only use MCP servers from --mcp-config, ignoring all other MCP configurations
  --system-prompt <prompt>                          System prompt to use for the session
  --tmux                                            Create a tmux session for the worktree (requires --worktree). Uses iTerm2 native panes when available; use --tmux=classic for
                                                    traditional tmux.
  --tools <tools...>                                Specify the list of available tools from the built-in set. Use "" to disable all tools, "default" to use all tools, or specify tool
                                                    names (e.g. "Bash,Edit,Read").
  --verbose                                         Override verbose mode setting from config
  -v, --version                                     Output the version number
  -w, --worktree [name]                             Create a new git worktree for this session (optionally specify a name)

Commands:
  agents [options]                                  List configured agents
  auth                                              Manage authentication
  auto-mode                                         Inspect auto mode classifier configuration
  doctor                                            Check the health of your Claude Code auto-updater. Note: The workspace trust dialog is skipped and stdio servers from .mcp.json are
                                                    spawned for health checks. Only use this command in directories you trust.
  install [options] [target]                        Install Claude Code native build. Use [target] to specify version (stable, latest, or specific version)
  mcp                                               Configure and manage MCP servers
  plugin|plugins                                    Manage Claude Code plugins
  setup-token                                       Set up a long-lived authentication token (requires Claude subscription)
  update|upgrade                                    Check for updates and install if available
```

#### 2.1 基础启动

```bash
# 交互式会话
claude

# 直接传入 prompt
claude "帮我重构这个文件"

# 非交互式（管道模式），打印结果后退出
claude -p "总结当前 git diff"

# 继续上一次对话
claude -c

# 恢复历史会话（交互式选择）
claude -r
```

#### 2.2 模型与 Effort

```bash
# 指定模型
claude --model opus
claude --model sonnet
claude --model claude-sonnet-4-6

# 指定 effort 级别（low/medium/high/xhigh/max）
claude --effort high

# 当默认模型过载时自动 fallback
claude --fallback-model sonnet -p "..."
```

#### 2.3 权限控制

```bash
# 指定权限模式
claude --permission-mode plan          # 仅规划，不执行
claude --permission-mode acceptEdits   # 自动接受编辑
claude --permission-mode bypassPermissions  # 跳过所有权限检查（慎用）

# 仅允许特定工具
claude --allowedTools "Bash(git *)" "Edit" "Read"

# 禁用特定工具
claude --disallowedTools "Bash(rm *)"

# 仅启用指定的内置工具
claude --tools "Bash,Edit,Read"
```

#### 2.4 上下文与目录

```bash
# 添加额外的允许访问目录
claude --add-dir ../shared-libs ../docs

# 创建 git worktree 进行隔离开发
claude -w feature-branch

# 配合 tmux 使用
claude -w feature-x --tmux
```

#### 2.5 输出格式（脚本化使用）

```bash
# JSON 单次结果
claude -p "审查这段代码" --output-format json

# 流式 JSON 输出（适合管道处理）
claude -p "..." --output-format stream-json --include-partial-messages

# 控制最大花费
claude -p "..." --max-budget-usd 0.50
```

#### 2.6 Bare 模式（极简）

```bash
# 跳过 hooks/LSP/插件同步/auto-memory 等
claude --bare -p "纯净环境运行"
```

---

### 三、交互式 Slash 命令

在交互式会话中输入 `/` 可以查看所有可用命令。常用如下：

| 命令 | 说明 |
|------|------|
| `/help` | 查看帮助 |
| `/clear` | 清除当前对话上下文 |
| `/compact` | 压缩对话上下文 |
| `/config` | 打开配置面板（主题、模型等） |
| `/model` | 切换模型 |
| `/agents` | 管理 subagents |
| `/mcp` | 管理 MCP 服务器 |
| `/plugin` | 管理插件 |
| `/resume` | 恢复历史会话 |
| `/review` | 审查 PR |
| `/security-review` | 对当前分支变更做安全审查 |
| `/init` | 为当前项目生成 CLAUDE.md |
| `/cost` | 查看本次会话花费 |
| `/login` `/logout` | 登录登出 |
| `/exit` | 退出 |

输入 `!<命令>` 可以在会话中直接执行 shell 命令，输出会带回上下文。

输入 `#<内容>` 可以快速添加内容到 CLAUDE.md。

---

### 四、CLAUDE.md：项目记忆

`CLAUDE.md` 是 Claude Code 在项目根目录读取的"长期记忆"文件，用来沉淀项目约定。

#### 4.1 文件查找顺序

1. 当前目录 `./CLAUDE.md`（项目级）
2. `~/.claude/CLAUDE.md`（用户级，所有项目共享）
3. `.claude/CLAUDE.md`（项目本地，可加入 .gitignore）

#### 4.2 推荐内容

```markdown
# 项目说明

## 技术栈
- Go 1.21
- MySQL 8.0
- Kubernetes + Helm

## 代码规范
- 错误处理使用 errors.Wrap
- 不要使用 panic，统一用 error 返回

## 常用命令
- 构建：make build
- 测试：make test
- 部署：./script/deploy.sh

## 注意事项
- 修改 schema 必须同时生成 migration
- 提交前运行 make lint
```

通过 `/init` 命令可以让 Claude 自动扫描项目并生成初版 CLAUDE.md。

---

### 五、settings.json 配置

配置文件位置：

- 用户级：`~/.claude/settings.json`
- 项目级：`.claude/settings.json`（团队共享，提交到 git）
- 本地级：`.claude/settings.local.json`（个人覆盖，不提交）

#### 5.1 常用配置示例

```json
{
  "model": "opus",
  "theme": "dark",
  "permissions": {
    "allow": [
      "Bash(git status)",
      "Bash(git diff:*)",
      "Bash(go test:*)",
      "Read",
      "Edit"
    ],
    "deny": [
      "Bash(rm -rf:*)"
    ]
  },
  "env": {
    "ANTHROPIC_LOG_LEVEL": "info"
  },
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [
          {
            "type": "command",
            "command": "gofmt -w $CLAUDE_FILE_PATHS"
          }
        ]
      }
    ]
  }
}
```

---

### 六、Hooks：事件钩子

Hooks 让你在工具调用前后、用户提交消息时等关键事件触发自定义脚本。

#### 6.1 主要事件

| 事件 | 触发时机 |
|------|----------|
| `PreToolUse` | 工具调用前（可拦截） |
| `PostToolUse` | 工具调用后 |
| `UserPromptSubmit` | 用户提交消息时 |
| `Stop` | Claude 完成响应时 |
| `Notification` | 系统通知时 |

#### 6.2 实战例子：保存前自动格式化 Go 代码

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Write|Edit",
        "hooks": [
          {
            "type": "command",
            "command": "for f in $CLAUDE_FILE_PATHS; do [[ $f == *.go ]] && gofmt -w $f; done"
          }
        ]
      }
    ]
  }
}
```

#### 6.3 阻止危险命令

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "echo \"$CLAUDE_TOOL_INPUT\" | grep -q 'rm -rf /' && exit 2 || exit 0"
          }
        ]
      }
    ]
  }
}
```

> 退出码 2 会阻断工具执行，并将 stderr 反馈给 Claude。

---

### 七、Subagents：专用子代理

Subagents 是为特定任务设计的独立子代理，拥有自己的系统提示词和工具集，主上下文不会被它们的搜索结果污染。

#### 7.1 内置 Subagents

- `general-purpose`：通用研究 / 多步任务
- `Explore`：只读快速代码搜索
- `Plan`：架构设计 / 实现规划

#### 7.2 自定义 Subagent

在 `.claude/agents/reviewer.md` 中定义：

```markdown
---
name: reviewer
description: Reviews Go code for idiomatic patterns and concurrency bugs
tools: Read, Grep, Bash
---

You are a senior Go reviewer. Focus on:
- goroutine leaks
- improper context usage
- error wrapping
Report findings concisely.
```

通过 `/agents` 命令管理，或在命令行用 `--agent reviewer` 指定。

---

### 八、MCP 服务器

MCP（Model Context Protocol）让 Claude 可以接入外部工具：数据库、Slack、GitHub、Linear 等。

#### 8.1 添加 MCP 服务器

```bash
# 查看已配置的 MCP 服务器
claude mcp list

# 添加（示例：GitHub）
claude mcp add github -- npx -y @modelcontextprotocol/server-github

# 通过配置文件加载
claude --mcp-config ./mcp.json
```

#### 8.2 MCP 配置示例

`.mcp.json`：

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"]
    },
    "postgres": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-postgres"],
      "env": {
        "DATABASE_URL": "postgres://localhost/mydb"
      }
    }
  }
}
```

---

### 九、Plugins：插件系统

插件可以批量分发 slash 命令、agents、hooks、MCP 配置。

```bash
# 查看插件
claude plugin list

# 加载本地目录的插件（仅当前会话）
claude --plugin-dir ./my-plugin
```

---

### 十、典型工作流示例

#### 10.1 让 Claude 接管一个 Bug 修复

```bash
cd my-project
claude
```

```
> 帮我分析 #123 这个 issue 的根因
> /review 看看现有 PR 的实现思路
> 实现修复并写一个回归测试
> 跑一下测试
> 提交并推送
```

#### 10.2 在 CI 中使用 Claude 做代码审查

`.github/workflows/claude-review.yml`：

```yaml
- name: Claude Review
  env:
    ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
  run: |
    claude -p "审查 PR #${{ github.event.pull_request.number }}，重点关注安全和性能" \
      --output-format json \
      --max-budget-usd 1.00 \
      --permission-mode plan > review.json
```

#### 10.3 隔离 worktree 做实验

```bash
# 在新 worktree 中尝试重构，不影响主分支
claude -w experiment-refactor --permission-mode acceptEdits
```

---

### 十一、最佳实践

1. **从 `/init` 开始**：让 Claude 先生成 CLAUDE.md，再人工补充团队约定。
2. **善用 `--permission-mode plan`**：在改动大的任务前，先让 Claude 出方案。
3. **配置 hooks 兜底**：自动 lint / fmt，防止 Claude 写出不规范代码。
4. **危险目录用 worktree 隔离**：`-w` 参数隔离实验性改动。
5. **定期 `/clear` 或 `/compact`**：长会话注意上下文管理，节省 token。
6. **生产 CI 用 `-p` + `--max-budget-usd`**：控制成本。
7. **敏感命令进 deny 列表**：通过 `permissions.deny` 兜底拦截 `rm -rf`、`git push --force` 等。
8. **CLAUDE.md 写"为什么"**：记录约定背后的动机，比单纯写规则更有效。

---

### 参考

- [官方文档 - 内置命令](https://code.claude.com/docs/zh-CN/commands)
- [Hooks 文档](https://code.claude.com/docs/zh-CN/hooks)
- [MCP 协议](https://modelcontextprotocol.io)
- [Subagents 指南](https://code.claude.com/docs/zh-CN/sub-agents)
