// Copyright 2016 Canonical Ltd.  All rights reserved.

package cmd

import (
	"github.com/juju/cmd"
	"github.com/juju/errors"
	"launchpad.net/gnuflag"
)

const suspendPlanDoc = `
suspend-plan is used to suspend plan for a set of charms
Example
suspend-plan foocorp/free cs:~foocorp/app-0 cs:~foocorp/app-1
	disables deploys of the two specified charms using the foocorp/free plan.
`

// suspendCommand suspends plan for a set of charms.
type suspendCommand struct {
	baseCommand

	PlanURL   string
	CharmURLs []string
	All       bool
}

// NewSuspendCommand creates a new suspendCommand.
func NewSuspendCommand() *suspendCommand {
	return &suspendCommand{}
}

// SetFlags implements Command.SetFlags.
func (c *suspendCommand) SetFlags(f *gnuflag.FlagSet) {
	c.baseCommand.ServiceURL = defaultServiceURL()
	c.baseCommand.SetFlags(f)
	f.BoolVar(&c.All, "all", false, "suspend plan for all charms")
}

// Info implements Command.Info.
func (c *suspendCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "suspend-plan",
		Args:    "<plan url> [<charm url>[...<charm url N>]]",
		Purpose: "suspends plan for specified charms",
		Doc:     suspendPlanDoc,
	}
}

// Init implements Command.Init.
func (c *suspendCommand) Init(args []string) error {
	if !c.All && len(args) < 2 {
		return errors.New("missing plan or charm url")
	} else if c.All && len(args) > 1 {
		return errors.New("cannot use --all and specify charm urls")
	}

	c.PlanURL, c.CharmURLs = args[0], args[1:]
	return nil
}

// Run implements Command.Run.
func (c *suspendCommand) Run(ctx *cmd.Context) error {
	defer c.Close()
	client, err := c.NewClient()
	if err != nil {
		return errors.Annotate(err, "failed to create an http client")
	}
	apiClient, err := newClient(c.ServiceURL, client)
	if err != nil {
		return errors.Annotate(err, "failed to create a plan API client")
	}
	return errors.Trace(apiClient.Suspend(c.PlanURL, c.All, c.CharmURLs...))
}
