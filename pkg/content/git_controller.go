package content

import (
	"context"
	"github.com/westcoastcode-se/gocms/pkg/event"
	"os/exec"
)

type GitController struct {
	bus      *event.Bus
	RootPath string
}

func (g *GitController) Update(ctx context.Context, commit string) error {
	err := g.Pull()
	if err != nil {
		return err
	}

	return g.Checkout(ctx, commit)
}

func (g *GitController) Save(ctx context.Context, message string) error {
	err := g.Commit(message)
	if err != nil {
		return err
	}

	return g.Push(ctx)
}

func (g *GitController) Pull() error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = g.RootPath
	return cmd.Run()
}

// Checkout the supplied commit and notify all listeners that
func (g *GitController) Checkout(ctx context.Context, commit string) error {
	cmd := exec.Command("git", "checkout", commit)
	cmd.Dir = g.RootPath
	if err := cmd.Run(); err != nil {
		return err
	}

	// NotifyAll next event that a checkout has happened
	if err := g.bus.NotifyAll(ctx, &event.Checkout{Commit: commit}); err != nil {
		return err
	}

	return nil
}

func (g *GitController) Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = g.RootPath
	return cmd.Run()
}

func (g *GitController) Push(ctx context.Context) error {
	cmd := exec.Command("git", "push")
	cmd.Dir = g.RootPath
	if err := cmd.Run(); err != nil {
		return err
	}

	// NotifyAll next event that a checkout has happened
	if err := g.bus.NotifyAll(ctx, &event.Push{}); err != nil {
		return err
	}

	return nil
}

func NewGitController(bus *event.Bus, rootPath string) *GitController {
	return &GitController{bus: bus, RootPath: rootPath}
}
