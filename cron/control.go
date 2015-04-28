package cron

import (
	"log"
	"os/exec"
)

func Control(workdir, arg string) (string, error) {
	cmd := exec.Command("./control", arg)
	cmd.Dir = workdir
	bs, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cd %s; ./control %s fail %v. output: %s", workdir, arg, err, string(bs))
	}
	return string(bs), err
}

func ControlStatus(workdir string) (string, error) {
	return Control(workdir, "status")
}

func ControlStart(workdir string) (string, error) {
	return Control(workdir, "start")
}

func ControlStop(workdir string) (string, error) {
	return Control(workdir, "stop")
}
