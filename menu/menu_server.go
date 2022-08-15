package menu

import (
	"log"

	"github.com/knightpp/sni/generated/d_bus_menu"

	"github.com/godbus/dbus/v5"
)

func NewMenuServer(tree ItemTree) *MenuServer {
	idToItem := make(map[int32]*Item)
	tree.root.forEach(func(i *Item) {
		idToItem[i.id] = i
	})
	return &MenuServer{
		tree:     tree,
		idToItem: idToItem,
	}
}

type MenuServer struct {
	*d_bus_menu.UnimplementedDbusmenu
	tree     ItemTree
	idToItem map[int32]*Item
}

type Layout = struct {
	// V0 is id
	V0 int32
	// V1 is properties
	V1 map[string]dbus.Variant
	// V2 is submenu items
	V2 []dbus.Variant
}

func (m *MenuServer) GetLayout(
	parentId int32,
	recursionDepth int32,
	propertyNames []string,
) (revision uint32, layout Layout, err *dbus.Error) {
	layout = m.tree.ToLayout()
	log.Printf("GetLayout(parentId = %d, recursionDepth = %d,"+
		"propertyNames = %+v) return %+v",
		parentId, recursionDepth, propertyNames, layout)
	return
}

// GetGroupProperties is com.canonical.dbusmenu.GetGroupProperties method.
func (m *MenuServer) GetGroupProperties(ids []int32, propertyNames []string) (properties []struct {
	V0 int32
	V1 map[string]dbus.Variant
}, err *dbus.Error,
) {
	log.Printf("GetGroupProperties(ids = %+v, propertyNames = %+v)", ids, propertyNames)
	return
}

// GetProperty is com.canonical.dbusmenu.GetProperty method.
func (m *MenuServer) GetProperty(id int32, name string) (value dbus.Variant, err *dbus.Error) {
	log.Printf("GetProperty(id = %d, name = %s)", id, name)
	return
}

// Event is com.canonical.dbusmenu.Event method.
func (m *MenuServer) Event(id int32, eventId string, data dbus.Variant, timestamp uint32) (err *dbus.Error) {
	log.Printf("Event(id = %d, eventId = %s, data = %s, timestamp = %d)",
		id, eventId, data, timestamp)
	if eventId == "clicked" {
		item, ok := m.idToItem[id]
		if ok && item.onClick != nil {
			item.onClick()
		}
	}
	return
}

// EventGroup is com.canonical.dbusmenu.EventGroup method.
func (m *MenuServer) EventGroup(events []struct {
	V0 int32
	V1 string
	V2 dbus.Variant
	V3 uint32
},
) (idErrors []int32, err *dbus.Error) {
	log.Printf("EventGroup(events = %+v)", events)
	return
}

// AboutToShow is com.canonical.dbusmenu.AboutToShow method.
func (m *MenuServer) AboutToShow(id int32) (needUpdate bool, err *dbus.Error) {
	log.Printf("AboutToShow(id = %d)", id)
	return
}

// AboutToShowGroup is com.canonical.dbusmenu.AboutToShowGroup method.
func (m *MenuServer) AboutToShowGroup(ids []int32) (updatesNeeded, idErrors []int32, err *dbus.Error) {
	log.Printf("AboutToShowGroup(ids = %+v)", ids)
	return
}
