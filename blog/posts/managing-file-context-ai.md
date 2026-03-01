# Managing File Context in AI Pair Programming

**How to stop losing track of files and keep your AI assistant focused**

---

## The Context Problem

You're building a feature. You ask your AI assistant:

> "Which files are relevant to the user authentication flow?"

And it delivers:

```
src/auth/login.go
src/auth/jwt.go
src/middleware/auth.go
src/models/user.go
src/handlers/auth.go
internal/crypto/password.go
config/auth.yaml
```

Eight files. Perfect. You open them... one by one.

Three hours later, you've lost track. The next day, you're asking the same question.

**The problem isn't finding files. It's remembering them.**

## Why This Happens

AI assistants are great at *discovery* but terrible at *persistence*. They can identify relevant files in your codebase, but they can't remember what you were working on yesterday.

You end up in a loop:

1. Ask AI for files
2. Open them manually
3. Work for a few hours
4. Close everything
5. Repeat tomorrow

This is exhausting. And it's not how AI pair programming should work.

## The Missing Piece: Working Sets

What you need is a **working set** — a focused, persistent list of files that are relevant to what you're doing right now.

Not your entire codebase. Not your "most recently used" files. Just the files that matter for this specific feature, on this specific branch.

**Enter ws.**

## ws: Branch-Scoped Working Sets

ws is a terminal tool that keeps track of your working set. It's:

- **Branch-scoped**: Each git branch has its own working set
- **Persistent**: Survives terminal sessions
- **Git-aware**: Shows inline status, auto-syncs modified files
- **Terminal-first**: Works with any editor or AI assistant

### The Workflow

```bash
# Add relevant files to your working set
ws add src/auth/login.go src/auth/jwt.go src/middleware/auth.go

# Open ws to see and navigate your working set
ws
```

That's it. Now you have a focused list of exactly the files you need.

Press `e` to open a file. Press `/` to fuzzy-search. Press `r` to refresh git status.

### The AI Workflow

Here's where it gets powerful.

**You:**
> "Find all files related to OAuth authentication"

**AI:**
> "I found these files:
> - src/auth/oauth.go
> - src/auth/providers.go
> - config/oauth.yaml
>
> Adding them to your working set..."

```bash
ws add src/auth/oauth.go src/auth/providers.go config/oauth.yaml
```

**You:** Run `ws` to see all your files in one place.

**AI:** Now has a focused context. It knows exactly which files you're working with.

## Branch-Scoped Context

This is the game-changer.

You're working on `feature/user-authentication`. Your working set has 8 auth-related files.

Then you switch to `feature/payment-processing`:

```bash
git checkout feature/payment-processing
ws
```

Now your working set shows 6 payment-related files.

**No manual cleanup. No cross-branch contamination. Context switches when you do.**

## Real-World Example

Let's say you're building OAuth login.

### Without ws

1. Ask AI: "Find OAuth files"
2. AI lists 6 files
3. You open them in tabs
4. Tabs get lost among 50 other tabs
5. You close the editor
6. Next day: repeat from step 1

### With ws

1. Ask AI: "Find OAuth files"
2. AI lists 6 files, adds them to ws
3. You run `ws` — see all 6 files in a focused TUI
4. Navigate with `e` (open), `/` (search)
5. Close the editor
6. Next day: run `ws` — all 6 files are there
7. Switch branches → ws shows different files

The difference is persistence and focus.

## How ws Helps AI Assistants

When your AI assistant knows your working set, it can:

### Provide Focused Answers

> "How does authentication work in this codebase?"

**AI knows your working set** and only analyzes those files:

> "Based on the 8 files in your working set, authentication works like this..."

### Suggest Relevant Files

> "I'm getting a JWT error"

**AI checks your working set:**

> "I see `jwt.go` and `login.go` in your working set. Let me check those first..."

### Maintain Context

> "Continue working on OAuth"

**AI:**

> "You have 6 OAuth-related files in your working set. Should I focus on those?"

## Integrations

ws works with all major AI assistants:

### Claude Code

```bash
/plugin marketplace add n-filatov/ws
/plugin install ws@n-filatov-ws
```

```
/ws:map OAuth authentication
```

Claude finds files and adds them to ws automatically.

### Cursor

Tell Cursor what you're working on, then add files to ws:

```bash
ws add src/auth/oauth.go config/oauth.yaml
```

Cursor now has focused context.

### ChatGPT/Claude (Web)

Copy-paste the file list from ChatGPT:

```bash
ws add <paste list>
```

Run `ws` to navigate.

### GitHub Copilot

```bash
ws list | xargs -I {} copilot analyze {}
```

## Features

- **Tree view TUI**: Built with Bubbletea, vim-style keybindings
- **Fuzzy search**: Press `/` to instantly filter
- **Git status**: See `M` (modified), `A` (staged), `?` (untracked) inline
- **Auto-refresh**: Modified files are auto-added
- **Stale cleanup**: Prompts to clean old branch sets

## Installation

**macOS/Linux:**
```bash
brew tap n-filatov/tap && brew install ws
```

**Debian/Ubuntu:**
```bash
curl -fsSL https://n-filatov.github.io/ws/gpg.key | sudo gpg --dearmor -o /usr/share/keyrings/ws.gpg
echo "deb [signed-by=/usr/share/keyrings/ws.gpg] https://n-filatov.github.io/ws ./" | sudo tee /etc/apt/sources.list.d/ws.list
sudo apt update && sudo apt install ws
```

**Go:**
```bash
go install github.com/n-filatov/ws@latest
```

**From source:**
```bash
git clone https://github.com/n-filatov/ws && cd ws && make install
```

## The Bottom Line

AI assistants are powerful. But they need context.

ws gives them that context by keeping a focused, branch-scoped list of files.

**Less time navigating. More time building.**

**Try it:** [github.com/n-filatov/ws](https://github.com/n-filatov/ws)

---

**Published:** 2026-03-01
**Tags:** #AI #pair-programming #developer-tools #golang #claude-code #cursor #copilot
**Author:** Nikita Filatov
**Length:** ~900 words
