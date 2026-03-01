# Newsletter Announcements

**Pre-written announcements for developer newsletters**

---

## System Override

**Subject:** CLI tool that solves AI file context problems

**Body:**

Hi team,

I built a tool to solve a problem I kept hitting when pair-programming with AI.

**The problem:** You ask Claude/Cursor to identify relevant files, and they do — but then you open them one by one, lose track, and repeat the next day.

**The solution:** ws — a branch-scoped working set manager.

```bash
ws add src/auth/login.go src/auth/jwt.go
ws  # Opens TUI with tree view, fuzzy search
```

Each branch gets its own working set. Switch branches → your file list switches too.

**Why it matters for AI:**
- AI assistants (Claude Code, Cursor, Copilot) get focused context
- No more "which files was I working on?"
- Branch-scoped = automatic context switching

**Built with:** Go + Bubbletea
**License:** MIT
**GitHub:** https://github.com/n-filatov/ws

Would love feedback from the devtools community!

Best,
Nikita

---

## Golang Weekly

**Subject:** ws - Branch-scoped working sets for Go projects

**Body:**

Hi Gophers,

I released ws (working set) — a CLI tool for managing files across git branches.

**What it does:**
- Maintains a focused list of files per branch
- Terminal UI with tree view and fuzzy search
- Git-aware (shows inline status)
- Integrates with AI assistants (Claude Code, Cursor)

**Example:**
```bash
ws add cmd/auth/main.go internal/auth/*.go
ws  # Opens TUI
```

**Why Go:**
Fast startup, single binary, no dependencies. Perfect for a terminal tool.

**GitHub:** https://github.com/n-filatov/ws

Would love feedback from the Go community!

Cheers,
Nikita

---

## Hacker Newsletter

**Subject:** Show HN: ws - Branch-scoped file manager I built

**Body:**

I built ws to solve a specific problem: losing file context when coding with AI.

**Problem:** AI assistants find relevant files, but there's no persistent way to remember them across sessions or branches.

**Solution:** ws — branch-scoped working sets.

```bash
ws add src/auth/*.go
ws  # See files in TUI
```

Each branch gets its own set. Switch branches → context switches.

**HN discussion:** https://news.ycombinator.com/item?id=[number]
**GitHub:** https://github.com/n-filatov/ws

Curious what the HN community thinks!

Best,
Nikita

---

## Pragmatic Engineer

**Subject:** Tool for managing file context in AI-assisted development

**Body:**

Hi Gergely,

I built a tool that might interest your readers working with AI assistants.

**The problem:**
When pair-programming with Claude/Cursor, AI identifies relevant files — but there's no persistent way to remember them. You open files manually, lose track, repeat.

**The solution:**
ws — a branch-scoped working set manager.

```bash
ws add src/auth/*.go
ws  # Terminal UI
```

**Why engineers should care:**
- Reduces context switching overhead
- Integrates with Claude Code, Cursor, Copilot
- Automatic per-branch context
- Open source (MIT)

**GitHub:** https://github.com/n-filatov/ws
**Show HN:** [link if applicable]

Thought your readers might find it useful for AI-assisted workflows.

Best,
Nikita

---

## Devtools Weekly

**Subject:** ws: AI-native working set manager

**Body:**

Hi [Name],

I released ws — an AI-native tool for managing working sets across git branches.

**What it does:**
- Maintains branch-scoped file lists
- Terminal UI (built with Bubbletea)
- Git-aware with inline status
- Integrates with Claude Code, Cursor, Copilot

**Why AI-native:**
Built specifically for AI pair programming. AI assistants can discover files and add them to ws, maintaining focused context.

**Install:**
```bash
brew tap n-filatov/tap && brew install ws
```

**GitHub:** https://github.com/n-filatov/ws

Would love to be featured in an upcoming issue!

Best,
Nikita

---

## CLI Quarterly

**Subject:** ws - Branch-scoped file manager for the terminal

**Body:**

Hi CLI Quarterly team,

I built ws — a branch-scoped working set manager for terminal users.

**Key features:**
- Branch-scoped file lists (each branch has its own set)
- Terminal UI with vim-style keybindings
- Fuzzy search, tree view, git status
- Integrates with AI assistants

**Use case:**
When working on feature/auth, you add auth files. When you switch to feature/payments, ws shows payment files instead.

**Tech stack:**
- Go 1.21+
- Bubbletea (TUI framework)
- MIT licensed

**GitHub:** https://github.com/n-filatov/ws

Hoping CLI Quarterly readers might find it useful!

Cheers,
Nikita

---

## Terminal Trove

**Subject:** ws - AI-native file manager for terminal workflows

**Body:**

Hi Terminal Trove,

I'd like to submit ws for consideration in your next issue.

**Description:**
ws is a branch-scoped working set manager. It keeps track of files relevant to your current git branch, with a terminal UI for navigation.

**Features:**
- Branch-scoped file lists
- TUI with tree view and fuzzy search
- Git-aware (shows inline status)
- Integrates with AI assistants

**Install:**
```bash
brew install n-filatov/tap/ws
# or
go install github.com/n-filatov/ws@latest
```

**GitHub:** https://github.com/n-filatov/ws
**Demo:** [link to GIF/video if available]

Hope terminal users find it useful!

Best,
Nikita

---

## Rusty Shell (if they cover non-Rust tools)

**Subject:** ws - Go-based working set manager

**Body:**

Hi Rusty Shell,

I built ws (in Go) and thought your readers might find it useful even though it's not Rust.

**What it does:**
- Branch-scoped working sets for git projects
- Terminal UI with vim-style navigation
- Git-aware, AI-native

**Why Go:**
Wanted fast startup and single binary deployment. Rust would work too, but Go was the right choice for this project.

**Install:**
```bash
brew install n-filatov/tap/ws
```

**GitHub:** https://github.com/n-filatov/ws

Even though it's Go, hope it's interesting to the shell/terminal community!

Best,
Nikita

---

## FLOSS Weekly

**Subject:** ws - Open source working set manager

**Body:**

Hi FLOSS Weekly,

I'd like to suggest ws for your show — an open source tool I built for managing file context in AI-assisted development.

**What it is:**
- Branch-scoped working set manager
- Terminal UI built with Bubbletea
- Integrates with AI assistants (Claude Code, Cursor)
- MIT licensed

**Why it matters:**
As AI pair programming becomes common, developers need better tools for managing context. ws addresses this gap.

**GitHub:** https://github.com/n-filatov/ws
**Episode idea:** Discuss the missing primitive of branch-scoped file lists, AI-assisted development workflows, terminal tool design with Bubbletea.

Would love to discuss on the show!

Best,
Nikita

---

## Open Source Insider

**Subject:** New tool: ws for branch-scoped file management

**Body:**

Hi Open Source Insider,

I released ws as open source (MIT) and thought your readers might be interested.

**Problem it solves:**
AI assistants find relevant files, but developers lose track of them across sessions and branches.

**Solution:**
Branch-scoped working sets with a terminal UI.

**Tech:**
- Go + Bubbletea
- Git-native design
- AI assistant integrations

**GitHub:** https://github.com/n-filatov/ws
**Show HN:** [link if applicable]

Hoping the open source community finds it useful!

Best,
Nikita

---

## General Template

**Subject:** [Tool Name] - [One-line description]

**Body:**

Hi [Newsletter Name],

I built [tool name] to solve [specific problem].

**[2-3 sentence problem description]**

**[2-3 sentence solution description]**

**Key features:**
- [Feature 1]
- [Feature 2]
- [Feature 3]

**[Code example if relevant]**

**[Why readers should care]**

**GitHub:** [Link]
**[Other links: Show HN, demo, etc.]**

**[Call to action: feedback, feature consideration, etc.]**

Best,
[Your Name]

---

## Submission URLs/Emails

- **System Override:** [Find submission form/email]
- **Golang Weekly:** https://golangweekly.com/
- **Hacker Newsletter:** Submit via HN discussion
- **Pragmatic Engineer:** gergely@omohayer.com
- **Devtools Weekly:** https://dev.tools/ (find contact)
- **CLI Quarterly:** [Find contact info]
- **Terminal Trove:** [Find submission form]
- **FLOSS Weekly:** https://flossweekly.com/ (find contact)
- **Open Source Insider:** [Find contact]

---

## Tips for Submission

1. **Personalize:** Don't copy-paste. Customize for each newsletter.
2. **Be brief:** Newsletter curators are busy. Get to the point.
3. **Show value:** Explain why *their* readers would care.
4. **Include links:** GitHub, demo, Show HN, etc.
5. **Follow up:** If no response in 2 weeks, polite follow-up.
6. **Timing:** Submit 2-3 weeks before you want it published.

---

*Last updated: 2026-03-01*
