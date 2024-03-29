package main

import (
	"flag"
	"log"
	"oasisTracker/api"
	"oasisTracker/cli"
	"oasisTracker/common/genesis"
	"oasisTracker/common/modules"
	"oasisTracker/conf"
	"oasisTracker/dao"
	"oasisTracker/services"
	"oasisTracker/services/scanners"
	"os"
	"os/signal"
	"syscall"

	"github.com/roylee0704/gron"
	"go.uber.org/zap"
)

var parserDisableFlag = flag.Bool("parser-disable", false, "disable cron for api tests")

func main() {
	flag.Parse()
	configFile := flag.String("conf", "./config.json", "Path to config file")
	cfg, err := conf.NewFromFile(configFile)
	if err != nil {
		log.Fatal("can`t read config from file", zap.Error(err))
	}

	d, err := dao.New(cfg)
	if err != nil {
		log.Fatal("dao.New", zap.Error(err))
	}

	args := os.Args[1:]
	if len(args) > 0 && !*parserDisableFlag {
		cli := cli.NewCli(d)

		err = cli.Setup(args)
		if err != nil {
			log.Fatal("cli.SetupGenesisJson", zap.Error(err))
		}
		return
	}

	gen, err := genesis.ReadGenesisFile(genesis.DefaultGenesisFileName)
	if err != nil {
		log.Fatal("ReadGenesisFile", zap.Error(err))
	}

	s := services.NewService(cfg, d.GetServiceDAO(), gen.GenesisHeight)
	s.D = d

	a := api.NewAPI(cfg, s)
	//mds := []modules.Module{a}
	s.Modules = []modules.Module{a}
	cron := gron.New()

	//services.AddToCron(cron, cfg, d)
	s.AddToCron(cron, cfg, d)

	if !*parserDisableFlag {
		sm := scanners.NewManager(cfg, d)

		wt, err := scanners.NewWatcher(cfg, d)
		if err != nil {
			log.Fatal("Watcher.New", zap.Error(err))
		}
		s.Modules = append(s.Modules, sm, wt)
	}

	cron.Start()
	defer cron.Stop()

	//s.Modules = mds
	modules.Run(s.Modules)

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	<-gracefulStop
	modules.Stop(s.Modules)
}
