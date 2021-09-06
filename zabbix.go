package main

import (
	"fmt"

	"github.com/spetr/go-zabbix-sender"
)

var zabbixMetrics []*zabbix.Metric

func startZabbixSender() {
	z := zabbix.NewSender(conf.Zabbix.Server)
	resActive, errActive, resTrapper, errTrapper := z.SendMetrics(zabbixMetrics)

	fmt.Printf("Agent active, response=%s, info=%s, error=%v\n", resActive.Response, resActive.Info, errActive)
	fmt.Printf("Trapper, response=%s, info=%s,error=%v\n", resTrapper.Response, resTrapper.Info, errTrapper)
}
