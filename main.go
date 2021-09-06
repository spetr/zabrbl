package main

import (
	"flag"
	"log"
	"os"

	"github.com/kardianos/service"
)

type program struct {
	exit chan struct{}
}

var (
	logger      service.Logger
	flagService *string
	flagReport  *string
)

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})
	ConfLoad()
	go func() {
		rblLookupList(conf.RBL.IPv4, conf.IP)
		if len(*flagReport) != 0 {
			os.Exit(1)
		}
	}()
	return nil
}

func (p *program) Stop(s service.Service) error {
	logger.Info("Stopping")
	close(p.exit)
	return nil
}

func init() {
	flagService = flag.String("service", "", "Control the system service.")
	flagReport = flag.String("report", "", "Report file (CSV format).")
	flag.Parse()
}

func main() {
	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"
	svcConfig := &service.Config{
		Name:        "ZabRBL",
		DisplayName: "ZabRBL monitoring service",
		Description: "RBL monitoring service for Zabbix.",
		Dependencies: []string{
			"Requires=network.target",
			"After=network-online.target syslog.target"},
		Option: options,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*flagService) != 0 {
		err := service.Control(s, *flagService)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
