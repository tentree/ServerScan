package service

import (
	"ServerScan/pkg/icmpcheck"
	"ServerScan/pkg/portscan"
	"ServerScan/pkg/vscan"
	"fmt"
	"log"

	"github.com/malfunkt/iprange"
)

type ServerScan struct {
	hostLists []string
	mode      string
	outFile   string
	port      string
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
		s.port = port
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
		aliveHosts = icmpcheck.ICMPRun(s.hostLists)
		for _, host := range aliveHosts {
			fmt.Printf("(ICMP) Target '%s' is alive\n", host)
		}
		aliveHosts, aliveAddress = portscan.TCPportScan(aliveHosts, s.port, s.mode, s.timeout)
	case "tcp":
		aliveHosts, aliveAddress = portscan.TCPportScan(s.hostLists, s.port, s.mode, s.timeout)
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

func (s *ServerScan) TcpPortScan() {

}

func (s *ServerScan) IcmpPortScan() {

}
