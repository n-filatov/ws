# ws vs Project-Scoped Tools

**Comparing branch-scoped vs project-scoped file management**

---

## Overview

**ws** maintains branch-scoped working sets. Project tools (project.nvim, .project2, vim-projectionist) maintain project-scoped file lists.

The key difference: **scope**.

---

## Feature Comparison

| Feature | ws | project.nvim | .project2 |
|---------|----|--------------|-----------|
| **Scope** | Branch-scoped | Project-scoped | Project-scoped |
| **Persistence** | Automatic | Manual | Manual |
| **Git-aware** | Yes (inline status) | No | No |
| **Switching cost** | Zero (auto) | High (manual) | High (manual) |
| **Cleanup** | Automatic (stale) | Manual | Manual |
| **Platform** | Terminal (any editor) | Neovim only | Vim/Neovim |
| **AI-integration** | Native | None | None |

---

## When to Use ws

**Choose ws if you:**
- Work across multiple feature branches
- Use AI assistants (Claude, Cursor, Copilot)
- Switch branches frequently
- Want automatic context switching
- Work in terminal or multiple editors

**Example workflow:**
```bash
git checkout feature/oauth
ws add src/auth/oauth.go config/oauth.yaml
# 6 files in working set

git checkout feature/payments
ws
# Now shows 8 payment files, not oauth files
```

---

## When to Use Project Tools

**Choose project.nvim if you:**
- Stay on one branch for long periods
- Work exclusively in Neovim
- Want editor integration (folds, status bar)
- Prefer manual file management
- Don't need branch-specific context

**Example workflow:**
```vim
:ProjectRoot
:Project Files **/*.go **/*.md
" Files persist across branches
:ProjectCD
```

---

## Key Differences

### 1. Scope

**ws:**
```
feature/auth → {auth files}
feature/payments → {payment files}
```

**project.nvim:**
```
All branches → {same project files}
```

**Trade-off:** ws requires more setup (per branch) but provides better isolation.

### 2. Context Switching

**ws:**
```bash
git checkout feature/payments
ws
# Working set automatically changes
```

**project.nvim:**
```vim
:ProjectFiles **/*.go
# Must manually change files for each branch
```

**Trade-off:** ws is automatic, project.nvim gives manual control.

### 3. Git Awareness

**ws:**
```
src/auth/login.go         M  # Modified
src/auth/jwt.go           A  # Staged
src/middleware/cors.go    ?  # Untracked
```

**project.nvim:**
```
src/auth/login.go
src/auth/jwt.go
src/middleware/cors.go
# No git status
```

**Trade-off:** ws provides more context, project.nvim is simpler.

### 4. AI Integration

**ws:**
- Claude Code plugin (`/ws:map`)
- Cursor workflow (`ws list`)
- Copilot integration (`ws list | xargs`)

**project.nvim:**
- No native AI integration
- Would need custom scripting

**Trade-off:** ws is AI-native, project.nvim is editor-focused.

---

## Workflow Comparison

### Scenario: Starting a New Feature

**With ws:**
```bash
git checkout feature/user-auth
# Working set is empty or shows previous branch's files
/ws:map user authentication
# AI finds files, adds to ws
ws
# See focused list
```

**With project.nvim:**
```vim
:ProjectRoot
:ProjectFiles **/*auth*.go **/*user*.go
# Manually specify patterns
" Files persist across all branches
```

**Difference:** ws is branch-aware, project.nvim is branch-agnostic.

### Scenario: Switching Branches

**With ws:**
```bash
git checkout feature/payments
ws
# Shows payment files (different from auth)
```

**With project.nvim:**
```vim
" Same files across branches
" Must manually remember which files matter
```

**Difference:** ws switches context automatically, project.nvim doesn't.

---

## Use Cases

### ws is Better For:

1. **Feature branch workflow**
   - Each feature = separate branch
   - Context switches frequently

2. **AI pair programming**
   - Claude Code, Cursor, Copilot
   - Need focused context

3. **Multi-editor setups**
   - Terminal + editor
   - Vim + VS Code
   - Any combination

4. **Git-heavy workflows**
   - Long-lived branches
   - Multiple developers
   - Code review branches

### project.nvim is Better For:

1. **Single-branch development**
   - Trunk-based development
   - Rarely switch branches

2. **Neovim power users**
   - Live in Neovim
   - Want editor integration
   - Use folds, status line, etc.

3. **Manual control**
   - Want to curate file list
   - Don't want automation
   - Prefer explicit management

---

## Can They Work Together?

**Yes!** Use both:

```bash
# Terminal
ws add src/auth/*.go
# Manage branch-specific files

# Neovim
:ProjectFiles **/*.go
# Manage editor-specific view
```

**Workflow:**
- ws for branch-scoped context
- project.nvim for editor features (folds, etc.)

---

## Migration Guide

### From project.nvim to ws

1. **Export current project:**
   ```vim
   :ProjectFiles **/*.go
   ```

2. **Create ws working set:**
   ```bash
   ws add [paste files]
   ```

3. **Switch branches:**
   ```bash
   git checkout feature/new-feature
   ws
   # Start fresh per branch
   ```

### From ws to project.nvim

1. **List current working set:**
   ```bash
   ws list
   ```

2. **Create project:**
   ```vim
   :ProjectFiles [paste files]
   ```

3. **Use across branches:**
   ```vim
   " Same files for all branches
   ```

---

## Summary

**Choose ws for:**
- Branch-scoped workflows
- AI pair programming
- Automatic context switching
- Terminal-first workflows

**Choose project.nvim for:**
- Project-scoped workflows
- Neovim integration
- Manual file curation
- Editor-first workflows

**Or use both:**
- ws for branch context
- project.nvim for editor features

---

**Related:**
- [ws vs File Managers](ws-vs-file-managers.md)
- [ws vs lazygit](ws-vs-lazygit.md)
- [Homepage](https://github.com/n-filatov/ws)

---

*Last updated: 2026-03-01*
