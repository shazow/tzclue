package main

import (
	"fmt"
	"log"
	"os"

	"github.com/godbus/dbus"
)

func main() {
	conn, err := dbus.SystemBus()
	if err != nil {
		exit(1, "failed to connect to system bus: %s\n", err)
	}

	var clientPath dbus.ObjectPath
	err = conn.Object("org.freedesktop.GeoClue2", "/org/freedesktop/GeoClue2/Manager").Call("org.freedesktop.GeoClue2.Manager.GetClient", 0).Store(&clientPath)
	if err != nil {
		exit(1, "failed to get list of owned names: %s\n", err)
	}

	log.Printf("Connected to client: %q", clientPath)

	client := conn.Object("org.freedesktop.GeoClue2", clientPath)

	var location interface{}
	err = client.Call("org.freedesktop.GeoClue2.Client.Location", 0).Store(&location)
	if err != nil {
		exit(1, "failed to get location: %s\n", err)
	}

	log.Printf("Location: %q", location)
}

func exit(code int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(code)
}
