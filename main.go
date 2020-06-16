package main

import (
	"ServerScan/cmd"
	"ServerScan/service"
)

func main() {

	scan := service.NewServerScan(service.WithHost(cmd.Host), service.WithMode(cmd.Mode), service.WithPort(cmd.Port), service.WithOutFile(cmd.OutFile))
	scan.PortScan()
}
