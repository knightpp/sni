package menu

import (
	"github.com/godbus/dbus/v5"
)

type ToggleType string
type Disposition string

const (
	ToggleTypeCheckmark = "checkmark"
	ToggleTypeRadio     = "radio"

	DispositionNormal      = "normal"
	DispositionInformative = "informative"
	DispositionWarning     = "warning"
	DispositionAlert       = "alert"
)

type ItemTree struct {
	root *Item
}

func (tree ItemTree) ToLayout() Layout {
	return tree.root.toLayout()
}

func (i *Item) toLayout() Layout {
	var layout Layout
	layout.V0 = i.id
	layout.V1 = i.properties
	for _, child := range i.children {
		layout.V2 = append(layout.V2,
			dbus.MakeVariant(child.toLayout()))
	}
	return layout

}
func (item *Item) forEach(fn func(*Item)) {
	fn(item)
	for _, child := range item.children {
		child.forEach(fn)
	}
}

type Item struct {
	id         int32
	onClick    func()
	children   []*Item
	properties map[string]dbus.Variant
}

func NewItem() *Item {
	m := make(map[string]dbus.Variant)
	return &Item{properties: m}
}

func (i *Item) Build() ItemTree {
	i.build(0)
	return ItemTree{root: i}
}
func (i *Item) build(id int32) {
	i.id = id
	for _, child := range i.children {
		id += 1
		child.build(id)
	}
}
func (i *Item) OnClick(fn func()) *Item {
	i.onClick = fn
	return i
}

func (i *Item) Separator(b bool) *Item {
	if b {
		i.properties["type"] = dbus.MakeVariant("separator")
	}
	return i
}
func (i *Item) Label(label string) *Item {
	i.properties["label"] = dbus.MakeVariant(label)
	return i
}
func (i *Item) CanBeActivated(b bool) *Item {
	i.properties["enabled"] = dbus.MakeVariant(b)
	return i
}
func (i *Item) Visible(b bool) *Item {
	i.properties["visible"] = dbus.MakeVariant(b)
	return i
}
func (i *Item) IconName(iconName string) *Item {
	i.properties["icon-name"] = dbus.MakeVariant(iconName)
	return i
}
func (i *Item) IconData(data []byte) *Item {
	i.properties["icon-data"] = dbus.MakeVariant(data)
	return i
}
func (i *Item) Shortcut(shortcut [][]string) *Item {
	i.properties["shortcut"] = dbus.MakeVariant(shortcut)
	return i
}
func (i *Item) ToggleType(tt ToggleType) *Item {
	i.properties["toggle-type"] = dbus.MakeVariant(tt)
	return i
}
func (i *Item) ToggleState(onoff bool) *Item {
	num := uint32(0)
	if onoff {
		num = 1
	}
	i.properties["toggle-state"] = dbus.MakeVariant(num)
	return i
}
func (i *Item) Submenu(children ...*Item) *Item {
	i.properties["children-display"] = dbus.MakeVariant("submenu")
	i.children = children
	return i
}
func (i *Item) Disposition(d Disposition) *Item {
	i.properties["disposition"] = dbus.MakeVariant(d)
	return i
}
