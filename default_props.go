package sni

import (
	"github.com/godbus/dbus/v5/prop"
	"github.com/knightpp/sni/sni"
)

func makePropsSni(id, title string) map[string]*prop.Prop {
	return map[string]*prop.Prop{
		"Category": {
			Value: sni.CategoryApplicationStatus,
		},
		"Id": {
			Value: id,
		},
		"Title": {
			Value: title,
		},
		"Status": {
			Value: sni.StatusActive,
		},
		"WindowId": {
			Value: 0,
		},
		"IconName": {
			Value: "face-cool",
		},
		"IconPixmap": {
			Value: Pixmap{},
		},
		"OverlayIconName": {
			Value: "",
		},
		"OverlayIconPixmap": {
			Value: Pixmap{},
		},
		"AttentionIconName": {
			Value: "",
		},
		"AttentionIconPixmap": {
			Value: Pixmap{},
		},
		"AttentionMovieName": {
			Value: "",
		},
		"ToolTip": {
			Value: ToolTip{},
		},
		"ItemIsMenu": {
			Value: false,
		},
		"Menu": {
			Value: MENU_PATH,
		},
		"IconThemePath": {
			Value: "",
		},
	}
}

func makePropsMenu() map[string]*prop.Prop {
	return map[string]*prop.Prop{
		// Provides the version of the DBusmenu API that this API is
		// implementing.
		"Version": {
			Value: uint32(3),
		},
		"TextDirection": {
			Value: "ltr",
		},
		"Status": {
			Value: "normal",
		},
		"IconThemePath": {
			Value: []string{},
		},
	}
}
