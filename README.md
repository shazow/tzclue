# tzclue

tzclue = geoclue (via https://github.com/godbus/dbus) + https://github.com/bradfitz/latlong


## Notes:


- https://www.freedesktop.org/software/geoclue/docs/libgeoclue/GClueSimple.html
- https://www.freedesktop.org/software/geoclue/docs/gdbus-org.freedesktop.GeoClue2.Location.html

```console
$ qdbus --system org.freedesktop.GeoClue2 /org/freedesktop/GeoClue2/Client/1
signal void org.freedesktop.DBus.Properties.PropertiesChanged(QString interface_name, QVariantMap changed_properties, QStringList invalidated_properties)
method QDBusVariant org.freedesktop.DBus.Properties.Get(QString interface_name, QString property_name)
method QVariantMap org.freedesktop.DBus.Properties.GetAll(QString interface_name)
method void org.freedesktop.DBus.Properties.Set(QString interface_name, QString property_name, QDBusVariant value)
method QString org.freedesktop.DBus.Introspectable.Introspect()
method QString org.freedesktop.DBus.Peer.GetMachineId()
method void org.freedesktop.DBus.Peer.Ping()
property read bool org.freedesktop.GeoClue2.Client.Active
property readwrite QString org.freedesktop.GeoClue2.Client.DesktopId
property readwrite uint org.freedesktop.GeoClue2.Client.DistanceThreshold
property read QDBusObjectPath org.freedesktop.GeoClue2.Client.Location
property readwrite uint org.freedesktop.GeoClue2.Client.RequestedAccuracyLevel
property readwrite uint org.freedesktop.GeoClue2.Client.TimeThreshold
signal void org.freedesktop.GeoClue2.Client.LocationUpdated(QDBusObjectPath old, QDBusObjectPath new)
method void org.freedesktop.GeoClue2.Client.Start()
method void org.freedesktop.GeoClue2.Client.Stop()
```


```
org.freedesktop.GeoClue2.Location
- Latitude     readable   d
- Longitude    readable   d
- Accuracy     readable   d
- Altitude     readable   d
- Speed        readable   d
- Heading      readable   d
- Description  readable   s
- Timestamp    readable   (tt)
```

## License

MIT
