package manager

import (
	"fmt"
	"time"
	"strings"
)

type Controller struct {
	exit bool
	mgr OrbiterManager
}

var commands = [11]string{
		"SHIP:FOCUS:OrbSpd",				// Orbital speed of the vehicle
		"SHIP:FOCUS:Alt",					// Altitude (dbl)
		"SHIP:FOCUS:Accel",					// Acceleration in m/s^2 along air speed vector (dbl)
		"SHIP:FOCUS:PropMass",				// Propellant mass (dbl)
		"SHIP:FOCUS:ShipAirspdVector",		// Air speed vector (vector)
		"SHIP:FOCUS:PropFlowRate",			// Propellant flow rate (dbl)
		"SHIP:FOCUS:AttitudeMode",			// Attitude mode (int)												// ERR02
		"SHIP:FOCUS:FltStatus",				// Flight status (int)
		"SHIP:FOCUS:Airspd",				// Air speed (dbl)													// ERR12
		"SHIP:FOCUS:VAccel",				// Vertical acceleration in m/s^2
		"SHIP:FOCUS:IndSpd",				// Indicated airspeed based on atmo conditions and flight regime	// ERR12
	}

func NewController() *Controller {
    ret := Controller {
		exit: false,
	}
	return &ret
}

func (c *Controller) SetManager(mgr *OrbiterManager) () {
	c.mgr = *mgr
}

func (c *Controller) SetExit() () {
	c.exit = true
}

/**
 * Pulls data from the Orbiter and sends it to Gemini.
 */
func (c *Controller) pullAndSend() () {
	fmt.Println("-> pullAndSend")
	// Pull down the data from Orbiter
	for _, cmd := range commands {
		message := c.mgr.Send(cmd) + "|"
		c.mgr.SendComPort(message)
		time.Sleep(time.Millisecond * 10)
	}
	fmt.Println("<- pullAndSend")
}

/**
 * Reads data from Gemini, and sends (or translates) to the Orbiter.
 */
func (c *Controller) dataRead() () {
	fmt.Println("-> dataRead")
	s := c.mgr.ReadyComPort()
	if s != "" {
		// Finish reading from the port
		s += c.mgr.ReadComPort()
		for _, elem := range strings.Split(s, "\n") {
			c.handleIncoming(elem)
		}
	}
	fmt.Println("<- dataRead")
}

func (c *Controller) handleIncoming(s string) () {
	fmt.Println("-- Received: " + s)
	if len(s) > 4 && s[:4] == "TEST" {
		// Test command, return everything after "TEST "
		reply := s[5:len(s)]
		fmt.Println("-- Writing: " + reply)
		c.mgr.SendComPort(reply)
	} else if c.mgr.IsOrbiterConnected() {
		// These received strings may or may not be direct OrbConnect commands.
		// For now let's assume they are.
		//c.mgr.Send(s)
	}
}

/**
 * The controller loop essentially consists of three stages:
 *  1. Read state from the Orbiter program
 *  2. Read commands from Gemini
 *  3. Push Orbiter state to Gemini
 *
 * I've built it in this order because Golang serial port code doesn't appear to
 * leverage the communication state of the port itself, namely whether there is
 * data to read or write. By ordering the loop this way, with the read operation
 * first, followed by a pseudo-atomic write operation, we minimize the risk of
 * cross-transmission of data.
 *
 * However, we may end up with a command coming in from Gemini that should
 * modify the state of the system in the write-out. A major to-do is to write the
 * Orbiter state to a map of values, and then during post-read command processing
 * update any modified mapped values before the write back.
 */
func (c *Controller) Run() () {
	fmt.Println("-> Controller loop")
	for {
		time.Sleep(time.Second * 1)
		if c.exit == true {
			fmt.Println("Cleaning up...")
			c.mgr.Cleanup()
			fmt.Println("Exiting...")
			break
		}
		
		// Read existing state from Orbiter and submit to Gemini
		c.pullAndSend()
		// Process any incoming commands from Gemini
		c.dataRead()
	}
	fmt.Println("<- Controller loop")
}
