package cron

import (
	"os/exec"
)

func Control(workdir, arg string) (string, error) {
	cmd := exec.Command("./control", arg)
	cmd.Dir = workdir
	bs, err := cmd.CombinedOutput()
	return string(bs), err
}
