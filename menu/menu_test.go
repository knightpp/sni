package menu_test

import (
	"testing"

	"github.com/knightpp/sni/menu"

	"github.com/godbus/dbus/v5"
	"github.com/stretchr/testify/assert"
)

func TestSimpleItem(t *testing.T) {
	assert := assert.New(t)
	item := menu.NewItem().Submenu(
		menu.NewItem().Label("Test button"),
		menu.NewItem().Separator(true),
	)
	layout := item.Build().ToLayout()

	assert.Equal(layout.V0, int32(0))
	assert.Equal(layout.V1, map[string]dbus.Variant{
		"children-display": dbus.MakeVariant("submenu"),
	})
	assert.Equal(len(layout.V2), int(2))
	label := layout.V2[0]
	assert.Equal(label, dbus.MakeVariant(menu.Layout{
		V0: 1,
		V1: map[string]dbus.Variant{
			"label": dbus.MakeVariant("Test button"),
		},
		V2: nil,
	}))
	sep := layout.V2[1]
	assert.Equal(sep, dbus.MakeVariant(menu.Layout{
		V0: 2,
		V1: map[string]dbus.Variant{
			"type": dbus.MakeVariant("separator"),
		},
		V2: nil,
	}))
}
