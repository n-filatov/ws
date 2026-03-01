# ws - AI-Native Working Set Manager

## One-Pager

### What is ws?

**ws** is a terminal UI tool that keeps track of files relevant to your current git branch. It's designed for AI-assisted development, helping you maintain context across sessions and branches.

### The Problem

When building features, developers work with the same 6-10 files constantly. But their editors show hundreds. They ask AI assistants (Claude, ChatGPT, Cursor) to identify relevant files — and get a list — but then what? They open them one by one, lose track, and repeat the next day.

### The Solution

**ws** maintains a focused, branch-scoped working set of files.

```bash
ws add src/auth/login.go src/auth/jwt.go
ws  # Opens TUI with tree view, fuzzy search, git status
```

Each git branch has its own working set. Switch branches, your context switches too.

### Key Features

- **Branch-scoped**: Each branch has its own working set
- **Persistent**: Survives terminal sessions
- **Git-aware**: Shows inline status, auto-syncs modified files
- **AI-native**: Built for AI pair programming workflows
- **Terminal-first**: Works with any editor or AI assistant

### Use Cases

- **AI pair programming**: Keep focused context for Claude Code, Cursor, Copilot
- **Feature development**: Track files for current feature
- **Branch switching**: Automatic context switching between branches
- **Code review**: See which files changed in PR/branch

### Integrations

- Claude Code (native plugin)
- Cursor AI
- GitHub Copilot
- VS Code (official extension)
- Zed (official extension)

### Tech Stack

Built with Go + Bubbletea (TUI framework). MIT licensed.

### Installation

```bash
brew tap n-filatov/tap && brew install ws
```

### Links

- **GitHub**: https://github.com/n-filatov/ws
- **Documentation**: https://github.com/n-filatov/ws#readme
- **Author**: Nikita Filatov
- **License**: MIT

---

*Version 1.1.0 | Last updated: 2026-03-01*
