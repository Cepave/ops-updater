package cron

import (
	"bytes"
	"fmt"
	"gitcafe.com/ops/common/model"
	"gitcafe.com/ops/updater/g"
	f "github.com/toolkits/file"
	"log"
	"os/exec"
	"path"
	"strings"
	"time"
)

func BuildHeartbeatRequest(hostname string, agentDirs []string) model.HeartbeatRequest {
	req := model.HeartbeatRequest{Hostname: hostname}

	realAgents := []*model.RealAgent{}
	now := time.Now().Unix()

	for _, agentDir := range agentDirs {
		// 如果目录下没有.version，我们认为这根本不是一个agent
		versionFile := path.Join(g.SelfDir, agentDir, ".version")
		if !f.IsExist(versionFile) {
			continue
		}

		version, err := f.ToTrimString(versionFile)
		if err != nil {
			log.Printf("read %s/.version fail: %v", agentDir, err)
			continue
		}

		controlFile := path.Join(g.SelfDir, agentDir, version, "control")
		if !f.IsExist(controlFile) {
			log.Printf("%s is nonexistent", controlFile)
			continue
		}

		status := ""

		cmd := exec.Command("./control", "status")
		cmd.Dir = path.Join(g.SelfDir, agentDir, version)

		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		err = cmd.Run()
		if err != nil {
			status = fmt.Sprintf("exec `./control status` fail: %s", err)
		} else {
			status = strings.TrimSpace(stdout.String())
		}

		realAgent := &model.RealAgent{
			Name:      agentDir,
			Version:   version,
			Status:    status,
			Timestamp: now,
		}

		realAgents = append(realAgents, realAgent)
	}

	req.RealAgents = realAgents
	return req
}

func ListAgentDirs() ([]string, error) {
	agentDirs, err := f.DirsUnder(g.SelfDir)
	if err != nil {
		log.Println("list dirs under", g.SelfDir, "fail", err)
	}
	return agentDirs, err
}
