# ws Claude Code Plugin

Manage your branch-scoped working sets directly from Claude Code.

## Overview

The `ws` plugin integrates the `ws` CLI tool with Claude Code, allowing you to manage working sets of files without leaving your conversation.

## What is ws?

`ws` is a terminal UI tool that keeps track of files relevant to your current git branch. It's designed for AI-assisted development, helping you maintain context across sessions and branches.

**Key features:**
- **Branch-scoped**: Each branch has its own working set
- **Persistent**: Survives terminal sessions
- **Git-aware**: Shows inline status, auto-syncs modified files
- **AI-native**: Built for AI pair programming workflows

## Installation

Install the plugin from the Claude Code marketplace:

```
/plugin marketplace add n-filatov/ws
/plugin install ws@n-filatov-ws
```

Then install the `ws` CLI:

```
/ws:install
```

This will detect your OS and install using the appropriate method (Homebrew, APT, or source).

## Available Commands

### /ws:map <feature description>

Search your codebase for files related to a feature and add them to your working set.

**Example:**
```
/ws:map user authentication flow
```

Claude will:
1. Search your codebase for relevant files
2. Identify authentication-related files
3. Run `ws add` for each file automatically
4. Report which files were added

**Output:**
```
I found 8 files related to user authentication:
- src/auth/login.go
- src/auth/jwt.go
- src/middleware/auth.go
- src/models/user.go
- src/handlers/auth.go
- internal/crypto/password.go
- config/auth.yaml
- tests/auth_test.go

Adding to your working set...
```

### /ws:install

Install the `ws` CLI tool.

Automatically detects your operating system and package manager, then installs using the appropriate method:

- **macOS/Linux with Homebrew**: `brew install n-filatov/tap/ws`
- **Debian/Ubuntu**: Adds repository and installs via apt
- **Other**: Builds from source

**Example:**
```
/ws:install
```

**Output:**
```
Detected: macOS
Installing via Homebrew...
brew tap n-filatov/tap
brew install ws
✓ ws installed successfully
```

### /ws:add <files...>

Manually add specific files to your working set.

**Example:**
```
/ws:add src/auth/login.go src/auth/jwt.go config/auth.yaml
```

### /ws:remove <file>

Remove a file from your working set.

**Example:**
```
/ws:remove src/auth/old-login.go
```

### /ws:open

Open the ws TUI (terminal UI) to view and navigate your working set.

**Example:**
```
/ws:open
```

This opens the interactive TUI where you can:
- Press `e` to open files in your editor
- Press `/` to fuzzy-search
- Press `r` to refresh git status
- Press `q` to quit

### /ws:clear

Clear all files from the current branch's working set.

**Example:**
```
/ws:clear
```

**Use with caution:** This removes all files from your working set for the current branch.

### /ws:list

List all files in your current working set.

**Example:**
```
/ws:list
```

**Output:**
```
Your working set contains 8 files:
src/auth/login.go
src/auth/jwt.go
src/middleware/auth.go
src/models/user.go
src/handlers/auth.go
internal/crypto/password.go
config/auth.yaml
tests/auth_test.go
```

## Usage Workflow

### Starting a New Feature

1. **Tell Claude what you're building**
   ```
   I'm building OAuth login functionality
   ```

2. **Map the feature**
   ```
   /ws:map OAuth login
   ```

3. **Claude finds and adds files**
   ```
   I found 6 files related to OAuth login.
   Adding to your working set...
   ```

4. **View your working set**
   ```
   /ws:open
   ```

5. **Navigate and work**
   - Use the TUI to open files with `e`
   - Search with `/`
   - Refresh with `r`

### Switching Branches

When you switch git branches, your working set automatically changes:

```
git checkout feature/payments
/ws:list
# Shows payment-related files, not auth files
```

### Growing Your Working Set

As you discover more files:

```
/ws:add src/payment/refund.go src/payment/webhook.go
```

## How It Works

The plugin wraps the `ws` CLI commands:

| Plugin Command | CLI Equivalent |
|----------------|----------------|
| `/ws:map <feature>` | Claude searches + `ws add <files>` |
| `/ws:add <files>` | `ws add <files>` |
| `/ws:remove <file>` | `ws rm <file>` |
| `/ws:list` | `ws list` |
| `/ws:open` | `ws` |
| `/ws:clear` | `ws clear` |
| `/ws:install` | Installation script |

## Integration with Claude Code

### Context Awareness

Claude Code can see your working set and provide more focused responses:

```
You: What files am I working with?

Claude: Let me check...
/ws:list
Claude: You have 8 files in your working set, all related to authentication.
```

### Better File Discovery

```
You: Find files for password reset

Claude: I found 3 files for password reset:
- src/auth/reset.go
- src/auth/token.go
- templates/reset.html

/ws:add src/auth/reset.go src/auth/token.go templates/reset.html
```

## Troubleshooting

### "ws command not found"

The `ws` CLI isn't installed. Run:
```
/ws:install
```

### "No working set found"

The current branch doesn't have a working set yet. Add files:
```
/ws:add <files>
```

### "Plugin not found"

Install the plugin:
```
/plugin marketplace add n-filatov/ws
/plugin install ws@n-filatov-ws
```

### "Permission denied"

The `ws` binary might not be executable. Try reinstalling:
```
/ws:install
```

## Tips and Best Practices

1. **Start with /ws:map**: Let Claude discover files automatically
2. **Keep sets focused**: Only add files relevant to current feature
3. **Use /ws:list**: Check what's in your working set before asking questions
4. **Let it grow organically**: Add files as you discover them
5. **Branch often**: Each branch gets its own set, switch freely

## Examples

### Example 1: Building a Feature

```
You: I'm working on user registration

/ws:map user registration

Claude: Found 6 files:
- src/auth/register.go
- src/auth/validator.go
- src/models/user.go
- templates/register.html
- tests/register_test.go
- config/email.yaml

Adding to ws...

You: Now explain how registration works

Claude: [Analyzes the 6 files in your working set]
```

### Example 2: Debugging

```
You: I'm getting a JWT error

/ws:list

Claude: You have 8 files. I see jwt.go in your set.
Let me check it for issues...

Claude: Found a bug on line 42 of jwt.go
```

### Example 3: Code Review

```
You: Review the authentication changes

/ws:list

Claude: Reviewing your 8 auth files...
[Provides focused review based on working set]
```

## Related Resources

- **GitHub**: https://github.com/n-filatov/ws
- **Documentation**: https://github.com/n-filatov/ws#readme
- **VS Code Extension**: `n-filatov/ws-vscode`
- **Zed Extension**: `n-filatov/ws-zed`

## Support

For issues or questions:
- GitHub Issues: https://github.com/n-filatov/ws/issues
- Author: Nikita Filatov

---

**Version:** 1.1.0
**Last updated:** 2026-03-01
