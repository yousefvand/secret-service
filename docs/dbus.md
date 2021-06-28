# dbus

## CLI Tools

- qdbus: watch and send messages via D-Bus
- dbus-monitor: view a bus activity
- dbus-send: send messages via D-Bus

## Tests

- Check D-Bus services and look for the org.freedesktop.secrets:

```bash
dbus-send --session --dest=org.freedesktop.DBus --type=method_call --print-reply /org/freedesktop/DBus org.freedesktop.DBus.ListNames | grep 'org.freedesktop.secrets'
```

or

```bash
qdbus | grep 'secrets'
```

get the service description:

```bash
dbus-send --session --print-reply --dest=org.freedesktop.DBus /org/freedesktop/secrets org.freedesktop.DBus.Introspectable.Introspect
```

or

```bash
qdbus org.freedesktop.secrets /org/freedesktop/secrets
```

Who is listening on the another side of the bus?

```bash
ps -p $(qdbus --session org.freedesktop.DBus / org.freedesktop.DBus.GetConnectionUnixProcessID org.freedesktop.secrets) -o comm=
```

Steps done:

- A call to the org.freedesktop.DBus service
- passing the / path
- Called the org.freedesktop.DBus.GetConnectionUnixProcessID method
- and passing the service org.freedesktop.secrets name to the method to get its PID

## Store and retrieve data

To store data in secret service:

```bash
secret-tool store --label=Sample attribute1 value1 attribute2 value2
```

To retrieve stored data:

```bash
secret-tool lookup attribute1 value1
secret-tool lookup attribute2 value2
```
