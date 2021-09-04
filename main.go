package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/kardianos/service"
	"gopkg.in/yaml.v2"
)

type program struct {
	exit chan struct{}
}

var (
	logger service.Logger
)

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})
	go func() {
		ConfLoad()
		for i := range conf.IP {
			// Single IP
			if ip := net.ParseIP(conf.IP[i]); ip != nil {
				r := rblLookup(conf.RBL.IPv4, ip)
				x, _ := yaml.Marshal(r)
				fmt.Println(string(x))
				continue
			}
			// CIDR
			if _, cidr, err := net.ParseCIDR(conf.IP[i]); err == nil {
				mask := binary.BigEndian.Uint32(cidr.Mask)
				start := binary.BigEndian.Uint32(cidr.IP)
				finish := (start & mask) | (mask ^ 0xffffffff)
				for i := start; i <= finish; i++ {
					ip := make(net.IP, 4)
					binary.BigEndian.PutUint32(ip, i)
					r := rblLookup(conf.RBL.IPv4, ip)
					x, _ := yaml.Marshal(r)
					fmt.Println(string(x))
				}
				continue
			}
		}
	}()
	return nil
}

func (p *program) Stop(s service.Service) error {
	logger.Info("Stopping")
	close(p.exit)
	return nil
}

func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

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

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
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
