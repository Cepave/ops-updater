package cron

import (
	"fmt"
	"gitcafe.com/ops/common/model"
	"gitcafe.com/ops/common/utils"
	"gitcafe.com/ops/updater/g"
	"github.com/toolkits/file"
	"log"
	"os/exec"
	"path"
	"strings"
	"time"
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
	} else {
		log.Println("unknown cmd", da)
	}
}

func StartDesiredAgent(da *model.DesiredAgent) {
	err := InsureAgentDirExists(da)
	if err != nil {
		return
	}

	err = DownloadNewVersion(da)
	if err != nil {
		log.Println("download tarball and md5 fail", err, da)
		return
	}

	err = StopOldVersion(da)
	if err != nil {
		log.Println("stop", da.Name, "fail", err)
		return
	}

	err = StartNewVersion(da)
	if err != nil {
		log.Println("start new version", da.Name, "fail", err)
		return
	}

	file.WriteString(path.Join(da.AgentDir, ".version"), da.Version)
}

func StartNewVersion(da *model.DesiredAgent) error {
	out, err := Control(da.AgentVersionDir, "status")
	if err != nil {
		log.Printf("cd %s; ./control status fail: %v, stdout: %s", da.AgentVersionDir, err, out)
		return err
	}

	if strings.Contains(out, "started") {
		return nil
	}

	out, err = Control(da.AgentVersionDir, "start")
	if err != nil {
		log.Printf("cd %s; ./control start fail: %v, stdout: %v", da.AgentVersionDir, err, out)
		return err
	}

	time.Sleep(time.Second)

	out, err = Control(da.AgentVersionDir, "status")
	if err != nil {
		log.Printf("cd %s; ./control status fail: %v, stdout: %s", da.AgentVersionDir, err, out)
		return err
	}

	if strings.Contains(out, "started") {
		return nil
	}

	return fmt.Errorf("cd %s; ./control start fail", da.AgentVersionDir)
}

func DownloadNewVersion(da *model.DesiredAgent) error {
	if FilesReady(da) {
		return nil
	}

	downloadTarballCmd := exec.Command("wget", "-q", da.TarballUrl, "-O", da.TarballFilename)
	downloadTarballCmd.Dir = da.AgentVersionDir
	err := downloadTarballCmd.Run()
	if err != nil {
		log.Println("wget -q", da.TarballUrl, "-O", da.TarballFilename, "fail")
		return err
	}

	downloadMd5Cmd := exec.Command("wget", "-q", da.Md5Url, "-O", da.Md5Filename)
	downloadMd5Cmd.Dir = da.AgentVersionDir
	err = downloadMd5Cmd.Run()
	if err != nil {
		log.Println("wget -q", da.Md5Url, "-O", da.Md5Filename, "fail")
		return err
	}

	if utils.Md5sumCheck(da.AgentVersionDir, da.Md5Filename) {
		return nil
	} else {
		return fmt.Errorf("md5sum -c fail")
	}
}

func FilesReady(da *model.DesiredAgent) bool {
	if !file.IsExist(da.Md5Filepath) {
		return false
	}

	if !file.IsExist(da.TarballFilepath) {
		return false
	}

	if !file.IsExist(da.ControlFilepath) {
		return false
	}

	return utils.Md5sumCheck(da.AgentVersionDir, da.Md5Filename)
}

func StopOldVersion(da *model.DesiredAgent) error {
	versionFile := path.Join(da.AgentDir, ".version")
	if !file.IsExist(versionFile) {
		log.Printf("WARN: %s is nonexistent. no need stop old agent.", versionFile)
		return nil
	}

	version, err := file.ToTrimString(versionFile)
	if err != nil {
		log.Printf("WARN: read %s fail. no need stop old agent.", version)
		return nil
	}

	oldVersionDir := path.Join(da.AgentDir, version)
	if !file.IsExist(oldVersionDir) {
		log.Printf("WARN: %s nonexistent. no need to stop old version", oldVersionDir)
		return nil
	}

	out, err := Control(oldVersionDir, "status")
	if err != nil {
		log.Println("./control status fail", err, "in", oldVersionDir)
		return err
	}

	if strings.Contains(out, "stoped") {
		// 已经是stop状态
		return nil
	}

	_, err = Control(oldVersionDir, "stop")
	if err != nil {
		log.Println("./control stop fail", err, "in", oldVersionDir)
		return err
	}

	// 杀死一个进程可能比较费劲，等两秒再检查一下
	time.Sleep(time.Second * 2)

	out, err = Control(oldVersionDir, "status")
	if err != nil {
		log.Println("./control status fail", err, "in", oldVersionDir)
		return err
	}

	if strings.Contains(out, "stoped") {
		// 已经是stop状态
		return nil
	}

	return fmt.Errorf("cannot stop %v", da.Name)
}

func StopDesiredAgent(da *model.DesiredAgent) {

}

func InsureAgentDirExists(da *model.DesiredAgent) error {
	err := file.InsureDir(da.AgentDir)
	if err != nil {
		return err
	}

	return file.InsureDir(da.AgentVersionDir)
}
