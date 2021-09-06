package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

type RBLResults struct {
	IP      net.IP    `json:"ip" yaml:"ip"`
	Results []*Result `json:"results" yaml:"results"`
}

type Result struct {
	Rbl    string `json:"rbl" yaml:"rbl"`
	Listed bool   `json:"listed" yaml:"listed"`
	Text   string `json:"text,omitempty" yaml:"text,omitempty"`
}

func rblQuery(ip net.IP, rbl string) (r *Result) {

	// Reverse IPv4 address
	splitAddress := strings.Split(ip.String(), ".")
	for i, j := 0, len(splitAddress)-1; i < len(splitAddress)/2; i, j = i+1, j-1 {
		splitAddress[i], splitAddress[j] = splitAddress[j], splitAddress[i]
	}

	lookup := fmt.Sprintf("%s.%s", strings.Join(splitAddress, "."), rbl)
	r = &Result{
		Listed: false,
		Rbl:    rbl,
	}
	regexpResponse, _ := regexp.Compile(`^127\.0\.0\.*`)
	res, _ := net.LookupHost(lookup)
	if len(res) > 0 {
		for i := range res {
			if regexpResponse.MatchString(res[i]) {
				r.Listed = true
				reportWriteLine(ip, rbl)
			}
		}
		txt, _ := net.LookupTXT(lookup)
		if len(txt) > 0 {
			r.Text = strings.Join(txt, "")
		}
	}
	return r
}

func rblLookup(rblList []string, ip net.IP) (res *RBLResults) {
	var wg sync.WaitGroup
	res = &RBLResults{
		IP: ip,
	}
	for i := range rblList {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			defer func() {
				wg.Done()
			}()
			r := rblQuery(ip, rblList[i])
			res.Results = append(res.Results, r)
		}(&wg, i)
	}
	wg.Wait()
	return res
}

func rblLookupList(rblList []string, ip []string) {
	for i := range ip {
		// Single IP
		if ip := net.ParseIP(ip[i]); ip != nil {
			r := rblLookup(rblList, ip)
			x, _ := yaml.Marshal(r)
			fmt.Println(string(x))
			continue
		}
		// CIDR
		if _, cidr, err := net.ParseCIDR(ip[i]); err == nil {
			mask := binary.BigEndian.Uint32(cidr.Mask)
			start := binary.BigEndian.Uint32(cidr.IP)
			finish := (start & mask) | (mask ^ 0xffffffff)
			for i := start; i <= finish; i++ {
				ip := make(net.IP, 4)
				binary.BigEndian.PutUint32(ip, i)
				r := rblLookup(rblList, ip)
				x, _ := yaml.Marshal(r)
				fmt.Println(string(x))
			}
			continue
		}
	}
}
