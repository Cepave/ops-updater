package cron

import (
	"gitcafe.com/ops/common/utils"
	"gitcafe.com/ops/updater/g"
	"log"
	"time"
)

func Heartbeat() {
	for {
		heartbeat()
		d := time.Duration(g.Config().Interval) * time.Second
		time.Sleep(d)
	}
}

func heartbeat() {
	agentDirs, err := ListAgentDirs()
	if err != nil {
		return
	}

	hostname, err := utils.Hostname(g.Config().Hostname)
	if err != nil {
		return
	}

	heartbeatRequest := BuildHeartbeatRequest(hostname, agentDirs)
	log.Println(heartbeatRequest)
}
