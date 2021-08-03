package sni

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/knightpp/sni/interfaces/d_bus"
	"github.com/knightpp/sni/interfaces/d_bus_menu"
	"github.com/knightpp/sni/interfaces/status_notifier_item"
	"github.com/knightpp/sni/interfaces/status_notifier_watcher"
	"github.com/knightpp/sni/menu"
	"github.com/knightpp/sni/sni"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

// instance increments every time when Tray::Setup method is called
var instance uint32

// Tray is trying to abstract tray functionality into one place.
//
// It's not safe to call methods from multiple goroutines simultaneously without
// synchronysation.
type Tray struct {
	// conn is connection to dbus session bus
	conn *dbus.Conn
	// propsSni maps StatusNotifierItem prop name to a value
	propsSni map[string]*prop.Prop
	// propsMenu maps dbusmenu prop name to a value
	propsMenu map[string]*prop.Prop
	// menuServer implements com.canonical.dbusmenu
	menuServer *menu.MenuServer
	// sniServer implements org.kde.StatusNotifierWatcher
	sniServer *sni.SniServer
}

// NewTray allocates new Tray. Note: this function doesn't communicate through
// dbus, to "start tray" you should call .Setup method.
//
// Returns error if it couldn't establish connection to dbus session bus.
func NewTray(id, title string, itemTree menu.ItemTree) (*Tray, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to session bus: %w", err)
	}
	return NewTrayWithConn(conn, id, title, itemTree), nil
}

// NewTray allocates new Tray. Note: this function doesn't communicate through
// dbus, to "start tray" you should call .Setup method
func NewTrayWithConn(conn *dbus.Conn, id, title string, itemTree menu.ItemTree) *Tray {
	return &Tray{
		conn:       conn,
		propsSni:   makePropsSni(id, title),
		propsMenu:  makePropsMenu(),
		menuServer: menu.NewMenuServer(itemTree),
		sniServer:  sni.NewSniServer(),
	}
}

// Close closes underlying dbus connection
func (t *Tray) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}

// Setup do necessary setups. Requests dbus name; exports servers; exports
// properties; registers with StatusNotifierWatcher; listens for
// OwnerNameChanged dbus signal etc.
func (t *Tray) Setup() error {
	inst := atomic.AddUint32(&instance, 1)
	name := NameBySpec(inst)
	reply, err := t.conn.RequestName(name,
		dbus.NameFlagReplaceExisting|dbus.NameFlagAllowReplacement)
	if err != nil {
		return err
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("name already taken")
	}
	err = status_notifier_item.ExportStatusNotifierItem(t.conn,
		SNI_PATH, t.sniServer)
	if err != nil {
		return err
	}
	err = d_bus_menu.ExportDbusmenu(t.conn, MENU_PATH, t.menuServer)
	if err != nil {
		return err
	}

	/*--------------- PROPS ---------------*/

	props := make(map[string]map[string]*prop.Prop)
	props[SNI_INTERFACE_NAME] = t.propsSni
	sniProp, err := prop.Export(t.conn, SNI_PATH, props)
	if err != nil {
		return err
	}
	props = make(map[string]map[string]*prop.Prop)
	props[DBUSMENU_INTERFACE_NAME] = t.propsMenu
	menuProp, err := prop.Export(t.conn, MENU_PATH, props)
	if err != nil {
		return err
	}
	/*--------------- END-PROPS ---------------*/
	/*--------------- INTROSPECTION ---------------*/
	sniNode := introspect.Node{
		Name: SNI_PATH,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			{
				Name: status_notifier_item.InterfaceStatusNotifierItem,
				Methods: introspect.Methods(
					status_notifier_item.StatusNotifierItemer(t.sniServer)),
				Properties: sniProp.Introspection(SNI_INTERFACE_NAME),
				Signals:    nil,
			},
		},
	}
	err = t.conn.Export(introspect.NewIntrospectable(&sniNode), SNI_PATH,
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		return err
	}

	menuNode := introspect.Node{
		Name: MENU_PATH,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			{
				Name: d_bus_menu.InterfaceDbusmenu,
				Methods: introspect.Methods(
					d_bus_menu.Dbusmenuer(t.menuServer),
				),
				Properties: menuProp.Introspection(DBUSMENU_INTERFACE_NAME),
				Signals:    nil,
			},
		},
	}
	err = t.conn.Export(introspect.NewIntrospectable(&menuNode), MENU_PATH,
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		return err
	}
	/*--------------- END-INTROSPECTION ---------------*/

	if err = register(t.conn, name); err != nil {
		return err
	}
	go func(name string) {
		if err := t.listen(name); err != nil {
			log.Print("NameOwnerChanged listener exitted with error: ", err)
		}

	}(name)
	return nil
}

// register registers service name with StatusNotifierWatcher
func register(conn *dbus.Conn, service string) error {
	ctx := context.Background()
	obj := conn.Object("org.kde.StatusNotifierWatcher",
		"/StatusNotifierWatcher")
	watcher := status_notifier_watcher.NewStatusNotifierWatcher(obj)
	err := watcher.RegisterStatusNotifierItem(ctx, service)
	return err
}

// SetId sets an id that should be unique for this application and consistent
// between sessions, such as the application name itself.
func (t *Tray) SetId(id string) *Tray {
	t.propsSni["Id"].Value = id
	return t
}

// SetTitle sets a name that describes the application, it can be more
// descriptive than Id.
func (t *Tray) SetTitle(title string) *Tray {
	t.propsSni["Title"].Value = title
	return t
}

// SetIconName sets StatusNotifierItem IconName property
func (t *Tray) SetIconName(name string) *Tray {
	t.propsSni["IconName"].Value = name
	return t
}

// GetIconName return StatusNotifierItem IconName property
func (t *Tray) GetIconName() string {
	v, ok := t.propsSni["IconName"]
	if !ok {
		panic("GetIconName(): no such value in a map")
	}
	s, ok := v.Value.(string)
	if !ok {
		panic("GetIconName(): value is not string")
	}
	return s
}

// SetSniStatus sets StatusNotifierItem Status property.
//
// It describes the status of this item or of the associated application.
//
// The allowed values for the Status property are:
//
// Passive: The item doesn't convey important information to the user, it
// can be considered an "idle" status and is likely that visualizations
// will chose to hide it.
//
// Active: The item is active, is more important that the item will be shown
// in some way to the user.
//
// NeedsAttention: The item carries really important information for the user,
// such as battery charge running out and is wants to incentive the direct user
// intervention. Visualizations should emphasize in some way the items with
// NeedsAttention status.
func (t *Tray) SetSniStatus(status sni.Status) *Tray {
	t.propsSni["Status"].Value = status
	return t
}

// Represents the way the text direction of the application. This
// allows the server to handle mismatches intelligently. For left-
// to-right the string is "ltr" for right-to-left it is "rtl".
func (t *Tray) SetMenuTextDirection(dir TextDirection) *Tray {
	t.propsMenu["TextDirection"].Value = dir
	return t
}

// Tells if the menus are in a normal state or they believe that they
// could use some attention. Cases for showing them would be if help
// were referring to them or they accessors were being highlighted.
// This property can have two values:
//
// - "normal" in almost all cases and
//
// - "notice" when they should have a higher priority to be shown.
func (t *Tray) SetMenuStatus(status MenuStatus) *Tray {
	t.propsMenu["Status"].Value = status
	return t
}

// A list of directories that should be used for finding icons using
// the icon naming spec. Ideally there should only be one for the icon
// theme, but additional ones are often added by applications for
// app specific icons.
func (t *Tray) SetMenuIconThemePath(path []string) *Tray {
	t.propsMenu["IconThemePath"].Value = path
	return t
}

// SignalNewIcon emits signal on dbus thus requesting re-rendering of its icon.
// You should emit this signal to reflect change of the icon visually.
func (t *Tray) SignalNewIcon() error {
	err := status_notifier_item.Emit(t.conn, &status_notifier_item.StatusNotifierItem_NewIconSignal{
		Path: SNI_PATH,
		Body: &status_notifier_item.StatusNotifierItem_NewIconSignalBody{},
	})
	return err
}

// listen blocks forever
func (t *Tray) listen(appName string) error {
	err := d_bus.AddMatchSignal(t.conn, &d_bus.DBus_NameOwnerChangedSignal{})
	if err != nil {
		return err
	}
	ch := make(chan *dbus.Signal)
	t.conn.Signal(ch)
	for sig := range ch {
		s, err := d_bus.LookupSignal(sig)
		if err != nil {
			return err
		}
		switch sig := s.(type) {
		case *d_bus.DBus_NameOwnerChangedSignal:
			name := sig.Body.V0
			// oldOwner := sig.Body.V1
			newOwner := sig.Body.V2
			if name == "org.kde.StatusNotifierWatcher" {
				if len(newOwner) == 0 {
					return fmt.Errorf("stop")
				} else {
					log.Printf("Registering !!!")
					if err := register(t.conn, appName); err != nil {
						return err
					}
				}
			}
		default:
			return fmt.Errorf("unknown signal: %+v", sig)
		}
	}
	return nil
}
