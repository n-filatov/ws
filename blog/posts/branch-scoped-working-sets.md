# Branch-Scoped Working Sets: A Missing Git Primitive

**A deep dive into the design and implementation of ws, a branch-scoped file manager for AI-assisted development**

---

## Introduction

Git gives us branches — isolated contexts for development. Editors give us file lists — all the files in a project. But nothing gives us **branch-scoped file lists**.

Enter `ws`: a working set manager that maintains a focused list of files per git branch. When you switch branches, your working set switches too.

This post explores the design decisions, implementation details, and philosophy behind ws.

## The Problem Space

### Context in Modern Development

Modern development involves:
1. **Multiple branches**: feature/auth, feature/payments, feature/oauth
2. **Many files**: projects have hundreds or thousands of files
3. **AI assistants**: Claude Code, Cursor, Copilot need context

### The Missing Primitive

Git knows about branches. Editors know about files. But nothing bridges the gap: **which files matter for this branch?**

Current solutions:
- **Project-scoped tools** (project.nvim, .project2): One set of files for entire project
- **Recent files** (editor MRU): Shows what you opened, not what's relevant
- **Manual lists** (harpoon marks): Requires manual management, no git awareness

**What's missing**: A branch-scoped, persistent, git-aware file list.

## Design Principles

### 1. Branch-Scoped

Each git branch has its own working set.

```
feature/auth → {auth files}
feature/payments → {payment files}
main → {different files or empty}
```

**Why**: Context switches when you switch branches. No manual cleanup.

### 2. Persistent

Working sets survive terminal sessions.

```
$ ws add src/auth/login.go
$ exit
$ # ... next day ...
$ ws
# Still shows login.go
```

**Why**: Don't make users remember what they were working on.

### 3. Git-Aware

Working sets know about git status.

```
src/auth/login.go         M  # Modified
src/auth/jwt.go           A  # Staged
src/middleware/cors.go    ?  # Untracked
```

**Why**: Developers need to know file state at a glance.

### 4. Terminal-First

Works in terminal, independent of editor.

```bash
$ ws add <files>
$ ws  # TUI
$ ws list  # Script-friendly
```

**Why**: Works with any editor, any AI assistant.

## Implementation

### Data Storage

Working sets live outside the repository:

```
~/.local/share/ws/<repo>/.workingset-<branch>
```

**Why outside repo?**
- Working sets are user-specific, not project-specific
- Don't pollute git history
- No `.gitignore` needed

**Why separate files per branch?**
- Simple and reliable
- No conflicts between branches
- Easy to clean up old branches

### File Format

Working sets are stored as newline-separated file paths:

```
src/auth/login.go
src/auth/jwt.go
src/middleware/auth.go
src/models/user.go
```

**Why plain text?**
- Human-readable (easy to debug)
- Easy to edit manually if needed
- No parsing overhead
- Git-friendly (can diff if you want)

### Git Integration

Detecting current branch:

```go
cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
branch, _ := cmd.Output()
// branch = "feature/auth\n"
```

Loading working set:

```go
workingSetFile := filepath.Join(baseDir, repoName, ".workingset-"+branch)
files := readLines(workingSetFile)
```

Getting git status:

```go
cmd := exec.Command("git", "status", "--porcelain", file)
status, _ := cmd.Output()
// Parse status: "M src/auth/login.go"
```

### The TUI (Bubbletea)

Built with [Bubbletea](https://github.com/charmbracelet/bubbletea):

```go
type model struct {
    files    []File
    cursor   int
    selected int
    filtering bool
    query    string
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "j", "down":
            m.cursor++
        case "k", "up":
            m.cursor--
        case "e":
            return m, openFile(m.files[m.cursor])
        case "/":
            m.filtering = true
        case "q":
            return m, tea.Quit
        }
    }
    return m, nil
}
```

**Key features:**
- Vim-style keybindings
- Tree view rendering
- Fuzzy search
- Git status inline

### Auto-Sync Modified Files

On startup, ws checks for modified files:

```go
cmd := exec.Command("git", "diff", "--name-only", "--cached")
staged, _ := cmd.Output()

cmd = exec.Command("git", "diff", "--name-only")
modified, _ := cmd.Output()

cmd = exec.Command("git", "ls-files", "--others", "--exclude-standard")
untracked, _ := cmd.Output()

// Add all to working set
workingSet.AddAll(staged, modified, untracked)
```

**Why auto-add?**
- Keeps working set in sync with reality
- Users don't forget to add files
- One less thing to manage manually

## Performance Considerations

### Large Working Sets

**Problem**: Some branches have 100+ files

**Solution**: Lazy loading and pagination

```go
type model struct {
    files    []File
    visible  []File  // Only visible files
    offset   int     // Scroll position
}

func (m model) View() string {
    // Only render visible files
    for _, file := range m.visible {
        // render
    }
}
```

### Large Repositories

**Problem**: git operations can be slow

**Solution**: Cache git status

```go
type Cache struct {
    status map[string]GitStatus
    expiry time.Time
}

func (c *Cache) Get(file string) GitStatus {
    if time.Now().Before(c.expiry) {
        return c.status[file]
    }
    // Refresh cache
    c.refresh()
    return c.status[file]
}
```

### Startup Time

**Target**: <100ms to open TUI

**Optimizations:**
- Read working set file once
- Lazy git status (only on visible files)
- Don't read file contents (just paths)

## Edge Cases

### Deleted Files

**Problem**: File in working set but deleted from disk

**Solution**: Mark as `[DELETED]` in TUI

```go
if _, err := os.Stat(file.Path); os.IsNotExist(err) {
    file.Status = "[DELETED]"
}
```

### Renamed Files

**Problem**: File renamed in git

**Solution**: Detect rename and update working set

```go
cmd := exec.Command("git", "diff", "--name-status", "HEAD")
// Output: "R100\told.go\tpath/to/new.go"
```

### Branch Deletion

**Problem**: Branch deleted, working set remains

**Solution**: Offer cleanup after N days

```go
if time.Since(branch.DeletedAt) > cleanupDays * 24 * time.Hour {
    // Prompt to clean up
    fmt.Printf("Clean up working set for %s? [y/N] ", branch.Name)
}
```

### Submodules

**Problem**: Submodules have their own git repos

**Solution**: Track submodule files separately

```go
if isInSubmodule(file) {
    submodulePath := getSubmodulePath(file)
    // Add to submodule's working set
} else {
    // Add to main working set
}
```

## Comparison to Alternatives

### project.nvim

| Feature | ws | project.nvim |
|---------|----|--------------|
| Scope | Branch-scoped | Project-scoped |
| Persistence | Automatic | Manual |
| Git-aware | Yes (inline status) | No |
| Editor | Terminal (any) | Neovim only |

**Trade-off**: ws requires terminal; project.nvim stays in editor.

### harpoon

| Feature | ws | harpoon |
|---------|----|----|
| Marks | Automatic (via AI) | Manual |
| Scope | Branch-scoped | Global |
| Persistence | File-based | Buffer-based |
| Git | Git-aware | Git-agnostic |

**Trade-off**: ws is more opinionated about git workflows.

### zoxide

| Feature | ws | zoxide |
|---------|----|----|
| Scope | Files | Directories |
| Context | Branch-scoped | Global |
| Workflow | AI-assisted | Directory jumping |

**Trade-off**: Different scope. ws is for files, zoxide is for directories.

## AI Integration Patterns

### Claude Code Plugin

```json
{
  "commands": [
    {
      "name": "ws:map",
      "description": "Search and add files",
      "parameters": [{"name": "feature", "required": true}]
    }
  ]
}
```

### Cursor Workflow

```bash
# Cursor finds files
# User adds to ws
ws add src/auth/*.go

# Cursor sees focused context
ws list
```

### Copilot Integration

```bash
# Give Copilot focused context
ws list | xargs -I {} copilot analyze {}
```

## Future Directions

### Remote Working Sets

Sync working sets across machines:

```
~/.local/share/ws/<repo>/.workingset-<branch>
→ Push to GitHub gist
→ Pull on other machines
```

### Collaborative Working Sets

Share working sets with team:

```
git ws push  # Push to repo
git ws pull  # Pull from repo
```

### Smart File Discovery

Use AST analysis to find related files:

```
ws add --smart src/auth/login.go
# Analyzes imports, adds related files
```

## Lessons Learned

### 1. Simple Storage Wins

Plain text files beat databases:
- Easier to debug
- No migrations
- Human-editable
- Tool-friendly

### 2. Git Native > Git Wrapper

Embrace git, don't fight it:
- Use git branches as isolation
- Use git status for file state
- Use git hooks for automation

### 3. Terminal First > Editor First

Terminal-first tools work everywhere:
- Any editor
- Any AI assistant
- Any platform

### 4. Opinionated > Flexible

Branch-scoped is opinionated but:
- Solves real problems
- Easy to understand
- Reduces cognitive load

## Conclusion

ws fills a gap in the development toolchain: **branch-scoped file lists**.

By maintaining working sets per branch, it gives developers:
- Automatic context switching
- Focused AI assistance
- Persistent file tracking

The design is simple:
- Plain text storage
- Git-native workflows
- Terminal-first architecture

Sometimes the missing primitive isn't complicated. It just needs someone to build it.

---

**Built with Go + Bubbletea**
**License: MIT**
**GitHub: https://github.com/n-filatov/ws**

---

*Published: 2026-03-01*
*Author: Nikita Filatov*
*Tags: #golang #devtools #git #terminal #design #architecture*
*Length: ~1800 words*
