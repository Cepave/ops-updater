package cron

import (
	"encoding/json"
	"fmt"
	"github.com/Cepave/ops-common/model"
	"github.com/Cepave/ops-common/utils"
	"github.com/Cepave/ops-updater/g"
	"github.com/toolkits/net/httplib"
	"log"
	"time"
)

func Heartbeat() {
	SleepRandomDuration()
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
	if g.Config().Debug {
		log.Println("====>>>>")
		log.Println(heartbeatRequest)
	}

	bs, err := json.Marshal(heartbeatRequest)
	if err != nil {
		log.Println("encode heartbeat request fail", err)
		return
	}

	url := fmt.Sprintf("https://%s/heartbeat", g.Config().Server)
	httpRequest := httplib.Post(url).SetTimeout(time.Second*10, time.Minute)
	httpRequest.Body(bs)
	httpResponse, err := httpRequest.Bytes()
	if err != nil {
		log.Printf("curl %s fail %v", url, err)
		return
	}

	var heartbeatResponse model.HeartbeatResponse
	err = json.Unmarshal(httpResponse, &heartbeatResponse)
	if err != nil {
		log.Println("decode heartbeat response fail", err)
		return
	}

	if g.Config().Debug {
		log.Println("<<<<====")
		log.Println(heartbeatResponse)
	}

	HandleHeartbeatResponse(&heartbeatResponse)

}
