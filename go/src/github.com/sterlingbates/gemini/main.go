package main

import (
	"fmt"
	"os"
	"os/signal"
	
	"github.com/marcturner/gemini/lib"
)

var controller *manager.Controller

func handleCtrlC() () {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Exit signal received")
		controller.SetExit()
	}()
}

func main() {
	fmt.Println("-> main")

	mgr := manager.NewOrbiterManager()
	controller = manager.NewController()
	controller.SetManager(mgr)

	// Set up a handler for Ctrl-C to clean things up properly
	handleCtrlC()
	// Run the controller loop
	controller.Run()
	fmt.Println("<- main")
}
