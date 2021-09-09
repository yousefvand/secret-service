package service

import (
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

// dbusSecretService creates dbusSecretService objects and interfaces on dbus
func dbusSecretService(service *Service) {

	////////////////////////////// Methods //////////////////////////////
	/*
		CreateSession ( IN String algorithm,
		                IN Variant input,
		                OUT Variant output,
		                OUT String serialnumber);
	*/
	createSession := []introspect.Arg{
		{
			Name:      "algorithm",
			Type:      "s",
			Direction: "in",
		},
		{
			Name:      "input",
			Type:      "v",
			Direction: "in",
		},
		{
			Name:      "output",
			Type:      "v",
			Direction: "out",
		},
		{
			Name:      "serialnumber",
			Type:      "s",
			Direction: "out",
		},
	}

	/*
		SetPassword ( IN  String      serialnumber
			            IN  Array<Byte> oldpassword,
									IN  Array<Byte> oldpassword_iv,
									IN  Array<Byte> newpassword,
									IN  Array<Byte> newpassword_iv,
									IN  Array<Byte> oldSalt,
									IN  Array<Byte> oldSalt_iv,
									IN  Array<Byte> newSalt,
									IN  Array<Byte> newSalt_iv
									OUT String result);
	*/
	setPassword := []introspect.Arg{
		{
			Name:      "serialnumber",
			Type:      "s",
			Direction: "in",
		},
		{
			Name:      "oldpassword",
			Type:      "ay",
			Direction: "in",
		},
		{
			Name:      "oldpassword_iv",
			Type:      "ay",
			Direction: "in",
		},
		{
			Name:      "newpassword",
			Type:      "ay",
			Direction: "in",
		},
		{
			Name:      "newpassword_iv",
			Type:      "ay",
			Direction: "in",
		},
		{
			Name:      "oldsalt",
			Type:      "ay",
			Direction: "in",
		},
		{
			Name:      "oldsalt_iv",
			Type:      "ay",
			Direction: "in",
		},
		{
			Name:      "newsalt",
			Type:      "ay",
			Direction: "in",
		},
		{
			Name:      "newsalt_iv",
			Type:      "ay",
			Direction: "in",
		},
		{
			Name:      "result",
			Type:      "s",
			Direction: "out",
		},
	}

	/*
		Login ( IN  String serialnumber,
		        IN  Array<Byte> passwordhash,
						IN  Array<Byte> passwordhash_iv,
						OUT Array<Byte> cookie,
						OUT Array<Byte> cookie_iv
						OUT String result);
	*/
	login := []introspect.Arg{
		{
			Name:      "serialnumber",
			Type:      "s",
			Direction: "in",
		},
		{
			Name:      "passwordhash",
			Type:      "ay",
			Direction: "in",
		},
		{
			Name:      "passwordhash_iv",
			Type:      "ay",
			Direction: "in",
		},
		{
			Name:      "cookie",
			Type:      "ay",
			Direction: "out",
		},
		{
			Name:      "cookie_iv",
			Type:      "ay",
			Direction: "out",
		},
		{
			Name:      "result",
			Type:      "s",
			Direction: "out",
		},
	}

	/*
		Command ( IN String command,
		          OUT String result);
	*/
	command := []introspect.Arg{
		{
			Name:      "serialnumber",
			Type:      "s",
			Direction: "in",
		},
		{
			Name:      "cookie",
			Type:      "ay",
			Direction: "in",
		},
		{
			Name:      "cookie_iv",
			Type:      "ay",
			Direction: "in",
		},
		{
			Name:      "command",
			Type:      "ay",
			Direction: "in",
		},
		{
			Name:      "command_iv",
			Type:      "ay",
			Direction: "in",
		},
		{
			Name:      "params",
			Type:      "ay",
			Direction: "in",
		},
		{
			Name:      "params_iv",
			Type:      "ay",
			Direction: "in",
		},
		{
			Name:      "result",
			Type:      "ay",
			Direction: "out",
		},
		{
			Name:      "result_iv",
			Type:      "ay",
			Direction: "out",
		},
	}

	////////////////////////////// Signals //////////////////////////////

	/*
		ServiceLocked;
	*/
	serviceLocked := []introspect.Arg{
		{
			Name: "collection",
		},
	}

	/*
		ServiceUnlocked;
	*/
	serviceUnlocked := []introspect.Arg{
		{
			Name: "collection",
		},
	}

	/////////////////////////////////// dbus ///////////////////////////////////

	introSecretService := &introspect.Node{
		Name: "/secretservice",
		Interfaces: []introspect.Interface{
			introspect.IntrospectData, prop.IntrospectData,
			{
				Name: "ir.remisa.SecretService", // interface name
				Methods: []introspect.Method{
					{
						Name: "CreateSession",
						Args: createSession,
					},
					{
						Name: "SetPassword",
						Args: setPassword,
					},
					{
						Name: "Login",
						Args: login,
					},
					{
						Name: "Command",
						Args: command,
					},
				},
				Signals: []introspect.Signal{
					{
						Name: "ServiceLocked",
						Args: serviceLocked,
					},
					{
						Name: "serviceUnlocked",
						Args: serviceUnlocked,
					},
				},
			},
		},
	}

	service.Connection.Export(service, "/secretservice",
		"ir.remisa.SecretService")

	service.Connection.Export(introspect.NewIntrospectable(introSecretService),
		"/secretservice", "org.freedesktop.DBus.Introspectable")

}
