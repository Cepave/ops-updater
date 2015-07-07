package cron

import (
	"github.com/Cepave/ops-common/model"
	"github.com/Cepave/ops-updater/g"
	"github.com/toolkits/file"
	"log"
	"path"
	"strings"
	"time"
)

func StopDesiredAgent(da *model.DesiredAgent) {
	if !file.IsExist(da.ControlFilepath) {
		return
	}

	ControlStopIn(da.AgentVersionDir)
}

func StopAgentOf(agentName, newVersion string) error {
	agentDir := path.Join(g.SelfDir, agentName)
	versionFile := path.Join(agentDir, ".version")

	if !file.IsExist(versionFile) {
		log.Printf("WARN: %s is nonexistent", versionFile)
		return nil
	}

	version, err := file.ToTrimString(versionFile)
	if err != nil {
		log.Printf("WARN: read %s fail %s", version, err)
		return nil
	}

	if version == newVersion {
		// do nothing
		return nil
	}

	versionDir := path.Join(agentDir, version)
	if !file.IsExist(versionDir) {
		log.Printf("WARN: %s nonexistent", versionDir)
		return nil
	}

	return ControlStopIn(versionDir)
}

func ControlStopIn(workdir string) error {
	if !file.IsExist(workdir) {
		return nil
	}

	out, err := ControlStatus(workdir)
	if err == nil && strings.Contains(out, "stoped") {
		return nil
	}

	_, err = ControlStop(workdir)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 3)

	out, err = ControlStatus(workdir)
	if err == nil && strings.Contains(out, "stoped") {
		return nil
	}

	return err
}
