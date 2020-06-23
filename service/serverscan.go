package service

import (
	"ServerScan/pkg/icmpcheck"
	"ServerScan/pkg/vscan"
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/malfunkt/iprange"
)

type ServerScan struct {
	hostLists []string
	mode      string
	outFile   string
	ports     []int
	timeout   int
}

func NewServerScan(params ...ParamOption) *ServerScan {
	s := &ServerScan{}

	for _, param := range params {
		param(s)
	}

	return s
}

type ParamOption func(*ServerScan)

func WithHost(host string) ParamOption {
	return func(s *ServerScan) {
		hostList, err := iprange.ParseList(host)
		if err != nil {
			log.Panic(err)
		}

		for _, v := range hostList.Expand() {
			s.hostLists = append(s.hostLists, v.String())
		}
	}
}

func WithPort(port string) ParamOption {
	return func(s *ServerScan) {
		s.ports = parsePort(port)
	}
}

func WithMode(mode string) ParamOption {
	return func(s *ServerScan) {
		s.mode = mode
	}
}

func WithTimeout(timeout int) ParamOption {
	return func(s *ServerScan) {
		s.timeout = timeout
	}
}

func WithOutFile(outFile string) ParamOption {
	return func(s *ServerScan) {
		s.outFile = outFile
	}
}

func (s *ServerScan) PortScan() {
	var aliveHosts []string
	var aliveAddress []string

	switch s.mode {
	case "icmp":
		s.hostLists = icmpcheck.ICMPRun(s.hostLists)
		for _, host := range aliveHosts {
			fmt.Printf("(ICMP) Target '%s' is alive\n", host)
		}
		aliveHosts, aliveAddress = s.TcpPortScan()
	case "tcp":
		aliveHosts, aliveAddress = s.TcpPortScan()
		for _, host := range aliveHosts {
			fmt.Printf("(TCP) Target '%s' is alive\n", host)
		}
		for _, addr := range aliveAddress {
			fmt.Println(addr)
		}
	}

	if len(aliveAddress) > 0 {
		vscan.GetProbes(aliveAddress)
	}
}

func (s *ServerScan) TcpPortScan() ([]string, []string) {
	var aliveHosts []string
	var aliveAddress []string

	for _, host := range s.hostLists {
		address := ProbeHosts(context.Background(), host, s.ports)
		if len(address) > 0 {
			aliveHosts = append(aliveHosts, host)
		}

		aliveAddress = append(aliveAddress, address...)
	}

	return aliveHosts, aliveAddress
}

func (s *ServerScan) IcmpPortScan() ([]string, []string) {
	var aliveHosts []string
	var aliveAddress []string

	for _, host := range s.hostLists {
		address := ProbeHosts(context.Background(), host, s.ports)
		if len(address) > 0 {
			aliveHosts = append(aliveHosts, host)
		}

		aliveAddress = append(aliveAddress, address...)
	}

	return aliveHosts, aliveAddress
}

func parsePort(ports string) []int {
	var scanPorts []int
	slices := strings.Split(ports, ",")

	for _, port := range slices {
		port = strings.Trim(port, " ")
		upper := port
		if strings.Contains(port, "-") {
			ranges := strings.Split(port, "-")
			if len(ranges) < 2 {
				continue
			}
			sort.Strings(ranges)
			port = ranges[0]
			upper = ranges[1]
		}
		start, _ := strconv.Atoi(port)
		end, _ := strconv.Atoi(upper)
		for i := start; i <= end; i++ {
			scanPorts = append(scanPorts, i)
		}
	}

	return scanPorts
}
