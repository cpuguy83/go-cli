package cli

import (
	"context"
	"testing"
)

func TestCommand(t *testing.T) {
	ctx := context.Background()

	var ran bool
	cmd := NewCmd("foo", func(ctx context.Context) error {
		ran = true
		return nil
	})
	if cmd.Name() != "foo" {
		t.Errorf("expected command name to be 'command_test', got %q", cmd.Name())
	}

	if err := cmd.Run(ctx, nil); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !ran {
		t.Error("expected command to run")
	}

	if err := cmd.Run(ctx, []string{"bar"}); err != ErrNoSuchCommand {
		t.Errorf("expected ErrNoSuchCommand, got %v", err)
	}
	cmd.NewCmd("bar", nil)
	if err := cmd.Run(ctx, []string{"bar"}); err != ErrNoSuchCommand {
		t.Errorf("expected ErrNoSuchCommand, got %v", err)
	}

	ran = false
	cmd.NewCmd("bar", func(ctx context.Context) error {
		ran = true
		return nil
	})
	if err := cmd.Run(ctx, []string{"bar"}); err != nil {
		t.Fatal(err)
	}
	if !ran {
		t.Error("expected command to run")
	}
}
