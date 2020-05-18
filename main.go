package main

import (
	"flag"
	"go.uber.org/zap"
	"log"
	"oasisTracker/api"
	"oasisTracker/common/modules"
	"oasisTracker/conf"
	"oasisTracker/dao"
	"oasisTracker/services"
	"oasisTracker/services/scanners"
	"os"
	"os/signal"
	"syscall"
)

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

	s := services.NewService(cfg, d)

	a := api.NewAPI(cfg, s)

	sm := scanners.NewManager(cfg, d)

	wt, err := scanners.NewWatcher(cfg, d)
	if err != nil {
		log.Fatal("Watcher.New", zap.Error(err))
	}

	mds := []modules.Module{wt, sm, a}

	modules.Run(mds)

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	<-gracefulStop
	modules.Stop(mds)
}
