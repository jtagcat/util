package std

import (
	"context"
	"fmt"
	"os/exec"
)

func RunCmdWithCtx(ctx context.Context, cmd *exec.Cmd) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting command: %w", err)
	}

	wait := make(chan error)
	go func() {
		wait <- cmd.Wait()
		close(wait)
	}()

	select {
	case <-ctx.Done():
		_ = cmd.Cancel()
		return ctx.Err()
	case err := <-wait:
		return err
	}
}
