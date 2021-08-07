package main

import (
	"github.com/godbus/dbus/v5/introspect"
)

var Node = introspect.Node{
	Name: "",
	Interfaces: []introspect.Interface{
		{
			Name: "",
			Methods: []introspect.Method{
				{
					Name: "",
					Args: []introspect.Arg{
						{
							Name:      "",
							Type:      "",
							Direction: "",
						},
					},
					Annotations: nil,
				},
			},
			Signals: []introspect.Signal{
				{
					Name:        "",
					Args:        []introspect.Arg{},
					Annotations: nil,
				},
			},
			Properties: []introspect.Property{
				{
					Name:        "",
					Type:        "",
					Access:      "",
					Annotations: nil,
				},
			},
			Annotations: []introspect.Annotation{
				{},
			},
		},
	},
}
