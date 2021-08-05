package sni

import (
	"context"
	"fmt"
	"image"
	"image/draw"
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
	menuServer d_bus_menu.Dbusmenuer
	// sniServer implements org.kde.StatusNotifierWatcher
	sniServer status_notifier_item.StatusNotifierItemer
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

func (t *Tray) SetSniServer(impl status_notifier_item.StatusNotifierItemer) *Tray {
	t.sniServer = impl
	return t
}

func (t *Tray) SetMenuServer(impl d_bus_menu.Dbusmenuer) *Tray {
	t.menuServer = impl
	return t
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
				if newOwner == "" {
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

func imageToArgb32(src image.Image) Pixmap {
	b := src.Bounds()
	width := b.Dx()
	height := b.Dy()
	m := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(m, m.Bounds(), src, b.Min, draw.Src)

	buf := m.Pix
	// _ = buf[len(buf)]
	// Swap RGBA to ARGB
	for i := 0; i < len(buf); i += 4 {
		buf[i], buf[i+3] = buf[i+3], buf[i]
		buf[i+1], buf[i+3] = buf[i+3], buf[i+1]
		buf[i+2], buf[i+3] = buf[i+3], buf[i+2]
	}
	return Pixmap{
		Width:  int32(width),
		Heigth: int32(height),
		Data:   buf,
	}
}
