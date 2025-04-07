package main

import (
	"log"
)

func main() {
	ports, err := GetPorts()
	if err != nil {
		log.Fatal(err)
	}
	CreateUI(ports)
}
