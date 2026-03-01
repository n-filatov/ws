# Infographic Specification

**How ws works — Visual representation**

---

## Overview

**Purpose:** Explain ws visually in a shareable format
**Target:** Twitter, Reddit, blog posts
**Format:** PNG/JPG for social, SVG for editing
**Size:** 1200x630 (Twitter card), 1080x1080 (Instagram/square)

---

## Concept: Before/After Workflow

### Panel 1: The Problem (Before)

**Layout:**
```
┌─────────────────────────────────┐
│  BEFORE: Working with AI        │
├─────────────────────────────────┤
│                                 │
│  🤖 AI: "Found 12 files:"       │
│                                 │
│  📁 src/auth/login.go           │
│  📁 src/auth/jwt.go             │
│  📁 src/middleware/auth.go      │
│  📁 src/models/user.go          │
│  📁 src/handlers/auth.go        │
│  📁 ... (7 more)                │
│                                 │
│  👤 You: *opens files one by one*│
│                                 │
│  ⏰ Next day:                    │
│  🤖 AI: "What were we working on?"│
│  👤 You: "I forgot, let me ask again"│
│                                 │
│  ❌ Lost context. Repeating.    │
└─────────────────────────────────┘
```

**Colors:**
- Background: #1a1a2e (dark)
- Text: #e0e0e0 (light)
- AI: #a855f7 (purple)
- Files: #3b82f6 (blue)
- Problem indicator: #ef4444 (red)

---

### Panel 2: The Solution (After)

**Layout:**
```
┌─────────────────────────────────┐
│  AFTER: Using ws                │
├─────────────────────────────────┤
│                                 │
│  🤖 AI: "Found 12 files:"       │
│                                 │
│  💡 AI: Adding to ws...         │
│  $ ws add src/auth/*.go         │
│                                 │
│  📺 TERMINAL (ws TUI):          │
│  ┌───────────────────────────┐  │
│  │ src/auth/                │  │
│  │ 📄 login.go          M   │  │
│  │ 📄 jwt.go            A   │  │
│  │ 📄 middleware.go     ?   │  │
│  │                         │  │
│  │ [e]dit [f]ind [q]uit    │  │
│  └───────────────────────────┘  │
│                                 │
│  ⏰ Next day:                    │
│  👤 You: $ ws                   │
│  📺 All 12 files still there!   │
│                                 │
│  ✅ Context preserved.          │
└─────────────────────────────────┘
```

**Colors:**
- Background: #1a1a2e (dark)
- Text: #e0e0e0 (light)
- AI: #a855f7 (purple)
- Terminal: #10b981 (green)
- Success indicator: #22c55e (green)

---

## Concept: Branch Switching

### Diagram Format

```
┌────────────────────────────────────────────────────────┐
│                   BRANCH-SCOPED CONTEXT                │
├────────────────────────────────────────────────────────┤
│                                                        │
│  feature/auth                    feature/payments       │
│  ┌──────────────────┐          ┌──────────────────┐  │
│  │ Working Set:     │          │ Working Set:     │  │
│  │                  │          │                  │  │
│  │ 📄 login.go      │          │ 📄 processor.go  │  │
│  │ 📄 jwt.go        │          │ 📄 stripe.go     │  │
│  │ 📄 middleware.go │          │ 📄 webhook.go    │  │
│  │ 📄 user.go       │          │ 📄 refund.go      │  │
│  │                  │          │                  │  │
│  │ (8 auth files)   │          │ (6 payment files)│  │
│  └──────────────────┘          └──────────────────┘  │
│           ↓                              ↓             │
│    git checkout              git checkout           │
│  feature/payments            feature/auth             │
│           ↓                              ↓             │
│  ┌──────────────────┐          ┌──────────────────┐  │
│  │ Working Set:     │          │ Working Set:     │  │
│  │                  │          │                  │  │
│  │ (6 payment files)│          │ (8 auth files)   │  │
│  └──────────────────┘          └──────────────────┘  │
│                                                        │
│  ✨ Context switches automatically                    │
└────────────────────────────────────────────────────────┘
```

**Animation (if video/GIF):**
1. Start on feature/auth
2. Show 8 auth files
3. `git checkout feature/payments`
4. Files morph into 6 payment files
5. `git checkout feature/auth`
6. Files morph back into 8 auth files

---

## Concept: Feature Comparison

### Table Format

```
┌──────────────────────────────────────────────────────────┐
│                   ws VS OTHER TOOLS                       │
├──────────────────────────────────────────────────────────┤
│                                                           │
│  ┌─────────────┬──────────┬──────────┬──────────┐      │
│  │  Feature    │    ws    │project.nv│ harpoon  │      │
│  ├─────────────┼──────────┼──────────┼──────────┤      │
│  │ Branch-scope│    ✅     │    ❌    │    ❌    │      │
│  │ Git-aware   │    ✅     │    ❌    │    ❌    │      │
│  │ AI-native   │    ✅     │    ❌    │    ❌    │      │
│  │ Auto-switch │    ✅     │    ❌    │    ❌    │      │
│  │ TUI         │    ✅     │    ❌    │    N/A    │      │
│  │ Manual setup│    ❌     │    ✅    │    ✅    │      │
│  │ Editor int. │    ❌     │    ✅    │    ✅    │      │
│  └─────────────┴──────────┴──────────┴──────────┘      │
│                                                           │
│  💡 ws: Branch-scoped, Git-aware, AI-native              │
└──────────────────────────────────────────────────────────┘
```

**Colors:**
- ✅ Yes: #22c55e (green)
- ❌ No: #ef4444 (red)
- N/A: #6b7280 (gray)

---

## Text Layout Guide

### Font Sizes

```
Title: 48px bold
Subtitle: 32px semi-bold
Body: 20px regular
Code: 16px monospace
Annotations: 14px regular
```

### Fonts

```
Headings: Inter, SF Pro Display, system-ui
Body: Inter, SF Pro Text, system-ui
Code: JetBrains Mono, Fira Code, monospace
```

### Color Palette

```
Background: #1a1a2e (dark blue)
Text: #e0e0e0 (off-white)
Accent: #3b82f6 (blue)
Success: #22c55e (green)
Warning: #f59e0b (amber)
Error: #ef4444 (red)
Purple: #a855f7 (AI)
Green: #10b981 (terminal)
```

---

## Sizes for Different Platforms

### Twitter/X
- **Size:** 1200x675 (16:9)
- **Format:** PNG < 5MB
- **Use:** Tweets, replies

### LinkedIn
- **Size:** 1200x627 (1.91:1)
- **Format:** PNG < 5MB
- **Use:** Article posts

### Instagram Square
- **Size:** 1080x1080 (1:1)
- **Format:** PNG/JPG < 1MB
- **Use:** Feed posts

### Instagram Story
- **Size:** 1080x1920 (9:16)
- **Format:** PNG/JPG < 1MB
- **Use:** Stories, highlights

### Reddit
- **Size:** 1920x1080 (16:9)
- **Format:** PNG
- **Use:** r/programming, r/devtools

---

## Alternative Concepts

### Concept A: Simple Flow

```
┌─────────────────────────────────────────┐
│  ws WORKFLOW IN 3 STEPS                │
├─────────────────────────────────────────┤
│                                         │
│  1️⃣  AI FINDS FILES                   │
│     🤖 "Found 8 auth files"            │
│                                         │
│       ↓                                 │
│                                         │
│  2️⃣  ADD TO ws                         │
│     $ ws add src/auth/*.go              │
│     ✅ Files added to working set       │
│                                         │
│       ↓                                 │
│                                         │
│  3️⃣  NAVIGATE WITH ws                   │
│     $ ws                                │
│     📺 See all files, press 'e' to open│
│                                         │
└─────────────────────────────────────────┘
```

### Concept B: Why ws?

```
┌─────────────────────────────────────────┐
│  3 REASONS TO USE ws                    │
├─────────────────────────────────────────┤
│                                         │
│  🌿 BRANCH-SCOPED                       │
│     Each branch has its own file list  │
│     Switch branches → files switch too │
│                                         │
│  🤖 AI-NATIVE                           │
│     Built for AI pair programming      │
│     Give AI focused context             │
│                                         │
│  💾 PERSISTENT                          │
│     Survives terminal sessions          │
│     Open ws tomorrow → files still there│
│                                         │
└─────────────────────────────────────────┘
```

---

## Production Tools

**Recommended tools:**
- Figma (design, collaboration)
- Canva (quick templates)
- Adobe Illustrator (professional)
- Excalidraw (hand-drawn style)
- Mermaid.js (diagrams, code-based)

**For developers:**
- Generate from code (diagram as code)
- VS Code extensions (Mermaid, PlantUML)
- GitHub Actions (auto-generate on commit)

---

## Export Checklist

- [ ] Correct aspect ratio for platform
- [ ] File size under limit
- [ ] Colors match brand palette
- [ ] Text is readable at small sizes
- [ ] GitHub link visible (if applicable)
- [ ] Logo included (if applicable)
- [ ] Call-to-action included (GitHub, install command)
- [ ] Alt text prepared (for accessibility)

---

## Usage

**Twitter:**
- Upload as image
- Add caption: "How ws works 🧵"
- Link to GitHub in thread

**Reddit:**
- Upload to r/devtools, r/programming
- Title: "Built ws to manage file context with AI"
- Link to GitHub in comments

**Blog:**
- Include in blog posts
- Break down into sections
- Use each concept as separate image

**Presentations:**
- Use in slides
- Explain ws visually
- Show before/after comparison

---

*Last updated: 2026-03-01*
