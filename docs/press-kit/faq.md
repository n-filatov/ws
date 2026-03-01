# Frequently Asked Questions (FAQ)

## General Questions

### What is ws?

**ws** is a terminal UI tool that maintains a focused list of files relevant to your current git branch. It's designed for AI-assisted development, helping you keep track of which files matter for the feature you're working on.

### Why did you build ws?

I built ws because I kept losing track of files when pair-programming with AI. AI assistants like Claude Code and Cursor are great at identifying relevant files, but there was no persistent way to remember those files across sessions or branches. ws solves this by maintaining branch-scoped working sets.

### Is ws open source?

Yes! ws is released under the MIT license. Source code is available on GitHub: https://github.com/n-filatov/ws

## Technical Questions

### What technology is ws built with?

ws is written in Go and uses Bubbletea for the terminal UI. It's designed to be fast, lightweight, and dependency-free.

### What platforms does ws support?

ws supports:
- macOS (via Homebrew)
- Linux (via Homebrew, APT, or source)
- Windows (via Scoop or WSL)
- Any platform with Go installed

### How does ws store data?

Working sets are stored in `~/.local/share/ws/<repo>/` — outside your git repository, so they're never committed. Each branch gets its own file: `.workingset-<branch-name>`.

### Is ws safe to use?

Yes. ws only stores file paths and timestamps. It doesn't read file contents except for git status checking. All data is stored locally on your machine.

## Usage Questions

### How is ws different from project.nvim or other project tools?

ws is **branch-scoped** while most project tools are project-scoped. This means:
- Each git branch has its own working set
- Context switches automatically when you switch branches
- No manual cleanup needed when finishing a feature

### How does ws work with AI assistants?

ws integrates with AI assistants in several ways:

**Claude Code:**
- Native plugin: `/plugin install ws@n-filatov-ws`
- Commands: `/ws:map`, `/ws:add`, `/ws:remove`, `/ws:open`

**Cursor:**
- Use `ws list` output to give Cursor focused context
- Cursor sees exactly which files you're working on

**GitHub Copilot:**
- Pipe `ws list` to xargs for analysis
- `ws list | xargs grep "TODO"`
- `ws list | xargs wc -l`

### Can I use ws with my editor?

Yes! ws is terminal-first and editor-agnostic. It works with:
- VS Code (official extension: `n-filatov/ws-vscode`)
- Zed (official extension: `n-filatov/ws-zed`)
- Vim/Neovim (set as `editor` in `~/.wsconfig`)
- Any editor that supports opening files from terminal

### How do I get started?

```bash
# Install
brew tap n-filatov/tap && brew install ws

# Add files to your working set
ws add src/auth/login.go src/auth/jwt.go

# Open ws TUI
ws
```

## Feature Questions

### What happens when I switch branches?

Your working set automatically switches to show files for the current branch. Each branch maintains its own independent working set.

### Can I use ws across multiple repositories?

Yes! ws maintains separate working sets for each git repository. You can have different working sets for different projects.

### Does ws sync with git?

Yes, ws is git-aware:
- Shows inline git status (M=modified, A=staged, ?=untracked)
- Auto-adds modified files on refresh
- Respects `.gitignore`

### Can I search within my working set?

Yes! Press `/` in the ws TUI to fuzzy-search your working set. Matches are highlighted as you type.

## Integration Questions

### Does ws replace lazygit or other git tools?

No, ws complements git tools like lazygit. Use lazygit for commits, branches, and git operations. Use ws for navigating files relevant to your current feature.

### Can I use ws with tmux?

Yes! ws works great with tmux. Run ws in a split pane while you work in another pane.

### Does ws work with git worktrees?

Yes! Each git worktree maintains its own working sets, independent of other worktrees.

## Business Questions

### Can I use ws in my company?

Yes! ws is released under the MIT license, which permits commercial use. No attribution required (though appreciated).

### Do you offer support?

Community support is available through GitHub Issues. For enterprise support or custom integrations, contact keddofilatov@gmail.com.

### Will ws remain free and open source?

Yes. ws is and will remain free and open source under the MIT license.

## Future Plans

### What's coming in ws?

Planned features include:
- Enhanced AI assistant integrations
- Remote working set synchronization
- Additional editor extensions
- Performance optimizations for large codebases

Check the GitHub issues and roadmap for updates: https://github.com/n-filatov/ws/issues

### Can I contribute?

Yes! Contributions are welcome. See CONTRIBUTING.md in the repository. Good areas for contributions:
- Additional editor extensions
- Performance improvements
- Bug fixes
- Documentation

## Press & Media

### Where can I get screenshots/demos?

Screenshots and demos are available in the press kit:
- Screenshots: `docs/press-kit/screenshots/`
- Demo video: [YouTube link]
- Logo: `docs/press-kit/logo/`

### Can I interview the author?

For interviews, speaking opportunities, or press inquiries, contact:
- Email: keddofilatov@gmail.com
- GitHub: https://github.com/n-filatov

### Do you have a media kit?

Yes! The complete press kit is available at `docs/press-kit/` including:
- One-pager overview
- Author bio
- Logo in multiple formats
- Screenshots
- FAQ (this document)

## Misc Questions

### Why the name "ws"?

**ws** stands for "working set" — a computer science term for the set of resources a process is actively working with. It's short, memorable, and describes exactly what the tool does.

### What does "AI-native" mean?

"AI-native" means ws was designed from the ground up for AI-assisted development workflows. Unlike tools adapted for AI, ws was built specifically to make AI pair programming more effective.

### Where can I learn more?

- **GitHub**: https://github.com/n-filatov/ws
- **Documentation**: https://github.com/n-filatov/ws#readme
- **Blog**: [Link to blog posts]
- **Twitter**: @nfilatov (if applicable)

---

*Last updated: 2026-03-01*
