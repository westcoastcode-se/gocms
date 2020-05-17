package content

import (
	"github.com/westcoastcode-se/gocms/event"
	"os/exec"
)

type Controller struct {
	bus      *event.Bus
	RootPath string
}

func (g *Controller) Pull() error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = g.RootPath
	return cmd.Run()
}

// Checkout the supplied commit and notify all listeners that
func (g *Controller) Checkout(commit string) error {
	cmd := exec.Command("git", "checkout", commit)
	cmd.Dir = g.RootPath
	if err := cmd.Run(); err != nil {
		return err
	}

	// NotifyAll next event that a checkout has happened
	if err := g.bus.NotifyAll(&event.Checkout{Commit: commit}); err != nil {
		return err
	}

	return nil
}

func (g *Controller) Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = g.RootPath
	return cmd.Run()
}

func (g *Controller) Push() error {
	cmd := exec.Command("git", "push")
	cmd.Dir = g.RootPath
	if err := cmd.Run(); err != nil {
		return err
	}

	// NotifyAll next event that a checkout has happened
	if err := g.bus.NotifyAll(&event.Push{}); err != nil {
		return err
	}

	return nil
}

func NewController(bus *event.Bus, rootPath string) *Controller {
	return &Controller{bus: bus, RootPath: rootPath}
}
