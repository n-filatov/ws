# Twitter Thread: Managing File Context in AI Pair Programming

**Based on blog post:** "Managing File Context in AI Pair Programming"

---

## Thread Structure

**Hook:** Focus on the pain of losing file context with AI assistants

---

## Tweet 1: The Problem

I built a tool because I kept losing track of files when working with AI.

AI: "Here are 12 files for the auth feature"
Me: *opens them one by one, loses track, repeats the next day*

Sound familiar? 🧵

---

## Tweet 2: The Loop

When building features, you jump between the same files constantly.

But your editor shows hundreds.

You ask Claude/Copilot "which files matter?" and get a list.

Then what? You open them manually. You lose track. Next day: repeat.

---

## Tweet 3: The Missing Piece

The problem isn't finding files. AI is great at that.

The problem is **remembering** them.

What you need is a working set — a focused, persistent list of files that matter for what you're doing right now.

---

## Tweet 4: Enter ws

```bash
ws add src/auth/login.go src/auth/jwt.go
ws
```

That's it.

Now you have a focused list of files. Branch-scoped. Persistent.

Press 'e' to open. Press '/' to search. Press 'q' to quit.

---

## Tweet 5: Why Branch-Scoped?

Each git branch gets its own working set.

Switch branches → your working set switches too.

feature/auth → {8 auth files}
feature/payments → {6 payment files}

No manual cleanup. Context switches when you do.

---

## Tweet 6: The AI Workflow

1. AI: "Find OAuth files"
2. AI: "Found 6 files. Adding to ws..."
3. `ws add src/auth/oauth.go config/oauth.yaml`
4. You: Run `ws` to navigate
5. AI: Now has focused context

---

## Tweet 7: How This Helps AI

When your AI knows your working set:

❌ "Analyze this codebase" (too broad)
✅ "Analyze these 8 files" (focused)

Better questions → Better answers

---

## Tweet 8: Real Example

Yesterday: "Find auth files"
AI lists 8 files
I add to ws

Today: Run `ws`
All 8 files are there
AI remembers context

No more repeating myself.

---

## Tweet 9: Features

🌳 Tree view TUI (built with Bubbletea)
🔍 Fuzzy search
📂 Git status inline
🔄 Auto-refresh
🧹 Stale cleanup

All terminal-based. Works with any editor.

---

## Tweet 10: Integrations

Works with all major AI assistants:

• Claude Code (native plugin)
• Cursor AI
• GitHub Copilot
• ChatGPT (manual)

Same tool, different assistants.

---

## Tweet 11: Installation

One command (macOS/Linux):

```bash
brew tap n-filatov/tap && brew install ws
```

Or Go:

```bash
go install github.com/n-filatov/ws@latest
```

---

## Tweet 12: Open Source

Built with Go + Bubbletea. MIT licensed.

Branch-scoped. Git-aware. AI-native.

GitHub: github.com/n-filatov/ws

---

## Tweet 13: TL;DR

AI assistants are powerful. But they need context.

ws gives them that context by keeping a focused list of files.

Less time navigating. More time building.

---

## Tweet 14: Try It

If you use AI to code, you need ws.

Keep your context tight. Stay focused.

github.com/n-filatov/ws

RTs appreciated 🙏

---

## Posting Strategy

**Best time:** 9-11 AM PST, Tuesday-Thursday

**Follow-up:**
- Reply with demo GIF (if available)
- Pin for 24 hours
- Engage with replies

**Tags:**
- @claude (if relevant)
- @cursor_ai (if relevant)
- #golang #devtools #AI #coding

**Engagement questions to add:**
- "How do you keep track of files?"
- "Ever lose context with AI assistants?"
