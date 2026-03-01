# Using ws with Cursor AI

A guide to integrating `ws` with Cursor AI for focused, context-aware development.

## Overview

Cursor AI is a powerful AI code editor. `ws` complements Cursor by maintaining a focused, branch-scoped list of files, giving Cursor better context and reducing the noise it needs to process.

## Quick Start

### 1. Install ws

```bash
# macOS/Linux
brew tap n-filatov/tap && brew install ws

# Go
go install github.com/n-filatov/ws@latest
```

### 2. Basic Workflow

```
You: Find all files related to user authentication

Cursor: [searches codebase]
I found these files:
- src/auth/login.go
- src/auth/jwt.go
- src/middleware/auth.go
- src/models/user.go

You: Add them to ws

You: ws add src/auth/login.go src/auth/jwt.go src/middleware/auth.go src/models/user.go
```

Now run `ws` in a separate terminal to see your focused file list.

## Key Workflows

### Workflow 1: Feature Discovery

**Step 1: Tell Cursor what you're building**

```
I'm building OAuth integration with Google
```

**Step 2: Cursor identifies files**

Cursor will search and identify:
- `src/auth/oauth.go`
- `src/auth/providers.go`
- `config/oauth.yaml`
- `internal/oauth/client.go`

**Step 3: Add to ws**

```bash
ws add src/auth/oauth.go src/auth/providers.go config/oauth.yaml internal/oauth/client.go
```

**Step 4: View in ws**

```bash
ws
```

Now you have a focused TUI with your OAuth files.

### Workflow 2: Using `ws list` Output

Cursor can work directly with `ws list` output:

```bash
# Count lines in your working set
ws list | xargs wc -l

# Grep for TODOs
ws list | xargs grep "TODO"

# Open all in your editor (if not using Cursor)
ws list | xargs cursor
```

### Workflow 3: Branch Switching

```bash
# Working on OAuth
git checkout feature/oauth
ws list
# Shows OAuth files

# Switch to payments
git checkout feature/payments
ws list
# Shows payment files (different set)
```

**Benefit:** Cursor's context changes with your branch. No confusion about which files are relevant.

### Workflow 4: Growing Context

As you work:

```
You: I need to handle OAuth errors

Cursor: You should add error handling in oauth.go

You: Add error handler file

You: ws add src/auth/oauth_error.go
```

Your working set grows organically. Cursor always knows what you're working on.

## Cursor + ws Synergies

### 1. Better Context

Cursor sees all your open files. `ws` keeps the important ones front and center.

**Before `ws`:**
- 50+ files in project
- Cursor searches everything
- Context gets diluted

**After `ws`:**
- 6-10 files in working set
- Cursor focuses on what matters
- Higher quality responses

### 2. Faster Navigation

**Without `ws`:**
- Search for file
- Wait for Cursor
- Open file
- Repeat

**With `ws`:**
- Run `ws`
- Press `/` to fuzzy-search
- Press `e` to open
- Instant

### 3. Branch Awareness

Cursor knows your git branch. `ws` knows your working set.

Together:
- Switch branches
- `ws` shows different files
- Cursor adapts to new context
- Seamless transition

## Advanced Patterns

### Pattern 1: Context Injection

```
You: Here's what I'm working on:

[Run: ws list]

Cursor: [receives 8 files]
Got it. You're working on OAuth. How can I help?
```

Cursor now has focused context.

### Pattern 2: Focused Analysis

```
You: Analyze the authentication flow

[Before: Cursor analyzes entire project]

[After: You run ws list first]

Cursor: Looking at your 6 auth files...
[Provides focused analysis]
```

### Pattern 3: Test-Driven Development

```bash
# Add test files to ws
ws add tests/auth_test.go tests/oauth_test.go

# Ask Cursor to test
ws list | xargs -I {} cursor --test {}
```

## Tips for Cursor Users

### Tip 1: Split Terminal

```
Left pane: Cursor (editing)
Right pane: ws (navigation)
```

Keep `ws` open while working. Press `e` to open files in Cursor.

### Tip 2: Use Fuzzy Search

```bash
ws
# Press /
# Type "auth"
# Press Enter to open match
```

Faster than Cursor's file finder for large projects.

### Tip 3: Refresh Often

```bash
ws
# Press 'r' to refresh git status
```

`ws` automatically adds modified files. Cursor sees the latest state.

### Tip 4: Integrate with Cursor Chat

```
You: Based on my working set, how does authentication work?

[Provide ws list output]
Cursor: Analyzing your 8 auth files...
```

Cursor provides focused, relevant answers.

## Common Scenarios

### Scenario 1: Legacy Codebase

**Problem:** Huge codebase, Cursor gets overwhelmed

```bash
# Start small
ws add main.go

# Ask Cursor to trace imports
You: What files does main.go import?

Cursor: [lists files]
You: [adds them to ws]
ws add handlers/ routes/ models/

# Now Cursor has focused context
```

### Scenario 2: Multiple Features

```bash
# Branch 1
git checkout feature/auth
ws add src/auth/*.go
# 8 files

# Branch 2
git checkout feature/payments
ws add src/payment/*.go
# 6 files

# Cursor adapts to each branch
```

### Scenario 3: Code Review

```bash
git checkout pr/add-oauth
ws
# See what's in the PR

You: Review these files for security issues

Cursor: [reviews your working set]
```

## Keyboard Shortcuts (ws TUI)

| Key | Action | Cursor Integration |
|-----|--------|-------------------|
| `e` | Open file in editor | Opens in Cursor if default editor |
| `/` | Fuzzy search | Find file, press `e` to open |
| `r` | Refresh git status | Cursor sees latest changes |
| `q` | Quit | Back to Cursor |

## Troubleshooting

### "Cursor doesn't open files from ws"

Set Cursor as your default editor:

```bash
# Add to ~/.wsconfig
editor=cursor
```

Or set environment variable:

```bash
export EDITOR=cursor
```

### "ws shows wrong files for this branch"

Check current branch:

```bash
git branch
ws
```

`ws` follows git branch. Switch to correct branch.

### "Cursor context is stale"

Refresh your working set:

```bash
ws
# Press 'r' to refresh
```

Modified files are auto-added.

## Best Practices

1. **Start with ws map**: Let Cursor discover files, add to ws
2. **Keep sets focused**: Only current feature files
3. **Use ws list**: Give Cursor focused context
4. **Switch branches freely**: ws handles context switching
5. **Review before committing**: Check ws for relevant files

## Example Session

```
You: I'm adding payment processing

Cursor: I found these files:
- src/payment/processor.go
- src/payment/stripe.go
- src/payment/webhook.go
- config/payment.yaml

You: [adds to ws]
ws add src/payment/processor.go src/payment/stripe.go src/payment/webhook.go config/payment.yaml

You: [opens ws in split pane]
ws

Cursor: What would you like to know?

You: How does the payment flow work?

Cursor: [analyzes your 4 payment files]
The payment flow works like this...
```

## Resources

- **ws GitHub**: https://github.com/n-filatov/ws
- **Cursor**: https://cursor.sh
- **Plugin**: `/plugin install ws@n-filatov-ws` (Claude Code)

---

**Last updated:** 2026-03-01
