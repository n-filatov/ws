# `ws` — Working Set Tool Requirements

## Purpose

A terminal TUI tool for managing a focused set of files during feature development. Designed to work alongside Claude Code (which adds files via CLI) and let the developer navigate and open them in a configured editor. Inspired by lazygit's file list UX.

---

## Tech Stack

- **Language:** Go
- **TUI:** Bubbletea + Bubbles + Lipgloss (Charm ecosystem)
- **Distribution:** Single binary installed to `~/.local/bin/ws`

---

## CLI Interface

```bash
ws                          # Open TUI (same as ws open)
ws add <file> [file...]     # Add one or more files to the working set
ws rm <file>                # Remove a file from the working set
ws list                     # Print all tracked files to stdout
ws clear                    # Clear all files from the working set
```

`ws add` is designed to be called by Claude Code from a separate terminal window. It writes to the same `.workingset-<branch>` file the TUI reads from.

---

## Scope & Storage

- Scope is determined by the **current git repo**, detected via `git rev-parse --show-toplevel`
- Running `ws` from anywhere inside the repo uses the same working set
- Working set is **per git branch** — switching branches gives you a fresh isolated scope
- File location: `<git-root>/.workingset-<branch-name>`
  - Examples: `.workingset-main`, `.workingset-feature-auth`
- Format: one absolute path per line, plain text
- All `.workingset-*` files should be added to `.gitignore`
- `ws add` deduplicates automatically — adding the same file twice is a no-op
- Accepts relative or absolute paths, always stores as absolute

---

## TUI Layout

Single full-screen panel — a navigable list of tracked files, similar to the file panel in lazygit.

Each row shows:
- Relative file path (relative to git root)
- Git status indicator (`M` modified, `A` added, `?` untracked)
- Visual indicator if the file no longer exists on disk (e.g. dimmed or struck through)

Bottom bar shows keybinding hints.

---

## TUI Keybindings

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `e` | Open selected file in configured editor |
| `a` | Show inline input prompt to type a file path to add |
| `d` | Remove selected file from working set (with git prompt if file has changes) |
| `r` | Refresh — re-sync git-modified files and re-render list |
| `q` / `Esc` | Quit |

---

## Feature Requirements

### 1. Add Files

- `ws add <file>` adds one or multiple files to the current branch's working set
- Deduplicates — adding an already-tracked file is silently ignored
- Also available inside TUI via `a` key, which shows an input field to type a path

### 2. Auto-add Git-modified Files

- On TUI open and on `r` (refresh), run `git status` and automatically add all files with any git status (`M`, `A`, `??`) to the working set
- Silent — no confirmation needed
- This ensures files you're actively editing always appear in the list

### 3. Remove File

- Press `d` on a selected file to remove it from the working set
- If the file has any git changes (`M`, `A`, `??`), show a prompt before acting:

  ```
  file.ts has uncommitted changes. Revert with git? [y/N/cancel]
  ```

  - `y` — runs `git checkout -- <file>`, then removes from working set
  - `N` — removes from working set only, file changes are kept
  - `cancel` / `Esc` — does nothing, returns to list

### 4. Refresh

- Press `r` to re-run `git status`, auto-add any newly modified files, and re-render the list
- When `ws add` is called from another terminal, press `r` to pick up the new file

### 5. Navigate

- `j` / `↓` moves selection down
- `k` / `↑` moves selection up
- Selection wraps at top and bottom of list

### 6. Open File in Editor

- Press `e` to open the selected file in the configured editor
- Default editor: `vim`
- Configurable via `~/.wsconfig`:
  ```
  editor=zed
  ```
- The editor command receives the absolute file path as its argument

### 7. Branch-scoped Working Sets

- Each git branch has its own `.workingset-<branch>` file at the git root
- Switching branches and running `ws` automatically loads the correct working set
- No bleed between branches — each feature has its own isolated file scope

---

## Claude Code Integration

The primary way files get added during development is Claude Code running `ws add` in the terminal. Suggested prompt to include at the start of a Claude Code session:

> "When you identify files related to this feature, run `ws add <filepath>` for each one."

Or explicitly:

> "Find all files related to the auth feature and add them to my working set using `ws add`."

---

## Non-Goals (v1)

- No file preview panel
- No multiple named working sets per branch
- No real-time file watching (refresh is manual via `r`)
- No sync across machines
- No VSCode / Zed extension
- No file tree view — flat list only
