package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nfilatov/ws/internal/config"
	"github.com/nfilatov/ws/internal/git"
	"github.com/nfilatov/ws/internal/store"
	"github.com/nfilatov/ws/internal/tui"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 || args[0] == "open" {
		runTUI()
		return
	}

	switch args[0] {
	case "add":
		if len(args) < 2 {
			fatalf("usage: ws add <file> [file...]\n")
		}
		wsPath, err := workingSetPath()
		if err != nil {
			fatalf("error: %v\n", err)
		}
		if err := store.Add(wsPath, args[1:]...); err != nil {
			fatalf("error: %v\n", err)
		}
		fmt.Printf("added %d file(s)\n", len(args)-1)

	case "rm":
		if len(args) < 2 {
			fatalf("usage: ws rm <file>\n")
		}
		wsPath, err := workingSetPath()
		if err != nil {
			fatalf("error: %v\n", err)
		}
		if err := store.Remove(wsPath, args[1]); err != nil {
			fatalf("error: %v\n", err)
		}
		fmt.Printf("removed %s\n", args[1])

	case "list":
		wsPath, err := workingSetPath()
		if err != nil {
			fatalf("error: %v\n", err)
		}
		files, err := store.Load(wsPath)
		if err != nil {
			fatalf("error: %v\n", err)
		}
		for _, f := range files {
			fmt.Println(f)
		}

	case "clear":
		wsPath, err := workingSetPath()
		if err != nil {
			fatalf("error: %v\n", err)
		}
		if err := store.Save(wsPath, []string{}); err != nil {
			fatalf("error: %v\n", err)
		}
		fmt.Println("working set cleared")

	default:
		fatalf("unknown command %q\n\nUsage:\n  ws                    open TUI\n  ws add <file>...      add files\n  ws rm <file>          remove a file\n  ws list               list tracked files\n  ws clear              clear working set\n", args[0])
	}
}

func runTUI() {
	root, err := git.RootDir()
	if err != nil {
		fatalf("error: %v\n", err)
	}
	branch, err := git.CurrentBranch()
	if err != nil {
		fatalf("error: %v\n", err)
	}

	store.MigrateIfNeeded(root, branch)
	store.WriteRepoPath(root)

	wsPath := store.WorkingSetPath(root, branch)

	cfg, err := config.Load()
	if err != nil {
		fatalf("error loading config: %v\n", err)
	}

	var stale []store.StaleCandidate
	if cfg.CleanupDays > 0 {
		stale = store.StaleCandidates(root, branch, cfg.CleanupDays)
	}

	m := tui.New(wsPath, root, cfg, stale)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fatalf("error: %v\n", err)
	}
}

func workingSetPath() (string, error) {
	root, err := git.RootDir()
	if err != nil {
		return "", err
	}
	branch, err := git.CurrentBranch()
	if err != nil {
		return "", err
	}
	store.MigrateIfNeeded(root, branch)
	store.WriteRepoPath(root)
	return store.WorkingSetPath(root, branch), nil
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
