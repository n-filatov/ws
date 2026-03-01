# Twitter/X Announcement Thread

**Thread Hook:** Focus on the problem of losing file context in AI pair programming

---

## Tweet 1: The Hook

I built a tool because I kept losing track of files when working with AI.

AI: "Here are 12 files for the auth feature"
Me: *opens them one by one, loses track, repeats the next day*

So I built ws. It keeps your working set branch-scoped and persistent. 🧵

---

## Tweet 2: The Problem

When building features, you jump between the same files constantly.

But your editor shows you hundreds. You ask Claude/Copilot "which files matter?" and get a list — but then what?

You open them manually. You lose track. Next day: repeat.

The problem isn't finding files. It's remembering them.

---

## Tweet 3: The Solution

Enter ws (working set).

Just add files to your current branch:
```
ws add src/auth/login.go src/auth/jwt.go
```

Run `ws` to see them in a TUI.

Each branch gets its own set. Switch branches → your working set switches too.

No more "what was I working on again?"

---

## Tweet 4: Why It's Different

Other tools: project-scoped or editor-specific

ws: branch-scoped + terminal-first

This means:
- ✅ Context switches automatically with git branches
- ✅ Works with Claude Code, Cursor, Copilot, anything
- ✅ Persistent across terminal sessions
- ✅ Git-aware (shows inline status)

---

## Tweet 5: AI-Native Design

Built specifically for AI pair programming.

When your AI assistant finds relevant files, it can add them to ws.

Now you have a focused list. You can navigate, review, and iterate — without losing context.

The AI sees the same files you do.

---

## Tweet 6: Quick Demo

```
$ ws add src/auth/* src/middleware/auth.go
$ ws
[opens TUI with tree view, git status, fuzzy search]
Press 'e' → opens file in editor
Press '/' → fuzzy search
Press 'r' → refresh git status
```

Less time navigating. More time building.

---

## Tweet 7: Installation

One command (macOS/Linux):
```bash
brew install n-filatov/tap/ws
```

Or Go:
```bash
go install github.com/n-filatov/ws@latest
```

---

## Tweet 8: Integrations

- Claude Code plugin (native integration)
- VS Code extension
- Zed extension
- Works with any AI assistant (ChatGPT, Cursor, Copilot)

---

## Tweet 9: Open Source

Built with Go + Bubbletea.

Branch-scoped. Git-aware. AI-native.

GitHub: github.com/n-filatov/ws

---

## Tweet 10: CTA

If you use AI to code, you need ws.

Keep your context tight. Stay focused on what matters.

Try it: github.com/n-filatov/ws

RTs appreciated 🙏

---

## Posting Strategy

**Best times to post:**
- Tuesday-Thursday, 9-11 AM PST
- Avoid weekends (lower tech engagement)

**Follow-up:**
- Reply to your own thread with demo GIF
- Pin the thread for 24 hours
- Engage with all replies for first 2 hours

**Tags to mention:**
- @claude (if relevant)
- @cursor_ai (if relevant)
- #golang #devtools #AI #coding

**Engagement questions:**
- "How do you keep track of files across branches?"
- "Ever lose context when switching features?"
