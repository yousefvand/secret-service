// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service

import "github.com/godbus/dbus/v5"

// OrgFreedesktopSecretErrorIsLocked
// "The object must be unlocked before this action can be carried out."
func ApiErrorIsLocked() *dbus.Error {
	return DbusErrorAccessDenied("The object must be unlocked before this action can be carried out")
}

// OrgFreedesktopSecretErrorNoSession
// "The session does not exist."
func ApiErrorNoSession() *dbus.Error {
	return DbusErrorUnknownObject("The session does not exist")
}

// OrgFreedesktopSecretErrorNoSuchObject
// "No such item or collection exists."
func ApiErrorNoSuchObject() *dbus.Error {
	return DbusErrorUnknownObject("No such item or collection exists")
}

// OrgFreedesktopDBusErrorNotSupported
// "Service does not support a specific set of algorithms for encryption."
func ApiErrorNotSupported() *dbus.Error {
	return DbusErrorNotSupported("Service does not support a specific set of algorithms for encryption")
}

////////////////////////////// D-Bus Errors //////////////////////////////

// DbusError is the low-level dbus error function
func DbusError(dbusError, message string) *dbus.Error {
	return &dbus.Error{
		Name: dbusError,
		Body: []interface{}{
			message,
		},
	}
}

// DbusErrorCallFailed means the call has failed
func DbusErrorCallFailed(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.Failed", message)
}

// DbusErrorNoMemory means system is out of memory
func DbusErrorNoMemory(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.NoMemory", message)
}

// DbusErrorServiceUnknown means the called service is not known
func DbusErrorServiceUnknown(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.ServiceUnknown", message)
}

// DbusErrorNoReply means the called method did not reply within the specified timeout
func DbusErrorNoReply(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.NoReply", message)
}

// DbusErrorBadAddress means the given address is not valid
func DbusErrorBadAddress(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.BadAddress", message)
}

// DbusErrorNotSupported means the call/operation is not supported
func DbusErrorNotSupported(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.NotSupported", message)
}

// DbusErrorLimitsExceeded means the limits allocated to this process/call/connection exceeded the pre-defined
func DbusErrorLimitsExceeded(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.LimitsExceeded", message)
}

// DbusErrorAccessDenied means the call/operation tried to access a resource it isn't allowed to
func DbusErrorAccessDenied(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.AccessDenied", message)
}

// DbusErrorNoServer means server is not listening on the address
func DbusErrorNoServer(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.NoServer", message)
}

// DbusErrorTimeout means operation has timed out
func DbusErrorTimeout(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.Timeout", message)
}

// DbusErrorNoNetwork means network is not available
func DbusErrorNoNetwork(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.NoNetwork", message)
}

// DbusErrorAddressInUse means D-Bus address is already taken
func DbusErrorAddressInUse(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.AddressInUse", message)
}

// DbusErrorDisconnected means connection is closed
func DbusErrorDisconnected(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.Disconnected", message)
}

// DbusErrorInvalidArgs means the arguments passed to this call/operation are not valid
func DbusErrorInvalidArgs(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.InvalidArgs", message)
}

// DbusErrorUnknownMethod means the method called was not found in this object/interface with the given parameters
func DbusErrorUnknownMethod(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.UnknownMethod", message)
}

// DbusErrorInvalidSignature means the type signature is not valid or compatible
func DbusErrorInvalidSignature(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.InvalidSignature", message)
}

// DbusErrorUnknownInterface means the interface is not known in this object
func DbusErrorUnknownInterface(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.UnknownInterface", message)
}

// DbusErrorUnknownObject means the object path points to an object that does not exist
func DbusErrorUnknownObject(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.UnknownObject", message)
}

// DbusErrorUnknownProperty means the property does not exist in this interface
func DbusErrorUnknownProperty(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.UnknownProperty", message)
}

// DbusErrorPropertyReadOnly means the property set failed because the property is read-only
func DbusErrorPropertyReadOnly(message string) *dbus.Error {
	return DbusError("org.freedesktop.DBus.Error.PropertyReadOnly", message)
}
