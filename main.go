package main

import (
	"flag"
	"fmt"
	"github.com/Cepave/ops-updater/cron"
	"github.com/Cepave/ops-updater/g"
	"github.com/Cepave/ops-updater/http"
	"github.com/toolkits/sys"
	"log"
	"os"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if err := g.ParseConfig(*cfg); err != nil {
		log.Fatalln(err)
	}

	g.InitGlobalVariables()

	CheckDependency()

	go http.Start()
	go cron.Heartbeat()

	select {}
}

func CheckDependency() {
	_, err := sys.CmdOut("wget", "--help")
	if err != nil {
		log.Fatalln("dependency wget not found")
	}

	_, err = sys.CmdOut("md5sum", "--help")
	if err != nil {
		log.Fatalln("dependency md5sum not found")
	}

	_, err = sys.CmdOut("tar", "--help")
	if err != nil {
		log.Fatalln("dependency tar not found")
	}
}
