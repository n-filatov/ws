# Using ws with GitHub Copilot

A guide to integrating `ws` with GitHub Copilot for context-aware AI assistance.

## Overview

GitHub Copilot is your AI pair programmer. `ws` enhances Copilot by maintaining a focused, branch-scoped list of files, giving Copilot better context and more relevant suggestions.

## Quick Start

### 1. Install ws

```bash
# macOS/Linux
brew tap n-filatov/tap && brew install ws

# Go
go install github.com/n-filatov/ws@latest
```

### 2. Basic Workflow

```bash
# Add files to your working set
ws add src/auth/login.go src/auth/jwt.go src/middleware/auth.go

# List files for Copilot context
ws list
```

## Key Integration Pattern: `ws list | xargs`

The most powerful pattern is piping `ws list` to xargs:

```bash
# Grep across your working set
ws list | xargs grep "TODO"

# Count lines
ws list | xargs wc -l

# Run linting
ws list | xargs golint

# Open in editor
ws list | xargs nvim
```

## Workflows

### Workflow 1: Give Copilot Focused Context

**Step 1: Build your working set**

```bash
ws add src/auth/login.go src/auth/jwt.go src/middleware/auth.go
```

**Step 2: List files for Copilot**

```bash
ws list
```

**Output:**
```
src/auth/login.go
src/auth/jwt.go
src/middleware/auth.go
```

**Step 3: Provide context to Copilot**

In your editor, select the output and paste it into Copilot Chat:

```
I'm working on these files:
src/auth/login.go
src/auth/jwt.go
src/middleware/auth.go

How does authentication work?
```

Copilot now provides focused, relevant answers.

### Workflow 2: Code Analysis

```bash
# Find TODOs in your working set
ws list | xargs grep "TODO"

# Paste results into Copilot Chat
```

```
I have these TODOs in my auth files:
- [TODO] Add rate limiting
- [TODO] Implement refresh tokens

How should I prioritize?
```

### Workflow 3: Test Coverage

```bash
# Check test coverage for working set
ws list | grep "_test.go" | xargs go test
```

Then ask Copilot:

```
My auth tests are failing. How can I fix them?
```

### Workflow 4: Refactoring

```bash
# List files you want to refactor
ws list
```

Tell Copilot:

```
Refactor these files to use the new error handling pattern:
[paste ws list output]
```

## Copilot Chat Prompts

### Prompt 1: Context Setting

```
I'm working on a feature with these files:
[paste output of: ws list]

Give me an overview of what this code does.
```

### Prompt 2: Code Review

```
Review these files for security issues:
[paste output of: ws list]

Focus on authentication and authorization.
```

### Prompt 3: Documentation

```
Generate documentation for these files:
[paste output of: ws list]

Include function descriptions and examples.
```

### Prompt 4: Testing

```
Write unit tests for these files:
[paste output of: ws list]

Focus on edge cases and error conditions.
```

### Prompt 5: Refactoring

```
Refactor these files to use dependency injection:
[paste output of: ws list]

Maintain backward compatibility.
```

## Advanced Patterns

### Pattern 1: Inline Comments

Use Copilot inline with ws:

```bash
# In your editor
// ws add src/auth/login.go

// Copilot suggests next line
// ws add src/auth/jwt.go

// Copilot suggests next line
// ws add src/middleware/auth.go
```

### Pattern 2: Script Integration

Create a script that combines ws and Copilot:

```bash
#!/bin/bash
# auth-context.sh

echo "I'm working on authentication:"
ws list
echo ""
echo "Key questions:"
echo "1. How does JWT validation work?"
echo "2. Where are auth middlewares used?"
echo "3. What's the error handling strategy?"
```

Run it and paste into Copilot Chat.

### Pattern 3: Branch-Aware Context

```bash
# Switch branches
git checkout feature/oauth

# ws shows different files
ws list

# Paste into Copilot
I'm now working on OAuth. Here are the files:
[paste output]
```

Copilot adapts to new context.

## Copilot Labs Features

If you have Copilot Labs:

### 1. Code Explanation

```bash
# Select files
ws list

# In Copilot Labs: "Explain this code"
# Paste your ws list output
```

### 2. Code Translation

```
Translate these files from Go to Rust:
[paste output of: ws list]
```

### 3. Brush Up

```
Brush up the code style in these files:
[paste output of: ws list]
```

## Tips for Copilot Users

### Tip 1: Always Set Context First

```bash
# Before asking questions, run
ws list

# Paste into Copilot
I'm working on: [paste list]
```

### Tip 2: Use Descriptive Prompts

Instead of:
```
What does this do?
```

Use:
```
Based on these files [paste ws list], how does authentication work?
```

### Tip 3: Leverage Branch Scoping

```bash
# Different context per branch
git checkout feature/auth && ws list
# Paste into Copilot

git checkout feature/payments && ws list
# Paste into Copilot (different context!)
```

### Tip 4: Combine with Git

```bash
# See changed files in working set
ws list | xargs git diff

# Paste into Copilot
Review these changes: [paste output]
```

### Tip 5: Iterative Refinement

```bash
# Start with broad set
ws add src/auth/*.go

# Ask Copilot questions
# Narrow down based on answers
ws rm src/auth/unused.go

# Ask more focused questions
```

## Common Scenarios

### Scenario 1: Onboarding

```
You: I'm new to this codebase. Help me understand authentication.

[Run: ws list]
[Create auth working set]

You: I'm working on these files:
[paste ws list]

Copilot: I'll explain how authentication works...
```

### Scenario 2: Bug Hunting

```bash
# Find potential issues
ws list | xargs grep -i "bug\|todo\|fixme"

# Paste into Copilot
I found these potential issues in my auth files:
[paste output]

How should I fix them?
```

### Scenario 3: Feature Addition

```
You: I need to add password reset.

[Run: ws list to see current auth files]

You: Based on these files [paste list], where should I add password reset?

Copilot: Based on the structure, add it here...
```

### Scenario 4: Code Review

```bash
# List files for review
git diff --name-only | ws add -

ws list

# Paste into Copilot
Review these changes [paste list] for:
1. Security
2. Performance
3. Code style
```

## Troubleshooting

### "Copilot gives generic answers"

**Solution:** Provide focused context with `ws list`

```
Instead of: "How do I authenticate users?"

Use: "How do I authenticate users in this codebase?
I'm working on these files:
[paste ws list output]"
```

### "Copilot loses context"

**Solution:** Re-establish context with `ws list`

```
Let me refocus you. I'm working on:
[paste ws list output]
```

### "Copilot suggests irrelevant code"

**Solution:** Your working set might be too broad

```bash
# Remove irrelevant files
ws rm unrelated/file.go

# Re-establish context
ws list
```

## Best Practices

1. **Always set context first**: Run `ws list` before asking questions
2. **Keep sets focused**: Only files for current feature
3. **Use prompts**: Paste `ws list` output into Copilot Chat
4. **Branch switch freely**: ws maintains context per branch
5. **Iterate**: Refine working set as you learn

## Integration with VS Code

If using Copilot in VS Code:

### 1. Integrated Terminal

```
Open: View > Terminal
Run: ws list
Select output → Copilot Chat: "Explain these files"
```

### 2. Code Actions

```go
// In your editor
// ws:map user authentication
// Copilot might suggest running ws add commands
```

### 3. Multi-file Editing

```bash
# List files
ws list

# In VS Code: Ctrl+Shift+P
# "Copilot: Edit in these files"
# Paste ws list output
```

## Example Session

```bash
# Start your feature
git checkout feature/oauth

# Build working set
ws add src/auth/oauth.go src/auth/providers.go config/oauth.yaml

# Get context
ws list

# Output:
src/auth/oauth.go
src/auth/providers.go
config/oauth.yaml

# Paste into Copilot Chat:
I'm implementing OAuth with these files:
src/auth/oauth.go
src/auth/providers.go
config/oauth.yaml

How should I handle token refresh?

Copilot: [provides focused, relevant answer based on your 3 files]
```

## Resources

- **ws GitHub**: https://github.com/n-filatov/ws
- **GitHub Copilot**: https://github.com/features/copilot
- **Copilot Chat**: https://docs.github.com/en/copilot

---

**Last updated:** 2026-03-01
