package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var mgr probeMgr

type probeMgr struct {
	mu              sync.Mutex
	unConfirmedHost map[string]*probeRequest
}

func init() {
	mgr.unConfirmedHost = make(map[string]*probeRequest)
}

type probeRequest struct {
	sync.RWMutex
	unConfirmPortNum int
	confirmAddress   []string      // 确认后的数据
	request          chan []string // 确认通道
}

func addProbeRequest(ctx context.Context, host string, ports []int) chan []string {

	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	_, ok := mgr.unConfirmedHost[host]
	if ok {
		log.Println("dumplicate host")
		return nil
	}

	requestChan := make(chan []string, 1)
	mgr.unConfirmedHost[host] = &probeRequest{
		unConfirmPortNum: len(ports),
		confirmAddress:   nil,
		request:          requestChan,
	}

	go func(host string) {
		// no deadline
		if _, ok := ctx.Deadline(); !ok {
			return
		}
		for {
			select {
			case <-ctx.Done():
				getResultNow(host)
			}
		}
	}(host)

	return requestChan
}

func recvResult(address string, result bool) {
	addrSplit := strings.Split(address, ":")
	if len(addrSplit) != 2 {
		return
	}
	host := addrSplit[0]

	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	value, ok := mgr.unConfirmedHost[host]
	if !ok {
		log.Println("cant find host")
		return
	}

	value.Lock()
	value.unConfirmPortNum--
	if result {
		if value.confirmAddress == nil {
			value.confirmAddress = append([]string{}, address)
		} else {
			value.confirmAddress = append(value.confirmAddress, address)
		}
	}

	if value.unConfirmPortNum == 0 {
		value.request <- value.confirmAddress
		delete(mgr.unConfirmedHost, host)
	}
	value.Unlock()
}

func getResultNow(host string) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	value, ok := mgr.unConfirmedHost[host]
	if !ok {
		// aleray exit
		return
	}

	value.Lock()
	value.request <- value.confirmAddress
	delete(mgr.unConfirmedHost, host)
	value.Unlock()

}

func ProbeHosts(ctx context.Context, host string, ports []int) []string {
	request := addProbeRequest(ctx, host, ports)

	dail := &net.Dialer{}
	for _, port := range ports {
		address := fmt.Sprintf("%s:%d", host, port)
		go func(address string) {
			dailCtx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			conn, err := dail.DialContext(dailCtx, "tcp4", address)
			if err != nil {
				log.Println(err)
				recvResult(address, false)
				return
			}
			conn.Close()

			recvResult(address, true)
		}(address)
	}

	for {
		select {
		case address, ok := <-request:
			if !ok {
				return nil
			}
			return address
		}
	}
}
