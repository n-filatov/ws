# Show HN Post

## Title

Show HN: ws – A terminal UI for managing working sets of files (branch-scoped, AI-native)

## Post Body

Hi HN,

I built `ws` because I kept losing track of files when pair-programming with AI.

**The problem:** When you're building a feature, you jump between the same 6-10 files constantly. You ask Claude/ChatGPT "which files are relevant?" and get a list — but then what? You open them one by one, lose track, and repeat the next day.

**The solution:** `ws` keeps a focused, branch-scoped working set of files.

```bash
ws add src/auth/login.go src/auth/jwt.go
ws  # Opens TUI with tree view, git status, fuzzy search
```

**Why it's different:**

- Branch-scoped: Each git branch has its own working set. Switch branches, your context switches too.
- AI-native: Built for AI pair programming workflows (Claude Code, Cursor, Copilot)
- Git-aware: Shows inline status, auto-adds modified files
- Persistent: Survives terminal sessions

**Workflow with AI:**

1. AI finds relevant files
2. AI runs `ws add` for each file
3. You run `ws` to navigate
4. AI has focused context on what you're working on

**Tech:** Go + Bubbletea (TUI), MIT licensed

**GitHub:** https://github.com/n-filatov/ws

Would love feedback from the HN community. This is my first real CLI tool — open to all suggestions!

**Screenshots/demo:** [Will add GIF before posting]

---

## Posting Tips

**Timing:**
- Best: Tuesday-Thursday, 8-10 AM PST
- Avoid: Friday afternoon, weekends
- Check: [HN Front Page Timing](https://hackerfront.com/) for optimal times

**Title variants (if rejected):**
- "Show HN: ws – Branch-scoped file manager for AI pair programming"
- "Show HN: I built a tool to manage file context when coding with AI"
- "Show HN: ws – Keep track of files across git branches"

**Engagement strategy:**
1. Post at optimal time
2. Stay on HN for 2 hours (reply to every comment)
3. Be authentic, humble, open to feedback
4. Answer questions quickly (builds momentum)
5. Update post with demo GIF if you have one

**Common questions to prepare for:**
- "Why not use [existing tool]?" — Have comparison ready
- "How does it differ from project.nvim?" — Explain branch-scoped vs project-scoped
- "Why another file manager?" — Emphasize AI-native design
- "Will it work with my editor?" — Yes, terminal-first, works with anything

**Follow-up comments:**
- "Thanks for the feedback! I hadn't considered that. Added to roadmap."
- "Great question — ws is branch-scoped, so each feature branch has its own set."
- "Yes, it works with Cursor/Copilot/any AI assistant."

**After 24 hours:**
- Thank top commenters
- Share interesting threads on Twitter/Mastodon
- Update README based on feedback
- Add requested features to roadmap

---

## Pre-Post Checklist

- [ ] GitHub repo is polished (README clear, installation works)
- [ ] Demo GIF or screenshot ready
- [ ] Prepare responses to common questions
- [ ] Schedule 2 hours free for engagement
- [ ] Test installation on clean machine (if possible)
- [ ] Check for typos in post
- [ ] Verify GitHub link works
- [ ] Have a thick skin (HN can be brutal)

---

## Success Metrics

- **Upvotes:** 50+ = good visibility, 100+ = front page likely
- **Comments:** 20+ = good engagement
- **GitHub stars:** Track before/after
- **Installs:** Check Homebrew/APT stats after 24 hours

---

## Template for Replies

**For positive feedback:**
> "Thanks! Glad you find it useful. Let me know if you have any suggestions!"

**For questions:**
> "Great question! [Answer]. Does that make sense? Happy to elaborate."

**For criticism:**
> "Fair point. I didn't consider [X]. I'll add that to the roadmap. Thanks for the feedback."

**For feature requests:**
> "Interesting idea! I'm tracking requests [here/ GitHub issues]. Feel free to open an issue."

**For comparisons to other tools:**
> "Yes, [tool] is great! ws is different because [key difference]. They're actually complementary — you can use both."
