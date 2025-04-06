package main

import (
	"log"
	"os"

	"gioui.org/app"
)

func main() {
	// ports, err := serial.GetPortsList()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, port := range ports {
	// 	fmt.Println(port)
	// }

	go func() {
		window := new(app.Window)
		err := CreateUI(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	go func() {
		window := new(app.Window)
		err := CreateUI(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
