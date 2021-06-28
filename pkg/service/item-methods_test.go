// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service_test

import (
	"reflect"
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/service"
)

func TestItem_SetProperties(t *testing.T) {
	tests := []struct {
		name       string
		in         map[string]dbus.Variant
		out        map[string]dbus.Variant
		attributes map[string]string
	}{
		{
			name: "wrong entries",
			in: map[string]dbus.Variant{
				"correct":                              dbus.MakeVariant("No"),
				"org.freedesktop.Secret.Item.Correct":  dbus.MakeVariant("Yes"),
				"org.freedesktop.Secret.Items.Correct": dbus.MakeVariant("No"),
			},
			out: map[string]dbus.Variant{
				"Correct": dbus.MakeVariant("Yes"),
				"Label":   dbus.MakeVariant(""),
			},
			attributes: map[string]string{},
		},
		{
			name: "reserved keys",
			in: map[string]dbus.Variant{
				"correct":                              dbus.MakeVariant("No"),
				"org.freedesktop.Secret.Item.Correct":  dbus.MakeVariant("Yes"),
				"org.freedesktop.Secret.Item.Locked":   dbus.MakeVariant(true),
				"org.freedesktop.Secret.Item.Created":  dbus.MakeVariant(1234),
				"org.freedesktop.Secret.Item.Modified": dbus.MakeVariant(5678),
				"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
					" a ": "b ",
				}),
			},
			out: map[string]dbus.Variant{
				"Correct": dbus.MakeVariant("Yes"),
				"Label":   dbus.MakeVariant(""),
			},
			attributes: map[string]string{"a": "b"},
		},
		{
			name: "non-string Label",
			in: map[string]dbus.Variant{
				"correct":                           dbus.MakeVariant("No"),
				"org.freedesktop.Secret.Item.Label": dbus.MakeVariant(5),
				"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
					"a": "b",
					"c": "d",
				}),
			},
			out: map[string]dbus.Variant{
				"Label": dbus.MakeVariant(""),
			},
			attributes: map[string]string{
				"a": "b",
				"c": "d",
			},
		},
		{
			name: "normal",
			in: map[string]dbus.Variant{
				"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("MyItem"),
				"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
					"a":  "b",
					"c":  "d",
					"e ": " f",
				}),
			},
			out: map[string]dbus.Variant{
				"Label": dbus.MakeVariant("MyItem"),
			},
			attributes: map[string]string{
				"a": "b",
				"c": "d",
				"e": "f",
			},
		},
	}

	collection := service.NewCollection(Service)

	for _, tt := range tests {

		item := service.NewItem(collection)

		t.Run(tt.name, func(t *testing.T) {

			item.SetProperties(tt.in)
			if !reflect.DeepEqual(tt.out, item.Properties) {
				t.Errorf("Expected: %v, got: %v", tt.out, item.Properties)
			}
			if !reflect.DeepEqual(tt.attributes, item.LookupAttributes) {
				t.Errorf("Expected: %v, got: %v", tt.out, item.Properties)
			}
		})
	}
}
