# Mastodon/Bluesky Announcement

**Post Type:** Long-form, single post (not a thread)

---

## Mastodon Post

I built a tool to solve a specific problem in AI-assisted development: **losing file context**.

## The Problem

When you're building a feature, you work with the same 6-10 files constantly. But your editor shows you hundreds. You ask your AI assistant (Claude, ChatGPT, Cursor) to identify relevant files — and it does! — but then what?

You open them one by one. You lose track. The next day, you repeat. When you switch branches, all context is lost.

## The Solution

**ws** (working set) — a branch-scoped file manager for the terminal.

```
ws add src/auth/login.go src/auth/jwt.go src/middleware/auth.go
ws
```

That's it. Now you have a focused, navigable list of exactly the files that matter for your current branch.

## Why ws Is Different

- **Branch-scoped**: Each git branch has its own working set. Switch branches, your context switches too.
- **Terminal-first**: Works with any editor, any AI assistant
- **AI-native**: Built specifically for AI pair programming workflows
- **Git-aware**: Shows inline status, auto-syncs modified files
- **Persistent**: Your working set survives terminal sessions

## The AI Workflow

1. You: "Find files for user authentication"
2. AI: [searches] "Found these files: login.go, jwt.go, auth.go"
3. AI: `ws add src/auth/login.go src/auth/jwt.go src/middleware/auth.go`
4. You: Run `ws` to see and navigate all files
5. AI: Now has focused context on exactly what you're working on

## Features

- 🌳 **Tree view TUI** — built with Bubbletea
- 🔍 **Fuzzy search** — instant filtering
- 📂 **Git status** — see modified/staged/untracked at a glance
- 🔄 **Auto-refresh** — modified files auto-added
- 🧹 **Stale cleanup** — prompts to clean old branch sets

## Integrations

- **Claude Code**: Native plugin (`/plugin install ws@n-filatov-ws`)
- **Cursor**: Works seamlessly with Cursor's AI
- **Copilot**: `ws list | xargs` for context
- **VS Code**: Official extension
- **Zed**: Official extension

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

## Built For

Developers who:
- Use AI assistants (Claude, ChatGPT, Cursor, Copilot)
- Work across multiple branches
- Need to maintain context across sessions
- Want terminal-first workflows

## Open Source

Built with Go + Bubbletea. MIT licensed.

**GitHub:** https://github.com/n-filatov/ws

If this sounds useful, give it a try! And if you like it, ⭐️ the repo — it helps with discoverability.

#golang #devtools #opensource #AI #coding #developerTools

---

## Bluesky Adaptation

**Same content**, but:
- Use Bluesky quotes instead of RTs
- Embed demo video if available
- Link to GitHub in first post
- Tag: @bsky.dev (if applicable)
- Shorter paragraphs (Bluesky readers prefer concise)

---

## Posting Tips

**Mastodon:**
- Post to @golang@mastodon.social, @devtools@fosstodon.org
- Use appropriate hashtags (max 5)
- Engage with replies for 24 hours
- Boost relevant responses

**Bluesky:**
- Post when US audience is active (9 AM - 12 PM EST)
- Use quote posts for follow-ups
- Embed media (GIFs, screenshots)

**Optimal times:**
- Mastodon: 9-11 AM PST (Tuesday-Thursday)
- Bluesky: 10 AM - 2 PM EST (weekdays)

**Follow-up posts (24 hours later):**
- "Quick demo of ws in action" [GIF]
- "How ws works with Cursor AI" [screenshots]
- "Branch-scoped working sets explained" [diagram]
