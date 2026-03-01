# Twitter Thread: Branch-Scoped Working Sets

**Based on blog post:** "Branch-Scoped Working Sets: A Missing Git Primitive"

---

## Thread Structure

**Hook:** Target systems programmers and tool builders

---

## Tweet 1: The Missing Primitive

Git gives us branches. Editors give us file lists.

But nothing gives us branch-scoped file lists.

That's what ws does.

A thread on design, implementation, and philosophy. 🧵

---

## Tweet 2: The Problem

Modern dev involves:
• Multiple branches (auth, payments, oauth)
• Many files (hundreds or thousands)
• AI assistants that need context

What's missing: Which files matter for THIS branch?

---

## Tweet 3: Current Solutions Fall Short

Project-scoped tools (project.nvim): One set for entire project
Recent files (editor MRU): What you opened, not what's relevant
Manual lists (harpoon): Requires management, no git awareness

The gap: Branch-scoped, persistent, git-aware file lists

---

## Tweet 4: Design Principles

1. Branch-scoped: Each branch has its own set
2. Persistent: Survives terminal sessions
3. Git-aware: Shows inline status
4. Terminal-first: Works with any editor

Simple principles, powerful results.

---

## Tweet 5: Data Storage

Where to store working sets?

~/.local/share/ws/<repo>/.workingset-<branch>

Why outside repo?
• User-specific, not project-specific
• Don't pollute git history
• No .gitignore needed

---

## Tweet 6: File Format

Working sets are newline-separated paths:

```
src/auth/login.go
src/auth/jwt.go
src/middleware/auth.go
```

Why plain text?
• Human-readable
• Easy to debug
• No parsing overhead
• Tool-friendly

---

## Tweet 7: Git Integration

Detect current branch:
```go
git rev-parse --abbrev-ref -- HEAD
```

Get file status:
```go
git status --porcelain <file>
```

Parse output: "M src/auth/login.go"

Simple, reliable, git-native.

---

## Tweet 8: The TUI (Bubbletea)

Built with Bubbletea (Go TUI framework):

```go
type model struct {
    files    []File
    cursor   int
    filtering bool
    query    string
}
```

Vim-style keybindings. Tree view. Fuzzy search.

---

## Tweet 9: Performance

Large working sets? Lazy loading.

Large repos? Cache git status.

Startup target: <100ms to open TUI.

Optimizations: Read once, lazy git status, no file contents.

---

## Tweet 10: Edge Cases

Deleted files → Mark [DELETED]
Renamed files → Detect via git diff --name-status
Branch deletion → Cleanup after N days
Submodules → Track separately

Handle reality, not just happy path.

---

## Tweet 11: AI Integration

Claude Code plugin:
```json
{
  "name": "ws:map",
  "description": "Search and add files"
}
```

Cursor workflow:
```bash
ws add src/auth/*.go
ws list  # Cursor sees focused context
```

---

## Tweet 12: Lessons Learned

1. Simple storage wins (plain text > database)
2. Git native > git wrapper (embrace git)
3. Terminal first > editor first (works everywhere)
4. Opinionated > flexible (reduces cognitive load)

---

## Tweet 13: Trade-offs

vs project.nvim: Terminal vs editor, branch vs project scope
vs harpoon: Automatic vs manual, branch vs global
vs zoxide: Files vs directories, different scopes

Different tools, different trade-offs.

---

## Tweet 14: Future Directions

• Remote working sets (sync across machines)
• Collaborative working sets (share with team)
• Smart file discovery (AST analysis)

Building in public, feedback welcome.

---

## Tweet 15: Philosophy

Sometimes the missing primitive isn't complicated.

It just needs someone to build it.

Branch-scoped file lists. Simple. Powerful. Missing no more.

---

## Tweet 16: Try It

Built with Go + Bubbletea. MIT licensed.

GitHub: github.com/n-filatov/ws

If you build tools, I'd love your feedback.

---

## Posting Strategy

**Target audience:** Systems programmers, tool builders, Go developers

**Best time:** 9-11 AM PST, Tuesday-Thursday

**Platform:** Twitter/X, also consider Mastodon for longer-form

**Follow-up:**
- Share code snippets from blog post
- Reply with Bubbletea examples
- Link to full technical deep-dive

**Tags:**
- @charmbracelet (Bubbletea maintainers)
- @golang
- #golang #systemsprogramming #devtools #cli #terminal

**Engagement questions:**
- "What's your missing primitive?"
- "What tooling gaps have you encountered?"
