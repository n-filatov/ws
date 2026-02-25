# ws Editor Extensions

This document is the reference for developers building editor integrations for `ws`.

## Overview

Two official extensions exist as separate repositories:

| Extension | Repo | Description |
|-----------|------|-------------|
| **ws-vscode** | `github.com/n-filatov/ws-vscode` | VS Code extension: sidebar panel, file decorations, status bar |
| **ws-zed** | `github.com/n-filatov/ws-zed` | Zed extension: `/ws` slash command for AI context, tasks |

Both extensions shell out to the `ws` CLI binary — they do **not** reimplement storage or git logic.

---

## CLI Contract

Extensions rely exclusively on these commands:

```bash
ws list              # stdout: one absolute path per line, exit 0 even when empty
ws add <file>...     # idempotent; adds absolute path(s) to current branch's working set
ws rm <file>         # removes by exact absolute path match
ws clear             # clears the entire working set for the current branch
```

**Critical**: always invoke `ws` with `cwd` set to the workspace/project root (any directory inside the git repo). The CLI runs `git rev-parse --show-toplevel` and `git rev-parse --abbrev-ref HEAD` internally to determine the correct store path. Extensions must not compute store paths themselves.

Exit behavior:
- `ws list` exits 0 and prints nothing if working set is empty
- `ws list` exits non-zero if not inside a git repo
- All write commands (`add`, `rm`, `clear`) exit 0 on success, non-zero on error

---

## Data Format

Working sets are stored at:

```
~/.local/share/ws/<repo-slug>/.workingset-<branch>
```

Where:
- `<repo-slug>` = `<repo-basename>-<sha1-of-repo-root[:8]>` (e.g., `ws-a3f2b1c4`)
- `<branch>` = branch name with `/` replaced by `-` (e.g., `feature-auth` for `feature/auth`)
- File contents: plain text, one **absolute path** per line, no trailing newlines required

Extensions that watch for external changes (e.g., files added from the terminal) should watch:
```
~/.local/share/ws/**/.workingset-*
```

The XDG base can be overridden via `$XDG_DATA_HOME`, but `~/.local/share/ws/` is the default.

---

## Architecture

```
┌─────────────────────────────┐
│   Editor Extension          │
│  (VSCode / Zed)             │
│                             │
│  • Sidebar / slash command  │
│  • File decorations         │
│  • Status bar / tasks       │
└────────────┬────────────────┘
             │ shells out (cwd = workspace root)
             ▼
┌─────────────────────────────┐
│   ws CLI binary             │
│                             │
│  ws add / rm / list / clear │
└────────────┬────────────────┘
             │ reads/writes
             ▼
┌─────────────────────────────┐
│   ~/.local/share/ws/        │
│   <repo-slug>/              │
│   .workingset-<branch>      │
└─────────────────────────────┘
```

---

## VSCode Extension (`ws-vscode`)

**Technology**: TypeScript, VSCode Extension API
**Activation**: `workspaceContains:.git`

### Features
- Sidebar tree view in the Explorer panel grouped by directory
- `WS` badge on files in the Explorer (via `FileDecorationProvider`)
- Status bar: `$(files) WS: N files (branch)`
- Commands: Add Current File, Remove Current File, Remove (from panel), Clear, Refresh
- Right-click context menu in Explorer: "Add to Working Set"
- File watcher on `~/.local/share/ws/**/.workingset-*` for auto-refresh from terminal changes
- Refresh on window focus

### Configuration
```json
{
  "ws.binaryPath": "ws"
}
```

Set to full path (e.g., `/home/user/.local/bin/ws`) if `ws` is not on VS Code's `$PATH`.

### Notable Implementation Details
- Uses `child_process.execFile` (not `exec`) to prevent shell injection on paths with spaces
- `FileItem` tree nodes use `resourceUri` for automatic language/theme icon rendering
- Badge is `'WS'` (VS Code silently truncates badges longer than 2 characters)
- Branch name is read from the built-in `vscode.git` extension API (`getAPI(1)`)
- File watcher uses `vscode.RelativePattern` with an absolute `vscode.Uri` base (requires VS Code 1.64+; engine constraint is 1.85+)
- Watcher callback is debounced ~100ms to avoid double-refresh when the extension itself writes

---

## Zed Extension (`ws-zed`)

**Technology**: Rust compiled to WASM (`wasm32-wasip1`)
**API crate**: `zed_extension_api`

### Features

**`/ws` slash command** (v1): Type `/ws` in the Zed AI panel. The extension runs `ws list`, reads each file's contents, and injects them into Claude's context as fenced code blocks with labeled sections.

**MCP context server** (v2): A standalone `ws-mcp-server` Go binary implementing the Model Context Protocol over stdio. Zed discovers it via `context_server_command()` in the extension. Exposes working set files as MCP resources (`resources/list`, `resources/read`).

**Tasks** (documented in README, user-paste):
```json
[
  { "label": "ws: add current file", "command": "ws add $ZED_FILE", "reveal": "no_focus" },
  { "label": "ws: remove current file", "command": "ws rm $ZED_FILE", "reveal": "no_focus" },
  { "label": "ws: list", "command": "ws list", "reveal": "always" },
  { "label": "ws: clear", "command": "ws clear", "reveal": "no_focus" }
]
```

**Keybindings** (documented in README, user-paste):
```json
[{ "context": "Workspace", "bindings": {
  "ctrl-alt-a": ["task::Spawn", { "task_name": "ws: add current file" }],
  "ctrl-alt-l": ["task::Spawn", { "task_name": "ws: list" }]
}}]
```

### Notable Implementation Details
- **Cannot use `std::process::Command` in WASM** — must use `zed_extension_api::process::Command` with declared `[capabilities]` in `extension.toml`
- `SlashCommandOutput` has `text: String` and `sections: Vec<SlashCommandOutputSection>` where each section has a `range: Range<usize>` into `text` and a `label: String`
- `context_server_command()` receives `&zed::Project` (not `&Worktree`)
- `schema_version = 1` is required in `extension.toml`
- A `LICENSE` file is required for the Zed marketplace (mandatory since October 2025)

---

## `ws-mcp-server` Binary

The Zed MCP context server is a separate Go binary `ws-mcp-server`. It is built from the `mcp-server/` subdirectory of the `ws-zed` repo, but distributed via the main `ws` release pipeline (added as a second build target in `.goreleaser.yaml`).

This means `brew install ws` (or the apt package) installs both `ws` and `ws-mcp-server` on the user's system.

### MCP Protocol (stdio, JSON-RPC 2.0)

```
← {"method":"initialize","params":{"capabilities":{}}}
→ {"result":{"capabilities":{"resources":{}}}}

← {"method":"resources/list"}
→ {"result":{"resources":[{"uri":"file:///abs/path","name":"filename.go"}]}}

← {"method":"resources/read","params":{"uri":"file:///abs/path"}}
→ {"result":{"contents":[{"uri":"file:///abs/path","text":"...file contents..."}]}}
```

---

## Versioning

Extensions should pin the `ws` CLI version they were tested against in their README. Breaking changes to the CLI contract (command names, exit codes, output format of `ws list`) will be noted in the main `ws` CHANGELOG with a `breaking:` prefix.
