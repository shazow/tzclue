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

	// Set DesktopID
	if call := client.Call("org.freedesktop.DBus.Properties.Set", 0, "org.freedesktop.GeoClue2.Client.DesktopId", appName); call.Err != nil {
		/*
			This errors with:
				call failed: Properties.Set: No such interface 'org.freedesktop.DBus.Properties' on object at path /org/freedesktop/GeoClue2/Client/1

			And yet,
			$ qdbus -r --system org.freedesktop.GeoClue2 /org/freedesktop/GeoClue2/Client/1
			method void org.freedesktop.DBus.Properties.Set(QString interface_name, QString property_name, QDBusVariant value)
			...
			property readwrite QString org.freedesktop.GeoClue2.Client.DesktopId
			$ qdbus --system org.freedesktop.GeoClue2 /org/freedesktop/GeoClue2/Client/1 "org.freedesktop.DBus.Properties.Set" org.freedesktop.GeoClue2.Client DesktopId "foo"
			Error: org.freedesktop.DBus.Error.AccessDenied
		*/
		exit(2, "call failed: Properties.Set: %s\n", call.Err)
	}

	// Set DistanceThreshold
	//client.Call("org.freedesktop.DBus.Properties.Set", 0, "org.freedesktop.GeoClue2.Client.DistanceThreshold", 0)
	/*
		ret_v = g_dbus_proxy_call_sync(
		geoclue_client,
		"org.freedesktop.DBus.Properties.Set",
		g_variant_new("(ssv)",
		"org.freedesktop.GeoClue2.Client",
		"DistanceThreshold",
		g_variant_new("u", 50000)),
	*/

	// Subscribe to signal
	/*
		const signalInterface = "org.freedesktop.GeoClue2.Client.LocationUpdated"
		if call := client.Call("org.freedesktop.DBus.AddMatch", 0, "type='signal',path='"+clientPath+"',interface='"+signalInterface+"'"); call.Err != nil {
			exit(2, "call failed: AddMatch: %s", call.Err)
		}
	*/

	if call := conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"type='signal',path='/org/freedesktop/DBus',interface='org.freedesktop.DBus',sender='org.freedesktop.DBus'"); call.Err != nil {
		exit(2, "call failed: AddMatch: %s", call.Err)

	}

	sig := make(chan *dbus.Signal, 10)
	conn.Signal(sig)

	log.Print("Subscribed to Location signals.")

	// Start client
	/*
		ret_v = g_dbus_proxy_call_sync(geoclue_client,
		"Start",
		NULL,
		G_DBUS_CALL_FLAGS_NONE,
		-1, NULL, &error);
	*/
	if call := client.Call("org.freedesktop.GeoClue2.Client.Start", 0); call.Err != nil {
		exit(2, "call failed: Client.Start: %s", call.Err)
	}

	log.Print("Client started")

	log.Print("Waiting for Location...")
	for v := range sig {
		fmt.Println(v)
	}

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
