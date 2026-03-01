# ws vs lazygit

**Comparing file navigation to git operations**

---

## Overview

**ws** manages file lists for development. **lazygit** manages git operations.

These are **complementary tools**, not competitors.

---

## Feature Comparison

| Feature | ws | lazygit |
|---------|----|----|
| **Purpose** | File navigation | Git operations |
| **Scope** | Working set | Git repository |
| **Operations** | Add, remove, list files | Commit, branch, push, pull |
| **Git awareness** | Shows status | Performs git actions |
| **Workflow** | Code navigation | Git workflow |
| **TUI focus** | File tree | Git status, commits, branches |

---

## What ws Does

**File management:**
```bash
ws add src/auth/*.go      # Add to working set
ws rm src/auth/old.go     # Remove from working set
ws                         # View files
ws list                    # List files
```

**Purpose:** Navigate files relevant to current feature.

---

## What lazygit Does

**Git management:**
```bash
lazygit
# Commit: Press 'c'
# Branch: Press 'b'
# Push: Press 'P'
# Stash: Press 's'
```

**Purpose:** Perform git operations efficiently.

---

## How They Work Together

### Typical Workflow

```bash
# 1. Use ws to navigate files
ws add src/auth/*.go
ws
# Press 'e' to open files

# 2. Make changes
# Edit files...

# 3. Use lazygit to commit
lazygit
# Stage files (Space)
# Commit (c)
# Push (P)

# 4. Back to ws
ws
# Press 'r' to refresh git status
```

### Key Integration

**ws shows git status:**
```
src/auth/login.go         M  # Modified
src/auth/jwt.go           A  # Staged
src/middleware/cors.go    ?  # Untracked
```

**lazygit performs git actions:**
- Stage files
- Commit changes
- Create branches
- Merge branches

---

## When to Use Each

### Use ws When You're:

**Navigating code:**
- "What files am I working on for this feature?"
- "Where was that authentication file?"
- "Let me see all OAuth-related files"

**Managing context:**
- Switching branches
- Starting new features
- Working with AI assistants

### Use lazygit When You're:

**Managing git:**
- "I need to commit these changes"
- "Create a new branch"
- "View commit history"
- "Resolve merge conflicts"

**Performing git operations:**
- Staging files
- Pushing/Pulling
- Rebasing/Cherry-picking
- Managing stashes

---

## Feature-by-Feature

### File Operations

| Task | ws | lazygit |
|------|----|----|
| List files | ✅ `ws list` | ❌ Not designed for this |
| Add to context | ✅ `ws add` | ❌ Stage for commit |
| Remove from context | ✅ `ws rm` | ❌ Not designed for this |
| Open files | ✅ Press `e` | ❌ Not designed for this |
| Search files | ✅ Press `/` | ❌ Not designed for this |

### Git Operations

| Task | ws | lazygit |
|------|----|----|
| Stage files | ❌ | ✅ Press `Space` |
| Commit | ❌ | ✅ Press `c` |
| Create branch | ❌ | ✅ Press `b` |
| View history | ❌ | ✅ Commits panel |
| Push/Pull | ❌ | ✅ Press `P`/`p` |
| Stash | ❌ | ✅ Press `s` |

### Git Status

| Task | ws | lazygit |
|------|----|----|
| Show status | ✅ Inline (M/A/?) | ✅ Full panel |
| Refresh status | ✅ Press `r` | ✅ Automatic |
| Auto-sync | ✅ Auto-adds modified | ❌ Manual staging |

---

## Workflow Examples

### Example 1: Start a Feature

```bash
# 1. Create branch
lazygit
# Press 'b' → 'n' → "feature/oauth"

# 2. Set up working set
ws add src/auth/*.go
ws
# 8 files ready to navigate

# 3. Work on feature
# Edit files...

# 4. Commit changes
lazygit
# Stage files (Space)
# Commit (c)
# Push (P)

# 5. Switch branches
git checkout main
ws
# Different files for main branch
```

### Example 2: Code Review

```bash
# 1. Checkout PR branch
git checkout pr/add-oauth

# 2. See what files changed
ws
# Shows files in PR (if added to ws)

# 3. Review with lazygit
lazygit
# See commits, changes, file diffs

# 4. Add to ws if needed
ws add src/auth/oauth.go
```

### Example 3: Bug Fix

```bash
# 1. ws helps find relevant files
ws add src/auth/*.go
ws
# See all auth files

# 2. Fix bug
# Edit files...

# 3. lazygit helps commit
lazygit
# Stage, commit, push
```

---

## Keyboard Shortcuts

### ws Keybindings

| Key | Action |
|-----|--------|
| `j`/`k` | Move down/up |
| `e` | Open file in editor |
| `/` | Fuzzy search |
| `r` | Refresh git status |
| `a` | Add file |
| `d` | Remove file |
| `q` | Quit |

### lazygit Keybindings

| Key | Action |
|-----|--------|
| `j`/`k` | Move down/up |
| `Space` | Stage/unstage file |
| `c` | Commit |
| `b` | Branch menu |
| `P` | Push |
| `p` | Pull |
| `s` | Stash |
| `q` | Quit |

---

## Complementary Use

### Recommended Setup

**Split terminal:**
```
Left pane: ws        (file navigation)
Right pane: lazygit  (git operations)
```

**Or use in sequence:**
```bash
# Navigate with ws
ws
# Press 'e' to edit files

# Commit with lazygit
lazygit
# Stage and commit

# Back to ws
ws
# Press 'r' to refresh status
```

### Key Bindings (Optional)

Create a custom command that opens both:

```bash
# ~/.bashrc or ~/.zshrc
ws-git() {
  ws
  lazygit
}
```

---

## Summary

**Use ws for:**
- File navigation
- Context management
- Branch-specific file lists
- AI pair programming

**Use lazygit for:**
- Git operations
- Commit management
- Branch management
- Git workflow

**Together:**
- ws helps you navigate files
- lazygit helps you commit changes
- Both improve developer workflow

---

## Recommendation

**Install both:**

```bash
# ws
brew tap n-filatov/tap && brew install ws

# lazygit
brew install jesseduffield/lazygit/lazygit
```

**Use both:**
- ws when coding
- lazygit when committing

They work better together.

---

**Related:**
- [ws vs Project Tools](ws-vs-project-tools.md)
- [ws vs File Managers](ws-vs-file-managers.md)
- [ws Homepage](https://github.com/n-filatov/ws)
- [lazygit Homepage](https://github.com/jesseduffield/lazygit)

---

*Last updated: 2026-03-01*
