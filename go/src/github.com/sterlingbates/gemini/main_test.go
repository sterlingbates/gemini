package main

import (
	"testing"
	"time"
	"fmt"
)

/*
func TestArduinoConversation(t *testing.T) {
	go main()
	// Give some time for initialization
	time.Sleep(3 * time.Second)
	// Start test mode
	controller.cl.SendComPort("testmode")
	// Wait for a few interactions
	time.Sleep(10 * time.Second)
	// End test mode
	controller.cl.SendComPort("testmode")
}
*/

func TestOrbiterCommands(t *testing.T) {
	go main()
	time.Sleep(3 * time.Second)
	s := controller.mgr.Send("SHIP:FOCUS:PropMass")
	fmt.Println(s)
}
