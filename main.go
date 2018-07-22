package main

import (
	"fmt"
	"log"
	"os"

	"github.com/godbus/dbus"
)

const appName = "tzclue"

func main() {
	conn, err := dbus.SystemBus()
	if err != nil {
		exit(1, "failed to connect to system bus: %s\n", err)
	}

	const geoclueInterface = "org.freedesktop.GeoClue2"

	var clientPath dbus.ObjectPath
	err = conn.Object(geoclueInterface, "/org/freedesktop/GeoClue2/Manager").Call("org.freedesktop.GeoClue2.Manager.GetClient", 0).Store(&clientPath)
	if err != nil {
		exit(1, "failed to get list of owned names: %s\n", err)
	}

	log.Printf("Connected to client: %q", clientPath)

	client := conn.Object(geoclueInterface, clientPath)

	// Setup client

	// Set DesktopId
	if call := client.Call("org.freedesktop.DBus.Properties.Set", 0, "org.freedesktop.GeoClue2.Client", "DesktopId", dbus.MakeVariant(appName)); call.Err != nil {
		exit(2, "call failed: Properties.Set(DesktopId): %s\n", call.Err)
	}

	// Set DistanceThreshold
	if call := client.Call("org.freedesktop.DBus.Properties.Set", 0, "org.freedesktop.GeoClue2.Client", "DistanceThreshold", dbus.MakeVariant(uint32(50000))); call.Err != nil {
		exit(2, "call failed: Properties.Set(DistanceThreshold): %s\n", call.Err)
	}

	if call := conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, "type='signal',interface='org.freedesktop.GeoClue2.Client',member='LocationUpdated'"); call.Err != nil {
		exit(2, "call failed: AddMatch: %s\n", call.Err)
	}

	sig := make(chan *dbus.Signal, 10)
	conn.Signal(sig)

	log.Print("Subscribed to Location signals.")

	// Start client
	// FIXME: Should we be expecting a reply?
	if call := client.Call("org.freedesktop.GeoClue2.Client.Start", dbus.FlagNoReplyExpected); call.Err != nil {
		exit(2, "call failed: Client.Start: %s\n", call.Err)
	}

	defer func() {
		client.Call("org.freedesktop.GeoClue2.Client.Stop", dbus.FlagNoReplyExpected)
	}()

	log.Print("Client started")

	log.Print("Waiting for Location...")
	for v := range sig {
		fmt.Println(v)
		break
	}

	/*
		location, err := client.GetProperty("org.freedesktop.GeoClue2.Client.Location")
		if err != nil {
			exit(1, "failed to get location: %s\n", err)
		}

		log.Printf("Location: %v", location.Value())
	*/
}

func exit(code int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(code)
}
