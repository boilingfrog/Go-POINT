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

**交互式 vs 非交互式**：默认进入交互式 REPL；加 `-p` 立刻执行后退出，适合管道。

```bash
# 1. 在当前目录开一个交互式会话（最常用）
cd ~/projects/my-app
claude

# 2. 直接把首条 prompt 写在命令行
claude "把 main.go 里的全局变量改成依赖注入"

# 3. 一次性任务：运行后立即退出（脚本/管道场景）
claude -p "总结当前 git diff 的改动"

# 4. 配合管道，把 stdin 内容当作上下文
git diff origin/master | claude -p "审查这次改动，按严重程度列出问题"

# 5. 继续上一次对话（在同一个目录下）
claude -c

# 5.1 继续上次对话并追加新 prompt
claude -c "接着把 reviewer 提的问题修一下"

# 6. 从历史会话中挑一个恢复（弹出交互式选择器）
claude -r

# 6.1 按 session UUID 直接恢复
claude -r 9f3b2e1a-...-d4

# 6.2 复用 session 时另起一个新 ID（避免覆盖原会话）
claude -c --fork-session

# 7. 给会话起一个好认的名字（显示在标题栏 / resume 列表中）
claude -n "fix-issue-123"

# 8. 从 GitHub PR 关联的会话恢复
claude --from-pr 456
```

#### 2.2 模型、Effort 与预算

不同模型在能力 / 速度 / 价格之间权衡，复杂重构选 `opus`，日常脚手架用 `sonnet`，琐碎补全可考虑 `haiku`。

```bash
# 1. 用别名切换（始终指向该系列最新模型）
claude --model opus                  # 推理最强
claude --model sonnet                # 综合最佳
claude --model haiku                 # 速度快、价格低

# 2. 用全名锁定具体版本（避免某天升级带来行为变化）
claude --model claude-sonnet-4-6

# 3. effort 控制思考深度，越高越"愿意"多步推理
claude --effort low                  # 短问答
claude --effort high                 # 复杂排查
claude --effort max -p "找出这段代码所有竞态条件"

# 4. 默认模型过载时自动降级（仅在 -p 下生效）
claude -p "review this PR" \
  --model opus \
  --fallback-model sonnet

# 5. CI 中限制单次运行最大花费（仅 -p 下生效）
claude -p "生成发布说明" --max-budget-usd 0.30

# 6. 强制提交时不超出预算（超出会被中断）
claude -p "全仓重构 logger" --max-budget-usd 5.00 --output-format json
```

#### 2.3 权限控制

权限决定 Claude 能否在不打扰你的前提下执行工具。**不同模式适用不同信任级别**：

| 模式 | 含义 | 典型场景 |
|------|------|----------|
| `default` | 危险操作弹窗确认 | 日常使用 |
| `plan` | 只读 + 出方案，不写文件 | 大改动前 review 思路 |
| `acceptEdits` | 自动接受文件编辑 | 已经批量明确的重构任务 |
| `dontAsk` | 静默执行已 allow 的工具 | 已配置好白名单的重复任务 |
| `bypassPermissions` | 全跳过（危险） | 临时沙箱、断网容器内 |

```bash
# 1. 大改动前先出方案，再决定是否真改
claude --permission-mode plan "梳理 user 模块拆成独立服务的步骤"

# 2. 已经明确改动范围，让它自动接受所有 Edit/Write
claude --permission-mode acceptEdits "把所有 fmt.Errorf 换成 errors.Wrap"

# 3. 完全沙箱（仅在断网容器/一次性 worktree 中使用！）
claude --dangerously-skip-permissions -p "..."

# 4. 临时白名单：本会话只允许这些工具，其他全部拒
claude --allowedTools "Read" "Grep" "Bash(git log:*)" "Bash(git diff:*)"

# 5. 临时黑名单：禁止某些破坏性命令
claude --disallowedTools "Bash(rm:*)" "Bash(git push:*)"

# 6. 完全限定可用工具集（最小权限）
claude --tools "Read,Grep"            # 只读模式
claude --tools ""                     # 完全禁用工具，纯对话
claude --tools "default"              # 启用所有内置工具

# 7. 命令匹配语法
#   "Bash(git *)"        → 允许任意 git 子命令
#   "Bash(git diff:*)"   → 允许 git diff 任意参数
#   "Edit"               → 允许全部 Edit 工具调用
#   "Read(/etc/**)"      → 允许读 /etc 下任意文件
```

#### 2.4 工作目录与 Worktree

```bash
# 1. 让 Claude 能访问父级 / 邻近目录（默认只允许 cwd）
claude --add-dir ../shared-libs ../docs ../../monorepo/proto

# 2. 在新建的 git worktree 中开干，不影响主分支
claude -w fix-login-bug
#   → 自动在 ../<repo>-fix-login-bug 创建 worktree
#   → 创建并 checkout 同名分支
#   → 进入该 worktree 启动 Claude

# 3. worktree + tmux：在 iTerm2 拆分面板里跑（适合并行多任务）
claude -w experiment-1 --tmux
claude -w experiment-2 --tmux        # 另开一个面板，互不干扰

# 4. 自动连接 IDE（仅当当前环境检测到一个 VS Code / JetBrains 时）
claude --ide

# 5. 指定一个特定的 session UUID（便于外部系统跟踪）
claude --session-id 11111111-2222-3333-4444-555555555555

# 6. 关闭会话持久化（不写 session 文件，仅 -p 下生效）
claude -p "查个文档" --no-session-persistence
```

#### 2.5 输出格式（脚本化与流式）

`-p` 模式下，配合 `--output-format` 可以把 Claude 当作 CLI 工具嵌入到任意脚本里。

```bash
# 1. 纯文本（默认）
claude -p "把这段 SQL 转成 GORM"

# 2. 单次 JSON：拿到结构化结果，便于 jq 解析
claude -p "审查这段代码" --output-format json | jq '.result'

# 输出大致结构：
# {
#   "type": "result",
#   "result": "...审查内容...",
#   "session_id": "...",
#   "total_cost_usd": 0.0123,
#   "duration_ms": 5421
# }

# 3. 流式 JSON：实时拿到每一条事件（适合做 UI 进度显示）
claude -p "重构这个函数" \
  --output-format stream-json \
  --include-partial-messages \
  --include-hook-events

# 4. stream-json 输入 + stream-json 输出（构造对话式 pipeline）
echo '{"type":"user","message":{"content":"hi"}}' | \
  claude -p --input-format stream-json --output-format stream-json

# 5. 使用 JSON Schema 强制结构化输出
claude -p "提取这段日志中的错误码与时间戳" \
  --json-schema '{"type":"object","properties":{"errors":{"type":"array","items":{"type":"object","properties":{"code":{"type":"string"},"ts":{"type":"string"}}}}},"required":["errors"]}'

# 6. 配合 --max-budget-usd 给批处理设上限
claude -p "$(cat task.txt)" \
  --output-format json \
  --max-budget-usd 0.20 > result.json
```

#### 2.6 Bare 模式与系统提示词

`--bare` 关掉一切"魔法"：不读 CLAUDE.md、不跑 hooks、不加载插件、不读 keychain，只走 `ANTHROPIC_API_KEY`。适合 CI / 复现 bug / 嵌入到其他工具。

```bash
# 1. 最干净的一次性运行
ANTHROPIC_API_KEY=sk-ant-xxx claude --bare -p "翻译这段文档"

# 2. bare 模式下显式指定上下文（不会自动发现 CLAUDE.md）
claude --bare \
  --add-dir ./src \
  --append-system-prompt "项目使用 Go 1.21 + Gin" \
  -p "为 user.go 写单元测试"

# 3. 完全替换默认系统提示词
claude --system-prompt "你是 SQL 优化专家，仅输出 SQL，不要解释" \
  -p "优化这条查询：SELECT * FROM orders WHERE ..."

# 4. 在默认提示词后追加（保留 Claude Code 行为，但加项目约定）
claude --append-system-prompt "所有改动必须有对应单测" \
  --permission-mode acceptEdits "修复这个 bug"

# 5. 改善多用户 prompt 缓存命中率（去除每机器差异部分）
claude --exclude-dynamic-system-prompt-sections -p "..."
```

#### 2.7 子命令速查

```bash
# 列出当前可用的 subagents
claude agents

# 管理认证（login / logout / status）
claude auth login
claude auth status

# 检查健康（更新器、MCP 服务器）
claude doctor

# 安装 / 升级 native 构建
claude install stable
claude install latest
claude update

# MCP 服务器管理
claude mcp list
claude mcp add github -- npx -y @modelcontextprotocol/server-github
claude mcp remove github

# 插件管理
claude plugin list
claude plugin install <plugin>

# 申请长期 token（适合 CI）
claude setup-token
```

---

### 三、交互式 Slash 命令

在交互式会话中输入 `/` 可以查看所有可用命令；输入 `/<字母>` 会做模糊筛选。

#### 3.1 命令速查

| 命令 | 说明 |
|------|------|
| `/help` | 查看帮助 |
| `/clear` | 清除当前对话上下文（重置 token，但保留 CLAUDE.md） |
| `/compact` | 压缩对话上下文（保留要点，节省 token） |
| `/config` | 打开配置面板（主题、模型等） |
| `/model` | 切换当前会话模型 |
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

#### 3.2 核心命令使用示例

**`/clear` vs `/compact`**：长会话的两种瘦身方式。

```
> /clear
# 清空全部历史，相当于重启会话。CLAUDE.md / settings 不会丢。
# 适合：完全切换到一个无关任务时使用。

> /compact
# 把已有对话压缩成摘要塞回上下文，保留关键决策与文件路径。
# 适合：同一任务做了很久、上下文逼近上限时延长会话。

> /compact 重点保留刚才讨论的 schema 设计决策
# 可以追加指令，告诉它哪些信息更值得保留。
```

**`/model`**：会话中途切换模型。

```
> /model
  ◯ claude-opus-4-7
  ● claude-sonnet-4-6   (current)
  ◯ claude-haiku-4-5

> /model opus
# 切到 Opus 处理一个复杂重构，做完再 /model sonnet 切回。
```

**`/init`**：为新项目一键生成 CLAUDE.md。

```
> /init
# Claude 会扫描项目结构、识别技术栈、读 README，
# 在根目录生成一份 CLAUDE.md 草稿，列出技术栈、构建命令、目录结构等。
# 生成后人工补充团队约定 / 注意事项。
```

**`/review`**：审查当前 PR / 分支改动。

```
> /review
# 自动检测当前分支对应的 PR，逐条点评。

> /review #123
# 指定 PR 号。

> /review --focus security
# 只看安全问题。
```

**`/security-review`**：聚焦安全的深度审查。

```
> /security-review
# 针对当前分支 vs main 的 diff，做 OWASP / 注入 / 鉴权 / 密钥泄露专项审查。
# 适合：合并前的最后一道关卡。
```

**`/agents`**：管理 subagents。

```
> /agents
  list      列出当前可用的 agents
  create    交互式创建一个新 agent
  edit      编辑现有 agent

> /agents create
  Name: db-migrator
  Description: Generate and review SQL migrations
  Tools: Read, Edit, Bash(make migrate:*)
  Prompt: 你是 DB migration 专家，遵循 expand-contract 模式...
# 创建后保存到 .claude/agents/db-migrator.md，团队共享。
```

**`/mcp`**：MCP 服务器管理。

```
> /mcp
  list      已加载的 MCP 服务器
  status    每个服务器的健康状态
  reload    重载 MCP 配置

> /mcp status
  github       ✓ connected    (12 tools)
  postgres     ✗ failed       Error: ECONNREFUSED 5432
  filesystem   ✓ connected    (5 tools)
```

**`/cost`**：查看花费。

```
> /cost
  Session cost:    $0.42
  Total tokens:    127,453 (in: 98,231 / out: 29,222)
  Cache hit rate:  68%
  Duration:        14m 22s
# 会话结束前看一眼，方便估算批处理任务成本。
```

**`/resume`**：恢复历史会话。

```
> /resume
# 弹出过去 7 天的会话列表，按时间倒序，可搜索。
  [1] fix-issue-123          2h ago     opus      $0.34
  [2] refactor-logger        yesterday  sonnet    $0.12
  [3] explore-auth-module    2d ago     opus      $0.89
```

#### 3.3 输入前缀

除了 `/`，交互式 prompt 还支持两个前缀：

```
# !<command>：直接执行 shell 命令，输出回到对话上下文
> !go test ./...
ok      myapp/internal/user     0.234s
FAIL    myapp/internal/order    0.456s
> 修一下上面 order 包失败的测试
# Claude 现在能看到 test 输出，可以直接定位问题。

# #<text>：把内容追加到 CLAUDE.md
> # 数据库迁移必须使用 expand-contract 模式
# 立刻被写入项目 CLAUDE.md，下次启动也能看到。

# @<file/dir>：把文件或目录显式拉入上下文（@ 触发自动补全）
> 帮我对比 @internal/user/service.go 和 @internal/order/service.go 的错误处理风格
```

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

#### 5.1 配置加载顺序与优先级

Claude Code 按以下顺序合并配置，**后加载的覆盖前面的**：

| 层级 | 路径 | 用途 | 是否提交 git |
|------|------|------|---------------|
| 1. 用户级 | `~/.claude/settings.json` | 个人偏好（主题、模型、全局白名单） | 否（在 home 目录） |
| 2. 项目级 | `<repo>/.claude/settings.json` | 团队共享约定 | **是** |
| 3. 本地级 | `<repo>/.claude/settings.local.json` | 个人对项目的覆盖 | **否**（加 .gitignore） |
| 4. CLI 临时 | `--settings <file-or-json>` | 单次会话覆盖 | — |

> 规则：`allow` 列表会取并集，`deny` 列表也是并集（任一层 deny 即 deny）；其他字段后覆盖前。

#### 5.2 完整配置字段速查

```json
{
  // ─── 基础 ─────────────────────────────────────────
  "model": "opus",                    // 默认模型：opus / sonnet / haiku / 全名
  "theme": "dark",                    // 主题：dark / light / dark-daltonized
  "agent": "reviewer",                // 默认 subagent
  "verbose": false,                   // 是否输出详细日志
  "effort": "medium",                 // 思考深度：low / medium / high / xhigh / max

  // ─── 权限 ─────────────────────────────────────────
  "permissions": {
    "allow": [
      "Read", "Edit", "Write",
      "Bash(git status)",
      "Bash(git diff:*)",
      "Bash(git log:*)",
      "Bash(go test:*)",
      "Bash(make:*)",
      "WebFetch(domain:docs.anthropic.com)"
    ],
    "deny": [
      "Bash(rm -rf:*)",
      "Bash(git push --force:*)",
      "Bash(git reset --hard:*)",
      "Read(/etc/**)",                // 拒绝读敏感目录
      "Read(.env*)"                   // 拒绝读环境变量文件
    ],
    "defaultMode": "default"          // default / plan / acceptEdits / dontAsk / bypassPermissions
  },

  // ─── 环境变量（注入到工具子进程）─────────────────
  "env": {
    "ANTHROPIC_LOG_LEVEL": "info",
    "GOFLAGS": "-mod=vendor",
    "EDITOR": "code --wait"
  },

  // ─── Hooks（见第六节）─────────────────────────────
  "hooks": { /* ... */ },

  // ─── MCP 服务器 ───────────────────────────────────
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-github"]
    }
  },

  // ─── 状态栏 ──────────────────────────────────────
  "statusLine": {
    "type": "command",
    "command": "echo \"$(git branch --show-current) | $(go version | awk '{print $3}')\""
  },

  // ─── 会话默认参数 ─────────────────────────────────
  "permissionMode": "default",        // 同 --permission-mode
  "maxBudgetUsd": 5.00,               // 单会话预算上限
  "addDirs": ["../shared-libs"],      // 额外允许访问的目录

  // ─── 排除文件 ─────────────────────────────────────
  "exclude": [
    "node_modules/**",
    "vendor/**",
    "**/*.min.js",
    "dist/**"
  ],

  // ─── 自动更新 ─────────────────────────────────────
  "autoUpdate": true,
  "autoUpdaterStatus": "enabled"
}
```

#### 5.3 团队 vs 个人的配置分工

```jsonc
// .claude/settings.json — 团队共享，提交 git
{
  "permissions": {
    "allow": ["Read", "Edit", "Bash(make:*)", "Bash(go test:*)"],
    "deny":  ["Bash(rm -rf:*)", "Bash(git push --force:*)"]
  },
  "hooks": {
    "PostToolUse": [
      { "matcher": "Edit|Write",
        "hooks": [{ "type": "command", "command": "make fmt" }] }
    ]
  }
}

// .claude/settings.local.json — 个人覆盖，加 .gitignore
{
  "model": "opus",                    // 我个人偏好用 opus
  "permissions": {
    "allow": ["Bash(docker:*)"]       // 只我本地需要 docker
  },
  "env": {
    "DEBUG": "true"
  }
}
```

#### 5.4 用 `--settings` 临时覆盖

```bash
# 1. 直接传 JSON 字符串
claude --settings '{"model":"opus","permissions":{"defaultMode":"plan"}}'

# 2. 传文件路径
claude --settings ./ci-settings.json -p "review this PR"

# 3. 控制加载哪些层级（CI 中只用项目级，忽略个人级）
claude --setting-sources project,local -p "..."
```

---

### 六、Hooks：事件钩子

Hooks 是 settings.json 中声明的**外部 shell 命令**，由 Claude Code harness（不是模型）在关键事件触发。这点很重要——如果你想要"每次都……"的强约束行为，必须用 hook 实现，让模型记住是不可靠的。

#### 6.1 全部事件

| 事件 | 触发时机 | 是否能阻断 |
|------|----------|------------|
| `SessionStart` | 会话启动时 | 否 |
| `UserPromptSubmit` | 用户提交消息时 | 是（exit 2 拒绝） |
| `PreToolUse` | 工具调用前 | **是**（exit 2 阻断） |
| `PostToolUse` | 工具调用后 | 否 |
| `Notification` | 系统通知时（如等用户输入） | 否 |
| `Stop` | Claude 完成响应时 | 是（exit 2 强制让 Claude 继续） |
| `SubagentStop` | subagent 完成时 | 是 |
| `SessionEnd` | 会话退出时 | 否 |
| `PreCompact` | 自动压缩上下文前 | 否 |

#### 6.2 Hook 配置结构

```json
{
  "hooks": {
    "<EventName>": [
      {
        "matcher": "<正则匹配工具名或 prompt>",
        "hooks": [
          {
            "type": "command",
            "command": "<shell 命令>",
            "timeout": 10000
          }
        ]
      }
    ]
  }
}
```

#### 6.3 注入的环境变量

Hook 命令可以读取以下环境变量：

| 变量 | 含义 |
|------|------|
| `$CLAUDE_TOOL_NAME` | 触发的工具名（Edit / Write / Bash …） |
| `$CLAUDE_TOOL_INPUT` | 工具输入（JSON 字符串） |
| `$CLAUDE_FILE_PATHS` | 受影响的文件路径（空格分隔） |
| `$CLAUDE_PROJECT_DIR` | 项目根目录 |
| `$CLAUDE_USER_PROMPT` | 用户最近提交的消息（仅 UserPromptSubmit） |
| `$CLAUDE_SESSION_ID` | 当前会话 ID |
| `$CLAUDE_NOTIFICATION` | 通知内容（仅 Notification） |

#### 6.4 退出码语义

| 退出码 | 行为 |
|--------|------|
| `0` | 通过，stdout 不会反馈给 Claude |
| `2` | **阻断**事件；stderr 内容作为反馈传给 Claude |
| 其他非零 | 失败但不阻断，stderr 记录到日志 |

#### 6.5 实战示例

**示例 1：保存 Go 文件后自动 gofmt**

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Write|Edit",
        "hooks": [
          {
            "type": "command",
            "command": "for f in $CLAUDE_FILE_PATHS; do [[ $f == *.go ]] && gofmt -w \"$f\"; done"
          }
        ]
      }
    ]
  }
}
```

**示例 2：阻止任何 `rm -rf /` 类的危险命令**

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "echo \"$CLAUDE_TOOL_INPUT\" | jq -r '.command' | grep -qE 'rm\\s+-rf?\\s+/' && { echo '禁止 rm -rf 根目录' >&2; exit 2; } || exit 0"
          }
        ]
      }
    ]
  }
}
```

**示例 3：UserPromptSubmit 自动注入项目状态**

```json
{
  "hooks": {
    "UserPromptSubmit": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "echo \"[当前分支: $(git branch --show-current) | 未提交改动: $(git status --porcelain | wc -l) 个文件]\""
          }
        ]
      }
    ]
  }
}
```

> stdout 内容会作为系统提醒注入给 Claude，让它每次回复都能看到 git 状态。

**示例 4：Stop 时桌面通知 + 播放提示音**

```json
{
  "hooks": {
    "Stop": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "osascript -e 'display notification \"Claude 已完成\" with title \"Claude Code\"' && afplay /System/Library/Sounds/Glass.aiff"
          }
        ]
      }
    ]
  }
}
```

**示例 5：写入前强制运行单测（PreToolUse 阻断）**

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Write|Edit",
        "hooks": [
          {
            "type": "command",
            "command": "echo \"$CLAUDE_FILE_PATHS\" | grep -q '_test.go$' || { go test ./... > /dev/null 2>&1 || { echo '当前测试未通过，先修测试再改代码' >&2; exit 2; }; }"
          }
        ]
      }
    ]
  }
}
```

**示例 6：SessionStart 把今日 TODO 注入上下文**

```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "[ -f ~/today.md ] && cat ~/today.md || true"
          }
        ]
      }
    ]
  }
}
```

#### 6.6 调试 Hooks

```bash
# 启动时开 hooks 调试
claude --debug hooks

# 在 stream-json 输出中包含 hook 事件
claude -p "..." --output-format stream-json --include-hook-events

# 查看 hook 实际收到的环境变量
# 在 hook 里加：env > /tmp/hook-env.log
```

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

#### 10.1 让 Claude 接管一个 Bug 修复（端到端）

```bash
cd my-project

# 用 worktree 隔离，避免污染主分支
claude -w fix-issue-123 --model opus
```

进入交互式会话后：

```
> 看一下 GitHub issue #123，分析根因
# Claude 会通过 gh CLI 拉 issue 内容，然后定位相关代码

> /agents
# 切到 Plan agent 让它先出方案，避免直接乱改

> @internal/order/service.go @internal/order/handler.go
> 这两个文件就是改动重点，请先给出修复方案

> 方案 OK，按方案 2 实现，并补一个能重现 bug 的回归测试

> !go test ./internal/order/... -run TestIssue123 -v
# 让 Claude 看到测试结果

> 测试过了。再跑一下完整测试和 lint
> !make test && make lint

> 提交：commit message 用 "fix: handle nil pointer in order calc (#123)"
# Claude 会调用 git add / commit

> /cost
# 看一下这次任务花了多少钱

> 推送并开 PR，PR 描述里包含根因分析和测试覆盖
```

#### 10.2 在 CI 中使用 Claude 做代码审查

`.github/workflows/claude-review.yml`：

```yaml
name: Claude Review
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Install Claude Code
        run: npm install -g @anthropic-ai/claude-code

      - name: Run review
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          git diff origin/main...HEAD > /tmp/diff.patch
          claude -p "$(cat <<EOF
          审查下列 diff，按以下格式输出：
          ## 必须修
          - <问题> (file:line)
          ## 建议
          - <问题> (file:line)
          ## 通过
          - <亮点>

          $(cat /tmp/diff.patch)
          EOF
          )" \
            --bare \
            --model sonnet \
            --output-format json \
            --max-budget-usd 0.50 \
            --permission-mode plan > review.json

      - name: Post comment
        run: |
          jq -r '.result' review.json | \
            gh pr comment ${{ github.event.pull_request.number }} -F -
```

#### 10.3 隔离 worktree 做实验

```bash
# 1. 起一个 worktree，用最激进的权限模式自由发挥
claude -w experiment-refactor \
  --model opus \
  --effort high \
  --permission-mode acceptEdits

# 2. 实验失败？直接删 worktree，主分支毫发无损
git worktree remove ../my-project-experiment-refactor
git branch -D experiment-refactor

# 3. 实验成功？merge 回主分支
cd ../my-project
git merge experiment-refactor
```

#### 10.4 批处理：一次性整理一批文件

```bash
# 用 shell 循环 + -p 模式，把几十个文件挨个让 Claude 处理
for f in $(find ./docs -name "*.md"); do
  claude -p "把 $f 中的中英文之间补上空格，结果直接覆盖原文件" \
    --add-dir ./docs \
    --permission-mode acceptEdits \
    --max-budget-usd 0.05 \
    --bare
done
```

#### 10.5 把 Claude 当作"命令行小工具"

```bash
# 起个 alias，把 Claude 当 SQL / 正则 / 文本处理器用
alias ai='claude -p --bare --model haiku'

# 用法
echo "SELECT * FROM users WHERE created > '2024-01-01'" | ai "把这条 SQL 改成 GORM 链式调用"

cat error.log | ai "提取所有 5xx 错误，按 path 分组统计"

git log --oneline -20 | ai "总结最近 20 个 commit 的主题分布"
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
