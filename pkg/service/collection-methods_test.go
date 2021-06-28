// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service_test

import (
	"reflect"
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/service"
)

func TestCollection_SetProperties(t *testing.T) {
	tests := []struct {
		name string
		in   map[string]dbus.Variant
		out  map[string]dbus.Variant
	}{
		{
			name: "wrong entries",
			in: map[string]dbus.Variant{
				"correct": dbus.MakeVariant("No"),
				"org.freedesktop.Secret.Collection.Correct":  dbus.MakeVariant("Yes"),
				"org.freedesktop.Secret.Collections.Correct": dbus.MakeVariant("No"),
			},
			out: map[string]dbus.Variant{
				"Label":   dbus.MakeVariant(""),
				"Correct": dbus.MakeVariant("Yes"),
			},
		},
		{
			name: "reserved keys",
			in: map[string]dbus.Variant{
				"correct": dbus.MakeVariant("No"),
				"org.freedesktop.Secret.Collection.Correct":  dbus.MakeVariant("Yes"),
				"org.freedesktop.Secret.Collection.Items":    dbus.MakeVariant([]string{"a", "b"}),
				"org.freedesktop.Secret.Collection.Locked":   dbus.MakeVariant(true),
				"org.freedesktop.Secret.Collection.Created":  dbus.MakeVariant(1234),
				"org.freedesktop.Secret.Collection.Modified": dbus.MakeVariant(5678),
			},
			out: map[string]dbus.Variant{
				"Label":   dbus.MakeVariant(""),
				"Correct": dbus.MakeVariant("Yes"),
			},
		},
		{
			name: "non-string Label",
			in: map[string]dbus.Variant{
				"correct": dbus.MakeVariant("No"),
				"org.freedesktop.Secret.Collection.Label": dbus.MakeVariant(5),
			},
			out: map[string]dbus.Variant{
				"Label": dbus.MakeVariant(""),
			},
		},
	}
	for _, tt := range tests {
		collection := service.NewCollection(Service)

		t.Run(tt.name, func(t *testing.T) {

			collection.SetProperties(tt.in)
			if !reflect.DeepEqual(tt.out, collection.Properties) {
				t.Errorf("Expected: %v, got: %v", tt.out, collection.Properties)
			}
		})
	}
}
