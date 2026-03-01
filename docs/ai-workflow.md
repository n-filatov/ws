# AI + ws Workflow Guide

A practical guide for using `ws` with AI assistants like Claude Code, Cursor, ChatGPT, and GitHub Copilot.

## Overview

`ws` keeps your working set of files organized and branch-scoped, making AI pair programming more effective. The AI always knows which files are relevant to your current work.

## Quick Start Workflow

### 1. Start a Feature

```bash
# Create or checkout your feature branch
git checkout -b feature/user-authentication

# Ask AI to identify relevant files
# AI discovers files and adds them to ws
ws add src/auth/login.go src/auth/jwt.go src/middleware/auth.go

# Open ws to navigate
ws
```

### 2. Work with AI

```bash
# AI can see your working set
ws list

# AI provides context-aware answers
AI: "I see you have 8 files in your working set. Let me analyze them..."
```

### 3. Switch Context

```bash
# Switch to another branch
git checkout feature/payment-processing

# ws automatically shows different files
ws
# Now shows payment-related files, not auth files
```

## Detailed Workflows

### Workflow 1: New Feature Discovery

**When you start a new feature:**

1. **Tell the AI what you're building**
   ```
   "I'm building user authentication with OAuth. Find relevant files."
   ```

2. **AI searches and adds to ws**
   ```bash
   AI: "Found these files:
       src/auth/login.go
       src/auth/oauth.go
       src/middleware/auth.go
       src/models/user.go

       Adding to ws..."
   ws add src/auth/login.go src/auth/oauth.go src/middleware/auth.go src/models/user.go
   ```

3. **Navigate with ws**
   ```
   Run `ws` to see your working set
   Press `e` to open files
   Press `/` to search
   ```

### Workflow 2: Growing Context

**As you discover more files:**

1. **AI finds additional files**
   ```
   AI: "I also found the config files for OAuth"
   ws add config/oauth.yaml config/providers.yaml
   ```

2. **Check your working set**
   ```bash
   ws
   # Now shows 6 files instead of 4
   ```

3. **Continue work**
   - Your working set grows organically
   - All files stay organized and accessible

### Workflow 3: Branch Switching

**When switching between features:**

1. **Switch branches**
   ```bash
   git checkout main
   ```

2. **ws shows different files**
   ```bash
   ws
   # Shows files for main branch (or empty set)
   ```

3. **Switch to feature branch**
   ```bash
   git checkout feature/user-authentication
   ```

4. **ws restores context**
   ```bash
   ws
   # Shows your 6 auth files
   ```

**Key benefit:** No manual cleanup. Context switches automatically.

### Workflow 4: Collaborative Problem Solving

**When debugging with AI:**

1. **Describe the problem**
   ```
   "Users can't log in after the latest commit"
   ```

2. **AI checks your working set**
   ```bash
   AI: "Let me see what you're working with:"
   ws list
   ```

3. **AI provides focused help**
   ```
   AI: "I see you have 6 auth-related files. Let me check
        login.go for issues. I found a bug on line 42."
   ```

4. **Fix and verify**
   ```bash
   # Press 'e' in ws to open login.go
   # Make the fix
   # Press 'r' in ws to refresh git status
   ```

## AI-Specific Workflows

### Claude Code

**Setup:**
```bash
/plugin marketplace add n-filatov/ws
/plugin install ws@n-filatov-ws
```

**Workflow:**
```
/ws:map user authentication flow
# Claude searches, finds files, runs ws add

ws
# Open in split terminal for navigation

# Ask Claude questions about your working set
"What do these files do?"
"How is login.go connected to jwt.go?"
```

### Cursor

**Workflow:**
```bash
# Tell Cursor what you're working on
"I'm working on payment processing"

# Cursor identifies files
# Add them to ws
ws add src/payment/processor.go src/payment/stripe.go

# Open ws in another terminal
ws

# Cursor now has focused context
# Ask questions about payment files
```

### ChatGPT/Claude (Web)

**Workflow:**
```bash
# ChatGPT/Claude can't run commands directly
# But you can copy-paste their suggestions

You: "Find OAuth-related files"
AI: "I found: src/auth/oauth.go, config/oauth.yaml"

You: ws add src/auth/oauth.go config/oauth.yaml

You: "Now what?"
AI: "Run 'ws' to navigate. I'll analyze these files..."
```

### GitHub Copilot

**Workflow:**
```bash
# Use ws list to give Copilot context
ws list | xargs -I {} cp {} /tmp/copilot-context/

# Or pipe to Copilot Chat
ws list
# Copilot sees your working set

# Ask Copilot questions
# "Explain how these files work together"
```

## Best Practices

### 1. Keep Working Sets Focused

**Good:**
```bash
# Only files for current feature
ws add src/auth/login.go src/auth/jwt.go
```

**Avoid:**
```bash
# Don't add entire codebase
ws add src/
```

### 2. Let Working Sets Grow Organically

```bash
# Start with core files
ws add main.go config.go

# Add more as needed
ws add handlers/auth.go middleware/auth.go

# Remove what you don't need
ws rm old-file.go
```

### 3. Use Git Status

```bash
ws
# Shows git status inline:
# M  Modified files
# A  Staged files
# ?  Untracked files
```

### 4. Refresh After Changes

```bash
# After you make changes
ws
# Press 'r' to refresh
# ws auto-adds any modified files
```

### 5. Leverage Search

```bash
ws
# Press '/' to fuzzy-search
# Type "auth" to filter
# Press Enter to open match
```

## Common Scenarios

### Scenario 1: Legacy Codebase

**Problem:** Huge codebase, don't know where to start

```bash
git checkout feature/refactor-auth
ws add src/auth/login.go  # Start with one file

AI: "I found related files"
ws add src/auth/user.go src/auth/session.go

ws
# Now navigate without overwhelm
```

### Scenario 2: Multiple Features

**Problem:** Juggling 3 features at once

```bash
# Branch 1: User auth
git checkout feature/user-auth
ws add src/auth/*.go
# 8 files in working set

# Branch 2: Payments
git checkout feature/payments
ws add src/payment/*.go
# 6 files in working set

# Branch 3: Admin
git checkout feature/admin
ws add src/admin/*.go
# 5 files in working set

# Switch between branches
# Each has its own working set
```

### Scenario 3: Code Review

**Problem:** Reviewing someone else's PR

```bash
git checkout pr/add-oauth
ws
# See what files were changed
ws add src/auth/oauth.go

# Ask AI questions
"What does oauth.go do?"
"How is it tested?"
```

### Scenario 4: Learning Codebase

**Problem:** New to project, need to understand architecture

```bash
git checkout main
ws add main.go

AI: "Let me trace the imports"
ws add handlers/ routes/ models/

ws
# Now you have a map of the architecture
```

## Pro Tips

### Tip 1: Use with FZF

```bash
# Fuzzy-find and open files
ws list | fzf | xargs nvim
```

### Tip 2: Git Worktree Integration

```bash
# Multiple branches, multiple directories
git worktree add ../ws-feature-auth feature/user-auth
cd ../ws-feature-auth
ws
# Each worktree has its own working set
```

### Tip 3: Backup Working Sets

```bash
# Working sets stored in:
~/.local/share/ws/<repo>/.workingset-<branch>

# Backup before major changes
cp -r ~/.local/share/ws ~/.local/share/ws.backup
```

### Tip 4: Clean Up Stale Sets

```bash
# ws offers to clean old sets
ws
# "Branch feature/old-auth was deleted 30 days ago. Clean working set? [y/N]"
```

## Troubleshooting

### "ws shows no files"

```bash
# Check if working set exists
ws list

# If empty, add files
ws add <files>
```

### "ws shows wrong branch"

```bash
# Check current branch
git branch

# ws follows git branch
# Switch to correct branch
git checkout <correct-branch>
ws
```

### "Files not showing up"

```bash
# Refresh git status
ws
# Press 'r'

# Or manually add
ws add <file>
```

## Advanced: Scripting with ws

### Count Lines in Working Set

```bash
ws list | xargs wc -l
```

### Grep Across Working Set

```bash
ws list | xargs grep "TODO"
```

### Open All in Editor

```bash
nvim $(ws list)
```

### Commit Working Set Files

```bash
git commit $(ws list) -m "Update auth flow"
```

## Summary

- **Start**: Add core files to ws
- **Grow**: Let AI discover and add more files
- **Navigate**: Use ws TUI to jump between files
- **Switch**: Change branches, ws changes context
- **Collaborate**: AI sees your working set, provides focused help

---

*Last updated: 2026-03-01*
