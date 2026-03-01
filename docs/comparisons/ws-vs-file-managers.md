# ws vs File Manager Tools

**Comparing branch-scoped file lists to manual file marking tools**

---

## Overview

**ws** maintains automatic, branch-scoped working sets. File manager tools (harpoon, marks, shortcuts) require manual file marking.

The key difference: **automation vs manual**.

---

## Feature Comparison

| Feature | ws | harpoon | vim-marks |
|---------|----|----|-----------|
| **Scope** | Branch-scoped | Global (buffer) | Global (buffer) |
| **Setup** | Automatic | Manual | Manual |
| **Persistence** | File-based | Buffer-based | Buffer-based |
| **Git-aware** | Yes | No | No |
| **Branch switching** | Automatic | Manual | Manual |
| **AI integration** | Native | None | None |
| **Platform** | Terminal (any) | Neovim only | Vim/Neovim |

---

## When to Use ws

**Choose ws if you:**
- Want automatic file tracking
- Work across multiple branches
- Use AI assistants
- Switch contexts frequently
- Work outside of Vim/Neovim

**Example workflow:**
```bash
/ws:map user authentication
# AI finds 8 files, adds to ws automatically
ws
# See all files in TUI
```

---

## When to Use File Managers

**Choose harpoon if you:**
- Want manual file selection
- Stay in one branch
- Work exclusively in Neovim
- Prefer quick access to fixed files
- Don't need git awareness

**Example workflow:**
```vim
:lua require("harpoon.mark").add_file()
" Manually mark important files
:lua require("harpoon.ui").toggle_quick_menu()
" Quick access menu
```

---

## Key Differences

### 1. Setup Effort

**ws:**
```bash
# One command to add many files
ws add src/auth/*.go config/auth.yaml
```

**harpoon:**
```vim
" Mark files one by one
:lua require("harpoon.mark").add_file()  " file 1
:lua require("harpoon.mark").add_file()  " file 2
:lua require("harpoon.mark").add_file()  " file 3
```

**Trade-off:** ws is faster to set up, harpoon gives manual control.

### 2. Branch Switching

**ws:**
```bash
git checkout feature/payments
ws
# Shows payment files (different from auth)
```

**harpoon:**
```vim
" Same marks across all branches
" Manual cleanup needed
```

**Trade-off:** ws is branch-aware, harpoon is branch-agnostic.

### 3. AI Integration

**ws:**
- Claude Code: `/ws:map feature`
- Cursor: `ws add [files]`
- Copilot: `ws list | xargs`

**harpoon:**
- No native AI integration
- Would need custom scripting

**Trade-off:** ws is AI-native, harpoon is editor-focused.

---

## Workflow Comparison

### Scenario: Starting a Feature

**With ws:**
```bash
git checkout feature/oauth
/ws:map OAuth
# AI finds files, adds to ws
ws
# 6 files ready to navigate
```

**With harpoon:**
```vim
" Must manually mark each file
:lua require("harpoon.mark").add_file()  " oauth.go
:lua require("harpoon.mark").add_file()  " providers.go
" ... repeat for each file
```

**Difference:** ws automates discovery, harpoon requires manual marking.

### Scenario: Context Switching

**With ws:**
```bash
git checkout main
ws
# Shows main branch files (or empty)

git checkout feature/oauth
ws
# Shows OAuth files
```

**With harpoon:**
```vim
" Same files in all branches
" Must manually unmark/mark
```

**Difference:** ws switches automatically, harpoon doesn't.

---

## Use Cases

### ws is Better For:

1. **Feature branch workflow**
   - Multiple branches, different files
   - Automatic context switching

2. **AI-assisted development**
   - Claude Code, Cursor, Copilot
   - AI discovers files, adds to ws

3. **Large working sets**
   - 10+ files per feature
   - Tedious to mark manually

4. **Terminal workflows**
   - Work in terminal + editor
   - TUI navigation

### harpoon is Better For:

1. **Fixed file set**
   - Always open config.yaml
   - Always open main.go
   - Small, stable set

2. **Quick access**
   - 2-3 keypresses to jump
   - No TUI needed

3. **Vim/Neovim integration**
   - Live in Neovim
   - Want buffer-level integration

---

## Can They Work Together?

**Yes!** Use both:

```bash
# Terminal
ws add src/auth/*.go
# Branch-specific files

# Neovim
:lua require("harpoon.mark").add_file()
# Quick access to 2-3 critical files
```

**Workflow:**
- ws for branch-specific context
- harpoon for critical files (config, main)

---

## Migration Guide

### From harpoon to ws

1. **Export current marks:**
   ```vim
   :lua require("harpoon.ui").toggle_quick_menu()
   :lua print(vim.inspect(require("harpoon").get_mark_config().marks))
   ```

2. **Create ws working set:**
   ```bash
   ws add [paste files]
   ```

3. **Use per branch:**
   ```bash
   git checkout feature/new-feature
   ws
   # Different set per branch
   ```

### From ws to harpoon

1. **List current files:**
   ```bash
   ws list
   ```

2. **Mark in harpoon:**
   ```vim
   :lua require("harpoon.mark").add_file()  " file 1
   :lua require("harpoon.mark").add_file()  " file 2
   ```

3. **Access quickly:**
   ```vim
   <leader>a  " Jump to mark 1
   <leader>s  " Jump to mark 2
   ```

---

## Summary

**Choose ws for:**
- Automatic file tracking
- Branch-scoped workflows
- AI pair programming
- Working sets >5 files

**Choose harpoon for:**
- Manual file marking
- Quick access to 2-3 files
- Neovim integration
- Fixed file sets

**Or use both:**
- ws for branch context
- harpoon for critical files

---

**Related:**
- [ws vs Project Tools](ws-vs-project-tools.md)
- [ws vs lazygit](ws-vs-lazygit.md)
- [Homepage](https://github.com/n-filatov/ws)

---

*Last updated: 2026-03-01*
