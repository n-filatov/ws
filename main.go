package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/n-filatov/ws/internal/config"
	"github.com/n-filatov/ws/internal/git"
	"github.com/n-filatov/ws/internal/store"
	"github.com/n-filatov/ws/internal/tui"
)

var version = "dev"

func main() {
	args := os.Args[1:]

	if len(args) == 0 || args[0] == "open" {
		runTUI()
		return
	}

	switch args[0] {
	case "version":
		fmt.Println("ws " + version)

	case "add":
		if len(args) < 2 {
			fatalf("usage: ws add <file> [file...]\n")
		}
		ws, err := initWorkspace()
		if err != nil {
			fatalf("error: %v\n", err)
		}
		if err := store.Add(ws.wsPath, args[1:]...); err != nil {
			fatalf("error: %v\n", err)
		}
		fmt.Printf("added %d file(s)\n", len(args)-1)

	case "rm":
		if len(args) < 2 {
			fatalf("usage: ws rm <file>\n")
		}
		ws, err := initWorkspace()
		if err != nil {
			fatalf("error: %v\n", err)
		}
		if err := store.Remove(ws.wsPath, args[1]); err != nil {
			fatalf("error: %v\n", err)
		}
		fmt.Printf("removed %s\n", args[1])

	case "list":
		ws, err := initWorkspace()
		if err != nil {
			fatalf("error: %v\n", err)
		}
		files, err := store.Load(ws.wsPath)
		if err != nil {
			fatalf("error: %v\n", err)
		}
		for _, f := range files {
			fmt.Println(f)
		}

	case "clear":
		ws, err := initWorkspace()
		if err != nil {
			fatalf("error: %v\n", err)
		}
		if err := store.Save(ws.wsPath, []string{}); err != nil {
			fatalf("error: %v\n", err)
		}
		fmt.Println("working set cleared")

	default:
		fatalf("unknown command %q\n\nUsage:\n  ws                    open TUI\n  ws add <file>...      add files\n  ws rm <file>          remove a file\n  ws list               list tracked files\n  ws clear              clear working set\n  ws version            print version\n", args[0])
	}
}

// workspace holds the resolved git context for the current invocation.
type workspace struct {
	root   string
	branch string
	wsPath string
}

// initWorkspace resolves the git root and branch, runs any pending migrations,
// records the repo path for future gc, and returns the resolved context.
func initWorkspace() (workspace, error) {
	root, err := git.RootDir()
	if err != nil {
		return workspace{}, err
	}
	branch, err := git.CurrentBranch()
	if err != nil {
		return workspace{}, err
	}
	store.MigrateIfNeeded(root, branch)
	store.WriteRepoPath(root)
	return workspace{
		root:   root,
		branch: branch,
		wsPath: store.WorkingSetPath(root, branch),
	}, nil
}

func runTUI() {
	ws, err := initWorkspace()
	if err != nil {
		fatalf("error: %v\n", err)
	}

	cfg, err := config.Load()
	if err != nil {
		fatalf("error loading config: %v\n", err)
	}

	var stale []store.StaleCandidate
	if cfg.CleanupDays > 0 {
		stale = store.StaleCandidates(ws.root, ws.branch, cfg.CleanupDays)
	}

	m := tui.New(ws.wsPath, ws.root, cfg, stale)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fatalf("error: %v\n", err)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
