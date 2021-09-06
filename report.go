package main

import (
	"fmt"
	"net"
	"os"
)

func reportWriteLine(ip net.IP, rbl string) {
	if len(*flagReport) == 0 {
		return
	}
	f, err := os.OpenFile(*flagReport, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("\"%s\",\"%s\"\n", ip.String(), rbl))
	if err != nil {
		logger.Error(err.Error())
		return
	}

}
