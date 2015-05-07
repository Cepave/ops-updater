package cron

import (
	"gitcafe.com/ops/common/model"
	"gitcafe.com/ops/updater/g"
	"log"
)

func HandleHeartbeatResponse(respone *model.HeartbeatResponse) {
	if respone.ErrorMessage != "" {
		log.Println("receive error message:", respone.ErrorMessage)
		return
	}

	das := respone.DesiredAgents
	if das == nil || len(das) == 0 {
		return
	}

	for _, da := range das {
		da.FillAttrs(g.SelfDir)
		HandleDesiredAgent(da)
	}
}

func HandleDesiredAgent(da *model.DesiredAgent) {
	if da.Cmd == "start" {
		StartDesiredAgent(da)
	} else if da.Cmd == "stop" {
		StopDesiredAgent(da)
	} else if da.Cmd == "none" || da.Cmd == "nil" || da.Cmd == "" {
		// do nothing
	} else {
		log.Println("unknown cmd", da)
	}
}
