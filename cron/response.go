package cron

import (
	"github.com/Cepave/ops-common/model"
	"github.com/Cepave/ops-updater/g"
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

		if g.Config().DesiredAgent == "" || g.Config().DesiredAgent == da.Name {
			HandleDesiredAgent(da)
		}
	}
}

func HandleDesiredAgent(da *model.DesiredAgent) {
	if da.Cmd == "start" {
		StartDesiredAgent(da)
	} else if da.Cmd == "stop" {
		StopDesiredAgent(da)
	} else {
		log.Println("unknown cmd", da)
	}
}
