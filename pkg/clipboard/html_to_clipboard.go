package clipboard

import (
	"fmt"
	"os/exec"
	"strings"
)

func CopyHTMLToClipboardAsRTF(html string) error {
	pandoc := exec.Command("pandoc", "-f", "html", "-t", "rtf", "-s")
	pandoc.Stdin = strings.NewReader(html)
	pandoc.Stderr = &strings.Builder{}
	stdin, err := pandoc.StdoutPipe()
	if err != nil {
		return err
	}
	pbcopy := exec.Command("pbcopy")
	pbcopy.Stdin = stdin
	pbcopy.Stderr = &strings.Builder{}

	if err := pandoc.Start(); err != nil {
		return err
	}
	if err := pbcopy.Start(); err != nil {
		return err
	}

	if err := pandoc.Wait(); err != nil {
		return reportError(pandoc, err)
	}
	if err := pbcopy.Wait(); err != nil {
		return reportError(pbcopy, err)
	}

	if err := reportExecError(pandoc); err != nil {
		return err
	}
	if err := reportExecError(pbcopy); err != nil {
		return err
	}

	return nil
}

func reportExecError(cmd *exec.Cmd) error {
	if cmd.ProcessState.ExitCode() != 0 {
		serr := cmd.Stderr.(*strings.Builder)
		return fmt.Errorf("%s exit code %d, stderr=%s", cmd.Path, cmd.ProcessState.ExitCode(), serr.String())
	}
	return nil
}

func reportError(cmd *exec.Cmd, err error) error {
	serr := cmd.Stderr.(*strings.Builder)
	return fmt.Errorf("%s %#v, exit code %d, stderr=%s", cmd.Path, err, cmd.ProcessState.ExitCode(), serr.String())
}
