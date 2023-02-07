package cli

import (
	"context"
	"errors"
	"flag"
)

var (
	ErrNoSuchCommand = errors.New("no such command")
)

type FlagSet interface {
	Parse([]string) error
	NArg() int
	Args() []string
	Arg(int) string
}

type FlagSetCreator[F FlagSet] func(name string) F

type Cmd[F FlagSet] struct {
	flags    F
	subs     map[string]*Cmd[F]
	h        func(context.Context) error
	newFlags FlagSetCreator[F]
	name     string
}

// NewCmd creates a new command with the given name and handler using the built-in flag.FlagSet type.
func NewCmd(name string, h func(context.Context) error) *Cmd[*flag.FlagSet] {
	f := func(name string) *flag.FlagSet {
		return flag.NewFlagSet(name, flag.ContinueOnError)
	}
	return NewCmdWithFlagSet(name, h, f)
}

// NewCmdWithFlagSet is a generic version of NewCmd that allows you to use any type that implements the FlagSet interface.
func NewCmdWithFlagSet[F FlagSet](name string, h func(context.Context) error, newFlags FlagSetCreator[F]) *Cmd[F] {
	return &Cmd[F]{flags: newFlags(name), h: h, newFlags: newFlags, name: name}
}

// Name returns the name of the command.
func (c *Cmd[F]) Name() string {
	return c.name
}

// Flags returns the flagset of the command.
func (c *Cmd[F]) Flags() F {
	return c.flags
}

// Run runs the command with the given args.
// If the command has no subcommands and no handler, it will print the usage and return ErrNoSuchCommand.
func (c *Cmd[F]) Run(ctx context.Context, args []string) error {
	if err := c.flags.Parse(args); err != nil {
		return err
	}

	if c.flags.NArg() == 0 {
		if c.h == nil {
			return ErrNoSuchCommand
		}
		return c.h(ctx)
	}

	sub, ok := c.subs[c.flags.Arg(0)]
	if !ok {
		return ErrNoSuchCommand
	}

	return sub.Run(ctx, c.flags.Args()[1:])
}

// NewCmd creates a new subcommand with the given name and handler.
func (c *Cmd[F]) NewCmd(name string, h func(context.Context) error) *Cmd[F] {
	if c.subs == nil {
		c.subs = make(map[string]*Cmd[F])
	}
	sub := NewCmdWithFlagSet(name, h, c.newFlags)
	c.subs[name] = sub
	return sub
}

// AddCmd adds the given command as a subcommand.
func (c *Cmd[F]) AddCmd(cmd *Cmd[F]) {
	if c.subs == nil {
		c.subs = make(map[string]*Cmd[F])
	}
	c.subs[cmd.Name()] = cmd
}

func (c *Cmd[F]) Commands() []*Cmd[F] {
	cmds := make([]*Cmd[F], 0, len(c.subs))
	for _, cmd := range c.subs {
		cmds = append(cmds, cmd)
	}
	return cmds
}

func (c *Cmd[F]) Cmd(name string) *Cmd[F] {
	return c.subs[name]
}
