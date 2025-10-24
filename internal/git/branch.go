package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

// Branch returns the current branch name of the repository.
func Branch(ctx context.Context, dirPath string) (string, error) {
	// ensure git is installed and available to run
	gitFilePath, lookErr := binPath()
	if lookErr != nil {
		return "", lookErr
	}

	// get the current branch name
	var cmd = exec.CommandContext(ctx, gitFilePath, "rev-parse",
		"--abbrev-ref",
		"HEAD",
	)

	cmd.Dir = dirPath
	cmd.Env = []string{
		"LC_ALL=C", "LANG=C", // forces the system to use the "C" (POSIX) locale, English-based output with no localization
		"NO_COLOR=1",            // disables colored output
		"GIT_CONFIG_NOSYSTEM=1", // do not use the system-wide configuration file
	}

	var stdOut, stdErr bytes.Buffer

	stdOut.Grow(256) //nolint:mnd // 256 bytes should be enough for branch name

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		if stdErr.Len() > 0 {
			err = fmt.Errorf("%s: %w", stdErrToString(stdErr.String()), err)
		}

		return "", fmt.Errorf("git rev-parse failed: %w", err)
	}

	return stdOut.String(), nil
}
