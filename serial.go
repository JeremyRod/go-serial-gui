package main

import (
	"go.bug.st/serial"
)

type SerialConf struct {
	Port     string
	BaudRate int
	Parity   serial.Parity
	DataBits int
	StopBits serial.StopBits
}

func GetPorts() ([]string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return []string{}, err
	}
	return ports, nil
}

func (c *SerialConf) OpenPort() (*serial.Port, error) {
	mode := serial.Mode{
		BaudRate: c.BaudRate,
		Parity:   c.Parity,
		DataBits: c.DataBits,
		StopBits: c.StopBits,
	}

	port, err := serial.Open(c.Port, &mode)
	if err != nil {
		return nil, err
	}
	return &port, nil
}

func (c *SerialConf) ClosePort(port *serial.Port) error {
	err := (*port).Close()
	if err != nil {
		return err
	}
	return nil
}
