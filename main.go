package main

import (
	"encoding/binary"
	"fmt"
	"net"

	"gopkg.in/yaml.v2"
)

func main() {
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

}
