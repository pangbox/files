package shell

import "os/exec"

func OpenFolder(folder string) error {
	return exec.Command("xdg-open", folder).Start()
}
